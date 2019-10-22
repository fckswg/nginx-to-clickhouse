package config

import (
	"flag"
	"github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
	"gopkg.in/go-playground/validator.v9"
	"os"
)

var (
	confPath string
	validate *validator.Validate
)

type Config struct {
	App struct {
		BatchSize int `mapstructure:"batch_size" validate:"required"`
		LogPath string `mapstructure:"log_path"`
	} `mapstructure:"app" validate:"required"`
	Clickhouse struct {
		Connection struct {
			Host string `mapstructure:"host" validate:"required"`
			Port int `mapstructure:"port" validate:"required"`
			Db string `mapstructure:"database" validate:"required"`
			Table string `mapstructure:"table" validate:"required"`
		} `mapstructure:"connection" validate:"required"`
		Credentials struct {
			User string `mapstructure:"user" validate:"required"`
			Password string `mapstructure:"password" validate:"required"`
		} `mapstructure:"credentials" validate:"required"`
	} `mapstructure:"clickhouse" validate:"required"`
	Nginx struct {
		WarnCount int `mapstructure:"warn_count" validate:"required"`
		LogPath string `mapstructure:"log_path" validate:"required"`
	} `mapstructure:"nginx" validate:"required"`
	Notifications struct {
		Telegram struct {
			BotToken string `mapstructure:"bot_token"`
			ChatId string `mapstructure:"chat_id"`
		} `mapstructure:"telegram"`
	} `mapstructure:"notifications"`
}


func init() {
	flag.StringVar(&confPath, "config_directory", "/etc/nginx-clickhouse/", "Config path.")
	flag.Parse()
}


func Read() *Config {

	c := Config{}

	v := viper.New()
	v.SetConfigName("config")
	v.AddConfigPath(confPath)

	logrus.Info("Reading config file: " + confPath)
	if err := v.ReadInConfig(); err != nil {
		logrus.Fatalf("couldn't load config: %s\n", err.Error())
	}

	if err := v.Unmarshal(&c); err != nil {
		logrus.Fatalf("couldn't unmarshal config: %s\n", err.Error())
	}

	return &c
}

func Validate(c *Config) error {
	validate := validator.New()
	logrus.Info("Validating config")
	err := validate.Struct(c)
	if err == nil {
		logrus.Info("Config validation completed")
	}
	return err
}

func SetLogger(c *Config) {
	if c.App.LogPath != "" {
		logger := logrus.New()
		f, err := os.OpenFile(c.App.LogPath, os.O_WRONLY | os.O_CREATE, 0755)
		if err != nil {
			logger.Fatalf("Cant open logfile for app: %s", err.Error())
		}
		logger.SetOutput(f)
	}
}
