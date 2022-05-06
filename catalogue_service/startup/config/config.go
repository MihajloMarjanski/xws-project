package config

import "os"

type Config struct {
	Port            string
	CatalogueDBHost string
	CatalogueDBPort string
}

func NewConfig() *Config {
	return &Config{
		Port:            os.Getenv("CATALOGUE_SERVICE_PORT"),
		CatalogueDBHost: os.Getenv("CATALOGUE_DB_HOST"),
		CatalogueDBPort: os.Getenv("CATALOGUE_DB_PORT"),
	}
}
