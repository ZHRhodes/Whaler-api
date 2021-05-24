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
	ID        string     `json:"id" gorm:"type:uuid;unique;primaryKey;default:uuid_generate_v4()"`
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

	// db.Migrator().DropTable(&Account{})
	// db.Debug().Migrator().CreateTable(&Account{})
	// db.Migrator().DropTable(&Contact{})
	// db.Migrator().CreateTable(&Contact{})
	// db.Migrator().DropTable(&ContactAssignmentEntry{})
	// db.Migrator().CreateTable(&ContactAssignmentEntry{})
	// db.Migrator().DropTable(&AccountAssignmentEntry{})
	// db.Migrator().CreateTable(&AccountAssignmentEntry{})
	// db.Migrator().DropTable(&Organization{})
	// db.Migrator().CreateTable(&Organization{})
	// db.Migrator().DropTable(&RefreshToken{})
	// db.Migrator().CreateTable(&RefreshToken{})
	// db.Migrator().DropTable(&User{})
	// db.Migrator().CreateTable(&User{})
	// db.Migrator().DropTable(&Workspace{})
	// db.Migrator().CreateTable(&Workspace{})

	db.AutoMigrate(&Account{})
	db.AutoMigrate(&Contact{})
	db.AutoMigrate(&User{})
	// db.AutoMigrate(&ContactAssignmentEntry{})
	// db.AutoMigrate(&AccountAssignmentEntry{})
	// db.AutoMigrate(&Organization{})
	// db.AutoMigrate(&RefreshToken{})
	db.AutoMigrate(&Note{})
	db.AutoMigrate(&Task{})
	// db.AutoMigrate(&Workspace{})

	// err2 := db.Debug().AutoMigrate(&Contact{}).Error
	// if err2 != nil {
	// 	fmt.Println(err2)
	// }
}

func DB() *gorm.DB {
	return db
}
