package main

import (
	"bytes"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func TestHandleLine(t *testing.T) {
	patterns := []*PatternConfig{
		&PatternConfig{
			Match:    ".*foo.*",
			Metric:   "lines_with_foo",
			Help:     "number of lines which contained 'foo'",
			Type:     "counter",
			Action:   "inc",
			Value:    "1",
			Continue: true,
		},
		&PatternConfig{
			Match:  "this sets.*bar.*",
			Metric: "current_bar",
			Help:   "current bar value",
			Type:   "gauge",
			Action: "set",
			Value:  "1",
		},
		&PatternConfig{
			Match:  "with labels.*",
			Metric: "with_labels",
			Help:   "with labels",
			Type:   "gauge",
			Action: "set",
			Value:  "1",
			Labels: map[string]string{
				"foo": "bar",
			},
		},
		&PatternConfig{
			Match:  "calculated_label=(?P<a>.+)",
			Metric: "calculated_label",
			Help:   "calculated label",
			Type:   "counter",
			Action: "inc",
			Value:  "1",
			Labels: map[string]string{
				"calculated": "$a",
			},
		},
		&PatternConfig{
			Match:    "request status: failed",
			Metric:   "request_status_failed",
			Help:     "test",
			Type:     "counter",
			Action:   "inc",
			Value:    "1",
			Continue: true,
		},
		&PatternConfig{
			Match:  "request status: .*",
			Metric: "request_status_total",
			Help:   "test",
			Type:   "counter",
			Action: "inc",
			Value:  "1",
		},

		&PatternConfig{
			Match:    "send status: failed",
			Metric:   "send_status_failed",
			Help:     "test",
			Type:     "counter",
			Action:   "inc",
			Value:    "1",
			Continue: false,
		},
		&PatternConfig{
			Match:  "send status: .*",
			Metric: "send_status_exceptfailed",
			Help:   "test",
			Type:   "counter",
			Action: "inc",
			Value:  "1",
		},

		&PatternConfig{
			Match:  "starting job",
			Metric: "jobs_inflight",
			Help:   "Currently running jobs",
			Type:   "gauge",
			Action: "inc",
			Value:  "1",
		},
		&PatternConfig{
			Match:  "finishing job",
			Metric: "jobs_inflight",
			Help:   "Currently running jobs",
			Type:   "gauge",
			Action: "dec",
			Value:  "1",
		},
	}
	err := registerMetrics(patterns)
	if err != nil {
		t.Fatalf("failed to register metrics %s", err)
	}

	handleLine("uncounted", patterns)

	handler := promhttp.Handler()

	want := func(wanted string) {
		req := httptest.NewRequest("GET", "/metrics", &bytes.Buffer{})
		rc := httptest.NewRecorder()
		handler.ServeHTTP(rc, req)
		got := rc.Body.String()
		if !strings.Contains(got, wanted) {
			t.Fatalf("did not find '%s' in output", wanted)
		}
	}

	handleLine("counted as 'foo' appears", patterns)
	want("\nlines_with_foo 1\n")

	handleLine("this sets current bar...", patterns)
	want("\ncurrent_bar 1\n")

	handleLine("with labels...", patterns)
	want("\nwith_labels{foo=\"bar\"} 1\n")

	handleLine("calculated_label=foo", patterns)
	want("\ncalculated_label{calculated=\"foo\"} 1\n")

	handleLine("request status: failed", patterns)
	want("\nrequest_status_total 1\n")
	want("\nrequest_status_failed 1\n")

	handleLine("send status: ok", patterns)
	handleLine("send status: failed", patterns)
	want("\nsend_status_failed 1\n")
	want("\nsend_status_exceptfailed 1\n")

	handleLine("starting job", patterns)
	want("\njobs_inflight 1\n")
	handleLine("starting job", patterns)
	want("\njobs_inflight 2\n")
	handleLine("finishing job", patterns)
	want("\njobs_inflight 1\n")
	handleLine("finishing job", patterns)
	want("\njobs_inflight 0\n")
}

func TestRunTailer(t *testing.T) {
	c := &Config{}
	tailers := NewSafeTailers()
	err := runTailer(c, "doc/example.yaml", tailers)
	if err != nil {
		t.Fatalf("Got error: %v, wanted: nil", err)
	}

}
