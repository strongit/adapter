package logger

import (
	"encoding/json"
	"fmt"
	"github.com/phachon/go-logger"
	"log"
	"os"
	"time"
)

var (
	project = "adapter/"
	log_file_path = "/home/data/logs/"
)

func Println(filename string, level string, msgess ...interface{}){
	Sendlog(filename, level, msgess)
}

/**
	记录日志
	基本用法：
 */
func Sendlog(filename string, level string, msgess ...interface{}){
	msg := fmt.Sprintln(msgess...)//格式化，支持多字段消息

	filepath := log_file_path + project + filename //存储路径
	logger := go_logger.NewLogger()
	logger.SetAsync()

	logger.Detach("console")

	// 命令行输出配置
	consoleConfig := &go_logger.ConsoleConfig{
		Color: true, // 命令行输出字符串是否显示颜色
		JsonFormat: true, // 命令行输出字符串是否格式化
		Format: "", // 如果输出的不是 json 字符串，JsonFormat: false, 自定义输出的格式
	}

	// 添加 console 为 logger 的一个输出
	logger.Attach("console", go_logger.LOGGER_LEVEL_DEBUG, consoleConfig)

	// 文件输出配置
	fileConfig := &go_logger.FileConfig {
		Filename : filepath, // 日志输出文件名，不自动存在
		// 如果要将单独的日志分离为文件，请配置LealFrimeNem参数。
		//LevelFileName : map[int]string {
		//	logger.LoggerLevel("error"): "./error.log",    // Error 级别日志被写入 error .log 文件
		//	logger.LoggerLevel("info"): "./info.log",      // Info 级别日志被写入到 info.log 文件中
		//	logger.LoggerLevel("debug"): "./debug.log",    // Debug 级别日志被写入到 debug.log 文件中
		//},
		MaxSize : 10 * 1024,  // 文件最大值（KB），默认值0不限
		MaxLine : 100000, // 文件最大行数，默认 0 不限制
		DateSlice : "d",  // 文件根据日期切分， 支持 "Y" (年), "m" (月), "d" (日), "H" (时), 默认 "no"， 不切分
		JsonFormat: false, // 写入文件的数据是否 json 格式化
		Format: "%timestamp_format% [%file%:%line%] %body%", // 如果写入文件的数据不 json 格式化，自定义日志格式
	}

	// 添加 file 为 logger 的一个输出
	logger.Attach("file", go_logger.LOGGER_LEVEL_DEBUG, fileConfig)

	switch level {
		case "debug":
			logger.Debug(msg)
			break
		case "error":
			logger.Error(msg)
			break
		case "warning":
			logger.Warning(msg)
			break
		default:
			logger.Info(msg)
	}
	SendToHaina(filename, level, msg)//发送到海纳日志系统

	// 程序结束前必须调用 Flush
	logger.Flush()
}

type msgess struct {
	Log_msg string    `json:"log_msg"`
	Log_plat string   `json:"log_plat"`
	Tbl string        `json:"tbl"`
	Log_time string   `json:"log_time"`
	Log_file string   `json:"log_file"`
	Log_level string  `json:"log_level"`
	Log_su string     `json:"log_su"`
}

/**
	发送一份到海纳实时日志
 */
func SendToHaina(filename string, level string, msg string){
	str := msgess{
		msg,
		"adapter",
		"t_lalalog",
		time.Now().Format("2006-01-02 15:04:05"),
		filename,
		level,
		"",
	}

	b, err := json.Marshal(str)
	if err != nil {
		log.Println("SendToHaina json error", err)
	}

	SendHaina(string(b))
}


/**
	记录日志(自己实现，暂时不用吧)
 */
func Sendlocaldebug(filename string, msg ...string){
	filepath := log_file_path + project + filename //存储路径

	_, err := os.Stat(filepath)
	if os.IsNotExist(err) {//返回true表示文件不存在
		file, err := os.Create(filepath)
		if err != nil {
			log.Println("fail to create file!", file)
		}
	}

	file, err := os.OpenFile(filepath, os.O_WRONLY|os.O_APPEND|os.O_CREATE|os.O_SYNC, 0666) //打开文件
	if err != nil {
		log.Println("fail to open file!", file)
	}

	logger := log.New(file, "", log.Llongfile)//设置文件路径
	logger.SetFlags(log.LstdFlags)//设置写入格式

	logger.Println(msg) //写入文件内容
}












