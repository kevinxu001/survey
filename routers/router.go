package routers

import (
	"github.com/astaxie/beego"
	"github.com/kevinxu001/survey/controllers"
)

func init() {
	beego.Router("/", &controllers.IndexController{})
	beego.Router("/dashboard", &controllers.IndexController{}, "get:Dashboard")
	//beego.Router("/blank", &controllers.IndexController{}, "get:Blank")

	beego.Router("/organizations", &controllers.OrganizationController{})
	beego.Router("/api/organizations", &controllers.OrganizationController{}, "*:GetOrganizations")
	beego.Router("/api/organizations/:id:int", &controllers.OrganizationController{}, "*:GetOrganizationsById")
	beego.Router("/api/organizations/add", &controllers.OrganizationController{}, "*:Add")
	beego.Router("/api/organizations/modify", &controllers.OrganizationController{}, "*:Modify")
	beego.Router("/api/organizations/delete", &controllers.OrganizationController{}, "*:Delete")

	beego.Router("/users", &controllers.UserController{})
	beego.Router("/api/users", &controllers.UserController{}, "*:GetUsers")
	beego.Router("/api/users/admin", &controllers.UserController{}, "*:AdminList")

	// beego.Router("/login", &controllers.LoginController{}, "get:Get;post:Post")
	// beego.Router("/logout", &controllers.LoginController{}, "*:Logout")
}
