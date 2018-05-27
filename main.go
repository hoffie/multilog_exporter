package main

import (
	log "github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	debug             = kingpin.Flag("debug", "Debug mode.").Bool()
	configFile        = kingpin.Flag("config.file", "path to the config file").Required().String()
	metricsListenAddr = kingpin.Flag("metrics.listen-addr", "listen address for metrics webserver").Default("127.0.0.1:9144").String()
)

var config *Config

func main() {
	kingpin.Parse()
	if *debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}
	log.WithFields(log.Fields{"configFile": *configFile}).Info("Loading config file")
	config, err := loadConfigFile(*configFile)
	if err != nil {
		log.Fatalf("failed to load config: %s", err)
	}

	if len(config.Logs) < 1 {
		log.Fatal("No Logs configured, cannot do anything")
	}

	for _, lc := range config.Logs {
		log.WithFields(log.Fields{"path": lc.Path}).Info("Registering relevant metrics")
		err := registerMetrics(lc.Patterns)
		if err != nil {
			log.Fatalf("failed to register metrics: %s", err)
		}
		log.WithFields(log.Fields{"path": lc.Path}).Info("Starting tailer")
		go runTailer(lc)
	}

	runServer(*metricsListenAddr)
}
