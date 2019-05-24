package logger

import (
	"bytes"
	"compress/zlib"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"os"
	"time"
)

var ports = [5]string{"57001", "57002", "57003", "57004", "57005"}
var hosts = map[string]string{"dev":"172.18.84.153", "prd":"haina-upload.myhll.cn"}

/**
	获取环境变量
 */
func Getenv() string{
	env := os.Getenv("DEVELOP")
	if env == ""{//不存在默认走prd
		env = "prd"
	}
	return env
}

/**
	获取udp请求地址
 */
func Getaddr() string{
	env := Getenv()
	host := hosts[env]

	rand.Seed(time.Now().UnixNano())
	key := rand.Intn(5)

	port := ports[key]

	addr := host + ":" + port
	return addr
}

/**
	udp发送给海纳
 */
func SendHaina(msg string){
	//获取udpaddr
	addr := Getaddr()
	udpaddr, err := net.ResolveUDPAddr("udp4", addr);
	chkError(err);

	//连接，返回udpconn
	udpconn, err2 := net.DialUDP("udp", nil, udpaddr);
	chkError(err2);

	//写入数据
	_, err3 := udpconn.Write(ZipBytes(msg));
	chkError(err3);
}

/**
	异常记录日志
 */
func chkError(err error) {
	if err != nil {
		log.Println(err);
	}
}

//压缩数据
func ZipBytes(input string) []byte {
	msg := []byte(input)
	var buf bytes.Buffer

	compressor, err := zlib.NewWriterLevel(&buf, zlib.BestCompression)
	if err != nil {
		chkError(err)
		return msg
	}

	compressor.Write(msg)
	compressor.Close()

	return buf.Bytes()
}

//解压数据
func UnzipBytes(input []byte) []byte {
	b := bytes.NewReader(input)
	r, err := zlib.NewReader(b)

	defer r.Close()
	if err != nil {
		panic(err)
	}

	data, _ := ioutil.ReadAll(r)

	return data
}






