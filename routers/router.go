package routers

import (
	"github.com/astaxie/beego"
	"github.com/kevinxu001/survey/controllers"
)

func init() {
	beego.Router("/", &controllers.MainController{})
}
