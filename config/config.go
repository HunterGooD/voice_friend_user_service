package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	App struct {
		CertFilePath string `yaml:"cert_file_path"`
	}
	// For rest server
	Server struct {
		Host    string `yaml:"host"`
		Port    string `yaml:"port"`
		Timeout struct {
			Server time.Duration `yaml:"server"`
			Write  time.Duration `yaml:"write"`
			Read   time.Duration `yaml:"read"`
			Idle   time.Duration `yaml:"idle"`
		} `yaml:"timeout"`
	} `yaml:"server"`
	Database struct {
		Host           string `yaml:"host"`
		Port           string `yaml:"port"`
		User           string `yaml:"user"`
		Password       string `yaml:"password"`
		DBName         string `yaml:"dbname"`
		SSLMode        string `yaml:"sslmode"`
		PoolConnection struct {
			MaxOpenConns int           `yaml:"max_open_conns"`
			MaxIdleConns int           `yaml:"max_idle_conns"`
			MaxLifeTime  time.Duration `yaml:"max_life_time"`
		} `yaml:"pool_connection"`
	} `yaml:"database"`
	JWT struct {
		AccessTokenDuration  time.Duration `yaml:"accessTokenDuration"`
		RefreshTokenDuration int           `yaml:"refreshTokenDuration"`
		Issuer               string        `yaml:"issuer"`
	} `yaml:"jwt"`
	Redis struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		DBIdx    int    `yaml:"dbIdx"`
	} `yaml:"redis"`
	Argon2 struct {
		Times   uint32 `yaml:"times"`
		Memory  uint32 `yaml:"memory"`
		KeyLen  uint32 `yaml:"keyLen"`
		SaltLen uint32 `yaml:"saltLen"`
		Threads uint8  `yaml:"threads"`
	} `yaml:"argon2"`
}

func NewConfig(file_name string) (*Config, error) {
	var config Config

	file, err := os.Open(file_name)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	d := yaml.NewDecoder(file)

	if err := d.Decode(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

// TODO: build on cfg param driver
func (c *Config) BuildDSN() string {
	return fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s",
		c.Database.User,
		c.Database.Password,
		c.Database.Host,
		c.Database.Port,
		c.Database.DBName,
		c.Database.SSLMode,
	)
}

func (c *Config) GetRedisAddr() string {
	return fmt.Sprintf("%s:%s", c.Redis.Host, c.Redis.Port)
}

func (c *Config) BuildRedisConn() string {
	return fmt.Sprintf("redis://%s:%s@%s:%s/%d",
		c.Redis.User,
		c.Redis.Password,
		c.Redis.Host,
		c.Redis.Port,
		c.Redis.DBIdx,
	)
}

func (c *Config) GetAddress() string {
	return fmt.Sprintf("%s:%s", c.Server.Host, c.Server.Port)
}

func (c *Config) GetRefreshTokenTime() time.Duration {
	// day * days in config
	return time.Duration(c.JWT.RefreshTokenDuration) * (time.Hour * 24)
}
