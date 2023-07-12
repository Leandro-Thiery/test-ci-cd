package config

import (
	"log"
	"sync"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

type Config struct {
	mu           sync.Mutex
	serverConfig ServerConfig
}

type ServerConfig struct {
	Port int    `koanf:"port"`
	Host string `koanf:"host"`
}

func InitConfig() *Config {
	k := koanf.New(".")
	f := file.Provider("config.yml")
	if err := k.Load(f, yaml.Parser()); err != nil {
		log.Fatalf("error loading config file: %v", err)
	}

	cfg := Config{}

	k.UnmarshalWithConf("server", &cfg.serverConfig, koanf.UnmarshalConf{Tag: "koanf"})

	f.Watch(func(event interface{}, err error) {
		if err != nil {
			log.Printf("watch error: %v", err)
			return
		}
		log.Println("config changed. Reloading ...")
		k = koanf.New(".")

		cfg.mu.Lock()
		defer cfg.mu.Unlock()

		if err := k.Load(f, yaml.Parser()); err != nil {
			log.Fatalf("error loading config file: %v", err)
		}
		k.UnmarshalWithConf("server", &cfg.serverConfig, koanf.UnmarshalConf{Tag: "koanf"})

		k.Print()
	})
	return &cfg
}

func (config *Config) GetConfig() *Config {
	config.mu.Lock()
	defer config.mu.Unlock()
	return config
}

func (config *Config) GetServerConfig() *ServerConfig {
	return &config.serverConfig
}
