package config

import (
	"time"
	"wallet/internal/infrastructure/broker/kafka"
	"wallet/internal/infrastructure/cache/redis"
	"wallet/internal/infrastructure/database/postgres"
	"wallet/internal/utils/httpserver"
	"wallet/internal/utils/metrics"
	"wallet/internal/utils/pprof"
)

type (
	Config struct {
		HTTPServer HTTPServerConfig
		Database   DBConfig
		PProf      PProfConfig
		Metrics    MetricsConfig
		Cache      CacheConfig
		Consumer   ConsumerConfig
		Producer   ProducerConfig
	}

	HTTPServerConfig struct {
		//Port         string        `env:"HTTP_SERVER_PORT" env-default:"8080"`
		ReadTimeout  time.Duration `env:"HTTP_SERVER_READ_TIMEOUT" env-default:"5s"`
		WriteTimeout time.Duration `env:"HTTP_SERVER_WRITE_TIMEOUT" env-default:"5s"`
	}

	DBConfig struct {
		Host     string `env:"DB_HOST" env-default:"localhost"`
		Port     string `env:"DB_PORT" env-default:""`
		Username string `env:"DB_USERNAME" env-default:"postgres"`
		Password string `env:"DB_PASSWORD" env-default:"postgres"`
		Database string `env:"DB_DATABASE" env-default:"postgres"`
	}

	CacheConfig struct {
		URI string `env:"CACHE_URI" env-default:"localhost:6379"`
	}

	ConsumerConfig struct {
		Addr    string `env:"CONSUMER_ADDR" env-default:"localhost:29092"`
		Topic   string `env:"CONSUMER_TOPIC" env-default:"test"`
		GroupID string `env:"CONSUMER_GROUP_ID" env-default:"test"`
	}

	ProducerConfig struct {
		Addr  string `env:"PRODUCER_ADDR" env-default:"localhost:29092"`
		Topic string `env:"PRODUCER_TOPIC" env-default:"test"`
	}

	PProfConfig struct {
		//Port string `env:"PPROF_PORT" env-default:"8081"`
	}

	MetricsConfig struct {
		//Port string `env:"METRICS_PORT" env-default:"8082"`
	}
)

func (srv HTTPServerConfig) Convert() httpserver.ServerConfig {
	return httpserver.ServerConfig{
		Addr:         ":8080",
		ReadTimeout:  srv.ReadTimeout,
		WriteTimeout: srv.WriteTimeout,
	}
}

func (db DBConfig) Convert() postgres.DBConfig {
	return postgres.DBConfig{
		Host:     db.Host,
		Port:     db.Port,
		User:     db.Username,
		Pass:     db.Password,
		Database: db.Database,
	}
}

func (c CacheConfig) Convert() redis.Config {
	return redis.Config{
		URI: c.URI,
	}
}

func (p PProfConfig) Convert() pprof.Config {
	return pprof.Config{
		Addr: ":8081",
	}
}

func (ms MetricsConfig) Convert() metrics.Config {
	return metrics.Config{
		Addr: ":8082",
	}
}

func (c ConsumerConfig) Convert() kafka.ConsumerConfig {
	return kafka.ConsumerConfig{
		Addr:    c.Addr,
		Topic:   c.Topic,
		GroupID: c.GroupID,
	}
}

func (p ProducerConfig) Convert() kafka.ProducerConfig {
	return kafka.ProducerConfig{
		Addr:  p.Addr,
		Topic: p.Topic,
	}
}
