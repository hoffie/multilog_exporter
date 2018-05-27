package main

import log "github.com/sirupsen/logrus"

type simpleLogger struct{}

func (sl simpleLogger) Debug(msg string, params ...interface{}) {
	log.Debugf(msg, params...)
}
