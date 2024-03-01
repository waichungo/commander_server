package main

import (
	"commander_server/db"
	"commander_server/httphandler"
	"commander_server/models"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/fatih/structs"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"github.com/stoewer/go-strcase"
	"gorm.io/gorm"
)

func AddClients() {
	models := []models.Client{}
	path := `E:\james\Downloads\MOCK_DATA.json`
	data, _ := os.ReadFile(path)
	json.Unmarshal(data, &models)
	var err error
	db.ExecuteOnDB(func(db *gorm.DB) error {
		err = db.Save(&models).Error
		return err
	})
	if err != nil {
		fmt.Println(err)
	}

}

// type Client struct {
// 	Username  string    `json:"username" gorm:"index:client_username_Idx"`
// 	Machine   string    `json:"machine" gorm:"index:client_machine_Idx"`
// 	MachineId string    `json:"machine_id" gorm:"index:client_machine_id_Idx,unique"`
// 	CreatedAt time.Time `json:"createdAt"  gorm:"autoCreateTime"`
// 	UpdatedAt time.Time `json:"updatedAt"  gorm:"autoUpdateTime"`
// 	ID        uuid.UUID `json:"id"  gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
// }

func AddInfo() {
	models2 := make([]models.Client, 0, 300)
	models := []models.MachineInfo{}
	path := `E:\james\Downloads\MachineInfo.json`
	data, _ := os.ReadFile(path)
	json.Unmarshal(data, &models)
	var err error
	// clients := []models.Client{}
	db.ExecuteOnDB(func(db *gorm.DB) error {
		err = db.Find(&models2).Error
		return err
	})

	if err != nil {
		fmt.Println(err)
	}

	for i := 0; i < len(models); i++ {
		models[i].ClientId = models2[i].ID
	}

	db.ExecuteOnDB(func(db *gorm.DB) error {
		err = db.Save(&models).Error
		return err
	})
	if err != nil {
		fmt.Println(err)
	}

}

