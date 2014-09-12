package routers

import (
	"github.com/astaxie/beego"
	"github.com/kevinxu001/survey/controllers"
)

func init() {
	beego.Router("/", &controllers.IndexController{})
	beego.Router("/dashboard", &controllers.IndexController{}, "get:Dashboard")
	//beego.Router("/blank", &controllers.IndexController{}, "get:Blank")

	beego.Router("/organizations", &controllers.OrganizationController{}, "*:Get")
	beego.Router("/api/organizations", &controllers.OrganizationController{}, "*:GetOrganizations")
	beego.Router("/api/organizations/:id:int", &controllers.OrganizationController{}, "*:GetOrganizationsById")
	beego.Router("/api/organizations/add", &controllers.OrganizationController{}, "*:Add")
	beego.Router("/api/organizations/modify", &controllers.OrganizationController{}, "*:Modify")
	beego.Router("/api/organizations/delete", &controllers.OrganizationController{}, "*:Delete")

	beego.Router("/users", &controllers.UserController{}, "*:Get")
	beego.Router("/api/users", &controllers.UserController{}, "*:GetUsers")
	beego.Router("/api/users/:id:int", &controllers.UserController{}, "*:GetUserById")
	beego.Router("/api/users/add", &controllers.UserController{}, "*:Add")
	beego.Router("/api/users/modify", &controllers.UserController{}, "*:Modify")
	beego.Router("/api/users/delete", &controllers.UserController{}, "*:Delete")

	beego.Router("/surveys", &controllers.SurveyController{}, "*:Get")
	beego.Router("/surveys/:id:int", &controllers.SurveyController{}, "*:GetSurveyById")

	// beego.Router("/login", &controllers.LoginController{}, "get:Get;post:Post")
	// beego.Router("/logout", &controllers.LoginController{}, "*:Logout")
}
