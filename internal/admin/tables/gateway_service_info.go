package tables

import (
	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/GoAdminGroup/go-admin/template/types/form"
	"github.com/lllllan-fv/gateway-admin/internal/admin/models"
)

func GetGatewayServiceInfoTable(ctx *context.Context) table.Table {

	gatewayServiceInfo := table.NewDefaultTable(table.DefaultConfigWithDriver("mysql").SetPrimaryKey("id", db.Bigint))

	setServiceInfo(gatewayServiceInfo.GetInfo())
	setServiceForm(gatewayServiceInfo.GetForm())

	return gatewayServiceInfo
}

func setServiceInfo(info *types.InfoPanel) {
	info.HideFilterArea()

	info.AddField("ID", "id", db.Bigint).FieldFilterable()
	info.AddField("服务名称", "service_name", db.Varchar)
	info.AddField("服务描述", "service_desc", db.Varchar)
	info.AddField("类型", "load_type", db.Tinyint).FieldDisplay(func(model types.FieldModel) interface{} {
		switch model.Value {
		case models.HttpLoadType:
			return "HTTP"
		case models.TcpLoadType:
			return "TCP"
		case models.GrpcLoadType:
			return "GRPC"
		}
		return "unknown"
	})
	// TODO
	info.AddField("服务地址", "address", db.Varchar).
		FieldDisplay(func(model types.FieldModel) interface{} { return "127.0.0.1" })
	info.AddField("QPS", "qps", db.Int).
		FieldDisplay(func(model types.FieldModel) interface{} { return "1" })
	info.AddField("日请求数", "qpd", db.Int).
		FieldDisplay(func(model types.FieldModel) interface{} { return "20" })
	info.AddField("节点数", "nodes_number", db.Int).
		FieldDisplay(func(model types.FieldModel) interface{} { return "2" })
	// info.AddField("Open_auth", "open_auth", db.Tinyint)
	// info.AddField("Black_list", "black_list", db.Varchar)
	// info.AddField("White_list", "white_list", db.Varchar)
	// info.AddField("White_host_name", "white_host_name", db.Varchar)
	// info.AddField("Clientip_flow_limit", "clientip_flow_limit", db.Int)
	// info.AddField("Service_flow_limit", "service_flow_limit", db.Int)
	// info.AddField("Port", "port", db.Int)
	// info.AddField("Header_transfor", "header_transfor", db.Varchar)
	// info.AddField("Rule_type", "rule_type", db.Tinyint)
	// info.AddField("Rule", "rule", db.Varchar)
	// info.AddField("Need_https", "need_https", db.Tinyint)
	// info.AddField("Need_strip_uri", "need_strip_uri", db.Tinyint)
	// info.AddField("Need_websocket", "need_websocket", db.Tinyint)
	// info.AddField("Url_rewrite", "url_rewrite", db.Varchar)
	// info.AddField("Check_method", "check_method", db.Tinyint)
	// info.AddField("Check_timeout", "check_timeout", db.Int)
	// info.AddField("Check_interval", "check_interval", db.Int)
	// info.AddField("Round_type", "round_type", db.Tinyint)
	// info.AddField("Ip_list", "ip_list", db.Varchar)
	// info.AddField("Weight_list", "weight_list", db.Varchar)
	// info.AddField("Forbid_list", "forbid_list", db.Varchar)
	// info.AddField("Upstream_connect_timeout", "upstream_connect_timeout", db.Int)
	// info.AddField("Upstream_header_timeout", "upstream_header_timeout", db.Int)
	// info.AddField("Upstream_idle_timeout", "upstream_idle_timeout", db.Int)
	// info.AddField("Upstream_max_idle", "upstream_max_idle", db.Int)
	// info.AddField("Created_at", "created_at", db.Timestamp)
	// info.AddField("Updated_at", "updated_at", db.Timestamp)

	info.SetTable("gateway_service_info").SetTitle("服务管理").SetDescription("服务列表")
}

