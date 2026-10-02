package tutor

import (
	"context"
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"math/big"
	"net/http"
	"strings"
	"time"
)

func hasPlanScope(scope string) bool {
	for _, s := range strings.Fields(scope) {
		if s == "chatgpt.tokens.use.direct" {
			return true
		}
	}
	return false
}

type openAIIdentity struct {
	Issuer   string          `json:"iss"`
	Subject  string          `json:"sub"`
	Email    string          `json:"email"`
	Audience json.RawMessage `json:"aud"`
	Expires  int64           `json:"exp"`
	Issued   int64           `json:"iat"`
	Nonce    string          `json:"nonce"`
}

// Validate only signed RS256 identity tokens from OpenAI's pinned issuer and JWKS.
func (t *OpenAIChatGPTPlanTutor) validateIDToken(ctx context.Context, token, clientID, nonce string) (*openAIIdentity, error) {
	invalid := errors.New("OpenAI ID token validation failed")
	parts := strings.Split(token, ".")
	if len(parts) != 3 || nonce == "" {
		return nil, invalid
	}
	decode := base64.RawURLEncoding.DecodeString
	headerBytes, err := decode(parts[0])
	if err != nil {
		return nil, invalid
	}
	var header struct {
		Alg string `json:"alg"`
		Kid string `json:"kid"`
	}
	if json.Unmarshal(headerBytes, &header) != nil || header.Alg != "RS256" || header.Kid == "" {
		return nil, invalid
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://auth.openai.com/.well-known/jwks.json", nil)
	if err != nil {
		return nil, err
	}
	resp, err := t.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, invalid
	}
	var jwks struct {
		Keys []struct {
			Kid string `json:"kid"`
			Kty string `json:"kty"`
			N   string `json:"n"`
			E   string `json:"e"`
		} `json:"keys"`
	}
	if json.NewDecoder(io.LimitReader(resp.Body, 1<<20)).Decode(&jwks) != nil {
		return nil, invalid
	}
	signature, err := decode(parts[2])
	if err != nil {
		return nil, invalid
	}
	hash := sha256.Sum256([]byte(parts[0] + "." + parts[1]))
	verified := false
	for _, key := range jwks.Keys {
		if key.Kid != header.Kid || key.Kty != "RSA" {
			continue
		}
		n, e1 := decode(key.N)
		e, e2 := decode(key.E)
		if e1 != nil || e2 != nil || len(e) > 4 || len(e) == 0 {
			continue
		}
		exponent := 0
		for _, b := range e {
			exponent = exponent*256 + int(b)
		}
		if rsa.VerifyPKCS1v15(&rsa.PublicKey{N: new(big.Int).SetBytes(n), E: exponent}, crypto.SHA256, hash[:], signature) == nil {
			verified = true
			break
		}
	}
	if !verified {
		return nil, invalid
	}
	payload, err := decode(parts[1])
	if err != nil {
		return nil, invalid
	}
	var identity openAIIdentity
	if json.Unmarshal(payload, &identity) != nil {
		return nil, invalid
	}
	var audience string
	var audiences []string
	_ = json.Unmarshal(identity.Audience, &audience)
	_ = json.Unmarshal(identity.Audience, &audiences)
	match := audience == clientID
	for _, a := range audiences {
		match = match || a == clientID
	}
	now := time.Now().Unix()
	if !match || identity.Issuer != "https://auth.openai.com" || identity.Subject == "" || identity.Nonce != nonce || identity.Expires <= now || identity.Issued <= 0 || identity.Issued > now+5 {
		return nil, invalid
	}
	return &identity, nil
}
