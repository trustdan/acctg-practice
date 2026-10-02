package tutor

import (
	"bufio"
	"bytes"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

// Default OpenAI OAuth and API endpoints
const (
	DefaultOpenAIAuthURL = "https://auth.openai.com/oauth"
	DefaultOpenAIAPIURL  = "https://api.openai.com/v1"
	DefaultChatGPTModel  = "gpt-4o-mini"
	DefaultDCRClientID   = "acctg-practice-client"
	DefaultOAuthScope    = "openid profile email model.request"
)

// OpenAIChatGPTPlanTutor implements the ChatGPT Plus subscription path.
// It uses official "Sign in with ChatGPT" OAuth PKCE and the streamed Responses API.
//
// Invariant (AGENT-CONTRACT):
// Subscription-route preview limitations strictly enforced: store=false, stream=true.
// Tokens are saved securely in local user data dir with 0600 permissions.
type OpenAIChatGPTPlanTutor struct {
	mu         sync.Mutex
	authStore  *AuthStore
	authURL    string
	apiURL     string
	model      string
	clientID   string
	httpClient *http.Client
	budget     *Budget
}

// ChatGPTPlanConfig configures OpenAIChatGPTPlanTutor.
type ChatGPTPlanConfig struct {
	AuthStore  *AuthStore
	AuthURL    string
	APIURL     string
	Model      string
	ClientID   string
	HTTPClient *http.Client
	Budget     *Budget
}

// NewOpenAIChatGPTPlanTutor creates an OpenAIChatGPTPlanTutor.
func NewOpenAIChatGPTPlanTutor(cfg ChatGPTPlanConfig) *OpenAIChatGPTPlanTutor {
	authURL := cfg.AuthURL
	if authURL == "" {
		authURL = DefaultOpenAIAuthURL
	}
	apiURL := cfg.APIURL
	if apiURL == "" {
		apiURL = DefaultOpenAIAPIURL
	}
	model := cfg.Model
	if model == "" {
		model = DefaultChatGPTModel
	}
	clientID := cfg.ClientID
	if clientID == "" {
		clientID = DefaultDCRClientID
	}
	httpClient := cfg.HTTPClient
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 30 * time.Second}
	}

	return &OpenAIChatGPTPlanTutor{
		authStore:  cfg.AuthStore,
		authURL:    authURL,
		apiURL:     apiURL,
		model:      model,
		clientID:   clientID,
		httpClient: httpClient,
		budget:     cfg.Budget,
	}
}

// Name implements Tutor.
func (t *OpenAIChatGPTPlanTutor) Name() string {
	return "OpenAIChatGPTPlanTutor"
}

// Hint implements Tutor.
func (t *OpenAIChatGPTPlanTutor) Hint(ctx context.Context, req Request) (Response, error) {
	return t.callResponsesAPI(ctx, req, true)
}

// Explain implements Tutor.
func (t *OpenAIChatGPTPlanTutor) Explain(ctx context.Context, req Request) (Response, error) {
	return t.callResponsesAPI(ctx, req, false)
}

// PKCEParams holds OAuth PKCE code challenge and verifier.
type PKCEParams struct {
	Verifier  string
	Challenge string
	State     string
	Nonce     string
}

// GeneratePKCE creates cryptographically secure PKCE verifier, challenge, state, and nonce.
func GeneratePKCE() (*PKCEParams, error) {
	verifierBytes := make([]byte, 32)
	if _, err := rand.Read(verifierBytes); err != nil {
		return nil, fmt.Errorf("failed generating random verifier: %w", err)
	}
	verifier := base64.RawURLEncoding.EncodeToString(verifierBytes)

	hash := sha256.Sum256([]byte(verifier))
	challenge := base64.RawURLEncoding.EncodeToString(hash[:])

	stateBytes := make([]byte, 16)
	if _, err := rand.Read(stateBytes); err != nil {
		return nil, fmt.Errorf("failed generating random state: %w", err)
	}
	state := base64.RawURLEncoding.EncodeToString(stateBytes)

	nonceBytes := make([]byte, 16)
	if _, err := rand.Read(nonceBytes); err != nil {
		return nil, fmt.Errorf("failed generating random nonce: %w", err)
	}
	nonce := base64.RawURLEncoding.EncodeToString(nonceBytes)

	return &PKCEParams{
		Verifier:  verifier,
		Challenge: challenge,
		State:     state,
		Nonce:     nonce,
	}, nil
}

// OAuthFlow handles initiating and completing the browser-based OAuth PKCE sign-in.
type OAuthFlow struct {
	tutor       *OpenAIChatGPTPlanTutor
	listener    net.Listener
	RedirectURI string
	PKCE        *PKCEParams
	codeChan    chan string
	errChan     chan error
	server      *http.Server
}

