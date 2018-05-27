package main

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

var metricsByName = make(map[string]prometheus.Collector)

func registerMetrics(pcs []*PatternConfig) error {
	for idx, pc := range pcs {
		log.WithFields(log.Fields{"metric": pc.Metric}).Info("Registering pattern")
		err := pc.compile()
		if err != nil {
			return fmt.Errorf("pattern %d/%d[%s]: %s'", idx+1, len(pcs), pc.Metric, err)
		}
		switch pc.Type {
		case "counter":
			err := registerCounter(pc)
			if err != nil {
				return err
			}
		case "gauge":
			err := registerGauge(pc)
			if err != nil {
				return err
			}
		default:
			return fmt.Errorf("unknown type")
		}
	}
	return nil
}

func registerCounter(pc *PatternConfig) error {
	opts := prometheus.CounterOpts{
		Name: pc.Metric,
		Help: pc.Help,
	}
	labelNames := []string{}
	for name := range pc.Labels {
		labelNames = append(labelNames, name)
	}
	c := prometheus.NewCounterVec(opts, labelNames)
	if len(pc.Labels) == 0 {
		_ = c.WithLabelValues()
	}
	pc.CounterVec = c
	err := prometheus.Register(pc.CounterVec)
	if are, ok := err.(prometheus.AlreadyRegisteredError); ok {
		// A counter for that metric has been registered before.
		// Use the old counter from now on.
		pc.CounterVec = are.ExistingCollector.(*prometheus.CounterVec)
		log.Debug("counter already registered, using existing instance")
	} else {
		// Something else went wrong!
		return err
	}
	return nil
}

func registerGauge(pc *PatternConfig) error {
	opts := prometheus.GaugeOpts{
		Name: pc.Metric,
		Help: pc.Help,
	}
	labelNames := []string{}
	for name := range pc.Labels {
		labelNames = append(labelNames, name)
	}
	pc.GaugeVec = prometheus.NewGaugeVec(opts, labelNames)
	err := prometheus.Register(pc.GaugeVec)
	if are, ok := err.(prometheus.AlreadyRegisteredError); ok {
		// A counter for that metric has been registered before.
		// Use the old counter from now on.
		pc.GaugeVec = are.ExistingCollector.(*prometheus.GaugeVec)
		log.Debug("gauge already registered, using existing instance")
	} else {
		// Something else went wrong!
		return err
	}
	return nil
}
