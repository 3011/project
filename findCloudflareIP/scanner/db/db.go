package db

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type EndPoint struct {
	gorm.Model
	IP                string `gorm:"unique"`
	Colo              string
	City              string
	Region            string
	Continent         string
	AddTime           int64
	LastDetectionTime int64
	Active            bool
}

var db *gorm.DB

func init() {
	exePath, _ := os.Executable()
	dirPath := filepath.Dir(exePath)
	configPath := filepath.Join(dirPath, "data.db")
	var err error
	db, err = gorm.Open(sqlite.Open(configPath), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&EndPoint{})
}

func CreateEndPoint(endPoint *EndPoint) {
	now := time.Now().Unix()
	var foundEndPoint EndPoint
	result := db.Where("ip = ?", endPoint.IP).Find(&foundEndPoint)
	if result.RowsAffected == 0 {
		endPoint.AddTime = now
		endPoint.LastDetectionTime = now
		endPoint.Active = true
		db.Create(endPoint)
	} else {
		endPoint.AddTime = foundEndPoint.AddTime
		endPoint.LastDetectionTime = now
		endPoint.Active = true
		fmt.Printf("\"else\": %v\n", "else")
		db.Model(&foundEndPoint).Updates(endPoint)
	}

}
