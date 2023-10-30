package httphandler

import (
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

type ModelType interface {
	models.Client | models.ClientProfile | models.DownloadProgress | models.Group | models.GroupToClient | models.MachineInfo | models.Runtime | models.Task
}

func HandlePing(c *gin.Context) {
	msg := models.JSONResponse{
		Success: true,
	}
	c.JSON(200, msg)
}
func GetErrorResponseWithMessage(message string) models.JSONResponse {
	if len(message) == 0 {
		message = "Unkknown error ocurred"
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

func HandleList[T ModelType](c *gin.Context, fields []string) {

	limit := MaxPageSize
	page := 1
	total := 0
	search := ""
	orderKey := "created_at"
	orderType := "desc"

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
	searchQ := strings.ToLower(c.Query("search"))
	if len(searchQ) > 0 {
		search = searchQ
	}
	var results = []T{}
	err := db.ExecuteOnDB(func(db *gorm.DB) error {
		offset := (page - 1) * limit
		dbCtx := db.Model(&models.Client{})

		if len(search) > 0 {
			var queryStr = ""
			find := "%" + search + "%"
			repeatInterface := make([]interface{}, 0, len(fields))
			for _, field := range fields {
				q := fmt.Sprintf("lower(%s) LIKE ? ", field)
				if len(queryStr) > 0 {
					queryStr += "OR " + q
				} else {
					queryStr += q
				}
				repeatInterface = append(repeatInterface, find)
			}

			dbCtx.Where(queryStr, repeatInterface...)
		}
		count := int64(0)
		var err = dbCtx.Count(&count).Error
		if err != nil {
			return err
		}
		total = int(count)
		if offset > 0 {
			dbCtx.Offset(offset)
		}
		dbCtx.Limit(limit).Order(orderKey + " " + orderType)
		dbCtx.Find(&results)
		return dbCtx.Error
	})
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