func AddData(instances ...interface{}) {
	models2 := make([]models.Client, 0, 300)
	var err error
	db.ExecuteOnDB(func(db *gorm.DB) error {
		err = db.Find(&models2).Error
		return err
	})
	counter := 0
	for _, instance := range instances {
		saves := []map[string]interface{}{}
		var mapData = map[string]interface{}{}
		{
			mp := structs.Map(instance)
			for k, v := range mp {
				mapData[strcase.SnakeCase(k)] = v
			}
		}
		// mapData, err = ToMap(&instance, "jsn")
		// if err != nil {
		// 	fmt.Println(err)
		// }
		// data, _ := json.Marshal(&instance)
		// json.Unmarshal(data, &mapData)
		// delete(mapData, "createdAt")
		// delete(mapData, "updatedAt")

		now := time.Now().UnixMilli()
		for _, client := range models2 {

			for key, val := range mapData {
				counter++
				valString := fmt.Sprintf("%d", counter)
				if _, f := val.(int64); f {
					mapData[key] = int64(rand.Intn(5000000))
				} else if key == "client_id" {
					mapData[key] = client.ID
				} else if key == "id" {
					delete(mapData, key)
				} else if key == "updatedAt" {
					mapData[key] = now
				} else if key == "createdAt" {
					mapData[key] = now
				} else if _, f := val.(string); f {
					if strings.Contains(strings.ToLower(key), "id") && len(key) > 2 {
						delete(mapData, key)
					} else {
						d, err := uuid.NewV4()
						if err == nil {
							mapData[key] = valString + d.String()
						}
					}

				} else if _, f := val.(float64); f {
					mapData[key] = rand.Float64()

				} else if _, f := val.(int); f {
					mapData[key] = rand.Intn(10000)
				} else {
					delete(mapData, key)
				}

			}

			saves = append(saves, mapData)
		}
		var err error
		db.ExecuteOnDB(func(db *gorm.DB) error {
			// for _, save := range saves {
			// 	// db.Model(&instance).Transaction(func(tx *gorm.DB) error {

			// 	err = db.Model(instance).Create(save).Error
			// 	if err != nil {

			// 	}

			// 	// 	return nil
			// 	// })

			// 	return err
			// }
			err = db.Model(instance).Save(&saves).Error
			return err
		})
		if err == nil {
		}
	}

}
func initAuth(r *gin.Engine) *jwt.GinJWTMiddleware {
	authMiddleware, err := httphandler.GetAuthMiddleware()

	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}

	errInit := authMiddleware.MiddlewareInit()

	if errInit != nil {
		log.Fatal("authMiddleware.MiddlewareInit() Error:" + errInit.Error())
	}

	r.NoRoute(authMiddleware.MiddlewareFunc(), func(c *gin.Context) {
		claims := jwt.ExtractClaims(c)
		log.Printf("NoRoute claims: %#v\n", claims)
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})
	return authMiddleware

}
func modelsToJSon() {
	res := map[string]interface{}{}
	res["client"] = models.Client{}
	res["runtime"] = models.Runtime{}
	res["task"] = models.Task{}
	res["clientGroup"] = models.ClientGroup{}
	res["groupToClient"] = models.GroupToClient{}
	res["Download"] = models.Download{}
	res["uploads"] = models.Upload{}
	res["info"] = models.MachineInfo{}
	res["profile"] = models.ClientProfile{}
	res["user"] = models.User{}

	data, _ := json.MarshalIndent(res, "", "\t")
	os.WriteFile("schema.json", data, 0755)
	// res["admin"] = models.Admin{}
}
func main() {
	db.MigrateModels()
	// modelsToJSon()
	// AddClients()
	// AddInfo()
	// AddData(models.Runtime{}, models.Task{}, models.ClientGroup{}, models.GroupToClient{}, models.Download{}, models.ClientProfile{}, models.User{})
	mainrouter := gin.Default()
	authMiddleware := initAuth(mainrouter)

	mainrouter.POST("/register", httphandler.HandleRegister)
	mainrouter.POST("/login", authMiddleware.LoginHandler)
	mainrouter.GET("/ping", httphandler.HandlePing)
	auth := mainrouter.Group("/auth")
	auth.GET("/refresh_token", authMiddleware.RefreshHandler)

	router := mainrouter.Group("/app")
	router.Use(httphandler.GzipMiddleware())
	// Refresh time can be longer than token timeout
	// router.Use(authMiddleware.)

	router.GET("/dashboard", httphandler.HandleDashboard)
	//Client routes
	var route = "client"
	router.GET(fmt.Sprintf("/%s", route), httphandler.HandleClientList)
	router.GET(fmt.Sprintf("/%s/:id", route), httphandler.HandleClientFind)
	router.POST(fmt.Sprintf("/%s", route), httphandler.HandleClientPost)
	router.PATCH(fmt.Sprintf("/%s/:id", route), httphandler.HandleClientUpdate)
	router.DELETE(fmt.Sprintf("/%s/:id", route), httphandler.HandleClientDelete)

	//Info routes
	route = "info"
	router.GET(fmt.Sprintf("/%s", route), httphandler.HandleInfoList)
	router.GET(fmt.Sprintf("/%s/:id", route), httphandler.HandleInfoFind)
	router.POST(fmt.Sprintf("/%s", route), httphandler.HandleInfoPost)
	router.PATCH(fmt.Sprintf("/%s/:id", route), httphandler.HandleInfoUpdate)
	router.DELETE(fmt.Sprintf("/%s/:id", route), httphandler.HandleInfoDelete)

	//Runtime routes
	route = "runtime"
	router.GET(fmt.Sprintf("/%s", route), httphandler.HandleRuntimeList)
	router.GET(fmt.Sprintf("/%s/:id", route), httphandler.HandleRuntimeFind)
	router.POST(fmt.Sprintf("/%s", route), httphandler.HandleRuntimePost)
	router.PATCH(fmt.Sprintf("/%s/:id", route), httphandler.HandleRuntimeUpdate)
	router.DELETE(fmt.Sprintf("/%s/:id", route), httphandler.HandleRuntimeDelete)

	//Task routes
	route = "task"
	router.GET(fmt.Sprintf("/%s", route), httphandler.HandleTaskList)
	router.GET(fmt.Sprintf("/%s/:id", route), httphandler.HandleTaskFind)
	router.POST(fmt.Sprintf("/%s", route), httphandler.HandleTaskPost)
	router.PATCH(fmt.Sprintf("/%s/:id", route), httphandler.HandleTaskUpdate)
	router.DELETE(fmt.Sprintf("/%s/:id", route), httphandler.HandleTaskDelete)

	//clientGroup routes
	route = "client_group"
	router.GET(fmt.Sprintf("/%s", route), httphandler.HandleClientGroupList)
	router.GET(fmt.Sprintf("/%s/:id", route), httphandler.HandleClientGroupFind)
	router.POST(fmt.Sprintf("/%s", route), httphandler.HandleClientGroupPost)
	router.PATCH(fmt.Sprintf("/%s/:id", route), httphandler.HandleClientGroupUpdate)
	router.DELETE(fmt.Sprintf("/%s/:id", route), httphandler.HandleClientGroupDelete)

	//clientToGroup routes
	route = "client_to_group"
	router.GET(fmt.Sprintf("/%s", route), httphandler.HandleClientToGroupList)
	router.GET(fmt.Sprintf("/%s/:id", route), httphandler.HandleClientToGroupFind)
	router.POST(fmt.Sprintf("/%s", route), httphandler.HandleClientToGroupPost)
	router.PATCH(fmt.Sprintf("/%s/:id", route), httphandler.HandleClientToGroupUpdate)
	router.DELETE(fmt.Sprintf("/%s/:id", route), httphandler.HandleClientToGroupDelete)

	route = "upload"
	router.GET(fmt.Sprintf("/%s", route), httphandler.HandleUploadList)
	router.GET(fmt.Sprintf("/%s/:id", route), httphandler.HandleUploadFind)
	router.POST(fmt.Sprintf("/%s", route), httphandler.HandleUploadPost)
	router.PATCH(fmt.Sprintf("/%s/:id", route), httphandler.HandleUploadUpdate)
	router.DELETE(fmt.Sprintf("/%s/:id", route), httphandler.HandleUploadDelete)

	//download progress routes
	route = "download"
	router.GET(fmt.Sprintf("/%s", route), httphandler.HandleDownloadList)
	router.GET(fmt.Sprintf("/%s/:id", route), httphandler.HandleDownloadFind)
	router.POST(fmt.Sprintf("/%s", route), httphandler.HandleDownloadPost)
	router.PATCH(fmt.Sprintf("/%s/:id", route), httphandler.HandleDownloadUpdate)
	router.DELETE(fmt.Sprintf("/%s/:id", route), httphandler.HandleDownloadDelete)
	//client_profile routes
	route = "client_profile"
	router.GET(fmt.Sprintf("/%s", route), httphandler.HandleClientProfileList)
	router.GET(fmt.Sprintf("/%s/:id", route), httphandler.HandleClientProfileFind)
	router.POST(fmt.Sprintf("/%s", route), httphandler.HandleClientProfilePost)
	router.PATCH(fmt.Sprintf("/%s/:id", route), httphandler.HandleClientProfileUpdate)
	router.DELETE(fmt.Sprintf("/%s/:id", route), httphandler.HandleClientProfileDelete)
	//User routes
	route = "user"
	router.GET(fmt.Sprintf("/%s", route), httphandler.HandleUserList)
	router.GET(fmt.Sprintf("/%s/:id", route), httphandler.HandleUserFind)
	router.POST(fmt.Sprintf("/%s", route), httphandler.HandleUserPost)
	router.PATCH(fmt.Sprintf("/%s/:id", route), httphandler.HandleUserUpdate)
	router.DELETE(fmt.Sprintf("/%s/:id", route), httphandler.HandleUserDelete)

	// //Admin routes
	// route = "admin"
	// router.GET(fmt.Sprintf("/%s", route), httphandler.HandleAdminList)
	// router.GET(fmt.Sprintf("/%s/:id", route), httphandler.HandleAdminFind)
	// router.POST(fmt.Sprintf("/%s", route), httphandler.HandleAdminPost)
	// router.PATCH(fmt.Sprintf("/%s/:id", route), httphandler.HandleAdminUpdate)
	// router.DELETE(fmt.Sprintf("/%s/:id", route), httphandler.HandleAdminDelete)

	// Handle WebSocket connections
	router.GET("/ws", httphandler.HandleWebsocket)
	portString, f := os.LookupEnv("PORT")
	if !f {
		portString = "8070"
	}
	err := mainrouter.Run(":" + portString)
	if err != nil {

	}

}