// StartOAuthFlow starts an ephemeral loopback HTTP listener on 127.0.0.1:0 for the OAuth callback.
func (t *OpenAIChatGPTPlanTutor) StartOAuthFlow() (*OAuthFlow, string, error) {
	pkce, err := GeneratePKCE()
	if err != nil {
		return nil, "", err
	}

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return nil, "", fmt.Errorf("failed binding ephemeral loopback port: %w", err)
	}

	port := listener.Addr().(*net.TCPAddr).Port
	redirectURI := fmt.Sprintf("http://127.0.0.1:%d/callback", port)

	authURL, err := url.Parse(t.authURL + "/authorize")
	if err != nil {
		listener.Close()
		return nil, "", err
	}

	q := authURL.Query()
	q.Set("response_type", "code")
	q.Set("client_id", t.clientID)
	q.Set("redirect_uri", redirectURI)
	q.Set("scope", DefaultOAuthScope)
	q.Set("state", pkce.State)
	q.Set("nonce", pkce.Nonce)
	q.Set("code_challenge", pkce.Challenge)
	q.Set("code_challenge_method", "S256")
	authURL.RawQuery = q.Encode()

	flow := &OAuthFlow{
		tutor:       t,
		listener:    listener,
		RedirectURI: redirectURI,
		PKCE:        pkce,
		codeChan:    make(chan string, 1),
		errChan:     make(chan error, 1),
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/callback", flow.handleCallback)

	flow.server = &http.Server{
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		_ = flow.server.Serve(listener)
	}()

	return flow, authURL.String(), nil
}

func (f *OAuthFlow) handleCallback(w http.ResponseWriter, r *http.Request) {
	receivedState := r.URL.Query().Get("state")
	if receivedState != f.PKCE.State {
		http.Error(w, "Invalid state parameter", http.StatusBadRequest)
		f.errChan <- errors.New("oauth callback state mismatch")
		return
	}

	code := r.URL.Query().Get("code")
	if code == "" {
		errDesc := r.URL.Query().Get("error_description")
		if errDesc == "" {
			errDesc = r.URL.Query().Get("error")
		}
		if errDesc == "" {
			errDesc = "missing authorization code"
		}
		http.Error(w, "Authentication error: "+errDesc, http.StatusBadRequest)
		f.errChan <- fmt.Errorf("oauth error: %s", errDesc)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, `<!DOCTYPE html>
<html>
<head><title>Authentication Successful</title></head>
<body style="font-family: sans-serif; text-align: center; padding-top: 50px;">
  <h2>Authentication Successful!</h2>
  <p>Your ChatGPT Plus plan has been connected to <strong>acctg-practice</strong>.</p>
  <p>You can close this tab and return to the terminal drill.</p>
</body>
</html>`)

	f.codeChan <- code
}

// WaitForCallback waits for the user to complete the browser authorization.
func (f *OAuthFlow) WaitForCallback(ctx context.Context) (string, error) {
	defer f.Close()

	select {
	case <-ctx.Done():
		return "", ctx.Err()
	case err := <-f.errChan:
		return "", err
	case code := <-f.codeChan:
		return code, nil
	}
}

// Close shuts down the loopback listener.
func (f *OAuthFlow) Close() {
	if f.server != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		_ = f.server.Shutdown(ctx)
	}
}

// ExchangeCode exchanges an authorization code for access and refresh tokens.
func (t *OpenAIChatGPTPlanTutor) ExchangeCode(ctx context.Context, code string, redirectURI string, verifier string) (*OAuthToken, error) {
	tokenURL := t.authURL + "/token"

	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("client_id", t.clientID)
	data.Set("code", code)
	data.Set("redirect_uri", redirectURI)
	data.Set("code_verifier", verifier)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := t.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("token exchange network request failed: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("token exchange rejected (%d): %s", resp.StatusCode, string(bodyBytes))
	}

	var tokenResp struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int    `json:"expires_in"`
		Scope        string `json:"scope"`
	}
	if err := json.Unmarshal(bodyBytes, &tokenResp); err != nil {
		return nil, fmt.Errorf("failed decoding token response: %w", err)
	}

	expiresAt := time.Now().UTC().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)
	token := &OAuthToken{
		AccessToken:  tokenResp.AccessToken,
		RefreshToken: tokenResp.RefreshToken,
		ExpiresAt:    expiresAt,
		ClientID:     t.clientID,
		PlanType:     "chatgpt_plus",
		Scope:        tokenResp.Scope,
	}

	if t.authStore != nil {
		_ = t.authStore.SetChatGPTPlanToken(token)
		_ = t.authStore.SetActiveProvider(ProviderChatGPTPlan)
	}

	return token, nil
}

