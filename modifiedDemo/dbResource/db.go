package dbResource

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	_ "github.com/go-sql-driver/mysql"
)

type DBResource interface{
	CreateDBData() (int, map[string]interface{}, error)
	RetrieveDBDataByID(uint) (int, map[string]interface{}, error)
	UpdateDBDataByID(uint,map[string]string) (int, map[string]interface{}, error)
	DeleteDBDataByID(uint) (int, map[string]interface{}, error)
	GetAttrNameList() ([]string)
	ModifyAttrWithParams(map[string]string)
}
var pDB *gorm.DB 
func GetDB() (*gorm.DB, error){
	var db gorm.DB
	var err error
	if pDB == nil {
		db, err = gorm.Open("mysql", "root:root@/demoDB?charset=utf8&parseTime=True&loc=Local")
		if err != nil {
			fmt.Println("连接数据库失败:", err)
		} else {
			//关闭复数化
			db.SingularTable(true)
			//同步表结构
			db.AutoMigrate(&User{})
			
			pDB = &db
		}
	}
	return pDB, err
}
