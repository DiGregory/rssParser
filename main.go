package main

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"github.com/DiGregory/s7testTask/storage"
	"fmt"
	"github.com/DiGregory/s7testTask/parser"
)

type config struct {
	Parser struct {
		Links    []string `yaml:"links" env:"PARSER_LINKS"`
		Threads  int      `yaml:"threads" env:"PARSER_THREADS"`
		Timeout  int      `yaml:"timeout" env:"PARSER_TIMEOUT"`
		KeyWords []string `yaml:"key_words" env:"PARSER_KEY_WORDS"`
	} `yaml:"parser"`
	DB struct {
		Host     string `yaml:"host" env:"POSTGRES_HOST" env-default:"db"  `
		Port     string `yaml:"port" env:"POSTGRES_PORT" env-default:"5432" `
		User     string `yaml:"username" env:"POSTGRES_USER" env-default:"db"  `
		Password string `yaml:"password" env:"POSTGRES_PASSWORD" env-default:"1234"`
		Name     string `yaml:"db-name" env:"POSTGRES_DB"  env-default:"dev"`
	} `yaml:"db"`
	observer struct{
		Host     string `yaml:"host"   `
		Port     string `yaml:"port"   `
	}`yaml:"observer"`
}

func main() {
	var cfg config

	err := cleanenv.ReadConfig("config.yaml", &cfg)
	if err != nil {
		log.Fatal("error: ", err)
	}

	conn, err := storage.NewConn("postgres",
		cfg.DB.Host, cfg.DB.Port, cfg.DB.User, cfg.DB.Password, cfg.DB.Name)
	if err != nil {
		fmt.Println(err)
	}
	rssParser := parser.NewPool(cfg.Parser.Threads, cfg.Parser.Timeout, cfg.Parser.KeyWords, conn)
	rssParser.Start(cfg.Parser.Links)

}
