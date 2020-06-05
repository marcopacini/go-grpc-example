package main

import (
	"context"
	"fmt"
	"go-grpc-example/cache"
	"log"
	"os"

	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
)

const address = ":50051"

func main() {
	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:    "put",
				Aliases: []string{"p"},
				Usage:   "put an entry into the cache",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "key",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "value",
						Required: true,
					},
				},
				Action: func(c *cli.Context) error {
					conn, err := grpc.Dial(address, grpc.WithInsecure())
					if err != nil {
						return err
					}
					defer conn.Close()

					client := cache.NewCacheClient(conn)

					e := &cache.Entry{
						Key:   c.String("key"),
						Value: c.String("value"),
					}

					v, err := client.Put(context.Background(), e)
					if err != nil {
						return err
					}
					fmt.Println(v.Value)

					return nil
				},
			},
			{
				Name:    "get",
				Aliases: []string{"g"},
				Usage:   "get an entry from the cache",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "key",
						Required: true,
					},
				},
				Action: func(c *cli.Context) error {
					conn, err := grpc.Dial(address, grpc.WithInsecure())
					if err != nil {
						return err
					}
					defer conn.Close()

					client := cache.NewCacheClient(conn)

					key := &cache.Key{
						Key: c.String("key"),
					}
					v, err := client.Get(context.Background(), key)
					if err != nil {
						return err
					}
					fmt.Println(v.Value)

					return nil
				},
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
