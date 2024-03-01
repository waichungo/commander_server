package httphandler

import (
	"commander_server/db"
	"commander_server/models"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DashboardData struct {
	UsersCount         int             `json:"usersCount"`
	MachinesCount      int             `json:"machinesCount"`
	LastClientActivity int64           `json:"lastClientActivity"`
	RecentClients      []models.Client `json:"recentClients"`
}

func getCount[T ModelType](instance T) int {
	res := 0
	db.ExecuteOnDB(func(db *gorm.DB) error {
		count := int64(0)
		err := db.Model(&instance).Count(&count).Error
		res = int(count)
		return err
	})
	return res
}

// func getRecentModels[T ModelType](instance T, limit int) []T {
// 	recents := []T{}
// 	if limit < 1 {
// 		limit = 10
// 	}

//		db.ExecuteOnDB(func(db *gorm.DB) error {
//			err := db.Model(&instance).Order("updated_at desc").Limit(limit).Find(&recents).Error
//			return err
//		})
//		return recents
//	}
func getRecentClients(limit int) []models.Client {
	recents := []models.Client{}
	if limit < 1 {
		limit = 10
	}

	db.ExecuteOnDB(func(db *gorm.DB) error {
		err := db.Model(&models.Client{}).Order("updated_at desc").Limit(limit).Find(&recents).Error
		return err
	})
	return recents
}

func HandleDashboard(c *gin.Context) {
	usersCount := 0
	machinesCount := 0
	lastClientActivity := int64(0)
	clients := []models.Client{}
	db.ExecuteOnDB(func(db *gorm.DB) error {
		var wg = sync.WaitGroup{}
		wg.Add(4)
		go func() {
			defer wg.Done()
			usersCount = getCount(models.User{})
		}()
		go func() {
			defer wg.Done()
			machinesCount = getCount(models.Client{})
		}()
		go func() {
			defer wg.Done()
			db.Model(&models.MachineInfo{}).Select("updated_at").Order("updated_at desc").First(&lastClientActivity)

		}()
		go func() {
			defer wg.Done()
			clients = getRecentClients(10)

		}()
		wg.Wait()
		return nil
	})
	data := DashboardData{
		UsersCount:         usersCount,
		MachinesCount:      machinesCount,
		LastClientActivity: lastClientActivity,
		RecentClients:      clients,
	}
	c.JSON(http.StatusOK, GetSuccessResponseWithData(data))

}
