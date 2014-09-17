// organization

// @description 组织机构模型定义

// @link

// @license     https://github.com/kevinxu001/houserent/blob/master/LICENSE

// @authors     kevinxu

package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type OrganizationUnit struct {
	Id       int
	Pid      int     `orm:"default(0)" form:"Pid" valid:"Numeric"`
	UnitName string  `orm:"size(100)" form:"UnitName" valid:"Required;MaxSize(100)"`
	SortRank uint8   `orm:"default(1)" form:"SortRank" valid:"Range(1,100)"`
	Remark   string  `orm:"null;size(200)" form:"Remark" valid:"MaxSize(200)"`
	Status   uint8   `orm:"default(1)" form:"Status" valid:"Range(0，1)"`
	Users    []*User `orm:"reverse(many)"`
	//ChildOrganizationUnits []*OrganizationUnit `orm:"-"`
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
	orm.RegisterModelWithPrefix(beego.AppConfig.String("tableprefix"), new(OrganizationUnit))
}
