package main

import (
	"encoding/json"
	"fmt"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"time"
)

type StorePacket struct {
	Host   string `json:"host"`
	PubKey string `json:"pub_key"`
}

func main() {
	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:      "send",
				Aliases:   []string{"s"},
				ArgsUsage: "[host] [path_to_public_key]",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "host",
						Usage:    "Host of this machine for other machines to communicate",
						EnvVars:  []string{"HOST"},
						Required: true,
					},
					&cli.StringFlag{
						Name:    "ip",
						Aliases: []string{"i"},
						Usage:   "IP Address to send to",
					},
					&cli.StringFlag{
						Name:    "port",
						Aliases: []string{"p"},
						Usage:   "Port to send to",
						Value:   "8888",
					},
					&cli.StringFlag{
						Name:    "store-port",
						Aliases: []string{"s"},
						Usage:   "Port of the Store",
						Value:   "5000",
					},
					&cli.StringFlag{
						Name:    "send-port",
						Aliases: []string{"s"},
						Usage:   "Port of the Store",
						Value:   "50505",
					},
					&cli.IntFlag{
						Name:    "delay",
						Aliases: []string{"d"},
						Usage:   "Port to send to",
						Value:   0,
					},
				},
				Usage: "Send UDP packet to the IP about our own nix-store information. Reads IP address from STDIN",
				Action: func(cCtx *cli.Context) error {

					delay := cCtx.Int("delay")
					time.Sleep(time.Duration(delay) * time.Second)

					// pub key
					publicKeyPath := cCtx.Args().Get(0)
					if publicKeyPath == "" {
						return cli.Exit("path_to_public_key is required", 1)
					}
					body, err := os.ReadFile(publicKeyPath)
					if err != nil {
						return cli.Exit("path_to_public_key cannot be found", 1)
					}
					pubKey := string(body)

					// stores
					host := cCtx.String("host")
					storePort := cCtx.String("store-port")
					storePath := fmt.Sprintf("http://%s:%s", host, storePort)

					packet := StorePacket{
						Host:   storePath,
						PubKey: pubKey,
					}

					// json encode
					b, err := json.Marshal(packet)
					if err != nil {
						return cli.Exit(err.Error(), 1)
					}

					// Send address
					ip := cCtx.String("ip")
					if ip == "" {
						_, errScan := fmt.Scanln(&ip)
						if errScan != nil {
							return cli.Exit("Cannot read IP from stdin", 1)
						}
					}
					port := cCtx.String("port")
					sendPort := cCtx.String("send-port")
					send(ip, sendPort, port, string(b))
					return nil
				},
			},
			{
				Name:    "receive",
				Aliases: []string{"r"},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "port",
						Aliases: []string{"p"},
						Usage:   "Port to send to",
						Value:   "8888",
					},
					&cli.StringFlag{
						Name:    "cache",
						Aliases: []string{"c"},
						Usage:   "Cache Path",
						Value:   "/substituters-track.json",
					},
					&cli.IntFlag{
						Name:    "timeout",
						Aliases: []string{"t"},
						Usage:   "Timeout in seconds before a store is considered dead",
						Value:   30,
					},
					&cli.IntFlag{
						Name:  "poll",
						Usage: "Delay between each iteration to clear outdated Stores",
						Value: 5,
					},
				},
				Usage: "start a server to listen UDP packets",
				Action: func(cCtx *cli.Context) error {
					port := cCtx.String("port")
					cachePath := cCtx.String("cache")
					timeout := cCtx.Int("timeout")
					poll := cCtx.Int("poll")

					c := make(chan int, 1)
					c <- 0
					go func(timeout, poll int, cachePath string, c chan int) {
						fmt.Println("Monitoring thread started...")
						for {
							// sleep for poll time
							time.Sleep(time.Duration(poll) * time.Second)
							cleanOutdated(cachePath, time.Now().Unix(), timeout, c)
						}
					}(timeout, poll, cachePath, c)
					fmt.Println("Server Started...")
					for {

						receive(port, cachePath, timeout, c)
					}
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
