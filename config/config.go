package config

import (
	"math/big"
	"time"
)

type Option func(*Config)

type ConfigurationProvider interface {
	LoadConfig(filename string) Configuration
	ParseString(cfgStr string) Configuration
}

type Configuration interface {
	GetBoolean(path string, defaultVal ...bool) bool
	GetByteSize(path string) *big.Int
	GetInt32(path string, defaultVal ...int32) int32
	GetInt64(path string, defaultVal ...int64) int64
	GetString(path string, defaultVal ...string) string
	GetFloat32(path string, defaultVal ...float32) float32
	GetFloat64(path string, defaultVal ...float64) float64
	GetTimeDuration(path string, defaultVal ...time.Duration) time.Duration
	GetTimeDurationInfiniteNotAllowed(path string, defaultVal ...time.Duration) time.Duration
	GetBooleanList(path string) []bool
	GetFloat32List(path string) []float32
	GetFloat64List(path string) []float64
	GetInt32List(path string) []int32
	GetInt64List(path string) []int64
	GetByteList(path string) []byte
	GetStringList(path string) []string
	GetConfig(path string) Configuration
	WithFallback(fallback Configuration) Configuration
	HasPath(path string) bool
	Keys() []string
	IsEmpty() bool
	IsObject(path string) bool
	IsArray(path string) bool

	String() string
}

type Config struct {
	ConfigFile   string
	ConfigString string

	Configuration
	configProvider ConfigurationProvider
}

func NewConfig(opts ...Option) *Config {
	conf := &Config{}
	conf.init(opts...)
	return conf
}

func (p *Config) init(opts ...Option) {
	for i := 0; i < len(opts); i++ {
		opts[i](p)
	}

	if p.configProvider == nil {
		p.configProvider = &HOCONConfigProvider{}
	}

	conf := p.Configuration
	if conf == nil {
		conf = p.configProvider.ParseString("")
	}

	var confString, confFile Configuration

	if len(p.ConfigFile) > 0 {
		confFile = p.configProvider.LoadConfig(p.ConfigFile)
		conf = conf.WithFallback(confFile)
	}

	if len(p.ConfigString) > 0 {
		confString = p.configProvider.ParseString(p.ConfigString)
		conf = conf.WithFallback(confString)
	}

	p.Configuration = conf
}

func (p *Config) String() string {

	if p == nil || p.Configuration == nil {
		return ""
	}

	return p.Configuration.String()
}

func (p *Config) WithFallback(fallback Configuration) Configuration {

	if fallback == nil {
		return p
	}

	if p.Configuration == nil {
		return p
	}

	switch v := fallback.(type) {
	case *Config:
		{
			return p.WithFallback(v.Configuration)
		}
	default:
		p.Configuration = p.Configuration.WithFallback(v)
	}

	return p
}

func ConfigFile(fn string) Option {
	return func(o *Config) {
		o.ConfigFile = fn
	}
}

func ConfigString(str string) Option {
	return func(o *Config) {
		o.ConfigString = str
	}
}

func WithConfig(conf Configuration) Option {
	return func(o *Config) {
		o.Configuration = conf
	}
}

func ConfigProvider(provider ConfigurationProvider) Option {
	return func(o *Config) {
		o.configProvider = provider
	}
}

func ConfigProviderByName(driverName string) Option {
	return func(o *Config) {

		provider, err := NewConfigurationProvider(driverName)
		if err != nil {
			panic("create configration provider of " + driverName + " failure. error:" + err.Error())
		}

		o.configProvider = provider
	}
}
