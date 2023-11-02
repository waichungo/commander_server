package db

import (
	"commander_server/models"
	"context"
	"errors"
	"fmt"

	"github.com/cockroachdb/cockroach-go/v2/crdb/crdbgorm"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func GetConn() (*gorm.DB, error) {
	// db, err := gorm.Open(postgres.Open(os.Getenv("DATABASE_URL")), &gorm.Config{})
	const conString = "postgres://postgres:root@localhost:5432/collector"
	db, err := gorm.Open(postgres.Open(conString), &gorm.Config{})
	return db, err
}

var NilError = errors.New("Cannot execute on a nil reference")

func ExecuteOnDB(execute func(db *gorm.DB) error) error {
	if execute == nil {
		return NilError
	}
	var db, err = GetConn()
	if err == nil {
		var internalDB, err = db.DB()
		if err != nil {
			return err
		}
		defer internalDB.Close()
		err = crdbgorm.ExecuteTx(context.Background(), db, nil,
			execute,
		)
	}
	return err
}
func MigrateModels() error {
	var err error
	ExecuteOnDB(func(db *gorm.DB) error {
		err = db.AutoMigrate(&models.Client{})
		if err == nil {
			err = db.AutoMigrate(&models.Runtime{})
		}
		if err == nil {
			err = db.AutoMigrate(&models.Task{})
		}
		if err == nil {
			err = db.AutoMigrate(&models.Client{})
		}
		if err == nil {
			err = db.AutoMigrate(&models.GroupToClient{})
		}
		if err == nil {
			err = db.AutoMigrate(&models.DownloadProgress{})
		}
		if err == nil {
			err = db.AutoMigrate(&models.MachineInfo{})
		}
		if err == nil {
			err = db.AutoMigrate(&models.ClientProfile{})
		}
		if err == nil {
			err = db.AutoMigrate(&models.User{})
		}
		// if err == nil {
		// 	err = db.AutoMigrate(&models.Admin{})
		// }
		if err != nil {
			fmt.Println(err)
		}
		return err
	})
	return err
}
