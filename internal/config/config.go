package config

import v "github.com/go-ozzo/ozzo-validation/v4"

type PostgresConfig struct {
	Addr     string `config:"addr" yaml:"addr"`
	DB       string `config:"db" yaml:"db"`
	User     string `config:"user" yaml:"user"`
	Password string `config:"password" yaml:"password"`
}

func (c PostgresConfig) Validate() error {
	return v.ValidateStruct(
		&c,
		v.Field(&c.Addr, v.Required),
		v.Field(&c.DB, v.Required),
		v.Field(&c.User, v.Required),
		v.Field(&c.Password, v.Required),
	)
}

type GRPCConfig struct {
	Addr string `config:"addr" yaml:"addr"`
}

func (c GRPCConfig) Validate() error {
	return v.ValidateStruct(
		&c,
		v.Field(&c.Addr, v.Required),
	)
}

type TechServerConfig struct {
	Addr string `config:"addr" yaml:"addr"`
}

func (c *TechServerConfig) Validate() error {
	return v.ValidateStruct(
		&c,
		v.Field(&c.Addr, v.Required),
	)
}

type OTELConfig struct {
	Addr string `config:"addr" yaml:"addr"`
}

func (c *OTELConfig) Validate() error {
	return v.ValidateStruct(
		&c,
		v.Field(&c.Addr, v.Required),
	)
}

type Config struct {
	Postgres   PostgresConfig   `config:"postgres" yaml:"postgres"`
	GRPC       GRPCConfig       `config:"grpc" yaml:"grpc"`
	OTEL       OTELConfig       `config:"otel" yaml:"otel"`
	TechServer TechServerConfig `config:"tech_server" yaml:"tech_server"`
}

func (c Config) Validate() error {
	return v.ValidateStruct(
		&c,
		v.Field(&c.Postgres),
		v.Field(&c.GRPC),
		v.Field(&c.TechServer),
	)
}
