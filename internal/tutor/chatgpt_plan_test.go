package tutor

import (
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestPKCEGeneration(t *testing.T) {
	pkce, err := GeneratePKCE()
	if err != nil {
		t.Fatalf("failed generating PKCE: %v", err)
	}

	if len(pkce.Verifier) < 32 {
		t.Errorf("verifier too short: %d", len(pkce.Verifier))
	}
	if pkce.Challenge == "" || pkce.State == "" || pkce.Nonce == "" {
		t.Errorf("empty PKCE parameter")
	}

	// Verify SHA-256 challenge
	hash := sha256.Sum256([]byte(pkce.Verifier))
	expectedChallenge := base64.RawURLEncoding.EncodeToString(hash[:])
	if pkce.Challenge != expectedChallenge {
		t.Errorf("challenge mismatch: got %s, expected %s", pkce.Challenge, expectedChallenge)
	}
}

func TestOAuthLoopbackServerCallbackSuccess(t *testing.T) {
	tutor := NewOpenAIChatGPTPlanTutor(ChatGPTPlanConfig{})
	flow, authURL, err := tutor.StartOAuthFlow()
	if err != nil {
		t.Fatalf("failed starting OAuth flow: %v", err)
	}
	defer flow.Close()

	if !strings.Contains(authURL, "code_challenge=") || !strings.Contains(authURL, "response_type=code") {
		t.Errorf("invalid auth URL: %s", authURL)
	}

	// Simulate browser callback
	callbackURL := fmt.Sprintf("%s?code=test-auth-code-123&client_id=oaiapp_test&state=%s", flow.RedirectURI, flow.PKCE.State)

	go func() {
		time.Sleep(20 * time.Millisecond)
		resp, err := http.Get(callbackURL)
		if err != nil {
			t.Errorf("callback GET failed: %v", err)
			return
		}
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		if !strings.Contains(string(body), "Authorization received") {
			t.Errorf("expected success HTML, got: %s", string(body))
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	code, err := flow.WaitForCallback(ctx)
	if err != nil {
		t.Fatalf("WaitForCallback failed: %v", err)
	}
	if code != "test-auth-code-123" {
		t.Errorf("expected code test-auth-code-123, got %s", code)
	}
}

func TestOAuthLoopbackServerCallbackStateMismatch(t *testing.T) {
	tutor := NewOpenAIChatGPTPlanTutor(ChatGPTPlanConfig{})
	flow, _, err := tutor.StartOAuthFlow()
	if err != nil {
		t.Fatalf("failed starting OAuth flow: %v", err)
	}
	defer flow.Close()

	// Send wrong state
	callbackURL := fmt.Sprintf("%s?code=test-code&state=wrong-state-xyz", flow.RedirectURI)

	go func() {
		time.Sleep(20 * time.Millisecond)
		_, _ = http.Get(callbackURL)
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	_, err = flow.WaitForCallback(ctx)
	if err == nil || !strings.Contains(err.Error(), "state mismatch") {
		t.Fatalf("expected state mismatch error, got: %v", err)
	}
}

func TestTokenExchangeAndRefreshMock(t *testing.T) {
	signingKey, _ := rsa.GenerateKey(rand.Reader, 2048)
	mockOAuthServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/oauth/token" {
			http.NotFound(w, r)
			return
		}

		_ = r.ParseForm()
		grantType := r.Form.Get("grant_type")

		if grantType == "authorization_code" {
			if r.Form.Get("code") != "valid-code" || r.Form.Get("code_verifier") == "" {
				http.Error(w, "invalid request", http.StatusBadRequest)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"access_token":  "mock-access-token-1",
				"refresh_token": "mock-refresh-token-1",
				"expires_in":    3600,
				"scope":         DefaultOAuthScope,
				"id_token":      signedIdentity(t, signingKey, "oaiapp_test", "nonce"),
			})
			return
		}

		if grantType == "refresh_token" {
			if r.Form.Get("refresh_token") != "mock-refresh-token-1" {
				http.Error(w, "invalid refresh token", http.StatusBadRequest)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"access_token":  "mock-refreshed-token-2",
				"refresh_token": "mock-refresh-token-2",
				"expires_in":    7200,
			})
			return
		}

		http.Error(w, "unsupported grant", http.StatusBadRequest)
	}))
	defer mockOAuthServer.Close()

	tutor := NewOpenAIChatGPTPlanTutor(ChatGPTPlanConfig{
		AuthURL:    mockOAuthServer.URL,
		HTTPClient: identityClient(mockOAuthServer.Client(), signingKey),
		ClientID:   "oaiapp_test",
	})

	tutor.pendingNonce = "nonce"
	ctx := context.Background()

	// 1. Exchange Code
	token, err := tutor.ExchangeCode(ctx, "valid-code", "http://localhost:8080/callback", "mock-verifier")
	if err != nil {
		t.Fatalf("exchange code failed: %v", err)
	}
	if token.AccessToken != "mock-access-token-1" {
		t.Errorf("unexpected access token: %s", token.AccessToken)
	}
	if token.RefreshToken != "mock-refresh-token-1" {
		t.Errorf("unexpected refresh token: %s", token.RefreshToken)
	}

	// 2. Refresh Token
	refreshed, err := tutor.RefreshToken(ctx, token)
	if err != nil {
		t.Fatalf("refresh token failed: %v", err)
	}
	if refreshed.AccessToken != "mock-refreshed-token-2" {
		t.Errorf("unexpected refreshed token: %s", refreshed.AccessToken)
	}
}

