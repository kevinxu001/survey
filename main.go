package main

import (
	"mime"
	"os"

	"github.com/astaxie/beego"
	"github.com/kevinxu001/survey/models"
	_ "github.com/kevinxu001/survey/routers"
)

const (
	VERSION = "0.1"
)

func main() {
	// 初始化系统参数
	beego.Info("Start init system...")
	initialize()
	beego.Info("Init system successful")

	beego.Info("Start system server...")
	beego.Run()
}

func initialize() {
	mime.AddExtensionType(".css", "text/css")
	// 判断命令行参数并执行初始化操作
	initArgs()

	models.ConnectDB()

	//beego.AddFuncMap("stringsToJson", StringsToJson)
}
func initArgs() {
	args := os.Args
	for _, v := range args {
		if v == "-syncdb" {
			models.SyncDB()
			os.Exit(0)
		}
	}
}
