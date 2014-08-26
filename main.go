package main

import (
	"fmt"
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
	beego.Info("Start init survey system...")
	initialize()
	beego.Info("Init system successful")

	beego.Info("Start survey system server...")
	beego.Run()
}

func initialize() {
	mime.AddExtensionType(".css", "text/css")
	// 判断命令行参数并执行初始化操作
	initArgs()
	//连接数据库
	models.ConnectDB()
	//初始化cache
	models.InitCache()
	//初始化session
	models.InitSession()
	//初始化logs
	models.InitLogs()
	//初始化i18n
	models.InitI18n()
	//初始化RBAC权限控制
	models.InitAccessControl()

	//beego.AddFuncMap("stringsToJson", StringsToJson)
}

var usage = `Survey is a on-line survey web application.

Usage:

	survey [-syncdb | -h]

Description:
	
	-syncdb	Create database.This will delete all tables and data.

	-h	Show help page.

Use "bee -h" for this help.
`

func initArgs() {
	args := os.Args
	for _, v := range args {
		if v == "-syncdb" {
			models.SyncDB()
			os.Exit(0)
		} else if v == "-h" {
			fmt.Println(usage)
			os.Exit(0)
		}
	}
}
