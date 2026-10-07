package tutor

import (
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
	DefaultOpenAIAuthURL = "https://auth.openai.com/api/accounts"
	DefaultOpenAIAPIURL  = "https://api.openai.com/v1"
	DefaultChatGPTModel  = "gpt-4o-mini"
	DefaultDCRClientID   = "dynamic_agent_client"
	DefaultOAuthScope    = "openid profile email offline_access resource.invoke chatgpt.tokens.use.direct"
)

// OpenAIChatGPTPlanTutor implements the ChatGPT Plus subscription path.
// It uses official "Sign in with ChatGPT" OAuth PKCE and the streamed Responses API.
//
// Invariant (AGENT-CONTRACT):
// Subscription-route preview limitations strictly enforced: store=false, stream=true.
// Tokens are saved securely in local user data dir with 0600 permissions.
type OpenAIChatGPTPlanTutor struct {
	mu           sync.Mutex
	authStore    *AuthStore
	authURL      string
	apiURL       string
	model        string
	clientID     string
	pendingNonce string
	httpClient   *http.Client
	budget       *Budget
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
	if model == "" && cfg.AuthStore != nil {
		model = cfg.AuthStore.ResolveModel(ProviderChatGPTPlan)
	}
	if model == "" {
		model = DefaultChatGPTModel
	}
	clientID := cfg.ClientID
	if clientID == "" && cfg.AuthStore != nil {
		clientID = cfg.AuthStore.ResolveChatGPTClientID()
	}
	if clientID == "" {
		clientID = DefaultDCRClientID
	}
	httpClient := cfg.HTTPClient
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 60 * time.Second}
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

// SetModel updates the active model dynamically.
func (t *OpenAIChatGPTPlanTutor) SetModel(model string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.model = model
}

// GetModel returns the current model name.
func (t *OpenAIChatGPTPlanTutor) GetModel() string {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.model
}

// SetClientID updates the OAuth client ID.
func (t *OpenAIChatGPTPlanTutor) SetClientID(clientID string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.clientID = clientID
}

// GetClientID returns the current OAuth client ID.
func (t *OpenAIChatGPTPlanTutor) GetClientID() string {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.clientID
}

// Name implements Tutor.
func (t *OpenAIChatGPTPlanTutor) Name() string {
	return "OpenAIChatGPTPlanTutor"
}

// Hint implements Tutor.
func (t *OpenAIChatGPTPlanTutor) Hint(ctx context.Context, req Request) (Response, error) {
	return t.callResponsesAPI(ctx, req, true, nil)
}

// Explain implements Tutor.
func (t *OpenAIChatGPTPlanTutor) Explain(ctx context.Context, req Request) (Response, error) {
	return t.callResponsesAPI(ctx, req, false, nil)
}

// HintStream implements StreamingTutor.
func (t *OpenAIChatGPTPlanTutor) HintStream(ctx context.Context, req Request, onUpdate func(StreamUpdate)) (Response, error) {
	return t.callResponsesAPI(ctx, req, true, onUpdate)
}

// ExplainStream implements StreamingTutor.
func (t *OpenAIChatGPTPlanTutor) ExplainStream(ctx context.Context, req Request, onUpdate func(StreamUpdate)) (Response, error) {
	return t.callResponsesAPI(ctx, req, false, onUpdate)
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
	ClientID    string
	once        sync.Once
	closeOnce   sync.Once
	closed      chan struct{}
}

// StartOAuthFlow starts an ephemeral loopback HTTP listener on 127.0.0.1:0 for the OAuth callback.
func (t *OpenAIChatGPTPlanTutor) StartOAuthFlow() (*OAuthFlow, string, error) {
	if t.authStore == nil {
		t.authStore, _ = NewAuthStore("")
	}
	hostID, err := t.authStore.EnsureHostID()
	if err != nil {
		return nil, "", err
	}
	pkce, err := GeneratePKCE()
	if err != nil {
		return nil, "", err
	}

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return nil, "", fmt.Errorf("failed binding ephemeral loopback port: %w", err)
	}

	port := listener.Addr().(*net.TCPAddr).Port
	redirectURI := fmt.Sprintf("http://127.0.0.1:%d/auth/callback", port)

	authURL, err := url.Parse(t.authURL + "/authorize")
	if err != nil {
		listener.Close()
		return nil, "", err
	}

	q := authURL.Query()
	q.Set("response_type", "code")
	clientID := t.authStore.ResolveChatGPTClientID()
	q.Set("client_id", clientID)
	q.Set("ext_agent_host_id", hostID)
	q.Set("resource", DefaultOpenAIAPIURL)
	if clientID == DefaultDCRClientID {
		q.Set("agent_name_hint", "AccountTutor 9000")
	}
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
		ClientID:    clientID,
		closed:      make(chan struct{}),
		codeChan:    make(chan string, 1),
		errChan:     make(chan error, 1),
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/auth/callback", flow.handleCallback)

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
	if r.Method != http.MethodGet {
		http.Error(w, "GET required", 405)
		return
	}
	q := r.URL.Query()
	if q.Get("state") != f.PKCE.State {
		http.Error(w, "Invalid state parameter", 400)
		select {
		case f.errChan <- errors.New("oauth callback state mismatch"):
		default:
		}
		return
	}
	f.once.Do(func() {
		fail := func(msg string) { http.Error(w, msg, 400); f.errChan <- errors.New(msg) }
		if q.Get("error") != "" {
			fail("ChatGPT authorization declined or failed; try again")
			return
		}
		id := q.Get("client_id")
		if f.ClientID == DefaultDCRClientID {
			if !strings.HasPrefix(id, "oaiapp_") {
				fail("registration incomplete: missing issued client ID")
				return
			}
		} else {
			if id != "" && id != f.ClientID {
				fail("callback client ID mismatch")
				return
			}
			id = f.ClientID
		}
		if q.Get("code") == "" {
			fail("missing authorization code")
			return
		}
		f.ClientID = id
		f.tutor.SetClientID(id)
		f.tutor.pendingNonce = f.PKCE.Nonce
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, "<h1>Authorization received</h1><p>Return to AccountTutor 9000 to finish verifying your connection.</p>")
		f.codeChan <- q.Get("code")
	})
}

