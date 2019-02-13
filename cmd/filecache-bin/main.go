package main

import (
	"fmt"
	"github.com/Chyroc/filecache"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/urfave/cli"
)

func cmdGet() cli.Command {
	var file string
	return cli.Command{
		Name:        "get",
		Description: "get from filecache file",
		Usage:       "filecache-bin get <key>",
		Action: func(c *cli.Context) error {
			if len(c.Args()) != 1 {
				return fmt.Errorf("invalid params count")
			} else if file == "" {
				return fmt.Errorf("invalid file path")
			}
			val, err := filecache.New(file).Get(c.Args()[0])
			if err != nil {
				return err
			}
			fmt.Printf("%q\n", val)
			return nil
		},
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:        "f",
				Destination: &file,
			},
		},
	}
}

func cmdSet() cli.Command {
	var file string
	return cli.Command{
		Name:        "set",
		Description: "set k-v to filecache file",
		Usage:       "filecache-bin set <key> <val> <ttl seconds>",
		Action: func(c *cli.Context) error {
			if len(c.Args()) != 3 {
				return fmt.Errorf("invalid params count")
			} else if file == "" {
				return fmt.Errorf("invalid file path")
			}
			ttl, err := strconv.Atoi(c.Args()[2])
			if err != nil {
				return fmt.Errorf("invalid ttl seconds param")
			}
			if err := filecache.New(file).Set(c.Args()[0], c.Args()[1], time.Duration(ttl)*time.Second); err != nil {
				return err
			}
			fmt.Println("OK")
			return nil
		},
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:        "f",
				Destination: &file,
			},
		},
	}
}

func cmdTTL() cli.Command {
	var file string
	return cli.Command{
		Name:        "set",
		Description: "get ttl from filecache file",
		Usage:       "filecache-bin get <key>",
		Action: func(c *cli.Context) error {
			if len(c.Args()) != 1 {
				return fmt.Errorf("invalid params count")
			} else if file == "" {
				return fmt.Errorf("invalid file path")
			}
			ttl, err := filecache.New(file).TTL(c.Args()[0])
			if err != nil {
				return err
			}
			fmt.Println("got:", ttl)
			return nil
		},
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:        "f",
				Destination: &file,
			},
		},
	}
}

func cmdDel() cli.Command {
	var file string
	return cli.Command{
		Name:        "del",
		Description: "del from filecache file",
		Usage:       "filecache-bin del <key>",
		Action: func(c *cli.Context) error {
			if len(c.Args()) != 1 {
				return fmt.Errorf("invalid params count")
			} else if file == "" {
				return fmt.Errorf("invalid file path")
			}

			if err := filecache.New(file).Del(c.Args()[0]); err != nil {
				return err
			}
			fmt.Println("done")
			return nil
		},
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:        "f",
				Destination: &file,
			},
		},
	}
}

func cmdRange() cli.Command {
	var file string
	return cli.Command{
		Name:        "range",
		Description: "get all vals from filecache file",
		Usage:       "filecache-bin range",
		Action: func(c *cli.Context) error {
			if file == "" {
				return fmt.Errorf("invalid file path")
			}

			kvs, err := filecache.New(file).Range()
			if err != nil {
				return err
			}
			for idx, v := range kvs {
				fmt.Printf(` %d) "%s:%s"`+"\n", idx+1, v.Key, v.Val)
			}
			return nil
		},
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:        "f",
				Destination: &file,
			},
		},
	}
}

func main() {
	app := cli.NewApp()
	app.Name = "filecache client"
	app.Action = func(c *cli.Context) error {
		return cli.ShowAppHelp(c)
	}
	app.Commands = []cli.Command{
		cmdGet(),
		cmdSet(),
		cmdTTL(),
		cmdDel(),
		cmdRange(),
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
