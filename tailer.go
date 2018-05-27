package main

import (
	"fmt"

	"github.com/fstab/grok_exporter/tailer"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

func runTailer(lc LogConfig) {
	readall := false
	failOnMissing := false
	t := tailer.RunFseventFileTailer(lc.Path, readall, failOnMissing, &simpleLogger{})

	for {
		line := <-t.Lines()
		log.WithFields(log.Fields{"line": line, "path": lc.Path}).Debug("new line")
		handleLine(line, lc.Patterns)
	}
}

func handleLine(line string, patterns []*PatternConfig) error {
	for _, pc := range patterns {
		matches := pc.MatchCompiled.FindStringSubmatch(line)
		if matches == nil {
			continue
		}
		err := callLineSubHandler(pc, matches)
		if err != nil {
			log.WithFields(log.Fields{"err": err}).Warn("could not evaluate")
		}
		if !pc.Continue {
			return nil
		}
	}
	return nil
}

func callLineSubHandler(pc *PatternConfig, matches []string) error {
	val, err := pc.eval(matches)
	if err != nil {
		return err
	}
	if pc.CounterVec != nil {
		return handleLineCounter(pc, matches, val)
	}
	if pc.GaugeVec != nil {
		return handleLineGauge(pc, matches, val)
	}
	return fmt.Errorf("unknown type")
}

func buildLabels(pc *PatternConfig, matches []string) prometheus.Labels {
	labels := prometheus.Labels{}
	// we are ignoring errors here as the only error case is checked during startup:
	evalLabels, _ := pc.getEvaluatedLabels(matches)
	for label, value := range evalLabels {
		labels[label] = value
	}
	return labels
}

func handleLineCounter(pc *PatternConfig, matches []string, val float64) error {
	labels := buildLabels(pc, matches)
	c, err := pc.CounterVec.GetMetricWith(labels)
	if err != nil {
		return fmt.Errorf("failed to get metric: %s", err)
	}
	if pc.Action == "inc" {
		c.Add(val)
		return nil
	}
	return fmt.Errorf("unknown action '%s'", pc.Action)
}

func handleLineGauge(pc *PatternConfig, matches []string, val float64) error {
	labels := buildLabels(pc, matches)
	c, err := pc.GaugeVec.GetMetricWith(labels)
	if err != nil {
		return fmt.Errorf("failed to get metric: %s", err)
	}
	if pc.Action == "inc" {
		c.Add(val)
		return nil
	}
	if pc.Action == "dec" {
		c.Sub(val)
		return nil
	}
	if pc.Action == "set" {
		c.Set(val)
		return nil
	}
	return fmt.Errorf("unknown action '%s'", pc.Action)
}
