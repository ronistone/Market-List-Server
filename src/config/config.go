package config

import (
	"fmt"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"log"
	"strings"
)

var k *koanf.Koanf

func Init() error {
	k = koanf.New(".")

	if err := k.Load(file.Provider("./config/config.yaml"), yaml.Parser()); err != nil {
		log.Fatalf("Fail to loading config: %v", err)
		return err
	}

	return k.Load(env.Provider("MARKETLIST_", ".", func(s string) string {
		return strings.Replace(strings.ToLower(
			strings.TrimPrefix(s, "MARKETLIST_")), "_", ".", -1)
	}), nil)

}

func GetDatabaseDSN() string {
	host := k.Get("database.host")
	username := k.Get("database.username")
	password := k.Get("database.password")
	port := k.Get("database.port")
	database := k.Get("database.name")
	schema := k.Get("database.schema")

	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s search_path=%s sslmode=disable timezone=UTC",
		host, port, username, password, database, schema)

}
