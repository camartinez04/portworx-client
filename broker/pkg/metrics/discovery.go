package metrics

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

const (
	inClusterTokenFile = "/var/run/secrets/kubernetes.io/serviceaccount/token"
	inClusterCAFile    = "/var/run/secrets/kubernetes.io/serviceaccount/ca.crt"
	kubeAPIHost        = "https://kubernetes.default.svc"
	cacheTTL           = 60 * time.Second
)

// urlCache caches discovered pod metrics URLs for cacheTTL to avoid
// hammering the Kubernetes API on every metrics poll.
var urlCache struct {
	sync.Mutex
	urls      []string
	expiresAt time.Time
}

// ResolveMetricsURLs returns the list of Portworx node metrics endpoints.
//
// If override is non-empty it is used directly (single URL or comma-separated
// list) – useful for local development with port-forwards.
//
// When override is empty the function auto-discovers all running Portworx pod
// IPs via the in-cluster Kubernetes API and constructs a metrics URL for each
// pod.  Results are cached for 60 s.
//
// Auto-discovery is controlled by three optional environment variables:
//
//	PORTWORX_NAMESPACE      – namespace of Portworx pods     (default: kube-system)
//	PORTWORX_METRICS_PORT   – pod metrics port               (default: 9001)
//	PORTWORX_LABEL_SELECTOR – pod label selector             (default: name=portworx)
func ResolveMetricsURLs(override string) ([]string, error) {
	if override != "" {
		return splitURLs(override), nil
	}
	return discoverPortworxPodURLs()
}

// discoverPortworxPodURLs queries the K8s API for running Portworx pods and
// returns an http://{podIP}:{port}/metrics URL for each pod.
func discoverPortworxPodURLs() ([]string, error) {
	urlCache.Lock()
	defer urlCache.Unlock()

	if time.Now().Before(urlCache.expiresAt) && len(urlCache.urls) > 0 {
		return urlCache.urls, nil
	}

	namespace := envOrDefault("PORTWORX_NAMESPACE", "kube-system")
	port := envOrDefault("PORTWORX_METRICS_PORT", "9001")
	selector := envOrDefault("PORTWORX_LABEL_SELECTOR", "name=portworx")

	ips, err := listPortworxPodIPs(namespace, selector)
	if err != nil {
		return nil, fmt.Errorf("Portworx pod discovery failed: %w", err)
	}
	if len(ips) == 0 {
		return nil, fmt.Errorf("no running Portworx pods found in namespace %q with selector %q", namespace, selector)
	}

	urls := make([]string, 0, len(ips))
	for _, ip := range ips {
		urls = append(urls, "http://"+ip+":"+port+"/metrics")
	}

	urlCache.urls = urls
	urlCache.expiresAt = time.Now().Add(cacheTTL)
	return urls, nil
}

// podList is the minimal K8s API /pods response we need to parse.
type podList struct {
	Items []struct {
		Status struct {
			PodIP string `json:"podIP"`
			Phase string `json:"phase"`
		} `json:"status"`
	} `json:"items"`
}

// listPortworxPodIPs queries the Kubernetes API server using in-cluster
// credentials (service account token + CA cert) and returns the pod IP of
// every Running pod that matches labelSelector in namespace.
func listPortworxPodIPs(namespace, labelSelector string) ([]string, error) {
	token, err := os.ReadFile(inClusterTokenFile)
	if err != nil {
		return nil, fmt.Errorf("not running in-cluster (service account token not found at %s): %w",
			inClusterTokenFile, err)
	}

	caBytes, err := os.ReadFile(inClusterCAFile)
	if err != nil {
		return nil, fmt.Errorf("cannot read in-cluster CA certificate: %w", err)
	}

	pool := x509.NewCertPool()
	if !pool.AppendCertsFromPEM(caBytes) {
		return nil, fmt.Errorf("failed to parse in-cluster CA certificate")
	}

	httpClient := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{RootCAs: pool}, //nolint:gosec
		},
	}

	apiURL := kubeAPIHost + "/api/v1/namespaces/" + namespace + "/pods" +
		"?labelSelector=" + url.QueryEscape(labelSelector) +
		"&fieldSelector=" + url.QueryEscape("status.phase=Running")

	req, err := http.NewRequest(http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("build K8s API request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+strings.TrimSpace(string(token)))

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("K8s API request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("K8s API returned HTTP %d for namespace=%q selector=%q",
			resp.StatusCode, namespace, labelSelector)
	}

	var pl podList
	if err := json.NewDecoder(resp.Body).Decode(&pl); err != nil {
		return nil, fmt.Errorf("parse K8s pod list response: %w", err)
	}

	ips := make([]string, 0, len(pl.Items))
	for _, item := range pl.Items {
		if item.Status.Phase == "Running" && item.Status.PodIP != "" {
			ips = append(ips, item.Status.PodIP)
		}
	}
	return ips, nil
}

func envOrDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
