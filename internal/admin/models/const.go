package models

const (
	HttpLoadType = "1"
	TcpLoadType  = "2"
	GrpcLoadType = "3"

	RandomLoadBalance           = "1"
	RoundRobinLoadBalance       = "2"
	WeightRoundRobinLoadBalance = "3"
	IpHashLoadBalance           = "4"

	UrlPrefixRuleType = "0"
	DomainRuleType    = "1"
)
