package main

import (
	"bytes"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func TestHTTP(t *testing.T) {
	handler := promhttp.Handler()
	err := registerMetrics([]*PatternConfig{
		&PatternConfig{
			Metric: "example_metric",
			Help:   "just an example",
			Type:   "counter",
			Action: "inc",
			Value:  "1",
		},
	})
	if err != nil {
		t.Fatalf("failed to register metrics %s", err)
	}
	req := httptest.NewRequest("GET", "/metrics", &bytes.Buffer{})
	rc := httptest.NewRecorder()
	handler.ServeHTTP(rc, req)
	got := rc.Body.String()
	wanted := "\nexample_metric 0\n"
	if !strings.Contains(got, wanted) {
		t.Fatalf("unexpected output, got=%s, wanted=%s", got, wanted)
	}
}
