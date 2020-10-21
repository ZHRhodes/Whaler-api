package models

import (
	"fmt"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/joho/godotenv"
)

type DBModel struct {
	ID        string     `json:"id" gorm:"column:id;" sql:"type:uuid;primary_key;default:uuid_generate_v4()"`
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

	conn, err := gorm.Open(postgres.New(postgres.Config{
		DSN: os.Getenv("DATABASE_URL"), // data source name, refer https://github.com/jackc/pgx
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		fmt.Print(err)
	}

	db = conn
	// db.Debug().AutoMigrate(&ContactAssignmentEntry{})
	db.AutoMigrate(&User{})

	db.Migrator().DropTable(&Account{})
	db.Migrator().CreateTable(&Account{})
	db.Migrator().DropTable(&Contact{})
	db.Migrator().CreateTable(&Contact{})
	// db.AutoMigrate(&Account{})
	// db.AutoMigrate(&Contact{})
	// err2 := db.Debug().AutoMigrate(&Contact{}).Error
	// if err2 != nil {
	// 	fmt.Println(err2)
	// }
	db.Debug().AutoMigrate(&ContactAssignmentEntry{})
	db.AutoMigrate(&Organization{})
	db.AutoMigrate(&Workspace{})
}

func DB() *gorm.DB {
	return db
}
