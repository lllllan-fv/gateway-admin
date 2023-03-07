package models

type GatewayApp struct {
	ID        uint64
	AppID     string
	Name      string
	Secret    string
	WhiteIPs  string
	QPD       int64
	QPS       int64
	CreatedAt string
	UpdatedAt string
}

func (GatewayApp) TableName() string {
	return "gateway_app"
}
