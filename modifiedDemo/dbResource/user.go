package dbResource

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	_ "github.com/go-sql-driver/mysql"
	//"errors"
	"time"
	"strconv"
)

//结构定义
type User struct {
	gorm.Model
	Name string
	Gender uint
	Birthday time.Time
}

//自定义string转uint方法
func StringToUint(str string) (i uint){
	tmpInt, _ := strconv.Atoi(str)
	return uint(tmpInt)
}
//自定义string转timestamp方法
func StringToTimestamp(str string) (t time.Time) {
	timestamp, _ := time.Parse("2006-01-02", str)
	return timestamp
}
//使用string类型的属性名和值设置属性
func (u *User) setAttr(attrName string, valueStr string){
	switch attrName {
	case "name":
		u.Name = valueStr
	case "gender":
		u.Gender = StringToUint(valueStr)
	case "birthday":
		u.Birthday = StringToTimestamp(valueStr)
	}
}

//根据传入param修改属性
func (u *User) ModifyAttrWithParams(strParams map[string]string){
	fmt.Println("Modify Func is Call")
	for attrName, valueStr := range strParams {
		//暂时不考虑实际需要设置为空的情况
		if valueStr != "" {
			u.setAttr(attrName, valueStr)
		}		
	}
}

//添加数据
func (u *User) CreateDBData() (status int,retJSON map[string]interface{}, err error) {
	fmt.Println("Create Func is Call",u.Name, u.Gender, u.Birthday)
	var pDB *gorm.DB
	pDB, err = GetDB()
	retJSON = make(map[string]interface{})
	if pDB == nil {
		status = 500
	} else {
		(*pDB).Create(u)
		status = 200
		retJSON["id"] = u.ID
		retJSON["name"] = u.Name
		retJSON["gender"] = u.Gender 
		retJSON["birthday"] = u.Birthday
	}
	return
}
//读取数据
func (u *User) RetrieveDBDataByID(id uint) (status int, retJSON map[string]interface{}, err error) {
	fmt.Println("Retrieve Func is Call")
	var pDB *gorm.DB
	pDB, err = GetDB()
	retJSON = make(map[string]interface{})
	if pDB == nil {
		status = 500
	} else {
		//初始化ID为无效值
		u.ID = 0
		(*pDB).First(u, id)
		if u.ID == 0 {
			status = 404
		} else {
			status = 200
			retJSON["id"] = u.ID
			retJSON["name"] = u.Name
			retJSON["gender"] = u.Gender
			retJSON["birthday"] = u.Birthday
		}
	}
	return
}
//更新数据
func (u *User) UpdateDBDataByID(id uint, strParams map[string]string) (status int, retJSON map[string]interface{}, err error) {
	fmt.Println("Update Func is Call")
	var pDB *gorm.DB
	pDB, err = GetDB()
	retJSON = make(map[string]interface{})
	if pDB == nil {
		status = 500
	} else {
		//设置ID为无效值
		u.ID = 0
		(*pDB).First(u, id)
		if u.ID == 0 {
			status = 404
		} else {
			//修改属性
			u.ModifyAttrWithParams(strParams)
			//保持数据
			(*pDB).Save(u)
			status = 200
			retJSON["id"] = u.ID
			retJSON["name"] = u.Name
			retJSON["gender"] = u.Gender
			retJSON["birthday"] = u.Birthday
		}
	}
	return
}
//删除数据
func (u *User) DeleteDBDataByID(id uint) (status int, retJSON map[string]interface{}, err error) {
	fmt.Println("Delete Func is Call")
	var pDB *gorm.DB
	pDB, err = GetDB()
	retJSON = make(map[string]interface{})
	if pDB == nil {
		status = 500
	} else {
		//设置ID为无效值
		u.ID = 0
		(*pDB).First(u, id)
		if u.ID == 0 {
			status = 404
		} else {
			(*pDB).Delete(u)
			status = 200
			retJSON["id"] = u.ID
			retJSON["name"] = u.Name
			retJSON["gender"] = u.Gender
			retJSON["birthday"] = u.Birthday
		}
	}
	return
}
//获取属性列表
func (u *User) GetAttrNameList() ([]string) {
	return []string{"name", "gender", "birthday"}
}

