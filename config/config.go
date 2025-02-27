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
		Addr         string        `env:"HTTP_SERVER_ADDR" envDefault:":8080"`
		ReadTimeout  time.Duration `env:"HTTP_SERVER_READ_TIMEOUT" envDefault:"5s"`
		WriteTimeout time.Duration `env:"HTTP_SERVER_WRITE_TIMEOUT" envDefault:"5s"`
	}

	DBConfig struct {
		Host     string `env:"DB_HOST" envDefault:"localhost"`
		Port     string `env:"DB_PORT" envDefault:""`
		Username string `env:"DB_USERNAME" envDefault:"postgres"`
		Password string `env:"DB_PASSWORD" envDefault:"postgres"`
		Database string `env:"DB_DATABASE" envDefault:"postgres"`
	}

	CacheConfig struct {
		URI string `env:"CACHE_URI" envDefault:"localhost:6379"`
	}

	ConsumerConfig struct {
		Addr    string `env:"CONSUMER_ADDR" envDefault:"localhost:29092"`
		Topic   string `env:"CONSUMER_TOPIC" envDefault:"test"`
		GroupID string `env:"CONSUMER_GROUP_ID" envDefault:"test"`
	}

	ProducerConfig struct {
		Addr  string `env:"PRODUCER_ADDR" envDefault:"localhost:29092"`
		Topic string `env:"PRODUCER_TOPIC" envDefault:"test"`
	}

	PProfConfig struct {
		Addr string `env:"PPROF_ADDR" envDefault:":8081"`
	}

	MetricsConfig struct {
		Addr string `env:"METRICS_ADDR" envDefault:":8082"`
	}
)

func (srv HTTPServerConfig) Convert() httpserver.ServerConfig {
	return httpserver.ServerConfig{
		Addr:         srv.Addr,
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
		Addr: p.Addr,
	}
}

func (ms MetricsConfig) Convert() metrics.Config {
	return metrics.Config{
		Addr: ms.Addr,
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
