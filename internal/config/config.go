package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	App struct {
		Name     string `mapstructure:"name"`
		Env      string `mapstructure:"env"`
		LogLevel string `mapstructure:"log_level"`
	} `mapstructure:"app"`
	DB struct {
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		User     string `mapstructure:"user"`
		Password string `mapstructure:"password"`
		Name     string `mapstructure:"name"`
		MaxConns int32  `mapstructure:"max_conns"`
	} `mapstructure:"db"`
	API struct {
		Host              string `mapstructure:"host"`
		Port              int    `mapstructure:"port"`
		JWTPrivateKeyPath string `mapstructure:"jwt_private_key_path"`
		JWTPublicKeyPath  string `mapstructure:"jwt_public_key_path"`
	} `mapstructure:"api"`
	SFTP struct {
		ExternalAddress string `mapstructure:"external_address"`
		InternalAddress string `mapstructure:"internal_address"`
		HostKeyPath     string `mapstructure:"host_key_path"`
		DataRoot        string `mapstructure:"data_root"`
	} `mapstructure:"sftp"`
	Sync struct {
		IntervalSeconds int  `mapstructure:"interval_seconds"`
		Workers         int  `mapstructure:"workers"`
		ChecksumEnabled bool `mapstructure:"checksum_enabled"`
	} `mapstructure:"sync"`
}

func Load(path string) (*Config, error) {
	v := viper.New()
	v.SetConfigFile(path)
	v.SetEnvPrefix("SFTP")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}
	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}
	return &cfg, nil
}
