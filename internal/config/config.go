package config

import (
	"errors"
	"fmt"
	"time"

	"github.com/spf13/viper"
)

var ErrMissingRequiredConfig = errors.New("missing required config")

type Config struct {
	AppVersion  string
	ServiceName string

	Server     Server
	PostgresDB PostgresDB
	Auth       Auth
}

type Auth struct {
	AccessTokenSecret string
	AccessTokenTTL    time.Duration
}

type PostgresDB struct {
	HostPort string
	User     string
	Password string
	DBName   string
}

type Server struct {
	GRPCPort        int
	HTTPPort        int
	ShutDownTimeout time.Duration
}

func Get(v *viper.Viper) (Config, error) {
	v.AutomaticEnv()

	const (
		appVersionKey  = "VERSION"
		serviceNameKey = "SERVICE_NAME"
	)

	if !v.IsSet(appVersionKey) {
		return Config{}, fmt.Errorf("%w: %s", ErrMissingRequiredConfig, appVersionKey)
	}

	if !v.IsSet(serviceNameKey) {
		return Config{}, fmt.Errorf("%w: %s", ErrMissingRequiredConfig, serviceNameKey)
	}

	postgresDB, err := getPostgresDB(v)
	if err != nil {
		return Config{}, err
	}

	server, err := getServer(v)
	if err != nil {
		return Config{}, err
	}

	authCfg, err := getAuth(v)
	if err != nil {
		return Config{}, err
	}

	return Config{
		AppVersion:  v.GetString(appVersionKey),
		ServiceName: v.GetString(serviceNameKey),
		Server:      server,
		PostgresDB:  postgresDB,
		Auth:        authCfg,
	}, nil
}

func getServer(v *viper.Viper) (Server, error) {
	const (
		grpcPortKey        = "GRPC_PORT"
		httpPortKey        = "HTTP_PORT"
		shutdownTimeoutKey = "SHUTDOWN_TIMEOUT"
	)

	var server Server

	if !v.IsSet(grpcPortKey) {
		return server, fmt.Errorf("%w: %s", ErrMissingRequiredConfig, grpcPortKey)
	}

	if !v.IsSet(httpPortKey) {
		return server, fmt.Errorf("%w: %s", ErrMissingRequiredConfig, httpPortKey)
	}

	if !v.IsSet(shutdownTimeoutKey) {
		return server, fmt.Errorf("%w: %s", ErrMissingRequiredConfig, shutdownTimeoutKey)
	}

	server.GRPCPort = v.GetInt(grpcPortKey)
	server.HTTPPort = v.GetInt(httpPortKey)
	server.ShutDownTimeout = time.Duration(v.GetInt(shutdownTimeoutKey)) * time.Second

	return server, nil
}

func getPostgresDB(v *viper.Viper) (PostgresDB, error) {
	const (
		hostPortKey = "DB_HOST_PORT"
		userKey     = "DB_USER"
		passwordKey = "DB_PASSWORD"
		nameKey     = "DB_NAME"
	)

	var postgresDB PostgresDB

	if !v.IsSet(hostPortKey) {
		return postgresDB, fmt.Errorf("%w: %s", ErrMissingRequiredConfig, hostPortKey)
	}
	if !v.IsSet(userKey) {
		return postgresDB, fmt.Errorf("%w: %s", ErrMissingRequiredConfig, userKey)
	}
	if !v.IsSet(passwordKey) {
		return postgresDB, fmt.Errorf("%w: %s", ErrMissingRequiredConfig, passwordKey)
	}
	if !v.IsSet(nameKey) {
		return postgresDB, fmt.Errorf("%w: %s", ErrMissingRequiredConfig, nameKey)
	}

	postgresDB.HostPort = v.GetString(hostPortKey)
	postgresDB.User = v.GetString(userKey)
	postgresDB.Password = v.GetString(passwordKey)
	postgresDB.DBName = v.GetString(nameKey)

	return postgresDB, nil
}

func getAuth(v *viper.Viper) (Auth, error) {
	const (
		//nolint:gosec // env var for viper is not a secret
		authAccessTokenSecretKey = "AUTH_ACCESS_TOKEN_SECRET"
		ttlSecKey                = "AUTH_ACCESS_TOKEN_TTL_SEC"
	)

	var auth Auth

	if !v.IsSet(authAccessTokenSecretKey) {
		return auth, fmt.Errorf("%w: %s", ErrMissingRequiredConfig, authAccessTokenSecretKey)
	}

	if !v.IsSet(ttlSecKey) {
		return auth, fmt.Errorf("%w: %s", ErrMissingRequiredConfig, ttlSecKey)
	}

	auth.AccessTokenSecret = v.GetString(authAccessTokenSecretKey)
	auth.AccessTokenTTL = time.Duration(v.GetInt(ttlSecKey)) * time.Second

	return auth, nil
}
