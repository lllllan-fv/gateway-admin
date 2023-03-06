package tables

import (
	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	"github.com/GoAdminGroup/go-admin/template/types/form"
	"github.com/lllllan-fv/gateway-admin/public/rand"
)

func GetGatewayAppTable(ctx *context.Context) table.Table {

	gatewayApp := table.NewDefaultTable(table.DefaultConfigWithDriver("mysql").SetPrimaryKey("id", db.Bigint))

	info := gatewayApp.GetInfo().HideFilterArea()

	info.AddField("ID", "id", db.Bigint).
		FieldFilterable()
	info.AddField("APP ID", "app_id", db.Varchar)
	info.AddField("租户", "name", db.Varchar)
	info.AddField("密钥", "secret", db.Varchar)
	info.AddField("白名单", "white_ips", db.Varchar)
	info.AddField("日请求数", "qpd", db.Bigint)
	info.AddField("Qps", "qps", db.Bigint)
	info.AddField("创建时间", "created_at", db.Timestamp)
	info.AddField("更新时间", "updated_at", db.Timestamp)

	info.SetTable("gateway_app").SetTitle("租户管理").SetDescription("租户列表")

	formList := gatewayApp.GetForm()
	formList.AddField("ID", "id", db.Bigint, form.Default)
	formList.AddField("APP ID", "app_id", db.Varchar, form.Text).FieldMust().FieldNotAllowEdit()
	formList.AddField("租户", "name", db.Varchar, form.Text).FieldMust()
	formList.AddField("密钥", "secret", db.Varchar, form.Default).FieldDefault(rand.RandStringBytesMaskImpr(30))
	formList.AddField("日请求量限流", "qpd", db.Bigint, form.Number).FieldMust().FieldDefault("0")
	formList.AddField("QPS 限流", "qps", db.Bigint, form.Number).FieldMust().FieldDefault("0")
	// formList.AddField("白名单", "white_ips", db.Varchar, form.Text)
	// formList.AddField("Created_at", "created_at", db.Timestamp, form.Datetime)
	// formList.AddField("Updated_at", "updated_at", db.Timestamp, form.Datetime)

	formList.SetTable("gateway_app").SetTitle("租户管理").SetDescription("修改租户")

	return gatewayApp
}
