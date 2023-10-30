package main

import (
	"commander_server/db"
	"commander_server/httphandler"
	"commander_server/models"
	"encoding/json"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AddClients() {
	clients := []models.Client{}
	path := `C:\Users\James\Downloads\MOCK_DATA.json`
	data, _ := os.ReadFile(path)
	json.Unmarshal(data, &clients)

	err := db.ExecuteOnDB(func(db *gorm.DB) error {
		return db.Save(&clients).Error
	})
	if err != nil {
		fmt.Println(err)
	}

}
func main() {
	db.MigrateModels()
	// AddClients()

	router := gin.Default()
	router.GET("/ping", httphandler.HandlePing)

	router.GET("/client", httphandler.HandleClientList)
	router.GET("/client/:id", httphandler.HandleClientFind)
	router.POST("/client", httphandler.HandleClientPost)
	router.PUT("/client/:id", httphandler.HandleClientUpdate)
	router.DELETE("/client/:id", httphandler.HandleClientDelete)

	// Handle WebSocket connections
	router.GET("/ws", httphandler.HandleWebsocket)
	portString, f := os.LookupEnv("PORT")
	if !f {
		portString = "8070"
	}
	err := router.Run(":" + portString)
	if err != nil {

	}

}
