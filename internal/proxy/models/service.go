package models

type GatewayServiceInfo struct {
	ID                     uint
	LoadType               int
	ServiceName            string
	ServiceDesc            string
	CreatedAt              string
	UpdatedAt              string
	OpenAuth               int
	BlackList              string
	WhiteList              string
	WhiteHostName          string
	ClientIPFlowLimit      int
	ServiceFlowLimit       int
	HeaderTransfor         string
	RuleType               int
	Rule                   string
	NeedHTTPS              int
	NeedStripURI           int
	NeedWebSocket          int
	URLRewrite             string
	Port                   int
	CheckMethod            int
	CheckTimeout           int
	CheckInterval          int
	RoundType              int
	IPList                 string
	WeightList             string
	ForbidList             string
	UpstreamConnectTimeout int
	UpstreamHeaderTimeout  int
	UpstreamIdleTimeout    int
	UpstreamMaxIdle        int
}

func (GatewayServiceInfo) TableName() string {
	return "gateway_service_info"
}
