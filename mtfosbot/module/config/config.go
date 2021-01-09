package config

import (
	"errors"
	"io/ioutil"
	"os"
	"path"

	"/golang/mtfosbot/module/utils"

	yaml "gopkg.in/yaml.v2"
)

// Config -
type Config struct {
	Port         int    `yaml:"port"`
	URL          string `yaml:"url"`
	SelfKey      string `yaml:"self_key"`
	ImageRoot    string `yaml:"image_root"`
	LogImageRoot string `yaml:"log_image_root"`
	Line         struct {
		Secret string `yaml:"secret"`
		Access string `yaml:"access"`
	} `yaml:"line"`
	Twitch struct {
		ClientID     string `yaml:"client_id"`
		ClientSecret string `yaml:"client_secret"`
		SubSecret    string `yaml:"sub_secret"`
		ChatHost     string `yaml:"chat_host"`
		BotOauth     string `yaml:"bot_oauth"`
		BotUser      string `yaml:"bot_user"`
	} `yaml:"twitch"`
	Google struct {
		APIKey string `yaml:"api_key"`
	} `yaml:"google"`
	Database struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
		User string `yaml:"user"`
		Pass string `yaml:"pass"`
		DB   string `yaml:"dbname"`
	} `yaml:"database"`
	Redis struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
	} `yaml:"redis"`
	Elasticsearch struct {
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Host     string `yaml:"host"`
		Index    string `yaml:"index"`
	} `yaml:"elasticsearch"`
}

var conf *Config

// LoadConfig -
func LoadConfig(p ...string) error {
	var fp string
	if len(p) > 0 && len(p[0]) > 0 {
		fp = p[0]
	} else {
		wd, err := os.Getwd()
		if err != nil {
			return err
		}
		fp = path.Join(wd, "config.yml")
	}
	fp = utils.ParsePath(fp)

	exists := utils.CheckExists(fp, false)
	if !exists {
		return errors.New("config file not exists")
	}

	data, err := ioutil.ReadFile(fp)
	if err != nil {
		return err
	}

	conf = &Config{}
	err = yaml.Unmarshal(data, conf)
	if err != nil {
		return err
	}

	return nil
}

// GetConf -
func GetConf() *Config {
	return conf
}
