package httphandler

import (
	"commander_server/db"
	"commander_server/models"
	"errors"
	"net/http"
	"net/mail"
	"strings"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserForm struct {
	Email    string `form:"email" json:"email"`
	Password string `form:"password" json:"password"`
	Name     string `form:"name" json:"name"`
}

//	type ClientForm struct {
//		ClientId string `form:"client_id" json:"client_id" binding:"required"`
//	}
type MachineForm struct {
	Username  string `form:"username" json:"username" binding:"required"`
	Machine   string `form:"machine" json:"machine" binding:"required"`
	MachineId string `form:"machine_id" json:"machine_id" binding:"required"`
}

func HandleRegister(c *gin.Context) {
	isClient := strings.ToLower(c.GetHeader("X-APP-TYPE")) == "client"
	if isClient {
		HandleClientRegister(c)
	} else {
		HandleUserRegister(c)
	}
}
func HandleUserRegister(c *gin.Context) {

	// userForm := models.User{}
	userForm := UserForm{}
	var err error

	err = c.ShouldBindJSON(&userForm)
	if err != nil {
		c.JSON(http.StatusBadRequest, GetErrorResponseWithMessage(err.Error()))
		return
	}
	userForm.Name = strings.TrimSpace(userForm.Name)
	userForm.Email = strings.TrimSpace(userForm.Email)
	userForm.Password = strings.TrimSpace(userForm.Password)

	_, err = mail.ParseAddress(userForm.Email)
	if len(userForm.Name) > 0 && len(userForm.Password) > 7 && err != nil {
		var hashBytes, err = bcrypt.GenerateFromPassword([]byte(userForm.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusBadRequest, GetErrorResponseWithMessage(err.Error()))
			return
		}
		userForm.Password = string(hashBytes)
		user := models.User{
			Name:     userForm.Name,
			Email:    userForm.Email,
			Password: userForm.Password,
		}
		db.ExecuteOnDB(func(db *gorm.DB) error {
			err = db.Model(&models.User{}).Create(&user).Error
			return err
		})
		if err == nil {
			c.JSON(http.StatusCreated, GetSuccessResponseWithData("User created successfully\n.Awaiting verification"))
			return
		} else {
			c.JSON(http.StatusBadRequest, GetErrorResponseWithMessage(err.Error()))
		}
	} else {
		c.JSON(http.StatusBadRequest, GetErrorResponseWithMessage("invalid form"))
	}
}
func HandleClientRegister(c *gin.Context) {
	var err error

	mDetails := MachineForm{}
	err = c.ShouldBindBodyWith(&mDetails, binding.JSON)

	if err == nil {
		client := models.Client{
			Username:  mDetails.Username,
			Machine:   mDetails.Machine,
			MachineId: mDetails.MachineId,
		}
		db.ExecuteOnDB(func(db *gorm.DB) error {
			err = db.Model(&client).Create(&client).Error
			return err
		})
		if err == nil {
			c.JSON(http.StatusCreated, GetSuccessResponseWithData("Machine created successfully."))
			return
		} else {
			c.JSON(http.StatusBadRequest, GetErrorResponseWithMessage(err.Error()))
		}
	} else {
		c.JSON(http.StatusBadRequest, GetErrorResponseWithMessage(err.Error()))
		return
		// c.JSON(http.StatusBadRequest, GetErrorResponseWithMessage("invalid form"))
	}

}
func AuthenticatorHandler(c *gin.Context) (interface{}, error) {
	var err error
	payload := Identifier{}
	isClient := strings.ToLower(c.GetHeader("X-APP-TYPE")) == "client"
	if !isClient {
		logDetails := UserForm{}
		err = c.ShouldBindBodyWith(&logDetails, binding.JSON)
		if err != nil {
			return nil, err
		}
		user := models.User{}
		db.ExecuteOnDB(func(db *gorm.DB) error {
			err = db.Model(&user).Where("email = ?", logDetails.Email).First(&user).Error
			return err
		})
		if err != nil {
			err = errors.New("user identity not found")
		} else {
			err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(logDetails.Password))
			if err == nil {
				payload.Identifier = user.ID
				payload.Name = user.Name
				payload.Type = 1

			} else {
				err = errors.New("User and password do not match")
			}
		}

	} else {

		mDetails := MachineForm{}
		client := models.Client{}
		clientId := strings.TrimSpace(c.GetHeader("X-CLIENT-ID"))
		if len(clientId) == 0 {
			err = c.ShouldBindBodyWith(&mDetails, binding.JSON)

		}
		if err == nil {
			db.ExecuteOnDB(func(db *gorm.DB) error {
				if len(clientId) > 0 {
					err = db.Model(&client).Where("id = ?", clientId).First(&client).Error
				} else {
					err = db.Model(&client).Where("machine_id = ? AND username = ?", mDetails.MachineId, mDetails.Username).First(&client).Error
				}
				return err
			})
			if err != nil {
				err = errors.New("user identity not found")
			} else {
				payload.Identifier = clientId
				payload.Name = client.Username
				payload.Type = 2

			}
		}
	}

	return &payload, err
}

