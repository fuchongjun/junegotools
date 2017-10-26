package study

import (
	"github.com/kardianos/service"
	"os"
	"time"
)

type program struct{}

var myLogger *MyLog

//start,run,stop里面的日志服务启动后
// 在C:\Windows\System32\mylogs目录下，在服务运行的目录
func (p *program) Start(s service.Service) error {
	go p.run()
	myLogger.WriteLog("service have started!", Info)
	return nil
}

func (p *program) run() {
	SoundServer(8000)
}

func (p *program) Stop(s service.Service) error {
	// Stop should not block. Return with a few seconds.
	// 可以做一些事情然后停止,程序停止会回调此方法
	<-time.After(time.Second * 3)
	myLogger.WriteLog("service have stoped", Info)
	return nil
}
func CreateWindowsService() {
	svcConfig := &service.Config{
		Name:        "GoServiceExampleSimple",         //服务名称
		DisplayName: "Go Service Example",             //服务展示名称
		Description: "This is an example Go service.", //服务描述
	}

	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		myLogger.WriteLog(err.Error(), Error)
	}
	//命令行参数 "start", "stop", "restart", "install", "uninstall"
	if len(os.Args) > 1 {
		err = service.Control(s, os.Args[1])
		if err != nil {
			myLogger.WriteLog(err.Error(), Error)
			return
		}
		myLogger.WriteLog(string(os.Args[1])+" success!", Info) //此日志在exe文件运行目录下
		return
	}
	if err != nil {
		myLogger.WriteLog(err.Error(), Error)

	}
	err = s.Run() //不带参数单独执行exe，运行后服务执行
	if err != nil {
		myLogger.WriteLog(err.Error(), Error)
	}
}
