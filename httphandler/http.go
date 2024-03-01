package httphandler

import (
	"bytes"
	"commander_server/db"
	"commander_server/models"
	"commander_server/utils"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var MaxPageSize = 200

type BodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (r BodyWriter) Write(b []byte) (int, error) {
	return r.body.Write(b)
}

type ModelType interface {
	models.Client | models.ClientProfile | models.Download | models.ClientGroup | models.User | models.GroupToClient | models.MachineInfo | models.Runtime | models.Task
}

func HandlePing(c *gin.Context) {
	msg := models.JSONResponse{
		Success: true,
	}
	c.JSON(200, msg)
}
func GzipMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		canEncode := false

		canEncode = strings.Contains(strings.ToLower(c.Request.Header.Get("Accept-Encoding")), "gzip")

		var wb *BodyWriter
		if canEncode {
			// writerOriginal = &c.Writer

			wb = &BodyWriter{
				body:           &bytes.Buffer{},
				ResponseWriter: c.Writer,
			}
			c.Writer = wb
			c.Header("Content-Encoding", "gzip")
		}
		c.Next()

		if canEncode {

			var compressed []byte
			var err error
			{
				data := wb.body.Bytes()
				compressed, err = utils.CompressGzip(data)
			}
			if err == nil {
				wb.body = &bytes.Buffer{}
				// wb.Write(compressed)
				_, err = wb.ResponseWriter.Write(compressed)
			}
			if err != nil {
				c.JSON(http.StatusBadRequest, GetErrorResponseWithMessage(err.Error()))
			}

		}

	}
}
func GetErrorResponseWithMessage(message string) models.JSONResponse {
	if len(message) == 0 {
		message = "Unknown error ocurred"
	}
	return models.JSONResponse{
		Success: false,
		Error:   message,
	}
}
func GetSuccessResponseWithData(data interface{}) models.JSONResponse {
	return models.JSONResponse{
		Success: true,
		Data:    data,
	}
}
func HandlePost[T ModelType](c *gin.Context, instance T) {
	err := c.BindJSON(&instance)
	if err == nil {
		db.ExecuteOnDB(func(db *gorm.DB) error {
			err = db.Create(&instance).Error
			return err
		})
	}
	if err == nil {
		c.JSON(http.StatusOK, GetSuccessResponseWithData(instance))
	} else {
		c.AbortWithStatusJSON(http.StatusBadRequest, GetErrorResponseWithMessage(err.Error()))
	}
}
func HandleUpdate[T ModelType](c *gin.Context, result T) {
	id := strings.Trim(c.Param("id"), "/")
	if len(id) == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, GetErrorResponseWithMessage("resource not found"))
	}
	body := map[string]interface{}{}
	err := c.BindJSON(&body)
	if err == nil {
		db.ExecuteOnDB(func(db *gorm.DB) error {
			err = db.Model(&result).Where("id = ?", id).Find(&result).Error
			if err != nil {
				return err
			}
			return db.Model(&result).Select("*").Updates(body).Error
		})
	}
	if err == nil {
		c.JSON(http.StatusOK, GetSuccessResponseWithData(body))
	} else {
		c.AbortWithStatusJSON(http.StatusBadRequest, GetErrorResponseWithMessage(err.Error()))
	}
}
func HandleDelete[T ModelType](c *gin.Context, instance T) {
	id := strings.Trim(c.Param("id"), "/")
	if len(id) == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, GetErrorResponseWithMessage("resource not found"))
		return
	}
	var err error
	db.ExecuteOnDB(func(db *gorm.DB) error {
		err = db.Unscoped().Where("id = ?", id).Delete(&instance).Error
		return err
	})

	if err == nil {
		c.JSON(http.StatusOK, GetSuccessResponseWithData(nil))
	} else {
		c.AbortWithStatusJSON(http.StatusBadRequest, GetErrorResponseWithMessage(err.Error()))
	}
}
func HandleFind[T ModelType](c *gin.Context, result T) {
	id := strings.Trim(c.Param("id"), "/")
	if len(id) == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, GetErrorResponseWithMessage("resource not found"))
	}
	var err error
	db.ExecuteOnDB(func(db *gorm.DB) error {
		err = db.Model(&result).Where("id = ?", id).First(&result).Error
		return err
	})
	if err == nil {
		c.JSON(http.StatusOK, GetSuccessResponseWithData(result))
	} else {
		c.AbortWithStatusJSON(http.StatusBadRequest, GetErrorResponseWithMessage(err.Error()))
	}
}
func HandleList[T ModelType](c *gin.Context, instance T, searchableFields []string) {

	var results = []T{}
	limit := MaxPageSize
	page := 1
	total := 0
	search := ""
	orderKey := "created_at"
	orderType := "desc"
	fromTime := int64(0)
	ids := []string{}
	var err error

	limitQ := strings.ToLower(c.Query("limit"))
	if len(limitQ) > 0 {
		if num, err := strconv.ParseInt(limitQ, 10, 64); err == nil {
			if num > 0 && num < int64(MaxPageSize) {
				limit = int(num)
			}
		}
	}
	pageQ := strings.ToLower(c.Query("page"))
	if len(pageQ) > 0 {
		if num, err := strconv.ParseInt(pageQ, 10, 64); err == nil {
			if num > 0 {
				page = int(num)
			}
		}
	}
	timeQ := strings.ToLower(c.Query("fromTime"))
	if len(timeQ) > 0 {
		if num, err := strconv.ParseInt(timeQ, 10, 64); err == nil {
			if num > 0 {
				fromTime = num
			}
		}
	}

	orderQ := strings.ToLower(c.Query("order"))
	if len(orderQ) > 0 {
		arr := strings.Split(orderQ, ":")
		if len(arr) == 2 {
			if utils.InSlice(ClientFields, arr[0], false) {

				orderKey = arr[0]
				if strings.Contains(strings.ToLower(arr[1]), "asc") {
					orderType = "asc"
				}
			}
		}
	}

	idsQ := strings.ToLower(c.Query("ids"))
	if len(idsQ) > 0 {
		arr := strings.Split(idsQ, ",")
		arr = utils.RemoveEmptyFromSlice(arr)
		if len(arr) > 0 {
			ids = arr
		}
	}
	if len(ids) > 0 {
		err = db.ExecuteOnDB(func(db *gorm.DB) error {
			dbCtx := db.Model(&instance)
			dbCtx = dbCtx.Where("id in ?", ids)
			err := dbCtx.Find(&results).Error
			return err
		})

	} else {
		searchQ := strings.ToLower(c.Query("search"))
		if len(searchQ) > 0 {
			search = searchQ
		}

		db.ExecuteOnDB(func(db *gorm.DB) error {
			offset := (page - 1) * limit
			dbCtx := db.Model(&instance)
			if fromTime > 0 {
				dbCtx = dbCtx.Where("updated_at = ?", fromTime)
			}

			if len(search) > 0 && len(searchableFields) > 0 {
				var queryStr = ""
				find := "%" + search + "%"
				repeatInterface := make([]interface{}, 0, len(searchableFields))
				for _, field := range searchableFields {
					q := fmt.Sprintf("lower(%s) LIKE ? ", field)
					if len(queryStr) > 0 {
						queryStr += "OR " + q
					} else {
						queryStr += q
					}
					repeatInterface = append(repeatInterface, find)
				}

				dbCtx = dbCtx.Where(queryStr, repeatInterface...)
			}
			count := int64(0)
			err = dbCtx.Count(&count).Error
			if err != nil {
				return err
			}
			total = int(count)
			if offset > 0 {
				dbCtx.Offset(offset)
			}
			dbCtx.Limit(limit).Order(orderKey + " " + orderType)
			dbCtx.Find(&results)
			err = dbCtx.Error
			return err
		})
	}
	if err == nil {
		totalPages := 1
		if total > 0 {
			totalPages = int(math.Ceil(float64(total) / float64(limit)))
		}
		var pageres = models.PageResult{
			Total:      total,
			TotalPages: totalPages,
			Page:       page,
			PerPage:    limit,
			Data:       results,
		}
		msg := models.JSONResponse{
			Success: true,
			Data:    pageres,
		}
		c.JSON(200, msg)
	} else {
		c.AbortWithStatusJSON(http.StatusBadRequest, GetErrorResponseWithMessage(err.Error()))
	}
}
