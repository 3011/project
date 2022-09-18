package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/jinzhu/configor"
)

var Config = struct {
	PortList []int  `required:"true"`
	DBPath   string `default:"/etc/x-ui/x-ui.db"`
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
