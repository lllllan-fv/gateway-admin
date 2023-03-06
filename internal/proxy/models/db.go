package models

import (
	"github.com/jinzhu/gorm"
	"github.com/lllllan-fv/gateway-admin/internal/admin/models"
)

func GetDB() *gorm.DB {
	return models.GetDB()
}
