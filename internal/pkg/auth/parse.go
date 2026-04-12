package auth

import (
	"errors"
	"strings"
)

const (
	BearerPrefix             = "Bearer "
	bearerCredentialMaxParts = 2
)

var ErrInvalidToken = errors.New("invalid access token")

// FormatBearerCredential добавляет схему Bearer к непрозрачной части токена (wire).
func FormatBearerCredential(wire string) string {
	wire = strings.TrimSpace(wire)
	if wire == "" {
		return BearerPrefix
	}

	return BearerPrefix + wire
}

// WireFromBearerCredential извлекает непрозрачную часть из значения вида «Bearer <wire>».
func WireFromBearerCredential(s string) (string, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return "", ErrInvalidToken
	}

	parts := strings.SplitN(s, " ", bearerCredentialMaxParts)
	if len(parts) != bearerCredentialMaxParts || !strings.EqualFold(strings.TrimSpace(parts[0]), "Bearer") {
		return "", ErrInvalidToken
	}

	wire := strings.TrimSpace(parts[1])
	if wire == "" {
		return "", ErrInvalidToken
	}

	return wire, nil
}
