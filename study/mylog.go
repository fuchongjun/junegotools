package study

import (
	"fmt"
	"log"
	"os"
	"time"
)

type LogType int

const (
	Info  LogType = iota //0
	Debug                //1
	Error
)

type MyLog struct {
}

func (mylog *MyLog) WriteLog(info string, logtype LogType) error {
	dirpath, err := creatDir("mylogs")
	if err != nil {
		return err
	}
	var fileName string
	switch logtype {
	case Debug:
		log.Println(info)
		fileName = dirpath + time.Now().Format("20060102") + "_debug.log"
	case Info:
		log.Println(info)
		fileName = dirpath + time.Now().Format("20060102") + "_info.log"
	case Error:
		log.Println(info)
		fileName = dirpath + time.Now().Format("20060102") + "_error.log"
	default:
		return fmt.Errorf("未知错误类型")

	}
	writelog(fileName, info)
	return nil
}

func creatDir(dirname string) (dirpath string, err error) {
	var path string
	if os.IsPathSeparator('\\') { //前边的判断是否是系统的分隔符
		path = "\\"
	} else {
		path = "/"
	}
	dir, _ := os.Getwd() //当前的目录
	dirpath = dir + path + dirname + path + time.Now().Format("2006-01-02")
	err = os.MkdirAll(dirpath, os.ModePerm) //在当前目录下生成多级目录
	if err != nil {
		return
	}
	dirpath += path
	return
}

func writelog(fileName string, info string) {
	frefix := time.Now().Format("2006-01-02 15:04:05")
	logFile, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	defer logFile.Close()
	if err != nil {
		log.Fatalln("open file error")
	}
	debugLog := log.New(logFile, frefix+"   ", log.Llongfile)
	debugLog.Println(info)
}
