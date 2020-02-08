package main

import (
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	debug             = kingpin.Flag("debug", "Debug mode.").Bool()
	configFile        = kingpin.Flag("config.file", "path to the config file").Required().String()
	metricsListenAddr = kingpin.Flag("metrics.listen-addr", "listen address for metrics webserver").Default("127.0.0.1:9144").String()
)

func main() {
	config := &Config{}
	tailers := NewSafeTailers()

	kingpin.Parse()
	if *debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGHUP)

		for {
			select {
			case <-c:
				log.WithFields(log.Fields{"configFile": *configFile}).Info("SIGHUP received, " +
					"reloading config file")
				err := config.load(*configFile)
				if err != nil {
					log.Warnf("Unable to load configuration file: '%s'. Skipping configuration reload and "+
						"keeping current configuration", *configFile)
					break
				}

				err = runTailer(config, *configFile, tailers)
				if err != nil {
					log.Fatal(err)
				}

			}
		}
	}()

	err := config.load(*configFile)
	if err != nil {
		log.Fatal(err)
	}

	err = runTailer(config, *configFile, tailers)
	if err != nil {
		log.Fatal(err)
	}

	runServer(*metricsListenAddr)
}
