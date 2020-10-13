package models

import (
	"fmt"
	"time"
	"os"

	"gorm.io/gorm"
	"gorm.io/driver/postgres"
	"gorm.io/gorm/logger"

	"github.com/joho/godotenv"
)

type DBModel struct {
	ID        int        `json:"id" gorm:"primary_key"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt"`
}

var db *gorm.DB

func init() {
	e := godotenv.Load()
	if e != nil {
		fmt.Print(e)
	}

	// dsn := "user=gorm password=gorm DB.name=gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai"
	dsn := "uri=" + os.Getenv("DATABASE_URL") + "port=5432"
	conn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		fmt.Print(err)
	}

	db = conn
	// db.Debug().AutoMigrate(&ContactAssignmentEntry{})
	db.AutoMigrate(&User{})
	db.AutoMigrate(&Account{})
	// err2 := db.Debug().AutoMigrate(&Contact{}).Error
	// if err2 != nil {
	// 	fmt.Println(err2)
	// }
	db.Debug().AutoMigrate(&ContactAssignmentEntry{})
	db.Debug().AutoMigrate(&Contact{})
	db.AutoMigrate(&Organization{})
	db.AutoMigrate(&Workspace{})
}

func DB() *gorm.DB {
	return db
}
