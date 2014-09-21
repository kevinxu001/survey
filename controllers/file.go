package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/kevinxu001/survey/lib"
	"github.com/kevinxu001/survey/models"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type FileController struct {
	CommonController
}

type FileRsp struct {
	Success  bool
	Msg      string
	Id       int
	FileName string
	FileSize int
}

// 获取文件大小的接口
type Size interface {
	Size() int64
}

func (this *FileController) UploadFiles() {
	taskid, err := this.GetInt("taskid")
	if err != nil {
		beego.Error(err)
		this.Data["json"] = &FileRsp{Success: false, Msg: "无法读取任务id！"}
		this.ServeJson()
		return
	}

	o := orm.NewOrm()
	qstask := o.QueryTable("survey_task")

	var surveytask models.SurveyTask
	err = qstask.Filter("id", taskid).One(&surveytask)
	if err != nil {
		beego.Error(err)
		this.Redirect("/errorpage/500", 302)
		return
	}

	tnow := time.Now()
	// year, month, day := t.Date()
	// t = time.Date(year, month, day, 0, 0, 0, 0, time.Local)
	if tnow.Before(surveytask.TaskStarted) {
		this.Data["json"] = &FileRsp{Success: false, Msg: "调研还未开始！"}
		this.ServeJson()
		return
	}
	if tnow.After(surveytask.TaskEnded) {
		this.Data["json"] = &FileRsp{Success: false, Msg: "调研已经结束！"}
		this.ServeJson()
		return
	}

	//判断是否是普通用户登录
	v := this.GetSession("currentUser")
	user, ok := v.(models.User)
	if !ok {
		this.Data["json"] = &FileRsp{Success: false, Msg: "只有注册用户才能参与填报，请退出并重新登录！"}
		this.ServeJson()
		return
	}

	//存储上传的文件
	file, fileHeader, err := this.GetFile("upload-file")
	if err != nil {
		beego.Error(err)
		this.Data["json"] = &FileRsp{Success: false, Msg: "无法读取表单上传文件数据！"}
		this.ServeJson()
		return
	}
	defer file.Close()
	curdir, err := os.Getwd()

	var tofile string
	t := time.Now()
	tdate := t.Format("20060102")
	tofilepath := "/" + beego.AppConfig.String("conf::uploads_dir") + "/" + tdate + "/"
	tofilename := lib.StrToMD5("survey" + t.String() + fileHeader.Filename)
	err = os.MkdirAll(curdir+tofilepath, 0755)
	if err != nil {
		beego.Error(err)
		this.Data["json"] = &FileRsp{Success: false, Msg: "无法创建文件夹！"}
		this.ServeJson()
		return
	}
	tofile = curdir + tofilepath + tofilename
	f, err := os.OpenFile(tofile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		beego.Error(err)
		this.Data["json"] = &FileRsp{Success: false, Msg: "无法创建文件！"}
		this.ServeJson()
		return
	}
	defer f.Close()
	io.Copy(f, file)
	//文件信息写入数据库
	// qsfile := o.QueryTable("file")
	taskfile := new(models.File)
	taskfile.FileName = fileHeader.Filename
	taskfile.FilePath = tofilepath + tofilename
	taskfile.FileExt = filepath.Ext(fileHeader.Filename)
	//获取文件大小
	if sizeInterface, ok := file.(Size); ok {
		taskfile.FileSize = int(sizeInterface.Size())
	}
	taskfile.SurveyTask = &surveytask
	taskfile.User = &user
	beego.Info(taskfile)
	fileid, err := o.Insert(taskfile)
	if err != nil {
		beego.Error(err)
		//入库出错，删除已上传文件
		err = os.Remove(tofile)

		this.Data["json"] = &FileRsp{Success: false, Msg: "文件写入失败！"}
		this.ServeJson()
		return
	}

	this.Data["json"] = &FileRsp{Success: true, Id: int(fileid), FileName: taskfile.FileName, FileSize: taskfile.FileSize}
	this.ServeJson()
}

func (this *FileController) DownloadFileById() {
	fileid, err := strconv.Atoi(this.Ctx.Input.Param(":id"))
	if err != nil {
		beego.Error(err)
		this.Data["json"] = &FileRsp{Success: false, Msg: "找不到文件！"}
		this.ServeJson()
		return
	}

	//判断是否是普通用户登录
	v := this.GetSession("currentUser")
	user, ok := v.(models.User)
	if !ok {
		this.Data["json"] = &FileRsp{Success: false, Msg: "只有注册用户才能参与填报，请退出并重新登录！"}
		this.ServeJson()
		return
	}

	//读取出文件信息
	o := orm.NewOrm()
	qsfile := o.QueryTable("file")
	var file models.File
	err = qsfile.Filter("user__id", user.Id).Filter("id", fileid).One(&file)
	if err != nil {
		beego.Error(err)
		this.Data["json"] = &FileRsp{Success: false, Msg: "找不到文件的数据库信息！"}
		this.ServeJson()
		return
	}
	//开始文件下载
	code := http.StatusOK

	var ctype string
	ctype = mime.TypeByExtension(file.FileExt)
	w := this.Ctx.ResponseWriter
	beego.Info(ctype)
	w.Header().Set("Content-Type", ctype)
	w.Header().Set("Content-disposition", "attachment;filename="+file.FileName)
	w.WriteHeader(code)
	//打开文件输入流
	curdir, err := os.Getwd()
	fromfile := curdir + file.FilePath
	f, err := os.OpenFile(fromfile, os.O_RDONLY, 0644)
	if err != nil {
		beego.Error(err)
		this.Data["json"] = &FileRsp{Success: false, Msg: "无法打开文件！"}
		this.ServeJson()
		return
	}
	defer f.Close()

	io.CopyN(w, f, int64(file.FileSize))

	// this.Data["json"] = &FileRsp{Success: false, Msg: "成功下载文件！"}
	// this.ServeJson()
}

func (this *FileController) DeleteFileById() {
	fileid, err := strconv.Atoi(this.Ctx.Input.Param(":id"))
	if err != nil {
		beego.Error(err)
		this.Data["json"] = &FileRsp{Success: false, Msg: "找不到文件！"}
		this.ServeJson()
		return
	}

	//判断是否是普通用户登录
	v := this.GetSession("currentUser")
	user, ok := v.(models.User)
	if !ok {
		this.Data["json"] = &FileRsp{Success: false, Msg: "只有注册用户才能参与填报，请退出并重新登录！"}
		this.ServeJson()
		return
	}

	//读取出文件信息
	o := orm.NewOrm()
	qsfile := o.QueryTable("file")
	var file models.File
	err = qsfile.Filter("user__id", user.Id).Filter("id", fileid).One(&file)
	if err != nil {
		beego.Error(err)
		this.Data["json"] = &FileRsp{Success: false, Msg: "找不到文件的数据库信息！"}
		this.ServeJson()
		return
	}
	//删除文件
	curdir, err := os.Getwd()
	err = os.Remove(curdir + file.FilePath)
	if err != nil {
		beego.Error(err)
		this.Data["json"] = &FileRsp{Success: false, Msg: "删除文件失败！"}
		this.ServeJson()
		return
	}
	//删除数据库中的文件信息
	_, err = o.Delete(&file)
	if err != nil {
		beego.Error(err)
		this.Data["json"] = &FileRsp{Success: false, Msg: "删除文件的数据库信息失败！"}
		this.ServeJson()
		return
	}

	this.Data["json"] = &FileRsp{Success: true, Msg: "成功删除文件！"}
	this.ServeJson()
}
