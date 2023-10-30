package httphandler

import (
	"commander_server/db"
	"commander_server/models"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var MaxPageSize = 200

var ClientFields = []string{
	"username",
	"machine",
	"machine_id",
}

func HandleClientPost(c *gin.Context) {
	client := models.Client{}
	err := c.BindJSON(&client)
	if err == nil {
		err = db.ExecuteOnDB(func(db *gorm.DB) error {
			return db.Create(&client).Error
		})
	}
	if err == nil {
		c.JSON(http.StatusOK, GetSuccessResponseWithData(client))
	} else {
		c.AbortWithStatusJSON(http.StatusBadRequest, GetErrorResponseWithMessage(err.Error()))
	}

}
func HandleClientFind(c *gin.Context) {
	id := strings.Trim(c.Param("id"), "/")
	if len(id) == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, GetErrorResponseWithMessage("resource not found"))
	}
	client := models.Client{}

	err := db.ExecuteOnDB(func(db *gorm.DB) error {
		return db.Where("id = ?", id).Find(&client).Error
	})
	if err == nil {
		c.JSON(http.StatusOK, GetSuccessResponseWithData(client))
	} else {
		c.AbortWithStatusJSON(http.StatusBadRequest, GetErrorResponseWithMessage(err.Error()))
	}

}
func HandleClientUpdate(c *gin.Context) {
	id := strings.Trim(c.Param("id"), "/")
	if len(id) == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, GetErrorResponseWithMessage("resource not found"))
	}
	client := models.Client{}
	err := c.BindJSON(&client)
	if err == nil {
		err = db.ExecuteOnDB(func(db *gorm.DB) error {
			return db.Where("id = ?", id).Save(&client).Error
		})
	}
	if err == nil {
		c.JSON(http.StatusOK, GetSuccessResponseWithData(client))
	} else {
		c.AbortWithStatusJSON(http.StatusBadRequest, GetErrorResponseWithMessage(err.Error()))
	}

}

func HandleClientDelete(c *gin.Context) {
	id := strings.Trim(c.Param("id"), "/")
	if len(id) == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, GetErrorResponseWithMessage("resource not found"))
		return
	}
	err := db.ExecuteOnDB(func(db *gorm.DB) error {
		return db.Unscoped().Where("id = ?", id).Delete(models.Client{}).Error
	})

	if err == nil {
		c.JSON(http.StatusOK, GetSuccessResponseWithData(nil))
	} else {
		c.AbortWithStatusJSON(http.StatusBadRequest, GetErrorResponseWithMessage(err.Error()))
	}
}
func HandleClientList(c *gin.Context) {
	HandleList[models.Client](c, ClientFields)
}

// func HandleClientGet(c *gin.Context) {
// 	var fields = ClientFields
// 	limit := MaxPageSize
// 	page := 1
// 	total := 0
// 	search := ""
// 	orderKey := "createdAt"
// 	orderType := "desc"

// 	limitQ := strings.ToLower(c.Query("limit"))
// 	if len(limitQ) > 0 {
// 		if num, err := strconv.ParseInt(limitQ, 10, 64); err == nil {
// 			if num > 0 && num < int64(MaxPageSize) {
// 				limit = int(num)
// 			}
// 		}
// 	}
// 	pageQ := strings.ToLower(c.Query("page"))
// 	if len(pageQ) > 0 {
// 		if num, err := strconv.ParseInt(pageQ, 10, 64); err == nil {
// 			if num > 0 {
// 				page = int(num)
// 			}
// 		}
// 	}

// 	orderQ := strings.ToLower(c.Query("order"))
// 	if len(orderQ) > 0 {
// 		arr := strings.Split(orderQ, ":")
// 		if len(arr) == 2 {
// 			if utils.InSlice(ClientFields, arr[0], false) {

// 				orderKey = arr[0]
// 				if strings.Contains(strings.ToLower(arr[1]), "asc") {
// 					orderType = "asc"
// 				}
// 			}
// 		}
// 	}
// 	var clients = []models.Client{}
// 	err := db.ExecuteOnDB(func(db *gorm.DB) error {
// 		offset := (page - 1) * limit
// 		dbCtx := db.Model(&models.Client{})

// 		if len(search) > 0 {
// 			var queryStr = ""
// 			find := "%" + search + "%"
// 			repeatInterface := make([]interface{}, len(fields))
// 			for _, field := range fields {
// 				q := fmt.Sprintf("lower(%s) LIKE ? ", field)
// 				if len(queryStr) > 0 {
// 					queryStr += "OR " + q
// 				} else {
// 					queryStr += q
// 				}
// 				repeatInterface = append(repeatInterface, find)
// 			}

// 			dbCtx.Where(queryStr, repeatInterface...)
// 		}
// 		count := int64(0)
// 		var err = dbCtx.Count(&count).Error
// 		if err != nil {
// 			return err
// 		}
// 		total = int(count)
// 		if offset > 0 {
// 			dbCtx.Offset(offset)
// 		}
// 		dbCtx.Limit(limit).Order(orderKey + " " + orderType)
// 		dbCtx.Find(&clients)
// 		return dbCtx.Error
// 	})
// 	if err == nil {
// 		totalPages := 1
// 		if total > 0 {
// 			totalPages = int(math.Ceil(float64(total) / float64(limit)))
// 		}
// 		var pageres = models.PageResult{
// 			Total:      total,
// 			TotalPages: totalPages,
// 			Page:       page,
// 			PerPage:    limit,
// 			Data:       clients,
// 		}
// 		msg := models.JSONResponse{
// 			Success: true,
// 			Data:    pageres,
// 		}
// 		c.JSON(200, msg)
// 	} else {
// 		c.AbortWithStatusJSON(http.StatusBadRequest, GetErrorResponseWithMessage(err.Error()))
// 	}
// }
