package models

import (
	"fmt"
	"os"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

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

	conn, err := gorm.Open("postgres", os.Getenv("DATABASE_URL"))
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
	db.Debug().DropTable(&ContactAssignmentEntry{})
	db.Debug().DropTable(&ContactAssignmentEntry{})
	db.AutoMigrate(&Organization{})
	db.AutoMigrate(&Workspace{})
}

func DB() *gorm.DB {
	return db
}
