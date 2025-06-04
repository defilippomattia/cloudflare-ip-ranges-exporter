// main_test.go
package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
)

func resetMetric() {
	cloudflareIpRangesChanged.Set(0)
}

func TestDetectIpRangesChange_NoChanges(t *testing.T) {
	resetMetric()

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := io.WriteString(w, `173.245.48.0/20
103.21.244.0/22
103.22.200.0/22
103.31.4.0/22
141.101.64.0/18
108.162.192.0/18
190.93.240.0/20
188.114.96.0/20
197.234.240.0/22
198.41.128.0/17
162.158.0.0/15
104.16.0.0/13
104.24.0.0/14
172.64.0.0/13
131.0.72.0/22
`)
		if err != nil {
			t.Fatalf("failed to write mock response: %v", err)
		}
	}))
	defer mockServer.Close()

	detectIpRangesChange(mockServer.URL)

	value := testutil.ToFloat64(cloudflareIpRangesChanged)
	assert.Equal(t, float64(0), value)
}

func TestDetectIpRangesChange_OnlyOneIpResponse(t *testing.T) {
	resetMetric()

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := io.WriteString(w, "173.245.48.0/20")
		if err != nil {
			t.Fatalf("failed to write mock response: %v", err)
		}
	}))
	defer mockServer.Close()

	detectIpRangesChange(mockServer.URL)

	value := testutil.ToFloat64(cloudflareIpRangesChanged)
	assert.Equal(t, float64(1), value)
}

func TestDetectIpRangesChange_ChangedLastIp(t *testing.T) {
	resetMetric()

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := io.WriteString(w, `173.245.48.0/20
103.21.244.0/22
103.22.200.0/22
103.31.4.0/22
141.101.64.0/18
108.162.192.0/18
190.93.240.0/20
188.114.96.0/20
197.234.240.0/22
198.41.128.0/17
162.158.0.0/15
104.16.0.0/13
104.24.0.0/14
172.64.0.0/13
131.0.73.0/22
`)
		if err != nil {
			t.Fatalf("failed to write mock response: %v", err)
		}
	}))
	defer mockServer.Close()

	detectIpRangesChange(mockServer.URL)

	value := testutil.ToFloat64(cloudflareIpRangesChanged)
	assert.Equal(t, float64(1), value)
}

func TestDetectIpRangesChange_DeleteFirstIp(t *testing.T) {
	resetMetric()

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := io.WriteString(w, `103.21.244.0/22
103.22.200.0/22
103.31.4.0/22
141.101.64.0/18
108.162.192.0/18
190.93.240.0/20
188.114.96.0/20
197.234.240.0/22
198.41.128.0/17
162.158.0.0/15
104.16.0.0/13
104.24.0.0/14
172.64.0.0/13
131.0.73.0/22
`)
		if err != nil {
			t.Fatalf("failed to write mock response: %v", err)
		}
	}))
	defer mockServer.Close()

	detectIpRangesChange(mockServer.URL)

	value := testutil.ToFloat64(cloudflareIpRangesChanged)
	assert.Equal(t, float64(1), value)
}

func TestDetectIpRangesChange_AddedOneIp(t *testing.T) {
	resetMetric()

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := io.WriteString(w, `173.245.48.0/20
103.21.244.0/22
103.22.200.0/22
103.31.4.0/22
141.101.64.0/18
108.162.192.0/18
190.93.240.0/20
188.114.96.0/20
197.234.240.0/22
198.41.128.0/17
162.158.0.0/15
104.16.0.0/13
104.24.0.0/14
172.64.0.0/13
131.0.72.0/22
135.0.72.0/22
`)
		if err != nil {
			t.Fatalf("failed to write mock response: %v", err)
		}
	}))
	defer mockServer.Close()

	detectIpRangesChange(mockServer.URL)

	value := testutil.ToFloat64(cloudflareIpRangesChanged)
	assert.Equal(t, float64(1), value)
}

func TestDetectIpRangesChange_StringOnlyResponse(t *testing.T) {
	resetMetric()

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := io.WriteString(w, "\ntest\n")
		if err != nil {
			t.Fatalf("failed to write mock response: %v", err)
		}
	}))
	defer mockServer.Close()

	detectIpRangesChange(mockServer.URL)

	value := testutil.ToFloat64(cloudflareIpRangesChanged)
	assert.Equal(t, float64(1), value)
}

func TestDetectIpRangesChange_NoResponse(t *testing.T) {
	resetMetric()

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	defer mockServer.Close()

	detectIpRangesChange(mockServer.URL)

	value := testutil.ToFloat64(cloudflareIpRangesChanged)
	assert.Equal(t, float64(1), value)
}
