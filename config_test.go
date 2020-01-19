package main

import (
	"strconv"
	"testing"
	"time"
)

func TestLoadExampleConfig(t *testing.T) {
	c, err := loadConfigFile("doc/example.yaml")
	if err != nil {
		t.Fatalf("error: %s", err)
	}
	_ = c
}

func TestPatternConfigValidate(t *testing.T) {
	invalidPcs := []PatternConfig{
		PatternConfig{}, // empty metric name, bad action
		PatternConfig{
			Metric: "foo",
			Type:   "invalid",
		},
		PatternConfig{
			Metric: "foo",
			Type:   "counter",
			Action: "invalid",
		},
		PatternConfig{
			Metric: "foo",
			Type:   "gauge",
			Action: "",
		},
		PatternConfig{
			Metric: "foo",
			Type:   "counter",
			Action: "inc",
			Value:  "1",
			Match:  "^foo(", // missing closing brace
		},
		PatternConfig{
			Metric: "foo",
			Type:   "counter",
			Action: "inc",
			Value:  "1",
			Match:  "^foo",
			Labels: map[string]string{
				"foo": "$invalid",
			},
		},
	}
	for _, pc := range invalidPcs {
		err := pc.compile()
		if err == nil {
			t.Fatalf("expected compile to fail on %v, but it didn't", pc)
		}
	}

	validPcs := []PatternConfig{
		PatternConfig{
			Metric: "foo",
			Type:   "counter",
			Action: "inc",
			Value:  "1",
		},
		PatternConfig{
			Metric: "foo",
			Type:   "gauge",
			Action: "set",
			Value:  "30",
		},
		PatternConfig{
			Metric: "last_success",
			Type:   "gauge",
			Action: "set",
			Value:  "now()",
		},
		PatternConfig{
			Match:  "(?P<foo>.)",
			Metric: "last_success",
			Type:   "gauge",
			Action: "set",
			Value:  "$foo",
		},
	}
	for _, pc := range validPcs {
		err := pc.compile()
		if err != nil {
			t.Fatalf("expected compile to succeed on %v, but it didn't: %s", pc, err)
		}
	}
}

func TestPatternEvalSimple(t *testing.T) {
	pc := &PatternConfig{
		Metric: "last_success",
		Type:   "gauge",
		Action: "set",
		Value:  "83",
	}
	err := pc.compile()
	if err != nil {
		t.Fatalf("compile failed: %s", err)
	}
	got, err := pc.eval([]string{})
	if err != nil {
		t.Fatalf("got err: %s", err)
	}
	wanted := 83.0
	if got != wanted {
		t.Fatalf("wrong value, got=%v, wanted=%v", got, wanted)
	}
}

func TestPatternEvalTime(t *testing.T) {
	pc := &PatternConfig{
		Metric: "last_success",
		Type:   "gauge",
		Action: "set",
		Value:  "now()",
	}
	err := pc.compile()
	if err != nil {
		t.Fatalf("compile failed: %s", err)
	}
	got, err := pc.eval([]string{})
	if err != nil {
		t.Fatalf("got err: %s", err)
	}

	now:=time.Now().Unix()

	if !(int64(got) <= now && int64(got) > now -60) {
		t.Fatalf("wrong value, got=%v", got)
	}
}

func TestPatternEvalLabelTime(t *testing.T) {
	pc := &PatternConfig{
		Metric: "last_success",
		Type:   "gauge",
		Action: "set",
		Value:  "83",
		Labels: map[string]string{
			"timestamp": "now()",
		},
	}
	err := pc.compile()
	if err != nil {
		t.Fatalf("compile failed: %s", err)
	}
	got, err := pc.getEvaluatedLabels([]string{})
	if err != nil {
		t.Fatalf("got err: %s", err)
	}

	now:=time.Now().Unix()
	gotLabelTs, err := strconv.ParseInt(got["timestamp"], 10, 64)
	if err != nil {
		t.Fatalf("got err: %s", err)
	}

	if !(gotLabelTs <= now && gotLabelTs > now-60) {
		t.Fatalf("wrong value, got=%v", got)
	}
}


func TestPatternEvalMatch(t *testing.T) {
	pc := &PatternConfig{
		Match:  "^foo(?P<firstgroup>x.+)bar",
		Metric: "last_success",
		Type:   "gauge",
		Action: "set",
		Value:  "$firstgroup",
	}
	err := pc.compile()
	if err != nil {
		t.Fatalf("compile failed: %s", err)
	}
	got, err := pc.eval([]string{"", "1"})
	if err != nil {
		t.Fatalf("got err: %s", err)
	}
	if int64(got) < time.Now().Unix() && int64(got) > time.Now().Unix()-60 {
		t.Fatalf("wrong value, got=%v", got)
	}
}
