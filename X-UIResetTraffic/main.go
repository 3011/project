package main

import (
	"log"

	"github.com/3011/project/X-UIResetTraffic/config"
	"github.com/3011/project/X-UIResetTraffic/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func init() {
	var err error
	db, err = gorm.Open(sqlite.Open(config.Config.DBPath), &gorm.Config{})
	if err != nil {
		log.Panic(err)
	}
}

func main() {
	db.AutoMigrate(&model.Inbound{})
	for _, port := range config.Config.PortList {
		var inbound model.Inbound
		db.First(&inbound, "Port = ?", port)
		db.Model(&inbound).Updates(map[string]interface{}{"Up": 0, "Down": 0, "Enable": true})
	}
	log.Println("Done!")
}
