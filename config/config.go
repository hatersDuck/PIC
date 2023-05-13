package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	TelegramToken    string
	Messages         Messages
	DatabasePassword string
	DbPasswordTrade  string

	TimeSleep int
	TestNet   bool
}

type Messages struct {
	Responses
	Errors
	Buttons
	Alert
}

type Responses struct {
	Start          string `mapstructure:"start"`
	Account        string `mapstructure:"account"`
	NewUser        string `mapstructure:"new_user"`
	AddApi         string `mapstructure:"add_api"`
	AddSecret      string `mapstructure:"add_secret"`
	DeleteKey      string `mapstructure:"delete_key"`
	Report         string `mapstructure:"report"`
	ChangeStrategy string `mapstructure:"change_strategy"`
	OnTrade        string `mapstructure:"on_trade"`
	OffTrade       string `mapstructure:"off_trade"`

	Warnings  string `mapstructure:"warnings"`
	NoApiKeys string `mapstructure:"no_api_keys"`
	NoSuccess string `mapstructure:"no_success"`
}

type Errors struct {
	UndefinedCommand string `mapstructure:"undefined_command"`
}

type Buttons struct {
	BtnAccount        string `mapstructure:"account"`
	BtnCurrentOrder   string `mapstructure:"current_order"`
	BtnChangeStrategy string `mapstructure:"change_strategy"`
	BtnOnTrade        string `mapstructure:"on_trade"`
	BtnOffTrade       string `mapstructure:"off_trade"`
	BtnApiKeys        string `mapstructure:"api_keys"`
	BtnApiKeyEmpty    string `mapstructure:"api_key_empty"`
	BtnSecretKeyEmpty string `mapstructure:"secret_key_empty"`
	BtnApiKeyReady    string `mapstructure:"api_key_ready"`
	BtnSecretKeyReady string `mapstructure:"secret_key_ready"`
	BtnBack           string `mapstructure:"back"`
	BtnPrevious       string `mapstructure:"previous"`
	BtnNext           string `mapstructure:"next"`
	BtnReport         string `mapstructure:"report"`
	BtnAccept         string `mapstructure:"accept"`
}

type Alert struct {
	UndefinedButton string `mapstructure:"undefinedButton"`
	NoStrategy      string `mapstructure:"noStrategy"`
	FailedOff       string `mapstructure:"failedOff"`
}

func Init() (*Config, error) {
	viper.AddConfigPath("config")
	viper.SetConfigName("main")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	if err := viper.UnmarshalKey("messages.responses", &cfg.Messages.Responses); err != nil {
		return nil, err
	}
	if err := viper.UnmarshalKey("messages.errors", &cfg.Messages.Errors); err != nil {
		return nil, err
	}
	if err := viper.UnmarshalKey("messages.buttons", &cfg.Messages.Buttons); err != nil {
		return nil, err
	}
	if err := viper.UnmarshalKey("messages.alert", &cfg.Messages.Alert); err != nil {
		return nil, err
	}

	if err := viper.UnmarshalKey("time_sleep", &cfg.TimeSleep); err != nil {
		return nil, err
	}
	if err := viper.UnmarshalKey("test_net", &cfg.TestNet); err != nil {
		return nil, err
	}

	if err := ParseEnv(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func ParseEnv(cfg *Config) error {
	viper.SetConfigName("hid")
	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	cfg.TelegramToken = viper.GetString("BOT_TOKEN")
	cfg.DatabasePassword = viper.GetString("PASSWORD_DB")
	cfg.DbPasswordTrade = viper.GetString("PASSWORD_DB_TRADE")
	return nil
}
