package process

import (
	//"github.com/BurntSushi/toml"
	"adapter/lib"
	"adapter/modules/conf"
	"adapter/modules/tikv"
	"log"
)

// Init is init data
func Init(PDHost string, TimeInterval int) {
	// init runtime
	//if _, err := toml.DecodeFile(confPath, &conf.RunTimeMap); err != nil {
	//	log.Println(err)
	//	return
	//}
	//conf.RunTimeInfo = conf.RunTimeMap[runTime]
	conf.RunTimeInfo.PDHost = PDHost
	conf.RunTimeInfo.TimeInterval = TimeInterval
	log.Println("runtimeinfo", conf.RunTimeInfo)
	// init log
	lib.InitLog()
	// init tikv
	tikv.Init()
}
