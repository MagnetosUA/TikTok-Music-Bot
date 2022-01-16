package config

import (
	"github.com/spf13/viper"
	"os"
)

type Config struct {
	// env
	TelegramToken string
	ShazamHost    string
	ShazamKey     string
	// yml
	ShazamResourceURL string `mapstructure:"shazam_resource_url"`
	TikTokHostURL     string `mapstructure:"tik_tok_host_url"`
	ZkHostURL         string `mapstructure:"zk_host_url"`
	ZkSearchURL       string `mapstructure:"zk_search_url"`
	ZkDownloadURL     string `mapstructure:"zk_download_url"`
}

func Init() (*Config, error) {
	viper.AddConfigPath("configs")
	viper.SetConfigName("main")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config Config

	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	if err := parseEnv(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

func parseEnv(config *Config) error {
	os.Setenv("TELEGRAM_TOKEN", "1675374273:AAHCAYwVhhka8Qf-bYWFGRMViV5t2eZcPAE")
	os.Setenv("SHAZAM_KEY", "0dcc409e08msh55fe6be19bff0bcp192cf3jsn0d091d334fda")
	os.Setenv("SHAZAM_HOST", "shazam.p.rapidapi.com")

	if err := viper.BindEnv("telegram_token"); err != nil {
		return err
	}
	if err := viper.BindEnv("shazam_host"); err != nil {
		return err
	}
	if err := viper.BindEnv("shazam_key"); err != nil {
		return err
	}

	config.TelegramToken = viper.GetString("telegram_token")
	config.ShazamHost = viper.GetString("shazam_host")
	config.ShazamKey = viper.GetString("shazam_key")

	return nil
}
