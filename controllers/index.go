package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/kevinxu001/survey/lib"
	"github.com/kevinxu001/survey/models"
	"strconv"
	"time"
)

type IndexController struct {
	CommonController
}

type EventsData struct {
	//某天的描述，“后天”，“今天”，“昨天”，“星期一”，“8月31日”，“2013年12月15日”
	Today  bool
	Day    string
	Events []*models.Event
}

func (this *IndexController) Get() {
	this.Redirect("/dashboard", 302)
}

func (this *IndexController) Dashboard() {
	this.TplNames = "index.html"
	this.Data["isDashboard"] = true

	//暂时先直接查找所有的调研任务
	o := orm.NewOrm()
	qs := o.QueryTable("event")

	var events []*models.Event

	t := time.Now()
	year, month, day := t.Date()
	t = time.Date(year, month, day, 0, 0, 0, 0, time.Local)
	//startTime := time.Date(year-1, month, day, 0, 0, 0, 0, time.Local)
	//endTime := time.Date(year, month, day+3, 0, 0, 0, 0, time.Local)
	startTime := t.AddDate(-1, 0, 0)
	endTime := t.AddDate(0, 0, 3)

	_, err := qs.Filter("eventtime__gte", startTime).Filter("eventtime__lte", endTime).OrderBy("-eventtime").All(&events)
	if err != nil {
		beego.Error(err)
		return
	}

	var eventsList []*EventsData
	var es []*models.Event = make([]*models.Event, 0)
	var (
		first         bool
		prev, current string
	)

	for _, event := range events {
		current = lib.TimeToDesc(event.EventTime)

		if current != prev {
			first = true
		} else {
			first = false
		}

		if first {
			if len(es) > 0 {
				if prev == "今天" {
					eventsList = append(eventsList, &EventsData{Today: true, Day: prev, Events: es})
				} else {
					eventsList = append(eventsList, &EventsData{Today: false, Day: prev, Events: es})
				}

			}
			//清空原有的es
			es = make([]*models.Event, 0)
		}
		es = append(es, event)

		prev = current
	}
	if len(es) > 0 {
		if prev == "今天" {
			eventsList = append(eventsList, &EventsData{Today: true, Day: prev, Events: es})
		} else {
			eventsList = append(eventsList, &EventsData{Today: false, Day: prev, Events: es})
		}
	}

	this.Data["EventsList"] = &eventsList
}

func (this *IndexController) ErrorPage() {
	errno, _ := strconv.Atoi(this.Ctx.Input.Param(":errno"))
	switch errno {
	case 404:
		this.TplNames = "error/404.html"
	case 500:
		this.TplNames = "error/500.html"
	default:
		this.TplNames = "error/404.html"
	}
}

// func (this *IndexController) Blank() {
// 	this.TplNames = "blank.html"
// }
