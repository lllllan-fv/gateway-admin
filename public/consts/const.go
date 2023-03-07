package consts

const (
	HttpLoadType = 1
	TcpLoadType  = 2
	GrpcLoadType = 3

	PrefixURLHTTPRuleType = 0
	DomainHTTPRuleType    = 1

	HTTPRuleTypePrefixURL = 0
	HTTPRuleTypeDomain    = 1

	FlowTotal         = "flow_total"
	FlowServicePrefix = "flow_service_"
	FlowAppPrefix     = "flow_app_"

	RedisFlowDayKey  = "flow_day_count"
	RedisFlowHourKey = "flow_hour_count"

	JwtSignKey = "my_sign_key"
	JwtExpires = 60 * 60
)
