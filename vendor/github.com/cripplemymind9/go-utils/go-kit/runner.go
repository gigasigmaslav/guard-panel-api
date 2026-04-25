package gokit

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	SHUTDOWN_TIMEOUT    = 30 * time.Second
	RECONFIGURE_TIMEOUT = 10 * time.Second
)

const (
	MaxMessageSize = 4 * 1024 * 1024
)

type runner struct {
	config Config

	streamInterceptors []grpc.StreamServerInterceptor
	unaryInterceptors  []grpc.UnaryServerInterceptor

	serverMuxOptions []runtime.ServeMuxOption
	httpMiddlewares  []func(http.Handler) http.Handler
}

// NewRunner создает экземпляр Runner с конфигурацией по умолчанию
// Возвращает: указатель на новый Runner
func NewRunner() *runner {
	cfg := GetConfig(viper.New())

	return &runner{
		config: cfg,
	}
}

// Run запускает приложение, обрабатывает системные сигналы и управляет жизненным циклом
// Принимает: реализацию интерфейса App
// Возвращает: ошибку запуска или выполнения
func (r *runner) Run(appImpl App) error {
	log.Info().
		Any("config", r.config).
		Msg("starting runner")

	c := make(chan os.Signal, 1)

	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	go func() {
		for {
			sig := <-c

			log.Warn().
				Msgf("got signal: %s", sig)

			switch sig {
			case syscall.SIGINT:
				appImpl.Shutdown(SHUTDOWN_TIMEOUT)
			case syscall.SIGTERM:
				appImpl.Shutdown(SHUTDOWN_TIMEOUT)
			case syscall.SIGHUP:
				appImpl.Reconfigure(RECONFIGURE_TIMEOUT)
			}
		}
	}()

	opts := []grpc.ServerOption{
		grpc.ChainStreamInterceptor(r.streamInterceptors...),
		grpc.ChainUnaryInterceptor(r.unaryInterceptors...),
		grpc.MaxRecvMsgSize(MaxMessageSize),
		grpc.MaxSendMsgSize(MaxMessageSize),
	}
	grpcServer := grpc.NewServer(opts...)

	if err := r.initServer(appImpl, grpcServer, r.config); err != nil {
		log.Error().
			Err(err).
			Msg("error of server initialization")

		return err
	}

	return appImpl.Run()
}

// initServer настраивает и запускает gRPC сервер и HTTP шлюз (gateway)
// Принимает: реализацию App, настроенный gRPC сервер, конфигурацию App
// Возвращает: ошибку инициализации
func (r *runner) initServer(
	appImpl App,
	grpcServer *grpc.Server,
	cfg Config,
) error {
	appImpl.RegisterGRPCServices(grpcServer)

	server, err := NewServer(grpcServer, cfg, r.serverMuxOptions...)
	if err != nil {
		return fmt.Errorf("server creation error: %w", err)
	}

	server.StartGRPCServer()

	mux := server.GetMux()

	endpoint := fmt.Sprintf("0.0.0.0:%d", cfg.GRPCPort)

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	ctx := context.Background()
	if err := appImpl.RegisterHandlersFromEndpoint(ctx, mux, endpoint, opts); err != nil {
		return fmt.Errorf("failed to register HTTP handlers: %w", err)
	}

	r.setHTTPHandler(server, mux)

	if err := server.StartGateway(appImpl); err != nil {
		return fmt.Errorf("HTTP gateway start error: %w", err)
	}

	return nil
}

func (r *runner) setHTTPHandler(server *Server, mux *runtime.ServeMux) {
	if len(r.httpMiddlewares) == 0 {
		return
	}

	var handler http.Handler = mux
	for i := len(r.httpMiddlewares) - 1; i >= 0; i-- {
		handler = r.httpMiddlewares[i](handler)
	}

	server.SetHTTPHandler(handler)
}

func (r *runner) ApplyOptions(opts ...Option) {
	for _, opt := range opts {
		opt(r)
	}
}

type Option func(*runner)

func WithStreamInterceptor(s grpc.StreamServerInterceptor) Option {
	return func(r *runner) {
		r.streamInterceptors = append(r.streamInterceptors, s)
	}
}

func WithUnaryInterceptor(u grpc.UnaryServerInterceptor) Option {
	return func(r *runner) {
		r.unaryInterceptors = append(r.unaryInterceptors, u)
	}
}

func WithHTTPMiddleware(m func(http.Handler) http.Handler) Option {
	return func(r *runner) {
		r.httpMiddlewares = append(r.httpMiddlewares, m)
	}
}

type Runtime interface {
	ApplyOptions(opts ...Option)
}
