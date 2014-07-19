package models

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

var o orm.Ormer

func SyncDB() {
	createDB()

	ConnectDB()

	o = orm.NewOrm()
	// 数据库别名
	name := "default"
	// 强制重新建数据库
	force := true
	// 打印执行过程
	verbose := true
	// 遇到错误立即返回
	err := orm.RunSyncdb(name, force, verbose)
	if err != nil {
		beego.Error(err)
	}
	// insertUser()
	// insertGroup()
	// insertRole()
	// insertNodes()
	beego.Info("Database init success.\nPlease restart application server.")
}

//连接数据库
func ConnectDB() {
	db_type := beego.AppConfig.String("db_type")
	db_host := beego.AppConfig.String("db_host")
	db_port := beego.AppConfig.String("db_port")
	db_user := beego.AppConfig.String("db_user")
	db_pass := beego.AppConfig.String("db_pass")
	db_name := beego.AppConfig.String("db_name")
	db_path := beego.AppConfig.String("db_path")
	db_sslmode := beego.AppConfig.String("db_sslmode")

	var datasource string

	switch db_type {
	case "mysql":
		orm.RegisterDriver("mysql", orm.DR_MySQL)
		datasource = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8", db_user, db_pass, db_host, db_port, db_name)
		break
	case "postgres":
		orm.RegisterDriver("postgres", orm.DR_Postgres)
		datasource = fmt.Sprintf("dbname=%s host=%s  user=%s  password=%s  port=%s  sslmode=%s", db_name, db_host, db_user, db_pass, db_port, db_sslmode)
	case "sqlite3":
		orm.RegisterDriver("sqlite3", orm.DR_Sqlite)
		if db_path == "" {
			db_path = "./"
		}
		datasource = fmt.Sprintf("%s%s.db", db_path, db_name)
		break
	default:
		beego.Critical("Database driver not support: ", db_type)
	}

	orm.RegisterDataBase("default", db_type, datasource)
}

//创建数据库
func createDB() {
	db_type := beego.AppConfig.String("db_type")
	db_host := beego.AppConfig.String("db_host")
	db_port := beego.AppConfig.String("db_port")
	db_user := beego.AppConfig.String("db_user")
	db_pass := beego.AppConfig.String("db_pass")
	db_name := beego.AppConfig.String("db_name")
	db_path := beego.AppConfig.String("db_path")
	db_sslmode := beego.AppConfig.String("db_sslmode")

	var datasource string
	var sqlstring string

	switch db_type {
	case "mysql":
		datasource = fmt.Sprintf("%s:%s@tcp(%s:%s)/?charset=utf8", db_user, db_pass, db_host, db_port)
		sqlstring = fmt.Sprintf("CREATE DATABASE  if not exists `%s` CHARSET utf8 COLLATE utf8_general_ci", db_name)
		break
	case "postgres":
		datasource = fmt.Sprintf("host=%s  user=%s  password=%s  port=%s  sslmode=%s", db_host, db_user, db_pass, db_port, db_sslmode)
		sqlstring = fmt.Sprintf("CREATE DATABASE %s", db_name)
		break
	case "sqlite3":
		if db_path == "" {
			db_path = "./"
		}
		datasource = fmt.Sprintf("%s%s.db", db_path, db_name)
		os.Remove(datasource)
		sqlstring = "create table init (n varchar(32));drop table init;"
		break
	default:
		beego.Critical("Database driver not support: ", db_type)
	}

	db, err := sql.Open(db_type, datasource)
	defer db.Close()
	if err != nil {
		panic(err.Error())
	}
	r, err := db.Exec(sqlstring)
	if err != nil {
		log.Println(err)
		log.Println(r)
	} else {
		log.Println("Database ", db_name, " created")
	}
}

func insertUser() {
	fmt.Println("insert user ...")
	// u := new(User)
	// u.Username = "admin"
	// u.Nickname = "admin"
	// u.Password = "admin"
	// u.Email = "admin@admin.com"
	// u.Remark = "I'm admin"
	// u.Status = 2
	// o = orm.NewOrm()
	// o.Insert(u)
	fmt.Println("insert user end")
}