func (f *OAuthFlow) Complete(ctx context.Context) (*OAuthToken, error) {
	code, err := f.WaitForCallback(ctx)
	if err != nil {
		return nil, err
	}
	token, err := f.tutor.ExchangeCode(ctx, code, f.RedirectURI, f.PKCE.Verifier)
	if err != nil {
		return nil, err
	}
	if old := f.tutor.authStore.GetConfig().ChatGPTPlanToken; old != nil && old.ClientID == f.ClientID && old.Subject != "" && old.Subject != token.Subject {
		return nil, errors.New("returning account identity changed")
	}
	if err := f.tutor.authStore.SetChatGPTPlanToken(token); err != nil {
		return nil, err
	}
	if err := f.tutor.authStore.SetActiveProvider(ProviderChatGPTPlan); err != nil {
		return nil, err
	}
	return token, nil
}

// WaitForCallback waits for the user to complete the browser authorization.
func (f *OAuthFlow) WaitForCallback(ctx context.Context) (string, error) {
	defer f.Close()

	select {
	case <-f.closed:
		return "", errors.New("ChatGPT sign-in cancelled")
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
	f.closeOnce.Do(func() { close(f.closed) })
	if f.server != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		_ = f.server.Shutdown(ctx)
	}
}

// ExchangeCode exchanges an authorization code for access and refresh tokens.
func (t *OpenAIChatGPTPlanTutor) ExchangeCode(ctx context.Context, code string, redirectURI string, verifier string) (*OAuthToken, error) {
	tokenURL := t.authURL + "/oauth/token"

	data := url.Values{}
	if t.GetClientID() == DefaultDCRClientID {
		return nil, errors.New("issued client ID required")
	}
	data.Set("resource", DefaultOpenAIAPIURL)
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
		var errObj struct {
			Error            string `json:"error"`
			ErrorDescription string `json:"error_description"`
		}
		detail := string(bodyBytes)
		if err := json.Unmarshal(bodyBytes, &errObj); err == nil {
			if errObj.ErrorDescription != "" {
				detail = fmt.Sprintf("%s (%s)", errObj.Error, errObj.ErrorDescription)
			} else if errObj.Error != "" {
				detail = errObj.Error
			}
		}
		return nil, fmt.Errorf("token exchange rejected by OpenAI (%d): %s. Verify client ID '%s' or use commercial API key [5]", resp.StatusCode, detail, t.clientID)
	}

	var tokenResp struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int    `json:"expires_in"`
		Scope        string `json:"scope"`
		IDToken      string `json:"id_token"`
	}
	if err := json.Unmarshal(bodyBytes, &tokenResp); err != nil {
		return nil, fmt.Errorf("failed decoding token response: %w", err)
	}

	if tokenResp.AccessToken == "" || tokenResp.ExpiresIn <= 0 {
		return nil, errors.New("invalid token response")
	}
	if !hasPlanScope(tokenResp.Scope) {
		return nil, errors.New("ChatGPT plan usage was not authorized")
	}
	identity, err := t.validateIDToken(ctx, tokenResp.IDToken, t.clientID, t.pendingNonce)
	if err != nil {
		return nil, err
	}
	expiresAt := time.Now().UTC().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)
	token := &OAuthToken{
		AccessToken:  tokenResp.AccessToken,
		RefreshToken: tokenResp.RefreshToken,
		ExpiresAt:    expiresAt,
		ClientID:     t.clientID,
		IDToken:      tokenResp.IDToken,
		Subject:      identity.Subject,
		Email:        identity.Email,
		PlanType:     "chatgpt_plus",
		Scope:        tokenResp.Scope,
	}

	return token, nil
}

