// main_test.go
package main

import (
	"fmt"
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

	// returns the exact same IPs as hardcoded list
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "173.245.48.0/20")
		fmt.Fprintln(w, "103.21.244.0/22")
		fmt.Fprintln(w, "103.22.200.0/22")
		fmt.Fprintln(w, "103.31.4.0/22")
		fmt.Fprintln(w, "141.101.64.0/18")
		fmt.Fprintln(w, "108.162.192.0/18")
		fmt.Fprintln(w, "190.93.240.0/20")
		fmt.Fprintln(w, "188.114.96.0/20")
		fmt.Fprintln(w, "197.234.240.0/22")
		fmt.Fprintln(w, "198.41.128.0/17")
		fmt.Fprintln(w, "162.158.0.0/15")
		fmt.Fprintln(w, "104.16.0.0/13")
		fmt.Fprintln(w, "104.24.0.0/14")
		fmt.Fprintln(w, "172.64.0.0/13")
		fmt.Fprintln(w, "131.0.72.0/22")
	}))
	defer mockServer.Close()

	detectIpRangesChange(mockServer.URL)

	value := testutil.ToFloat64(cloudflareIpRangesChanged)
	assert.Equal(t, float64(0), value)
}

func TestDetectIpRangesChange_OnlyOneIpResponse(t *testing.T) {
	resetMetric()

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "173.245.48.0/20")
	}))
	defer mockServer.Close()

	detectIpRangesChange(mockServer.URL)

	value := testutil.ToFloat64(cloudflareIpRangesChanged)
	assert.Equal(t, float64(1), value)
}

func TestDetectIpRangesChange_ChangedLastIp(t *testing.T) {
	resetMetric()

	// change in the last ip
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "173.245.48.0/20")
		fmt.Fprintln(w, "103.21.244.0/22")
		fmt.Fprintln(w, "103.22.200.0/22")
		fmt.Fprintln(w, "103.31.4.0/22")
		fmt.Fprintln(w, "141.101.64.0/18")
		fmt.Fprintln(w, "108.162.192.0/18")
		fmt.Fprintln(w, "190.93.240.0/20")
		fmt.Fprintln(w, "188.114.96.0/20")
		fmt.Fprintln(w, "197.234.240.0/22")
		fmt.Fprintln(w, "198.41.128.0/17")
		fmt.Fprintln(w, "162.158.0.0/15")
		fmt.Fprintln(w, "104.16.0.0/13")
		fmt.Fprintln(w, "104.24.0.0/14")
		fmt.Fprintln(w, "172.64.0.0/13")
		fmt.Fprintln(w, "131.0.73.0/22") // changed from 131.0.72.0/22
	}))
	defer mockServer.Close()

	detectIpRangesChange(mockServer.URL)

	value := testutil.ToFloat64(cloudflareIpRangesChanged)
	assert.Equal(t, float64(1), value)
}

func TestDetectIpRangesChange_DeleteFirstIp(t *testing.T) {
	resetMetric()

	// deleted first ip
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "103.21.244.0/22")
		fmt.Fprintln(w, "103.22.200.0/22")
		fmt.Fprintln(w, "103.31.4.0/22")
		fmt.Fprintln(w, "141.101.64.0/18")
		fmt.Fprintln(w, "108.162.192.0/18")
		fmt.Fprintln(w, "190.93.240.0/20")
		fmt.Fprintln(w, "188.114.96.0/20")
		fmt.Fprintln(w, "197.234.240.0/22")
		fmt.Fprintln(w, "198.41.128.0/17")
		fmt.Fprintln(w, "162.158.0.0/15")
		fmt.Fprintln(w, "104.16.0.0/13")
		fmt.Fprintln(w, "104.24.0.0/14")
		fmt.Fprintln(w, "172.64.0.0/13")
		fmt.Fprintln(w, "131.0.73.0/22")
	}))
	defer mockServer.Close()

	detectIpRangesChange(mockServer.URL)

	value := testutil.ToFloat64(cloudflareIpRangesChanged)
	assert.Equal(t, float64(1), value)
}

func TestDetectIpRangesChange_AddedOneIp(t *testing.T) {
	resetMetric()

	// added one ip
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "173.245.48.0/20")
		fmt.Fprintln(w, "103.21.244.0/22")
		fmt.Fprintln(w, "103.22.200.0/22")
		fmt.Fprintln(w, "103.31.4.0/22")
		fmt.Fprintln(w, "141.101.64.0/18")
		fmt.Fprintln(w, "108.162.192.0/18")
		fmt.Fprintln(w, "190.93.240.0/20")
		fmt.Fprintln(w, "188.114.96.0/20")
		fmt.Fprintln(w, "197.234.240.0/22")
		fmt.Fprintln(w, "198.41.128.0/17")
		fmt.Fprintln(w, "162.158.0.0/15")
		fmt.Fprintln(w, "104.16.0.0/13")
		fmt.Fprintln(w, "104.24.0.0/14")
		fmt.Fprintln(w, "172.64.0.0/13")
		fmt.Fprintln(w, "131.0.72.0/22")
		fmt.Fprintln(w, "135.0.72.0/22") //this was added
	}))
	defer mockServer.Close()

	detectIpRangesChange(mockServer.URL)

	value := testutil.ToFloat64(cloudflareIpRangesChanged)
	assert.Equal(t, float64(1), value)
}

func TestDetectIpRangesChange_StringOnlyResponse(t *testing.T) {
	resetMetric()

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "\ntest")
	}))
	defer mockServer.Close()

	detectIpRangesChange(mockServer.URL)

	value := testutil.ToFloat64(cloudflareIpRangesChanged)
	assert.Equal(t, float64(1), value)
}

func TestDetectIpRangesChange_NoResponse(t *testing.T) {
	resetMetric()

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	}))
	defer mockServer.Close()

	detectIpRangesChange(mockServer.URL)

	value := testutil.ToFloat64(cloudflareIpRangesChanged)
	assert.Equal(t, float64(1), value)
}
