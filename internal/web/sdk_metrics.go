package web

import "sync"

type sdkMetricKey struct {
	Version string
	Package string
}

type sdkDownloadMetrics struct {
	mu          sync.Mutex
	acceptances map[sdkMetricKey]uint64
	redirects   map[sdkMetricKey]uint64
	errors      uint64
}

func newSDKDownloadMetrics() *sdkDownloadMetrics {
	return &sdkDownloadMetrics{acceptances: map[sdkMetricKey]uint64{}, redirects: map[sdkMetricKey]uint64{}}
}

func (m *sdkDownloadMetrics) recordAcceptance(version, packageSlug string) {
	if m == nil {
		return
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.acceptances[sdkMetricKey{Version: version, Package: packageSlug}]++
}

func (m *sdkDownloadMetrics) recordRedirect(version, packageSlug string) {
	if m == nil {
		return
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.redirects[sdkMetricKey{Version: version, Package: packageSlug}]++
}

func (m *sdkDownloadMetrics) recordError() {
	if m == nil {
		return
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.errors++
}

func (m *sdkDownloadMetrics) snapshot() (map[sdkMetricKey]uint64, map[sdkMetricKey]uint64, uint64) {
	if m == nil {
		return nil, nil, 0
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	acceptances := make(map[sdkMetricKey]uint64, len(m.acceptances))
	redirects := make(map[sdkMetricKey]uint64, len(m.redirects))
	for key, value := range m.acceptances {
		acceptances[key] = value
	}
	for key, value := range m.redirects {
		redirects[key] = value
	}
	return acceptances, redirects, m.errors
}
