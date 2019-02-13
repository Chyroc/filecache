# filecache

cache library that store data in local file.

## usage

```go
func Example_Cache() {
	cache := filecache.New("./cache.data")
	defer os.Remove("./cache.data")

	_, err := cache.Get("not-exist")
	fmt.Println(err)

	fmt.Println(cache.Set("k", "v", time.Minute))

	v, err := cache.Get("k")
	fmt.Println(v, err)

	ttl, err := cache.TTL("k")
	fmt.Println(int(math.Ceil(ttl.Seconds())), err)

	time.Sleep(time.Second)

	ttl, err = cache.TTL("k")
	fmt.Println(int(math.Ceil(ttl.Seconds())), err)

	fmt.Println(cache.Del("k"))

	_, err = cache.Get("k")
	fmt.Println(err)

	// output:
	// not found
	// <nil>
	// v <nil>
	// 60 <nil>
	// 59 <nil>
	// <nil>
	// not found
}
```

## benchmark

```
chyroc			"github.com/Chyroc/filecache"
dannyBen		"github.com/DannyBen/filecache"
fabiorphp		"github.com/fabiorphp/cachego"
gadelkareem		"github.com/gadelkareem/cachita"
gookit			"github.com/gookit/cache"
huntsman		"github.com/huntsman-li/go-cache"
miguelmota		"github.com/miguelmota/go-filecache"
```

```
pkg: github.com/Chyroc/filecache-benchmark
BenchmarkFileCache/chyroc-8         	     200	   7135164 ns/op
BenchmarkFileCache/huntsman-8       	       2	 537270751 ns/op
BenchmarkFileCache/fabiorphp-8      	       3	 405211720 ns/op
BenchmarkFileCache/dannyBen-8       	       3	 385529791 ns/op
BenchmarkFileCache/gadelkareem-8    	       3	 403604442 ns/op
BenchmarkFileCache/miguelmota-8     	       1	3222506519 ns/op
```