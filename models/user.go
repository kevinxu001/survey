// user

// @description 用户模型定义

// @link

// @license     https://github.com/kevinxu001/houserent/blob/master/LICENSE

// @authors     kevinxu

package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"time"
)

type User struct {
	Id       int
	UserName string `orm:"size(50)" form:"UserName" valid:"Required;MaxSize(50)"`
	PassWord string `orm:"size(40)" form:"PassWord" valid:"Required;MaxSize(40)"`
	RealName string `orm:"size(50);null" form:"RealName" valid:"Required;MaxSize(50)"`
	Mobile   string `orm:"size:(20);null" form:"Mobile" valid:"Required;Mobile"`
	Phone    string `orm:"size:(20);null" form:"Phone" valid:"Phone"`
	IdCard   string `orm:"size:(20);null" form:"IdCard" valid:"Required;Match(/\d{17}[\dXx]{1}/)"`
	//Avatar    string    `orm:"size(100)" form:"Avatar" valid:"MaxSize(100)"`
	Created   time.Time         `orm:"auto_now_add;type(datetime)" form:"Created"`
	Updated   time.Time         `orm:"auto_now;type(datetime)" form:"Created"`
	LastLogin time.Time         `orm:"type(datetime)"`
	Status    uint8             `orm:"default(0)" form:"Status" valid:"Range(0，1)"`
	OrgUnit   *OrganizationUnit `orm:"rel(fk)"`
}

// func (o *OrganizationUnit) TableName() string {
// 	return "organization_unit"
// }

// func (o *OrganizationUnit) TableIndex() [][]string {
// 	return [][]string{
// 		[]string{"Id","UnitName"},
// 	}
// }

// func (o *OrganizationUnit) TableUnique() [][]string {
// 	return [][]string{
// 		[]string{"UnitName"}
// 	}
// }

func init() {
	orm.RegisterModelWithPrefix(beego.AppConfig.String("tableprefix"), new(User))
}
