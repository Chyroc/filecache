package filecache

import (
	"bytes"
	"encoding/binary"
	"errors"
	"os"
	"time"

	"github.com/huichen/murmur"

	mmap "github.com/Chyroc/filecache/internal/gommap"
)

type Cache interface {
	Get(key string) (string, error)
	Set(key, val string, ttl time.Duration) error
	TTL(key string) (time.Duration, error)
}

func unixMs(ttl time.Duration) int64 {
	return int64(time.Now().Add(ttl).UnixNano() / int64(1000000))
}

func binaryInt(buf []byte) (int, error) {
	x, n := binary.Varint(buf)
	if n == 0 {
		return 0, errors.New("buf too small")
	} else if n < 0 {
		return 0, errors.New("value larger than 64 bits (overflow) and -n is the number of bytes read")
	}

	return int(x), nil
}

var NotFound = errors.New("not found")
var HashConflict = errors.New("hash conflict")
var KeyTooShort = errors.New("key too short")
var KeyTooLong = errors.New("key too long")
var ValueTooShort = errors.New("value too short")
var ValueTooLong = errors.New("value too long")

const entryCount = 512 // mod
const docLength = 1280
const docCount = 8
const docHeaderLength = 1 + 2 + 2 + 7
const minFileSize = 5242880

const MaxLengthKey = 244
const MaxLengthValue = 1024

// 文件初始大小是5M，空间不够就扩大
// 5M大小分成512个entry，4096个doc，每个doc大小是1280B，1个entry有8个doc
// doc的结构是 flag(1), key_len(2), val_len(2), ttl(7,13ms), key, val (k+v: 1268)

func New(filepath string) Cache {
	c := &CacheImpl{
		filepath:    filepath,
		CurrentSize: minFileSize, // B
	}

	c.file, c.err = os.OpenFile(filepath, os.O_CREATE|os.O_RDWR, 0600)
	if c.err != nil {
		return c
	}

	info, err := c.file.Stat()
	c.err = err
	if err != nil {
		return c
	}

	fill := make([]byte, c.CurrentSize-int(info.Size()))
	for i := 0; i < c.CurrentSize-int(info.Size()); i++ {
		fill[i] = 0
	}

	if _, c.err = c.file.WriteAt(fill, info.Size()); c.err != nil {
		return c
	}

	//c.mmap, c.err = mmap.Map(c.file, 0, mmap.RDWR)
	c.mmap, c.err = mmap.Map(c.file)
	return c
}

type CacheImpl struct {
	err         error
	filepath    string
	file        *os.File
	CurrentSize int
	mmap        mmap.MMap
}

func (r *CacheImpl) region(key string) int {
	return int(murmur.Murmur3([]byte(key)) % uint32(entryCount))
}

func (r *CacheImpl) Get(key string) (string, error) {
	keyLen := len(key)

	if r.err != nil {
		return "", r.err
	} else if keyLen > MaxLengthKey {
		return "", KeyTooLong
	} else if keyLen == 0 {
		return "", KeyTooShort
	}

	region := r.region(key) // 0 ~ mod-1
	regionOffset := region * docCount * docLength
	keyBytes := []byte(key)

	for i := 0; i < docCount; i++ {
		currentOffset := regionOffset + docLength*i
		if r.mmap[currentOffset] == 1 {
			// 当前有数据，判断key是否和给定的key重合
			keyLen, err := binaryInt(r.mmap[currentOffset+1 : currentOffset+3])
			if err != nil {
				return "", err
			}
			keyBytesFromMM := r.mmap[currentOffset+docHeaderLength : currentOffset+docHeaderLength+keyLen]
			if !bytes.Equal(keyBytes, keyBytesFromMM) {
				continue
			}

			expiredTimestampMS, err := binaryInt(r.mmap[currentOffset+5 : currentOffset+docHeaderLength])
			if int(time.Now().UnixNano()/int64(1000000)) > expiredTimestampMS {
				// 过期了
				// TODO: delete
				return "", NotFound
			}

			valLen, err := binaryInt(r.mmap[currentOffset+3 : currentOffset+5])
			if err != nil {
				return "", err
			}
			return string(r.mmap[currentOffset+docHeaderLength+keyLen : currentOffset+docHeaderLength+keyLen+valLen]), nil
		} else {
			// 全部flag为1，但是没有key重合的
			// 或者flag为0
		}
	}

	return "", NotFound
}

func (r *CacheImpl) Set(key, val string, ttl time.Duration) error {
	keyLen := len(key)
	valLen := len(val)

	if r.err != nil {
		return r.err
	} else if keyLen > MaxLengthKey {
		return KeyTooLong
	} else if len(val) > MaxLengthValue {
		return ValueTooLong
	} else if keyLen == 0 {
		return KeyTooShort
	} else if len(val) == 0 {
		return ValueTooShort
	}

	region := r.region(key) // 0 ~ 511
	regionOffset := region * docCount * docLength
	keyBytes := []byte(key)

	offset := -1
	for i := 0; i < docCount; i++ {
		currentOffset := regionOffset + docLength*i
		if r.mmap[currentOffset] == 1 {
			// 当前有数据，判断key是否和给定的key重合
			keyLen, err := binaryInt(r.mmap[currentOffset+1 : currentOffset+3])
			if err != nil {
				return err
			}
			keyBytesFromMM := r.mmap[currentOffset+docHeaderLength : currentOffset+docHeaderLength+keyLen]
			if bytes.Equal(keyBytes, keyBytesFromMM) {
				offset = currentOffset
				break
			}
		} else if offset < 0 { // r.mmap[currentOffset] == 0
			// 将第一个遇见的0doc给offset
			offset = currentOffset
		} else {
			// 全部flag为1，但是没有key重合的
		}
	}

	if offset < 0 {
		return HashConflict
	}

	docLen := docHeaderLength + keyLen + valLen

	buf := make([]byte, docLen) // TODO: use sync.Pool
	buf[0] = 1
	binary.PutVarint(buf[1:3], int64(keyLen))
	binary.PutVarint(buf[3:5], int64(valLen))
	binary.PutVarint(buf[5:docHeaderLength], unixMs(ttl))
	copy(buf[docHeaderLength:docHeaderLength+keyLen], key)
	copy(buf[docHeaderLength+keyLen:docLen], val)

	copy(r.mmap[offset:offset+docLen], buf[:])

	return nil
}

func (r *CacheImpl) TTL(key string) (time.Duration, error) {
	if r.err != nil {
		return 0, r.err
	}
	panic("implement me")
}
