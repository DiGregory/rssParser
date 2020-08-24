package main

import (
	"fmt"
	"github.com/DiGregory/s7testTask/parser"
	"github.com/DiGregory/s7testTask/storage"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/sirupsen/logrus"
	"github.com/DiGregory/s7testTask/observer"
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
	Observer struct {
		Host string `yaml:"host"   `
		Port string `yaml:"port"   `
	} `yaml:"observer"`
}

func main() {
	var cfg config
	err := cleanenv.ReadConfig("config.yaml", &cfg)
	if err != nil {
		logrus.WithError(err).Fatal("Read config error")
	}

	conn, err := storage.NewConn("postgres",
		cfg.DB.Host, cfg.DB.Port, cfg.DB.User, cfg.DB.Password, cfg.DB.Name)
	if err != nil {
		logrus.WithError(err).Fatal("DB connection error")
	}
	logrus.Info("Connection with db was set up")
	newsStorage := storage.NewNewsStorage(conn)

	rssParser := parser.NewPool(
		cfg.Parser.Threads,
		cfg.Parser.Timeout,
		cfg.Parser.KeyWords,
		storage.NewNewsStorage(conn),
	)
	go rssParser.Start(cfg.Parser.Links)

	err = observer.Start(
		fmt.Sprintf("%s:%s", cfg.Observer.Host, cfg.Observer.Port),
		observer.NewNewsService(newsStorage),
	)
	if err != nil {
		logrus.WithError(err).Fatal("Observer start error")
	}
	//observer.ClientStart(fmt.Sprintf("%s:%s", cfg.Observer.Host, cfg.Observer.Port))

}
