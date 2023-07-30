package initserver

import (
	"log"
	"time"

	"github.com/dominiclet/golang-base/init_server/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitGormDB(config *config.Config) *gorm.DB {
	var db *gorm.DB
	var err error
	for {
		db, err = gorm.Open(mysql.Open(config.DB))
		if err == nil {
			break
		}
		log.Printf("Waiting for DB connection: %s (%v)\n", config.DB, err)
		time.Sleep(2 * time.Second)
	}
	return db
}