func TestStreamedResponsesAPIPreviewLimitationsEnforced(t *testing.T) {
	// Verify that mock server checks strictly for store=false and stream=true
	mockAPIServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/responses" {
			http.NotFound(w, r)
			return
		}

		// Check Authorization
		authHeader := r.Header.Get("Authorization")
		if authHeader != "Bearer test-valid-jwt" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		// Check Body for Subscription Preview Invariants: store=false, stream=true
		bodyBytes, _ := io.ReadAll(r.Body)
		var reqBody map[string]interface{}
		_ = json.Unmarshal(bodyBytes, &reqBody)

		if reqBody["store"] != false {
			http.Error(w, "store must be false for subscription route", http.StatusBadRequest)
			return
		}
		if reqBody["stream"] != true {
			http.Error(w, "stream must be true for subscription route", http.StatusBadRequest)
			return
		}

		// Stream SSE events
		w.Header().Set("Content-Type", "text/event-stream")
		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "flusher not supported", http.StatusInternalServerError)
			return
		}

		events := []string{
			`{"type":"response.output_text.delta","delta":"Consider the economic "}`,
			`{"type":"response.output_text.delta","delta":"reality: cash was received before "}`,
			`{"type":"response.output_text.delta","delta":"work was performed."}`,
		}

		for _, ev := range events {
			fmt.Fprintf(w, "data: %s\n\n", ev)
			flusher.Flush()
		}
		fmt.Fprintf(w, "data: {\"type\":\"response.completed\"}\n\n")
		flusher.Flush()
	}))
	defer mockAPIServer.Close()

	authStore, _ := NewAuthStore("")
	_ = authStore.SetChatGPTPlanToken(&OAuthToken{
		AccessToken: "test-valid-jwt",
		Subject:     "user", Scope: DefaultOAuthScope, ClientID: "oaiapp_test",
		ExpiresAt: time.Now().UTC().Add(1 * time.Hour),
	})

	tutor := NewOpenAIChatGPTPlanTutor(ChatGPTPlanConfig{
		AuthStore:  authStore,
		APIURL:     mockAPIServer.URL,
		HTTPClient: mockAPIServer.Client(),
	})

	ctx := context.Background()
	req := Request{ProblemPrompt: "Customer pays $500 advance."}

	resp, err := tutor.Hint(ctx, req)
	if err != nil {
		t.Fatalf("hint failed: %v", err)
	}

	expected := "Consider the economic reality: cash was received before work was performed."
	if resp.Text != expected {
		t.Errorf("got text %q, expected %q", resp.Text, expected)
	}
	if resp.Provider != "chatgpt_plan" {
		t.Errorf("expected provider chatgpt_plan, got %s", resp.Provider)
	}
}

