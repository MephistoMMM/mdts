package discovery

// BrokerInfo is used to register broker into etcd
type BrokerInfo struct {
	IP        string
	Port      string
	Type      string
	Hash      string
	StartTime int64
}
