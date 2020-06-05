package server

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	Server *server `yaml:"server"`
	Log    *log    `yaml:"logger"`
	DB     *db     `yaml:"db"`
	Queue  *queue  `yaml:"queue"`
}

type server struct {
	Addr string `yaml:"port"`
}

type log struct {
	Level string `yaml:"level"`
}

type db struct {
	Driver   string `yaml:"driver"`
	Host     string `yaml:"host"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
	Port     string `yaml:"port"`
}

type queue struct {
	Network string `yaml:"network"`
	Host    string `yaml:"host"`
	Port    string `yaml:"port"`
}

func NewConfig() (*Config, error) {
	body, err := ioutil.ReadFile("configs/app.yml")
	if err != nil {
		return nil, err
	}
	config := &Config{}
	if err = yaml.Unmarshal(body, config); err != nil {
		return nil, err
	}
	return config, nil
}
