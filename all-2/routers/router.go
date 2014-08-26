package routers

import (
	"github.com/astaxie/beego"
	"github.com/kevinxu001/survey/controllers"
)

func init() {
	beego.Router("/", &controllers.IndexController{})

	beego.Router("/organizations", &controllers.OrganizationController{})
	beego.Router("/organizations/add", &controllers.OrganizationController{}, "*:Add")

	beego.Router("/users", &controllers.UserController{})
	beego.Router("/users/admin", &controllers.UserController{}, "get:AdminList")

	beego.Router("/login", &controllers.LoginController{}, "get:Get;post:Post")
	beego.Router("/logout", &controllers.LoginController{}, "*:Logout")
}
