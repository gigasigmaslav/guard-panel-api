package grpc

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/gigasigmaslav/guard-panel-api/internal/pkg/auth"
)

//nolint:gochecknoglobals // fixed set of RPCs without authorization
var publicMethods = map[string]bool{
	"/guard.v1.GuardPanelService/SignUp": true,
	"/guard.v1.GuardPanelService/SignIn": true,
}

const employeeIDKey = "employee_id"

// AuthInterceptor создаёт unary interceptor для проверки Bearer токенов
// Защищает все методы кроме публичных (SignUp, SignIn)
//
// Извлекает employee_id из токена и сохраняет в контексте.
func AuthInterceptor(tokenCodec auth.TokenCodec) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		if isPublic, exists := publicMethods[info.FullMethod]; exists && isPublic {
			return handler(ctx, req)
		}

		bearerToken, err := AuthorizationHeaderFromContext(ctx)
		if err != nil {
			return nil, err
		}

		employeeID, err := tokenCodec.ParseEmployeeIDFromBearer(bearerToken, time.Now())
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "invalid access token")
		}

		//nolint:revive,staticcheck // context-keys-type / SA1029 key - private package constant
		ctx = context.WithValue(ctx, employeeIDKey, employeeID)

		return handler(ctx, req)
	}
}

// GetEmployeeIDFronCtx извлекает employee_id из контекста
//
// Возвращает 0 если employee_id не найден (например, для публичных методов).
func GetEmployeeIDFronCtx(ctx context.Context) int64 {
	if employeeID, ok := ctx.Value(employeeIDKey).(int64); ok {
		return employeeID
	}
	return 0
}
