package routers

import (
	"github.com/kevinxu001/survey1/controllers"
	"github.com/astaxie/beego"
)

func init() {
    beego.Router("/", &controllers.MainController{})
}
