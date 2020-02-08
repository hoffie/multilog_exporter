package main

import (
	"fmt"
	"sync"

	"github.com/fstab/grok_exporter/tailer/fswatcher"
	"github.com/fstab/grok_exporter/tailer/glob"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

type SafeTailers struct {
	sync.Mutex
	Tailers []fswatcher.FileTailer
}

func NewSafeTailers() *SafeTailers {
	return &SafeTailers{}
}

func (tg *SafeTailers) start(lcs []LogConfig) error {
	tg.Lock()
	defer tg.Unlock()
	for _, lc := range lcs {
		log.WithFields(log.Fields{"path": lc.Path}).Info("Registering relevant metrics")
		err := registerMetrics(lc.Patterns)
		if err != nil {
			return fmt.Errorf("failed to register metrics: %s", err)
		}
		log.WithFields(log.Fields{"path": lc.Path}).Info("Starting tailer")
		readall := false
		failOnMissing := false
		t, err := fswatcher.RunFileTailer([]glob.Glob{glob.Glob(lc.Path)}, readall, failOnMissing, log.New())
		if err != nil {
			return fmt.Errorf("failed to start file tailer for path %s: %s", lc.Path, err)
		}
		tg.Tailers = append(tg.Tailers, t)

		go waitForTailerLines(t, lc)
	}
	return nil
}

func (tg *SafeTailers) stop() {
	tg.Lock()
	defer tg.Unlock()
	for _, tailer := range tg.Tailers {
		tailer.Close()
	}
	tg.Tailers = nil
}

func runTailer(config *Config, configPath string, tailers *SafeTailers) error {
	if tailers.Tailers != nil {
		log.WithFields(log.Fields{"configFile": configPath}).Debug("Stopping current tailers...")
		tailers.stop()
	}

	err := tailers.start(config.Logs)
	if err != nil {
		return err
	}

	return nil
}

func waitForTailerLines(t fswatcher.FileTailer, lc LogConfig) {
	for {
		line, open := <-t.Lines()
		if !open {
			break
		}
		log.WithFields(log.Fields{"line": line.Line, "path": lc.Path}).Debug("new line")
		handleLine(line.Line, lc.Patterns)
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
