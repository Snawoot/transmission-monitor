package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/spf13/viper"
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


	return 0
}

func main() {
	log.Default().SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)
	os.Exit(run())
}