// func HandleLogin(c *gin.Context) {

// 	authHeader := c.Request.Header.Get("Authorization")
// 	if len(authHeader) > 0 {
// 		if strings.Contains(authHeader, "Basic") {
// 			pwds := strings.Split(authHeader, " ")
// 			if len(pwds) != 2 {
// 				c.JSON(http.StatusBadRequest, GetErrorResponseWithMessage("Invalid authorization value"))
// 				return
// 			}
// 			authFields, err := base64.RawURLEncoding.DecodeString(pwds[0])
// 			if err != nil {
// 				c.JSON(http.StatusBadRequest, GetErrorResponseWithMessage("Invalid authorization value"))
// 				return
// 			}
// 			userAndPass := strings.Split(string(authFields), ":")
// 			if len(pwds) != 2 {
// 				c.JSON(http.StatusBadRequest, GetErrorResponseWithMessage("Invalid authorization value"))
// 				return
// 			}
// 			user := models.User{}
// 			db.ExecuteOnDB(func(db *gorm.DB) error {
// 				err = db.Model(&user).Where("email = ?", userAndPass[0]).First(&user).Error
// 				return err
// 			})
// 			if err != nil {
// 				c.JSON(http.StatusBadRequest, GetErrorResponseWithMessage("user identity not found"))
// 				return
// 			}
// 			err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userAndPass[1]))
// 			if err == nil {
// 				var payload = map[string]interface{}{
// 					"suceess": "Authenticated successfully",
// 				}
// 				c.JSON(http.StatusOK, GetSuccessResponseWithData(payload))

// 			} else {
// 				c.JSON(http.StatusBadRequest, GetErrorResponseWithMessage("User and password do not match"))

// 			}

// 		} else {
// 			c.JSON(http.StatusBadRequest, GetErrorResponseWithMessage("Unsuported authorization"))
// 		}
// 	}
// }

var identityKey = "id"

// User demo
type Identifier struct {
	Name       string `json:"name"`
	Identifier string `json:"identifier"`
	Type       int    `json:"type"`
}

func GetAuthMiddleware() (*jwt.GinJWTMiddleware, error) {
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "Commander",
		Key:         []byte("%KS^%HY~n'_zEG_KlBH;ma1)FASzhxr("),
		Timeout:     time.Hour * 24 * 7,
		MaxRefresh:  time.Hour * 24,
		IdentityKey: identityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*Identifier); ok {
				return jwt.MapClaims{
					"identifier": v.Identifier,
					"name":       v.Name,
					"type":       v.Type,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &Identifier{
				Identifier: claims["identifier"].(string),
				Name:       claims["name"].(string),
				Type:       int(claims["type"].(float64)),
			}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			return AuthenticatorHandler(c)
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			// identifier, ok := data.(*Identifier)
			// if ok && identifier.Type == 1 {
			// 	return true
			// }

			// return false
			return true
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, GetErrorResponseWithMessage(message))
		},
		// TokenLookup is a string in the form of "<source>:<name>" that is used
		// to extract token from the request.
		// Optional. Default value "header:Authorization".
		// Possible values:
		// - "header:<name>"
		// - "query:<name>"
		// - "cookie:<name>"
		// - "param:<name>"
		TokenLookup: "header: Authorization, query: token, cookie: jwt",
		// TokenLookup: "query:token",
		// TokenLookup: "cookie:token",

		// TokenHeadName is a string in the header. Default value is "Bearer"
		TokenHeadName: "Bearer",

		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,
	})
	return authMiddleware, err
}

// func main() {
// 	port := os.Getenv("PORT")
// 	r := gin.Default()

// 	if port == "" {
// 		port = "8000"
// 	}

// 	// the jwt middleware
// 	authMiddleware, err := GetAuthMiddleware()

// 	if err != nil {
// 		log.Fatal("JWT Error:" + err.Error())
// 	}

// 	// When you use jwt.New(), the function is already automatically called for checking,
// 	// which means you don't need to call it again.
// 	errInit := authMiddleware.MiddlewareInit()

// 	if errInit != nil {
// 		log.Fatal("authMiddleware.MiddlewareInit() Error:" + errInit.Error())
// 	}

// 	r.POST("/login", authMiddleware.LoginHandler)

// 	r.NoRoute(authMiddleware.MiddlewareFunc(), func(c *gin.Context) {
// 		claims := jwt.ExtractClaims(c)
// 		log.Printf("NoRoute claims: %#v\n", claims)
// 		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
// 	})

// 	auth := r.Group("/auth")
// 	// Refresh time can be longer than token timeout
// 	auth.GET("/refresh_token", authMiddleware.RefreshHandler)
// 	auth.Use(authMiddleware.MiddlewareFunc())
// 	{
// 		auth.GET("/hello", helloHandler)
// 	}

// 	if err := http.ListenAndServe(":"+port, r); err != nil {
// 		log.Fatal(err)
// 	}
// }
