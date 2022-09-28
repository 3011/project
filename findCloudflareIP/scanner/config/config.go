package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/jinzhu/configor"
)

var Config = struct {
	BotToken string `required:"true"`
	ChatID   string `required:"true"`
}{}

func init() {
	exePath, _ := os.Executable()
	dirPath := filepath.Dir(exePath)
	configPath := filepath.Join(dirPath, "config.yml")
	err := configor.Load(&Config, configPath)
	if err != nil {
		log.Panic(err.Error())
	}
}
