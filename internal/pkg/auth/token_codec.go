package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

const signedTokenWireParts = 2

type tokenClaims struct {
	Sub int64 `json:"sub"`
	Exp int64 `json:"exp"`
	Pos int32 `json:"pos"`
}

type TokenCodec struct {
	secret []byte
	ttl    time.Duration
}

func NewTokenCodec(secret []byte, ttl time.Duration) (TokenCodec, error) {
	if len(secret) == 0 {
		return TokenCodec{}, errors.New("access token secret is empty")
	}

	if ttl <= 0 {
		return TokenCodec{}, errors.New("access token ttl must be positive")
	}

	return TokenCodec{
		secret: append([]byte(nil), secret...),
		ttl:    ttl,
	}, nil
}

// IssueBearer возвращает полное значение для заголовка Authorization и момент истечения.
func (c *TokenCodec) IssueBearer(
	employeeID int64,
	position int32,
	now time.Time,
) (string, time.Time, error) {
	expiresAt := now.Add(c.ttl)

	payload := tokenClaims{
		Sub: employeeID,
		Exp: expiresAt.Unix(),
		Pos: position,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("marshal access token claims: %w", err)
	}

	mac := hmac.New(sha256.New, c.secret)
	if _, err = mac.Write(body); err != nil {
		return "", time.Time{}, fmt.Errorf("sign access token: %w", err)
	}

	sig := mac.Sum(nil)
	wire := base64.RawURLEncoding.EncodeToString(body) + "." +
		base64.RawURLEncoding.EncodeToString(sig)

	return FormatBearerCredential(wire), expiresAt, nil
}

// ParseEmployeeIDFromBearer разбирает Bearer-строку и возвращает идентификатор сотрудника из claims.
func (c *TokenCodec) ParseEmployeeIDFromBearer(
	raw string,
	now time.Time,
) (int64, error) {
	wire, err := WireFromBearerCredential(raw)
	if err != nil {
		return 0, err
	}

	parts := strings.SplitN(wire, ".", signedTokenWireParts)
	if len(parts) != signedTokenWireParts {
		return 0, ErrInvalidToken
	}

	body, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return 0, ErrInvalidToken
	}

	sig, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return 0, ErrInvalidToken
	}

	mac := hmac.New(sha256.New, c.secret)
	if _, err = mac.Write(body); err != nil {
		return 0, fmt.Errorf("verify access token: %w", err)
	}

	expected := mac.Sum(nil)
	if !hmac.Equal(sig, expected) {
		return 0, ErrInvalidToken
	}

	var claims tokenClaims
	if err = json.Unmarshal(body, &claims); err != nil {
		return 0, ErrInvalidToken
	}

	if claims.Exp <= now.Unix() {
		return 0, ErrInvalidToken
	}

	return claims.Sub, nil
}
