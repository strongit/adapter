package main

import (
	"adapter/config"
	"adapter/modules/logger"
	"adapter/modules/process"
	"adapter/modules/simpleHTTP"
	"gopkg.in/alecthomas/kingpin.v2"
	_ "net/http/pprof"
)
var (
	env = kingpin.Flag("env", "Running environment dev or prd").Default("prd").String()
	pdhost = kingpin.Flag("pdhost", "Connect pdserver format IP:port").Default("").String()
	port = kingpin.Flag("port", "Adapter start port").Default("12350").String()
)

func main() {
	kingpin.Parse()

	var PDHost string
	var TimeInterval int

	if *env == "dev" {
		PDHost = config.PDHost_dev //tikv连接地址
	}else{
		PDHost = config.PDHost_prd
	}
	TimeInterval = config.TimeInterval

	if *pdhost == "" {
		*pdhost = PDHost
	}

	logger.Sendlog("appstart.log", "info", *env, *pdhost, TimeInterval, *port)
	// init
	process.Init(*pdhost, TimeInterval)
	// http server
	simpleHTTP.Server(*port)
}