// RefreshToken refreshes an expired access token using the refresh token.
func (t *OpenAIChatGPTPlanTutor) RefreshToken(ctx context.Context, token *OAuthToken) (*OAuthToken, error) {
	if token == nil || token.RefreshToken == "" {
		return nil, errors.New("cannot refresh without refresh token")
	}

	tokenURL := t.authURL + "/oauth/token"
	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	if token.ClientID == "" || token.ClientID == DefaultDCRClientID {
		return nil, errors.New("saved issued client ID missing; reconnect ChatGPT")
	}
	data.Set("client_id", token.ClientID)
	data.Set("resource", DefaultOpenAIAPIURL)
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
		var errObj struct {
			Error            string `json:"error"`
			ErrorDescription string `json:"error_description"`
		}
		detail := string(bodyBytes)
		if err := json.Unmarshal(bodyBytes, &errObj); err == nil {
			if errObj.ErrorDescription != "" {
				detail = fmt.Sprintf("%s (%s)", errObj.Error, errObj.ErrorDescription)
			} else if errObj.Error != "" {
				detail = errObj.Error
			}
		}
		return nil, fmt.Errorf("token refresh rejected by OpenAI (%d): %s. Re-authenticate in Tutor Settings ('t')", resp.StatusCode, detail)
	}

	var tokenResp struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int    `json:"expires_in"`
	}
	if err := json.Unmarshal(bodyBytes, &tokenResp); err != nil {
		return nil, fmt.Errorf("failed decoding refresh response: %w", err)
	}

	if tokenResp.AccessToken == "" || tokenResp.ExpiresIn <= 0 {
		return nil, errors.New("invalid refresh response")
	}
	copyToken := *token
	token = &copyToken
	token.AccessToken = tokenResp.AccessToken
	if tokenResp.RefreshToken != "" {
		token.RefreshToken = tokenResp.RefreshToken
	}
	token.ExpiresAt = time.Now().UTC().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)

	if t.authStore != nil {
		if err := t.authStore.SetChatGPTPlanToken(token); err != nil {
			return nil, err
		}
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

	if !hasPlanScope(token.Scope) || token.Subject == "" {
		return "", errors.New("ChatGPT connection needs verified plan permission; reconnect in Tutor Settings")
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
// Text deltas are reported to onUpdate when it is non-nil; reasoning events
// only mark the reply as thinking.
//
// Strictly enforces subscription preview limitations:
// - store: false
// - stream: true
func (t *OpenAIChatGPTPlanTutor) callResponsesAPI(ctx context.Context, req Request, isHint bool, onUpdate func(StreamUpdate)) (Response, error) {
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
		"model":  t.GetModel(),
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

	client := t.httpClient
	if onUpdate != nil {
		client = streamingClient(client)
	}
	httpResp, err := client.Do(httpReq)
	if err != nil {
		return Response{}, fmt.Errorf("chatgpt api request failed: %w", err)
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		return Response{}, fmt.Errorf("chatgpt api error (%d): %s", httpResp.StatusCode, httpErrorBody(httpResp))
	}

	reply := &streamReply{label: "chatgpt responses api", provider: "chatgpt_plan", budget: t.budget, onUpdate: onUpdate}
	defer reply.recordBudget()
	completed := false
	err = readSSE(httpResp.Body, "ChatGPT", &completed, func(data string) (bool, error) {
		if data == "[DONE]" {
			return true, nil
		}
		var event struct {
			Type  string `json:"type"`
			Delta string `json:"delta"`
			Code  string `json:"code"`
			Error struct {
				Code    string `json:"code"`
				Message string `json:"message"`
			} `json:"error"`
			Response struct {
				Error struct {
					Code    string `json:"code"`
					Message string `json:"message"`
				} `json:"error"`
				Usage struct {
					OutputTokens int `json:"output_tokens"`
				} `json:"usage"`
			} `json:"response"`
		}
		if err := json.Unmarshal([]byte(data), &event); err != nil {
			return false, errors.New("invalid Responses stream event")
		}
		switch {
		case event.Type == "response.output_text.delta":
			if event.Delta != "" {
				reply.add(event.Delta, false)
			}
		case strings.HasPrefix(event.Type, "response.reasoning") && strings.HasSuffix(event.Type, ".delta"):
			reply.add("", true)
		case event.Type == "response.completed":
			completed = true
			reply.usage = event.Response.Usage.OutputTokens
			return true, nil
		case event.Type == "error", event.Type == "response.failed", event.Type == "response.incomplete":
			return false, fmt.Errorf("ChatGPT response failed or incomplete: %s %s %s", event.Code, event.Error.Code, event.Response.Error.Code)
		}
		return false, nil
	})
	if err != nil {
		return reply.partial(err)
	}
	if !completed {
		return reply.partial(errors.New("ChatGPT stream ended before response.completed"))
	}
	return reply.result()
}
