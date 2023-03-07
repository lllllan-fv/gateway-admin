package dao

import "github.com/lllllan-fv/gateway-admin/internal/proxy/models"

func ListApp() (list []*models.GatewayApp) {
	models.GetDB().Find(&list)
	return
}
