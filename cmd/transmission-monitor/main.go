package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/spf13/viper"
	"github.com/hekmon/transmissionrpc/v2"
)

var version = "undefined"

var (
	configFilename = flag.String("conf", "transmission-monitor.yaml", "path to configuration file")
	showVersion  = flag.Bool("version", false, "show program version and exit")
)

func run() int {
	flag.Parse()

	if *showVersion {
		fmt.Println(version)
		return 0
	}

	viper.SetConfigType("yaml")
	viper.SetConfigFile(*configFilename)
	setDefaults(viper.GetViper())
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("unable to read config file: %s", err)
	}

	trpc, err := transmissionrpc.New(
		viper.GetString("rpc.host"),
		viper.GetString("rpc.user"),
		viper.GetString("rpc.password"),
		&transmissionrpc.AdvancedConfig{
			HTTPS: viper.GetBool("rpc.https"),
			Port: uint16(viper.GetUint32("rpc.port")),
			RPCURI: viper.GetString("rpc.uri"),
			HTTPTimeout: viper.GetDuration("rpc.httptimeout"),
			UserAgent: viper.GetString("rpc.useragent"),
			Debug: viper.GetBool("rpc.debug"),
		},
	)
	if err != nil {
		log.Fatalf("unable to construct transmission RPC client: %v", err)
	}

	torrents, err := trpc.TorrentGetAll(context.Background())
	if err != nil {
		log.Fatalf("unable to get torrents: %v", err)
	}

	seenError := false
	for _, torrent := range torrents {
		if torrent.ErrorString != nil && *torrent.ErrorString != "" {
			fmt.Printf("torrent hash=%q, error=%q", *torrent.HashString, *torrent.ErrorString)
			seenError = true
		}
	}

	if seenError {
		return 3
	}

	return 0
}

func main() {
	log.Default().SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)
	os.Exit(run())
}