// RefreshToken refreshes an expired access token using the refresh token.
func (t *OpenAIChatGPTPlanTutor) RefreshToken(ctx context.Context, token *OAuthToken) (*OAuthToken, error) {
	if token == nil || token.RefreshToken == "" {
		return nil, errors.New("cannot refresh without refresh token")
	}

	tokenURL := t.authURL + "/token"
	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("client_id", t.clientID)
	data.Set("refresh_token", token.RefreshToken)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := t.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("token refresh network request failed: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("token refresh rejected (%d): %s", resp.StatusCode, string(bodyBytes))
	}

	var tokenResp struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int    `json:"expires_in"`
	}
	if err := json.Unmarshal(bodyBytes, &tokenResp); err != nil {
		return nil, fmt.Errorf("failed decoding refresh response: %w", err)
	}

	token.AccessToken = tokenResp.AccessToken
	if tokenResp.RefreshToken != "" {
		token.RefreshToken = tokenResp.RefreshToken
	}
	token.ExpiresAt = time.Now().UTC().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)

	if t.authStore != nil {
		_ = t.authStore.SetChatGPTPlanToken(token)
	}

	return token, nil
}

func (t *OpenAIChatGPTPlanTutor) getValidAccessToken(ctx context.Context) (string, error) {
	if t.authStore == nil {
		return "", errors.New("auth store not initialized for chatgpt plan")
	}

	cfg := t.authStore.GetConfig()
	token := cfg.ChatGPTPlanToken
	if token == nil || token.AccessToken == "" {
		return "", errors.New("chatgpt plus subscription not connected (sign in via 't' settings)")
	}

	if token.IsExpired() {
		refreshed, err := t.RefreshToken(ctx, token)
		if err != nil {
			return "", fmt.Errorf("chatgpt token expired and refresh failed: %w", err)
		}
		return refreshed.AccessToken, nil
	}

	return token.AccessToken, nil
}

// callResponsesAPI executes a streaming request to the OpenAI Responses API.
//
// Strictly enforces subscription preview limitations:
// - store: false
// - stream: true
func (t *OpenAIChatGPTPlanTutor) callResponsesAPI(ctx context.Context, req Request, isHint bool) (Response, error) {
	if t.budget != nil {
		if err := t.budget.Check(); err != nil {
			return Response{}, err
		}
	}

	accessToken, err := t.getValidAccessToken(ctx)
	if err != nil {
		return Response{}, err
	}

	endpoint := t.apiURL + "/responses"

	// Prepare payload enforcing store=false and stream=true
	payload := map[string]interface{}{
		"model":  t.model,
		"store":  false, // STRICT INVARIANT: Subscription preview requirement
		"stream": true,  // STRICT INVARIANT: Subscription preview requirement
		"input": []map[string]string{
			{"role": "system", "content": SystemPrompt},
			{"role": "user", "content": FormatUserPrompt(req, isHint)},
		},
	}

	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		return Response{}, fmt.Errorf("failed encoding request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(jsonBytes))
	if err != nil {
		return Response{}, err
	}
	httpReq.Header.Set("Authorization", "Bearer "+accessToken)
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "text/event-stream")

	httpResp, err := t.httpClient.Do(httpReq)
	if err != nil {
		return Response{}, fmt.Errorf("chatgpt api request failed: %w", err)
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		errBytes, _ := io.ReadAll(httpResp.Body)
		return Response{}, fmt.Errorf("chatgpt api error (%d): %s", httpResp.StatusCode, string(errBytes))
	}

	// Stream reader for SSE
	var accumulated strings.Builder
	scanner := bufio.NewScanner(httpResp.Body)
	for scanner.Scan() {
		select {
		case <-ctx.Done():
			return Response{}, ctx.Err()
		default:
		}

		line := scanner.Text()
		if !strings.HasPrefix(line, "data: ") {
			continue
		}
		data := strings.TrimPrefix(line, "data: ")
		if data == "[DONE]" {
			break
		}

		var event struct {
			Type     string `json:"type"`
			Delta    string `json:"delta"`
			Text     string `json:"text"`
			Response struct {
				Output []struct {
					Content []struct {
						Text string `json:"text"`
					} `json:"content"`
				} `json:"output"`
			} `json:"response"`
		}
		if err := json.Unmarshal([]byte(data), &event); err == nil {
			if event.Delta != "" {
				accumulated.WriteString(event.Delta)
			} else if event.Text != "" {
				accumulated.WriteString(event.Text)
			}
		}
	}

	if err := scanner.Err(); err != nil && !errors.Is(err, io.EOF) {
		return Response{}, fmt.Errorf("error reading streamed response: %w", err)
	}

	resultText := strings.TrimSpace(accumulated.String())
	if resultText == "" {
		return Response{}, errors.New("empty response received from chatgpt responses api")
	}

	tokens := (len(resultText) + 3) / 4
	if t.budget != nil {
		_ = t.budget.RecordUsage(tokens)
	}

	return Response{
		Text:        resultText,
		Provider:    "chatgpt_plan",
		TokensUsed:  tokens,
		Fallback:    false,
		GeneratedAt: time.Now().UTC(),
	}, nil
}
