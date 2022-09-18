package config

import (
	"log"

	"github.com/jinzhu/configor"
)

var Config = struct {
	PortList []int  `required:"true"`
	DBPath   string `default:"/etc/x-ui/x-ui.db"`
}{}

func init() {
	err := configor.Load(&Config, "config.yml")
	if err != nil {
		log.Panic(err.Error())
	}
}
