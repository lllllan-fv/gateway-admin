package loadbalance

type LoadBalance interface {
	Add(...string) error
	Get(string) (string, error)
	Update()
}