func TestChatGPTPlanFallbackOnUnauthenticated(t *testing.T) {
	authStore, _ := NewAuthStore("") // empty token
	chatgpt := NewOpenAIChatGPTPlanTutor(ChatGPTPlanConfig{
		AuthStore: authStore,
	})
	offline := NewOfflineTutor()
	wrapper := NewFallbackTutor(chatgpt, offline, 1*time.Second)

	req := Request{
		ProblemPrompt: "Acme pays rent.",
		CausalHint:    "Rent is an expense.",
	}

	resp, err := wrapper.Hint(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !resp.Fallback {
		t.Errorf("expected Fallback=true when chatgpt plan is unauthenticated")
	}
	if resp.Provider != "offline" {
		t.Errorf("expected offline provider, got %s", resp.Provider)
	}
}

func signedIdentity(t *testing.T, key *rsa.PrivateKey, client, nonce string) string {
	t.Helper()
	h, _ := json.Marshal(map[string]string{"alg": "RS256", "kid": "test"})
	p, _ := json.Marshal(map[string]interface{}{"iss": "https://auth.openai.com", "sub": "user", "aud": client, "nonce": nonce, "exp": time.Now().Add(time.Hour).Unix(), "iat": time.Now().Unix()})
	message := base64.RawURLEncoding.EncodeToString(h) + "." + base64.RawURLEncoding.EncodeToString(p)
	hash := sha256.Sum256([]byte(message))
	sig, err := rsa.SignPKCS1v15(rand.Reader, key, crypto.SHA256, hash[:])
	if err != nil {
		t.Fatal(err)
	}
	return message + "." + base64.RawURLEncoding.EncodeToString(sig)
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }
func identityClient(base *http.Client, key *rsa.PrivateKey) *http.Client {
	return &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		if r.URL.Host != "auth.openai.com" {
			return base.Transport.RoundTrip(r)
		}
		body, _ := json.Marshal(map[string]interface{}{"keys": []map[string]string{{"kid": "test", "kty": "RSA", "n": base64.RawURLEncoding.EncodeToString(key.N.Bytes()), "e": base64.RawURLEncoding.EncodeToString(big.NewInt(int64(key.E)).Bytes())}}})
		return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(string(body))), Header: make(http.Header)}, nil
	})}
}

func TestDynamicRegistrationHostAndCallback(t *testing.T) {
	path := filepath.Join(t.TempDir(), "auth.json")
	store, err := NewAuthStore(path)
	if err != nil {
		t.Fatal(err)
	}
	_ = store.SetChatGPTClientID("acctg-practice")
	provider := NewOpenAIChatGPTPlanTutor(ChatGPTPlanConfig{AuthStore: store})
	flow, raw, err := provider.StartOAuthFlow()
	if err != nil {
		t.Fatal(err)
	}
	defer flow.Close()
	u, _ := url.Parse(raw)
	q := u.Query()
	if u.Path != "/api/accounts/authorize" || q.Get("client_id") != DefaultDCRClientID || q.Get("resource") != DefaultOpenAIAPIURL || q.Get("agent_name_hint") != "AccountTutor 9000" || !hasPlanScope(q.Get("scope")) {
		t.Fatalf("incorrect registration request: %s", raw)
	}
	reloaded, err := NewAuthStore(path)
	if err != nil {
		t.Fatal(err)
	}
	host, err := reloaded.EnsureHostID()
	if err != nil || host != q.Get("ext_agent_host_id") || !strings.HasPrefix(host, "urn:uuid:") {
		t.Fatal("host identity not persisted")
	}
	resp, err := http.Get(flow.RedirectURI + "?state=" + flow.PKCE.State + "&code=code&client_id=oaiapp_issued")
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	_, err = flow.WaitForCallback(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if provider.GetClientID() != "oaiapp_issued" || store.IsConfigured(ProviderChatGPTPlan) {
		t.Fatal("callback should retain issued ID without trusting credentials yet")
	}
}

func TestIDTokenValidationRejectsInvalidIdentity(t *testing.T) {
	key, _ := rsa.GenerateKey(rand.Reader, 2048)
	provider := NewOpenAIChatGPTPlanTutor(ChatGPTPlanConfig{HTTPClient: identityClient(http.DefaultClient, key)})
	valid := signedIdentity(t, key, "oaiapp_test", "nonce")
	if _, err := provider.validateIDToken(context.Background(), valid, "oaiapp_test", "nonce"); err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct{ token, client, nonce string }{
		{valid, "wrong-client", "nonce"}, {valid, "oaiapp_test", "wrong-nonce"}, {valid + "x", "oaiapp_test", "nonce"}, {"unsigned", "oaiapp_test", "nonce"},
	} {
		if _, err := provider.validateIDToken(context.Background(), tc.token, tc.client, tc.nonce); err == nil {
			t.Fatal("invalid identity accepted")
		}
	}
}

func TestOAuthCancellationAndMissingRegistration(t *testing.T) {
	provider := NewOpenAIChatGPTPlanTutor(ChatGPTPlanConfig{})
	flow, _, err := provider.StartOAuthFlow()
	if err != nil {
		t.Fatal(err)
	}
	flow.Close()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if _, err := flow.WaitForCallback(ctx); err == nil {
		t.Fatal("closed sign-in did not cancel")
	}
	flow, _, err = provider.StartOAuthFlow()
	if err != nil {
		t.Fatal(err)
	}
	defer flow.Close()
	resp, err := http.Get(flow.RedirectURI + "?state=" + flow.PKCE.State + "&code=code")
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if _, err := flow.WaitForCallback(ctx); err == nil {
		t.Fatal("missing issued ID accepted")
	}
}
