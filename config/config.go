package config

type tikvCfg struct {
	PDHost string
	TimeInterval int
}

var devtikvCfg = tikvCfg{
	PDHost: "192.168.56.102:2379",
	TimeInterval: 5,
}

var prdtikvCfg = tikvCfg{
	PDHost: "172.18.84.169:32769",
	TimeInterval: 5,
}

const (
	PDHost_dev = "192.168.56.102:2379"
	//PDHost_prd = "172.18.84.169:32769"
	PDHost_prd = "172.18.191.94:2379"
	TimeInterval = 5
)


