package configs

import (
	"errors"
	"github.com/spf13/viper"
	"github.com/subosito/gotenv"
	"time"
)

var (
	ErrPostgresEmptyHost            = errors.New("postgres: empty host")
	ErrPostgresInvalidPort          = errors.New("postgres: invalid port (from 0 to 65535)")
	ErrPostgresEmptyUsername        = errors.New("postgres: empty username")
	ErrPostgresEmptyPassword        = errors.New("postgres: empty password")
	ErrPostgresEmptyDBName          = errors.New("postgres: empty db name")
	ErrPostgresInvalidSSLMode       = errors.New("postgres: invalid sslmode (disable, allow, prefer, require, verify-ca, verify-full)")
	ErrPostgresMaxConns             = errors.New("postgres: max conns must be positive")
	ErrPostgresMinConns             = errors.New("postgres: min conns must be positive")
	ErrPostgresMaxConnLifetime      = errors.New("postgres: max conn lifetime must be positive")
	ErrPostgresMaxConnIdleTime      = errors.New("postgres: max conn idle time must be positive")
	ErrPostgresParseMaxConnLifetime = errors.New("postgres: max conn lifetime must be represented as 1h2m3s (hours, minutes, seconds)")
	ErrPostgresParseMaxConnIdleTime = errors.New("postgres: max conn idle time must be represented as 1h2m3s (hours, minutes, seconds)")
	ErrServerParseReadTimeout       = errors.New("server: read timeout must be represented as 1h2m3s (hours, minutes, seconds)")
	ErrServerParseWriteTimeout      = errors.New("server: write timeout must be represented as 1h2m3s (hours, minutes, seconds)")
	ErrServerEmptyHost              = errors.New("server: empty host")
	ErrServerInvalidPort            = errors.New("server: invalid port (from 0 to 65535)")
	ErrServerReadTimeout            = errors.New("server: read timeout must be positive")
	ErrServerWriteTimeout           = errors.New("server: write timeout must be positive")
)

type Config struct {
	Postgres Postgres
	Server   Server
}

type Postgres struct {
	Host            string
	Port            int
	Username        string
	Password        string
	DBName          string
	SSLMode         string
	MaxConns        int
	MinConns        int
	MaxConnLifetime time.Duration
	MaxConnIdleTime time.Duration
}

type Server struct {
	Host         string
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

func NewConfig() (*Config, error) {
	viper.AddConfigPath("configs/")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	if err := gotenv.Load("configs/.env"); err != nil {
		return nil, err
	}
	viper.SetEnvPrefix("advert")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	postgres, err := newPostgresConf()
	if err != nil {
		return nil, err
	}
	if err := validatePostgresConf(postgres); err != nil {
		return nil, err
	}

	server, err := newServerConf()
	if err != nil {
		return nil, err
	}
	if err := validateServerConf(server); err != nil {
		return nil, err
	}

	config := Config{
		Postgres: postgres,
		Server:   server,
	}

	return &config, nil
}

func newPostgresConf() (Postgres, error) {
	host := viper.GetString("postgres.host")
	port := viper.GetInt("postgres.port")
	username := viper.GetString("POSTGRES_USERNAME")
	password := viper.GetString("POSTGRES_PASSWORD")
	DBName := viper.GetString("postgres.db_name")
	sslMode := viper.GetString("postgres.sslmode")
	maxConns := viper.GetInt("postgres.max_conns")
	minConns := viper.GetInt("postgres.min_conns")
	maxConnLifetime := viper.GetString("postgres.max_conn_lifetime")
	parsedMaxConnLifetime, err := time.ParseDuration(maxConnLifetime)
	if err != nil {
		return Postgres{}, ErrPostgresParseMaxConnLifetime
	}

	maxConnIdleTime := viper.GetString("postgres.max_conn_idle_time")
	parsedMaxConnIdleTime, err := time.ParseDuration(maxConnIdleTime)
	if err != nil {
		return Postgres{}, ErrPostgresParseMaxConnIdleTime
	}

	return Postgres{
		Host:            host,
		Port:            port,
		Username:        username,
		Password:        password,
		DBName:          DBName,
		SSLMode:         sslMode,
		MaxConns:        maxConns,
		MinConns:        minConns,
		MaxConnLifetime: parsedMaxConnLifetime,
		MaxConnIdleTime: parsedMaxConnIdleTime,
	}, nil
}

func validatePostgresConf(cfg Postgres) error {
	if cfg.Host == "" {
		return ErrPostgresEmptyHost
	}
	if cfg.Port < 0 || cfg.Port > 65535 {
		return ErrPostgresInvalidPort
	}
	if cfg.Username == "" {
		return ErrPostgresEmptyUsername
	}
	if cfg.Password == "" {
		return ErrPostgresEmptyPassword
	}
	if cfg.DBName == "" {
		return ErrPostgresEmptyDBName
	}
	sslModes := map[string]struct{}{
		"disable":     {},
		"allow":       {},
		"prefer":      {},
		"require":     {},
		"verify-ca":   {},
		"verify-full": {},
	}
	if _, ok := sslModes[cfg.SSLMode]; !ok {
		return ErrPostgresInvalidSSLMode
	}
	if cfg.MaxConns <= 0 {
		return ErrPostgresMaxConns
	}
	if cfg.MinConns <= 0 {
		return ErrPostgresMinConns
	}
	if cfg.MaxConnLifetime <= 0 {
		return ErrPostgresMaxConnLifetime
	}
	if cfg.MaxConnIdleTime <= 0 {
		return ErrPostgresMaxConnIdleTime
	}

	return nil
}

func newServerConf() (Server, error) {
	host := viper.GetString("server.host")
	port := viper.GetInt("server.port")
	readTimeout := viper.GetString("server.read_timeout")
	parsedReadTimeout, err := time.ParseDuration(readTimeout)
	if err != nil {
		return Server{}, ErrServerParseReadTimeout
	}

	writeTimeout := viper.GetString("server.write_timeout")
	parsedWriteTimeout, err := time.ParseDuration(writeTimeout)
	if err != nil {
		return Server{}, ErrServerParseWriteTimeout
	}

	return Server{
		Host:         host,
		Port:         port,
		ReadTimeout:  parsedReadTimeout,
		WriteTimeout: parsedWriteTimeout,
	}, nil
}

func validateServerConf(cfg Server) error {
	if cfg.Host == "" {
		return ErrServerEmptyHost
	}
	if cfg.Port < 0 || cfg.Port > 65535 {
		return ErrServerInvalidPort
	}
	if cfg.ReadTimeout <= 0 {
		return ErrServerReadTimeout
	}
	if cfg.WriteTimeout <= 0 {
		return ErrServerWriteTimeout
	}

	return nil
}
