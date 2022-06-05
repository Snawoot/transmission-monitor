package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/hekmon/transmissionrpc/v2"
	"github.com/spf13/viper"

	"github.com/Snawoot/transmission-monitor/db"
	"github.com/Snawoot/transmission-monitor/health"
	"github.com/Snawoot/transmission-monitor/monitor"
	"github.com/Snawoot/transmission-monitor/notifier"
)

var version = "undefined"

var (
	configFilename = flag.String("conf", "transmission-monitor.yaml", "path to configuration file")
	showVersion    = flag.Bool("version", false, "show program version and exit")
	clearDB        = flag.Bool("clear-db", false, "clear database")
	clearKey       = flag.String("clear-key", "", "delete specified hash from database")
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

	dbPath := viper.GetString("db.path")
	ensureDir(dbPath)
	dbInstance, err := db.Open(dbPath)
	defer dbInstance.Close()

	if *clearDB {
		err := dbInstance.Clear()
		if err != nil {
			log.Fatalf("failed to clear database: %v", err)
		}
		return 0
	}

	if *clearKey != "" {
		err := dbInstance.Delete(*clearKey)
		if err != nil {
			log.Fatalf("failed to remove key from database: %v", err)
		}
		return 0
	}

	trpc, err := transmissionrpc.New(
		viper.GetString("rpc.host"),
		viper.GetString("rpc.user"),
		viper.GetString("rpc.password"),
		&transmissionrpc.AdvancedConfig{
			HTTPS:       viper.GetBool("rpc.https"),
			Port:        uint16(viper.GetUint32("rpc.port")),
			RPCURI:      viper.GetString("rpc.uri"),
			HTTPTimeout: viper.GetDuration("rpc.httptimeout"),
			UserAgent:   viper.GetString("rpc.useragent"),
			Debug:       viper.GetBool("rpc.debug"),
		},
	)
	if err != nil {
		log.Fatalf("unable to construct transmission RPC client: %v", err)
	}

	torrents, err := trpc.TorrentGetAll(context.Background())
	if err != nil {
		log.Fatalf("unable to get torrents: %v", err)
	}

	t := make([]*transmissionrpc.Torrent, len(torrents))
	for i := range torrents {
		t[i] = &torrents[i]
	}

	checker := health.NewErrorCheck()
	var notifier monitor.Notifier = notifier.NewLogNotifier()
	mon := monitor.NewMonitor(dbInstance, checker, notifier)
	if err := mon.Process(t); err != nil {
		log.Fatalf("monitor returned error: %v", err)
	}

	return 0
}

func ensureDir(path string) {
	if err := os.MkdirAll(path, 0700); err != nil {
		log.Fatalf("failed to create database directory: %v", err)
	}
}

func main() {
	log.Default().SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)
	os.Exit(run())
}
