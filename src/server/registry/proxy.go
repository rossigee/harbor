// Copyright Project Harbor Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package registry

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	commonhttp "github.com/goharbor/harbor/src/common/http"
	"github.com/goharbor/harbor/src/lib/config"
	"github.com/goharbor/harbor/src/lib/log"
)

var proxy = newProxy()

type cachedToken struct {
	token   string
	expires time.Time
}

var tokenCache struct {
	mu   sync.RWMutex
	data *cachedToken
}

var detectedAuthType atomic.Value

// Override in tests to control the HTTP client used for registry probe.
var probeHTTPClient = defaultProbeClient()

func defaultProbeClient() *http.Client {
	return &http.Client{
		Timeout:   10 * time.Second,
		Transport: commonhttp.GetHTTPTransport(),
	}
}

// Override in tests to control the token service endpoint.
var getTokenServiceURL = func() string {
	return config.InternalTokenServiceEndpoint()
}

func newProxy() http.Handler {
	regURL, _ := config.RegistryURL()
	u, err := url.Parse(regURL)
	if err != nil {
		panic(fmt.Sprintf("failed to parse URL of registry: %v", err))
	}
	p := httputil.NewSingleHostReverseProxy(u)
	if commonhttp.InternalTLSEnabled() {
		p.Transport = commonhttp.GetHTTPTransport()
	}

	p.Director = authDirector(p.Director)
	return p
}

// authDirector returns a Director that authenticates to the upstream registry.
// If the upstream request already has a valid Authorization header (set by the
// v2auth middleware from the user's credentials), it passes it through.
// Otherwise it falls back to the shared registry credential.
func authDirector(d func(*http.Request)) func(*http.Request) {
	return func(r *http.Request) {
		d(r)
		if r == nil {
			return
		}
		// If the user already provided auth (validated by v2auth middleware),
		// pass it through to the upstream registry.
		if r.Header.Get("Authorization") != "" {
			return
		}
		switch detectRegistryAuthType() {
		case "token":
			if tk := getRegistryToken(r); tk != "" {
				r.Header.Set("Authorization", "Bearer "+tk)
			}
		default:
			u, p := config.RegistryCredential()
			r.SetBasicAuth(u, p)
		}
	}
}

// detectRegistryAuthType probes the upstream registry to determine which
// authentication scheme it expects (bearer token or basic auth). It first
// tries basic auth with the shared registry credential; if the registry
// responds with a Bearer challenge it returns "token".  The result is cached
// on success; on probe failure "basic" is returned as a safe default and the
// probe is retried on the next request.
func detectRegistryAuthType() string {
	if v := detectedAuthType.Load(); v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}

	authType, err := probeRegistry()
	if err != nil {
		log.Warningf("registry auth probe failed: %v, using basic auth as default", err)
		return "basic"
	}

	detectedAuthType.Store(authType)
	return authType
}

// probeRegistry makes a request to the registry's /v2/ endpoint with the
// shared registry credential as basic auth.  If the registry accepts the
// credential (any non-401 response) it returns "basic".  If the registry
// returns 401 with a Www-Authenticate: Bearer challenge it returns "token".
func probeRegistry() (string, error) {
	regURL, err := config.RegistryURL()
	if err != nil {
		return "", fmt.Errorf("failed to get registry URL: %w", err)
	}

	req, err := http.NewRequest(http.MethodGet, strings.TrimSuffix(regURL, "/")+"/v2/", nil)
	if err != nil {
		return "", fmt.Errorf("failed to create probe request: %w", err)
	}
	u, p := config.RegistryCredential()
	if u != "" {
		req.SetBasicAuth(u, p)
	}

	resp, err := probeHTTPClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to probe registry: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		wwwAuth := resp.Header.Get("Www-Authenticate")
		if strings.HasPrefix(wwwAuth, "Bearer") {
			return "token", nil
		}
	}
	return "basic", nil
}

// getRegistryToken obtains a bearer token for the upstream registry by
// exchanging the shared registry credential with Harbor's /service/token
// endpoint. The token is cached for 30 minutes.
func getRegistryToken(r *http.Request) string {
	tokenCache.mu.RLock()
	if tokenCache.data != nil && tokenCache.data.token != "" && time.Now().Before(tokenCache.data.expires) {
		tk := tokenCache.data.token
		tokenCache.mu.RUnlock()
		return tk
	}
	tokenCache.mu.RUnlock()

	tokenCache.mu.Lock()
	defer tokenCache.mu.Unlock()

	if tokenCache.data != nil && tokenCache.data.token != "" && time.Now().Before(tokenCache.data.expires) {
		return tokenCache.data.token
	}

	scope := scopeFromRequest(r)
	tokenURL := fmt.Sprintf("%s?service=harbor-registry&scope=%s", getTokenServiceURL(), url.QueryEscape(scope))

	req, err := http.NewRequest(http.MethodGet, tokenURL, nil)
	if err != nil {
		log.Warningf("failed to create token request: %v", err)
		return ""
	}

	username, password := config.RegistryCredential()
	if username != "" && password != "" {
		req.SetBasicAuth(username, password)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Warningf("failed to get registry token: %v", err)
		return ""
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Warningf("token service returned status %d", resp.StatusCode)
		return ""
	}

	var tokenResp struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		log.Warningf("failed to decode token response: %v", err)
		return ""
	}
	if tokenResp.Token == "" {
		return ""
	}

	tokenCache.data = &cachedToken{
		token:   tokenResp.Token,
		expires: time.Now().Add(30 * time.Minute),
	}
	return tokenResp.Token
}

// scopeFromRequest extracts the Docker registry scope from the request path.
// For a path like /v2/library/nginx/manifests/latest it returns
// repository:library/nginx:pull,push.
func scopeFromRequest(r *http.Request) string {
	if r == nil || r.URL == nil {
		return "repository:*:pull,push"
	}
	path := r.URL.Path
	if !strings.HasPrefix(path, "/v2/") {
		return "repository:*:pull,push"
	}
	parts := strings.SplitN(strings.TrimPrefix(path, "/v2/"), "/", 3)
	if len(parts) < 2 {
		return "repository:*:pull,push"
	}
	return fmt.Sprintf("repository:%s/%s:pull,push", parts[0], parts[1])
}

func basicAuthDirector(d func(*http.Request)) func(*http.Request) {
	return func(r *http.Request) {
		d(r)
		if r != nil {
			u, p := config.RegistryCredential()
			r.SetBasicAuth(u, p)
		}
	}
}
