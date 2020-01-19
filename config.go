package main

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"gopkg.in/yaml.v2"
)

type PatternConfig struct {
	Match         string
	MatchCompiled *regexp.Regexp `yaml:-`
	Metric        string
	Type          string
	Help          string
	Action        string
	Value         string
	Continue      bool
	Labels        map[string]string
	CounterVec    *prometheus.CounterVec
	GaugeVec      *prometheus.GaugeVec
}

type LogConfig struct {
	Path     string
	Patterns []*PatternConfig
}

type Config struct {
	Logs []LogConfig
}

func loadConfigFile(path string) (Config, error) {
	var c Config
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return c, fmt.Errorf("cannot read file '%s': %s", path, err)
	}
	err = yaml.UnmarshalStrict(content, &c)
	if err != nil {
		return c, fmt.Errorf("yaml parse error: %s", err)
	}
	return c, nil
}

func (pc *PatternConfig) compile() error {
	err := pc.validateTypeAction()
	if err != nil {
		return err
	}
	err = pc.compileMatch()
	if err != nil {
		return err
	}
	err = pc.validateValue()
	if err != nil {
		return err
	}
	err = pc.validateLabels()
	if err != nil {
		return err
	}
	return nil
}

func (pc *PatternConfig) validateTypeAction() error {
	switch pc.Type {
	case "counter":
		switch pc.Action {
		case "inc":
		default:
			return fmt.Errorf("action '%s' is unsupported for counters", pc.Action)
		}
	case "gauge":
		switch pc.Action {
		case "inc":
		case "dec":
		case "set":
		default:
			return fmt.Errorf("action '%s' is unsupported for counters", pc.Action)
		}
	default:
		return fmt.Errorf("unsupported type '%s'", pc.Type)
	}
	return nil
}

func (pc *PatternConfig) validateValue() error {
	numGroups := pc.MatchCompiled.NumSubexp()
	dummyMatches := []string{}
	for x := 0; x < numGroups+1; x++ {
		dummyMatches = append(dummyMatches, "0")
	}
	_, err := pc.eval(dummyMatches)
	return err
}

func (pc *PatternConfig) validateLabels() error {
	numGroups := pc.MatchCompiled.NumSubexp()
	dummyMatches := []string{}
	for x := 0; x < numGroups+1; x++ {
		dummyMatches = append(dummyMatches, "0")
	}
	_, err := pc.getEvaluatedLabels(dummyMatches)
	return err
}

func (pc *PatternConfig) compileMatch() error {
	r, err := regexp.Compile(pc.Match)
	if err != nil {
		return fmt.Errorf("failed to compile match expression: %s", err)
	}
	pc.MatchCompiled = r
	return nil
}

func (pc *PatternConfig) eval(matches []string) (float64, error) {
	var err error
	val := pc.Value
	if val == "now()" {
		return float64(time.Now().Unix()), nil
	}
	if strings.HasPrefix(val, "$") {
		val, err = pc.getGroupMatch(val[1:], matches)
		if err != nil {
			return 0, err
		}
	}
	return strconv.ParseFloat(val, 64)
}

func (pc *PatternConfig) getEvaluatedLabels(matches []string) (map[string]string, error) {
	ret := map[string]string{}
	for label, val := range pc.Labels {
		if val == "now()" {
			ret[label] = strconv.FormatInt(time.Now().Unix(), 10)
			continue
		}
		if !strings.HasPrefix(val, "$") {
			ret[label] = val
			continue
		}
		val, err := pc.getGroupMatch(val[1:], matches)
		if err != nil {
			return nil, err
		}
		ret[label] = val
	}
	return ret, nil
}

func (pc *PatternConfig) getGroupMatch(name string, matches []string) (string, error) {
	if len(name) < 1 {
		return "", fmt.Errorf("empty group name")
	}
	for idx, existingName := range pc.MatchCompiled.SubexpNames() {
		if existingName == name {
			if len(matches) < idx+1 {
				return "", fmt.Errorf("missing regexp match data at index %d", idx)
			}
			return matches[idx], nil
		}
	}
	return "", fmt.Errorf("named group not found")
}
