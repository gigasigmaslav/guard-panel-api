package grpc

import (
	"context"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// AuthorizationHeaderFromContext возвращает сырое значение заголовка Authorization
// (ожидается «Bearer <access_token>»). Разбор Bearer и проверка токена — в internal/pkg/auth и usecase/auth.
func AuthorizationHeaderFromContext(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", status.Error(codes.Unauthenticated, "missing authorization metadata")
	}

	var vals []string
	for _, k := range []string{"authorization", "grpcgateway-authorization"} {
		vals = md.Get(k)
		if len(vals) > 0 {
			break
		}
	}

	if len(vals) == 0 {
		return "", status.Error(codes.Unauthenticated, "missing authorization header")
	}

	raw := strings.TrimSpace(vals[0])
	if raw == "" {
		return "", status.Error(codes.Unauthenticated, "empty authorization header")
	}

	return raw, nil
}
