package app

import (
	"context"
	"errors"
	"fmt"
	"time"

	gokit "github.com/cripplemymind9/go-utils/go-kit"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"

	auth "github.com/gigasigmaslav/guard-panel-api/internal/pkg/auth"
	grpcpkg "github.com/gigasigmaslav/guard-panel-api/internal/pkg/grpc"
	middleware "github.com/gigasigmaslav/guard-panel-api/internal/pkg/http"
	"github.com/gigasigmaslav/guard-panel-api/internal/pkg/postgres"
	"github.com/gigasigmaslav/guard-panel-api/migrations"

	"github.com/gigasigmaslav/guard-panel-api/internal/adapter/repo"
	"github.com/gigasigmaslav/guard-panel-api/internal/config"
	"github.com/gigasigmaslav/guard-panel-api/internal/server"

	"github.com/gigasigmaslav/guard-panel-api/internal/domain/usecase/analytics"
	authuc "github.com/gigasigmaslav/guard-panel-api/internal/domain/usecase/auth"
	"github.com/gigasigmaslav/guard-panel-api/internal/domain/usecase/comment"
	"github.com/gigasigmaslav/guard-panel-api/internal/domain/usecase/employee"
	"github.com/gigasigmaslav/guard-panel-api/internal/domain/usecase/office"
	"github.com/gigasigmaslav/guard-panel-api/internal/domain/usecase/refund"
	"github.com/gigasigmaslav/guard-panel-api/internal/domain/usecase/task"
	vuddecision "github.com/gigasigmaslav/guard-panel-api/internal/domain/usecase/vud-decision"
)

type App struct {
	gokit.App

	ctx    context.Context
	cancel context.CancelFunc
	cfg    config.Config

	server *server.Server
}

func New(ctx context.Context, runtime gokit.Runtime, cfg config.Config) (*App, error) {
	repoStorage, err := getPostgresStorage(ctx, cfg.PostgresDB)
	if err != nil {
		return nil, err
	}

	tokenCodec, err := auth.NewTokenCodec(
		[]byte(cfg.Auth.AccessTokenSecret),
		cfg.Auth.AccessTokenTTL,
	)
	if err != nil {
		return nil, fmt.Errorf("auth access token codec: %w", err)
	}

	passwordHasher := auth.BcryptHasher{}

	serverDependencies := getGRPCServerDependencies(
		repoStorage,
		tokenCodec,
		passwordHasher,
	)

	server := server.New(cfg, serverDependencies)

	runtimeOpts := []gokit.Option{
		gokit.WithUnaryInterceptor(grpcpkg.AuthInterceptor(tokenCodec)),
		gokit.WithHTTPMiddleware(middleware.CORS(cfg.AllowedCORSOrigin)),
	}

	runtime.ApplyOptions(runtimeOpts...)

	ctx, cancel := context.WithCancel(ctx)

	return &App{
		ctx:    ctx,
		cancel: cancel,
		server: server,
		cfg:    cfg,
	}, nil
}

func (a *App) Run() error {
	<-a.ctx.Done()
	log.Info().Msg("bye")

	if errors.Is(a.ctx.Err(), context.Canceled) {
		return fmt.Errorf("context cancelled run app err: %w", a.ctx.Err())
	}

	return nil
}

func (a *App) Shutdown(dur time.Duration) error {
	time.Sleep(dur)
	a.cancel()

	return nil
}

func (a *App) RegisterGRPCServices(server grpc.ServiceRegistrar) {
	a.server.RegisterServices(server)
}

func (a *App) RegisterHandlersFromEndpoint(
	ctx context.Context,
	mux *runtime.ServeMux,
	endpoint string,
	opts []grpc.DialOption,
) error {
	return a.server.RegisterHandlersFromEndPoint(ctx, mux, endpoint, opts)
}

func getPostgresStorage(ctx context.Context, cfg config.PostgresDB) (*repo.Storage, error) {
	db, err := postgres.New(ctx, postgres.Config{
		DBName:   cfg.DBName,
		HostPort: cfg.HostPort,
		Username: cfg.User,
		Password: cfg.Password,
	})
	if err != nil {
		return nil, fmt.Errorf("new db instance err: %w", err)
	}

	stdDB := stdlib.OpenDBFromPool(db.Pool)

	if err = migrations.Up(ctx, stdDB); err != nil {
		return nil, fmt.Errorf("sheme migrations err: %w", err)
	}

	return repo.NewStorage(db), err
}

func getGRPCServerDependencies(
	repoStorage *repo.Storage,
	tokenCodec auth.TokenCodec,
	passwordHasher auth.BcryptHasher,
) *server.Dependencies {
	createCommentUC := comment.NewCreateCommentUseCase(repoStorage)
	deleteCommentUC := comment.NewDeleteCommentUseCase(repoStorage)
	createRefundUC := refund.NewCreateRefundUseCase(repoStorage)
	createVUDDecisionUC := vuddecision.NewCreateVUDDecisionUseCase(repoStorage)
	updateVUDDecisionUC := vuddecision.NewUpdateVUDDecisionUseCase(repoStorage)
	createEmployeeUC := employee.NewCreateEmployeeUseCase(repoStorage)
	updateEmployeeUC := employee.NewUpdateEmployeeUseCase(repoStorage)
	deleteEmployeeUC := employee.NewDeleteEmployeeUseCase(repoStorage)
	searchEmployeeUC := employee.NewSearchEmployeeUseCase(repoStorage)
	createOfficeUC := office.NewCreateOfficeUseCase(repoStorage)
	updateOfficeUC := office.NewUpdateOfficeUseCase(repoStorage)
	deleteOfficeUC := office.NewDeleteOfficeUseCase(repoStorage)
	searchOfficeUC := office.NewSearchOfficeUseCase(repoStorage)
	createTaskUC := task.NewCreateTaskUseCase(repoStorage, repoStorage, repoStorage)
	updateTaskUC := task.NewUpdateTaskUseCase(repoStorage)
	searchTasksUC := task.NewSearchTasksUseCase(repoStorage, repoStorage)
	signUpUC := authuc.NewSignUpUseCase(repoStorage, repoStorage, tokenCodec, passwordHasher)
	signInUC := authuc.NewSignInUseCase(repoStorage, repoStorage, tokenCodec, passwordHasher)
	whoAmIUC := authuc.NewWhoAmIUseCase(repoStorage)
	getAnalyticsUC := analytics.NewGetAnalyticsUseCase(repoStorage)
	getTaskByIDUC := task.NewGetTaskByIDUseCase(
		repoStorage,
		repoStorage,
		repoStorage,
		repoStorage,
		repoStorage,
		repoStorage,
	)

	return server.NewDependencies(
		createCommentUC,
		deleteCommentUC,
		createRefundUC,
		createVUDDecisionUC,
		updateVUDDecisionUC,
		createEmployeeUC,
		updateEmployeeUC,
		deleteEmployeeUC,
		searchEmployeeUC,
		createOfficeUC,
		updateOfficeUC,
		deleteOfficeUC,
		searchOfficeUC,
		createTaskUC,
		updateTaskUC,
		getTaskByIDUC,
		searchTasksUC,
		signUpUC,
		signInUC,
		whoAmIUC,
		getAnalyticsUC,
	)
}