func setServiceForm(formList *types.FormPanel) {
	formList.AddField("ID", "id", db.Bigint, form.Default)
	formList.AddField("服务名称", "service_name", db.Varchar, form.Text).FieldMust().FieldNotAllowEdit()
	formList.AddField("服务类型", "service_desc", db.Varchar, form.Text).FieldMust()

	formList.AddField("类型", "load_type", db.Tinyint, form.SelectSingle).
		FieldOptions(types.FieldOptions{
			{Text: "HTTP", Value: models.HttpLoadType},
			{Text: "TCP", Value: models.TcpLoadType},
			{Text: "GRPC", Value: models.GrpcLoadType},
		}).
		FieldOnChooseHide(models.HttpLoadType, "port").
		FieldOnChooseShow(models.HttpLoadType, "rule_type", "rule", "need_https", "need_strip_uri", "need_websocket", "url_rewrite",
			"upstream_connect_timeout", "upstream_header_timeout", "upstream_idle_timeout", "upstream_max_idle").
		FieldOnChooseHide(models.TcpLoadType, "header_transfor").
		FieldOnChooseShow(models.TcpLoadType).
		FieldOnChooseShow(models.GrpcLoadType).
		FieldDefault(models.HttpLoadType).
		FieldDivider("")

	formList.AddField("端口", "port", db.Int, form.Number).FieldHelpMsg("需要设置 8001-8999 范围内数字").FieldDefault("8001").FieldNotAllowEdit()
	formList.AddField("Header/metadata 转换", "header_transfor", db.Varchar, form.TextArea).
		FieldHelpMsg("header 转换支持增加(add)、删除(del)、修改(edit)格式:add headnarme headvalue 多条换行")

	formList.AddRow(func(pa *types.FormPanel) {
		formList.AddField("接入类型", "rule_type", db.Tinyint, form.SelectSingle).
			FieldOptions(types.FieldOptions{
				{Text: "路径", Value: models.UrlPrefixRuleType},
				{Text: "域名", Value: models.DomainRuleType},
			}).FieldDefault(models.UrlPrefixRuleType).FieldRowWidth(1).FieldNotAllowEdit()
		formList.AddField("", "rule", db.Varchar, form.Text).FieldHelpMsg("路径格式:/user/, 域名格式:www.test.com").
			FieldRowWidth(9).FieldNotAllowEdit()
	})
	formList.AddRow(func(pa *types.FormPanel) {
		formList.AddField("支持 https", "need_https", db.Tinyint, form.Switch).
			FieldOptions(types.FieldOptions{
				{Text: "", Value: "1"},
				{Text: "", Value: "0"},
			}).FieldDefault("0").FieldRowWidth(2)
		formList.AddField("支持 strip_uri", "need_strip_uri", db.Tinyint, form.Switch).
			FieldOptions(types.FieldOptions{
				{Text: "", Value: "1"},
				{Text: "", Value: "0"},
			}).FieldDefault("0").FieldRowWidth(3).FieldHeadWidth(3)
		formList.AddField("支持 websocket", "need_websocket", db.Tinyint, form.Switch).
			FieldOptions(types.FieldOptions{
				{Text: "", Value: "1"},
				{Text: "", Value: "0"},
			}).FieldDefault("0").FieldRowWidth(3).FieldHeadWidth(3)
	})
	formList.AddField("URL 重写", "url_rewrite", db.Varchar, form.TextArea).FieldHelpMsg("格式 /gatekeeper/test_service(.*)$1 多条换行")

	formList.AddField("开启验证", "open_auth", db.Tinyint, form.Switch).
		FieldOptions(types.FieldOptions{
			{Text: "", Value: "1"},
			{Text: "", Value: "0"},
		}).FieldDefault("0").FieldDivider("")
	formList.AddField("IP 白名单", "white_list", db.Varchar, form.TextArea).FieldHelpMsg("格式 127.0.0.1 多条换行，白名单优先级高于黑名单")
	formList.AddField("IP 黑名单", "black_list", db.Varchar, form.TextArea).FieldHelpMsg("格式 127.0.0.1 多条换行")
	formList.AddField("客户端限流", "clientip_flow_limit", db.Int, form.Number).FieldHelpMsg("零表示无限制").FieldDefault("0")
	formList.AddField("服务端限流", "service_flow_limit", db.Int, form.Number).FieldHelpMsg("零表示无限制").FieldDefault("0")
	formList.AddField("轮询方式", "round_type", db.Tinyint, form.SelectSingle).
		FieldOptions(types.FieldOptions{
			{Text: "random", Value: models.RandomLoadBalance},
			{Text: "round-robin", Value: models.RoundRobinLoadBalance},
			{Text: "weight_round-robin", Value: models.WeightRoundRobinLoadBalance},
			{Text: "ip_hash", Value: models.IpHashLoadBalance},
		}).FieldDefault(models.HttpLoadType)
	formList.AddField("IP 列表", "ip_list", db.Varchar, form.TextArea).FieldHelpMsg("格式 127.0.0.1:80 多条换行")
	formList.AddField("权重列表", "weight_list", db.Varchar, form.TextArea).FieldHelpMsg("格式 50 多条换行")
	formList.AddField("建立连接超时", "upstream_connect_timeout", db.Int, form.Number).FieldHelpMsg("单位s 零表示无限制").FieldDefault("0")
	formList.AddField("获取 header 超时", "upstream_header_timeout", db.Int, form.Number).FieldHelpMsg("单位s 零表示无限制").FieldDefault("0")
	formList.AddField("链接最大空闲时间", "upstream_idle_timeout", db.Int, form.Number).FieldHelpMsg("单位s 零表示无限制").FieldDefault("0")
	formList.AddField("最大空闲链接数", "upstream_max_idle", db.Int, form.Number).FieldHelpMsg("零表示无限制").FieldDefault("0")

	// formList.AddField("Check_method", "check_method", db.Tinyint, form.Number)
	// formList.AddField("Check_timeout", "check_timeout", db.Int, form.Number)
	// formList.AddField("Check_interval", "check_interval", db.Int, form.Number)
	// formList.AddField("Forbid_list", "forbid_list", db.Varchar, form.Text)
	// formList.AddField("White_host_name", "white_host_name", db.Varchar, form.Text)
	// formList.AddField("Created_at", "created_at", db.Timestamp, form.Datetime)
	// formList.AddField("Updated_at", "updated_at", db.Timestamp, form.Datetime)

	formList.SetTable("gateway_service_info").SetTitle("服务管理").SetDescription("服务管理")
}
