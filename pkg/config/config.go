package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	App      App      `yaml:"app"`
	Redis    Redis    `yaml:"redis"`
	Database Database `yaml:"database"`
	Allows   Allows   `yaml:"allows"`
	Kafka    Kafka    `yaml:"kafka"`
	Mongo    Mongo    `yaml:"mongo"`
	Consumer Consumer `yaml:"consumer"`
}

type App struct {
	Name      string `yaml:"name"`
	Port      string `yaml:"port"`
	Host      string `yaml:"host"`
	BaseUrl   string `yaml:"base_url"`
	JwtIssuer string `yaml:"jwt_issuer"`
	JwtSecret string `yaml:"jwt_secret"`
	JwtExpire int    `yaml:"jwt_expire"`
	ClientID  string `yaml:"client_id"`
}

type Consumer struct {
	Name string `yaml:"name"`
	Port string `yaml:"port"`
	Host string `yaml:"host"`
}

type Kafka struct {
	Brokers []string `yaml:"brokers"`
	Topics  struct {
		ATopic string `yaml:"a_topic"`
		BTopic string `yaml:"b_topic"`
		CTopic string `yaml:"c_topic"`
		DTopic string `yaml:"d_topic"`
		ETopic string `yaml:"e_topic"`
		FTopic string `yaml:"f_topic"`
	} `yaml:"topics"`
	MaxRetry       int  `yaml:"max_retry"`
	MaxMessageSize int  `yaml:"max_message_size"`
	ReturnErrors   bool `yaml:"return_errors"`
	ReturnSucces   bool `yaml:"return_succes"`
}

type Redis struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
	Pass string `yaml:"pass"`
}

type Database struct {
	Host    string `yaml:"host"`
	Port    string `yaml:"port"`
	User    string `yaml:"user"`
	Pass    string `yaml:"pass"`
	Name    string `yaml:"name"`
	SslMode string `yaml:"sslmode"`
}

type Mongo struct {
	User     string `yaml:"user"`
	Pass     string `yaml:"pass"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Database string `yaml:"database"`
}
type Allows struct {
	Methods []string `yaml:"methods"`
	Origins []string `yaml:"origins"`
	Headers []string `yaml:"headers"`
}

func InitConfig() *Config {

	//write exist path

	// file_name, _ := filepath.Abs("./config.yaml")

	path, _ := os.Executable()
	fmt.Println("os working directory: ", path)
	var configs Config
	file_name, _ := filepath.Abs("./config.yaml")
	yaml_file, _ := os.ReadFile(file_name)
	yaml.Unmarshal(yaml_file, &configs)
	return &configs
}

var configs *Config

func ReadValue() *Config {
	if configs != nil {
		return configs
	}
	filename, _ := filepath.Abs("./config.yaml")
	// Sanitize the destination path using filepath.Clean
	cleanedDst := filepath.Clean(filename)
	yamlFile, _ := os.ReadFile(cleanedDst)
	err := yaml.Unmarshal(yamlFile, &configs)
	if err != nil {
		log.Fatal("error loading config.yaml ", err)
	}
	return configs
}
