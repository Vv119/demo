package main
import (
	"strconv"
	//"errors"
	"fmt"
	"time"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	_ "github.com/go-sql-driver/mysql"
)

//定义User结构
type User struct {
	gorm.Model
	Name string `form:"user" binding:"required"`
	Gender uint `form:"gender" binding:"required"`	
	Birthday time.Time `form:"birthday" binding:"required"`
}

//定义操作结果
const (
	//操作成功
	success = iota
	//参数不合法
	paramInvalid
	//找不到记录
	recordNotExists
)

//数据库连接
var database gorm.DB

func initDatabase() (gorm.DB, bool) {
	var err error
	var db gorm.DB
	var initOK bool
	//连接数据库
        db, err = gorm.Open("mysql", "root:root@/demoDB?charset=utf8&parseTime=True&loc=Local")
        if err != nil {
        	fmt.Println("连接数据库失败:",err)
		initOK = false
        }else{
		//关闭复数化
		db.SingularTable(true)
		//同步数据库表结构
		db.AutoMigrate(&User{})
		initOK = true
	}
	return	db, initOK
}

//添加用户
func addUser(c *gin.Context){
	fmt.Println("addUser调用")
	opStatus := success
	//取FORM参数
	nameParam := c.PostForm("name")
	genderParam := c.PostForm("gender")
	birthdayParam := c.PostForm("birthday")
	//格式转换
	genderInt, _ := strconv.Atoi(genderParam)
	gender := uint(genderInt)
	birthday, _ := time.Parse("2006-01-02", birthdayParam)
	//创建对象
	user := User{}
	user.Name = nameParam
	user.Gender = gender
	user.Birthday = birthday
	//数据库操作
	database.Create(&user)
	//返回值封装
	c.JSON(http.StatusOK, gin.H{
		"op": "C",
		"status": opStatus,
		"ID": user.ID, 
		"name": user.Name, 
		"gender": user.Gender,
		"birthday": user.Birthday,
		"created_at": user.CreatedAt,
	})
}

//删除用户
func deleteUserByID(c *gin.Context){
	fmt.Println("deleteUserByID调用")
	opStatus := success
	//取URL参数
	idParam := c.Param("id")
	//格式转换
	idInt, _ := strconv.Atoi(idParam)
	id := uint(idInt)
	//创建对象
	user := User{}
	//数据库操作
	if id > 0 {
		user.ID = id
	       	ret := database.Delete(&user)
		fmt.Println("Delete ret:", ret)
	} else {
		opStatus = paramInvalid
	}
	//返回值封装
	c.JSON(http.StatusOK, gin.H{
		"op": "D",
		"status": opStatus,
		"ID": user.ID,
		"deleted_at": user.DeletedAt,
	})
}

//更新用户
func updateUserByID(c *gin.Context){
	fmt.Println("updateUserByID调用")
	opStatus := success
	//取URL参数
	idParam := c.Param("id")
	//取FORM参数
	nameParam := c.PostForm("name")
	genderParam := c.PostForm("gender")
	birthdayParam := c.PostForm("birthday")
	//格式转换
	idInt, _ := strconv.Atoi(idParam)
	id := uint(idInt)
	genderInt,_ := strconv.Atoi(genderParam)
	gender := uint(genderInt)
	birthday, _ := time.Parse("2006-01-02", birthdayParam)
	//创建对象
	user := User{}
	//数据库操作
	if id > 0 {
		user.ID = id
		database.First(&user)
		if nameParam != "" {
			user.Name = nameParam
		}
		if genderParam != "" {
			user.Gender = gender
		}
		if birthdayParam != "" {
			user.Birthday = birthday
		}
		database.Save(&user)
	} else {
		opStatus = paramInvalid
	}
	//返回封装值
	c.JSON(http.StatusOK, gin.H{
		"op": "U",
		"opStatus": opStatus,
		"ID": user.ID,
		"name": user.Name,
		"gender": user.Gender,
		"birthday": user.Birthday,
		"updated_at": user.UpdatedAt,
	})
}

//查询用户列表
func getUserList(c *gin.Context){
	fmt.Println("getUserList调用")
	opStatus := success
	//取URL参数
	idParam := c.Query("id")
	nameParam := c.Query("name")
	genderParam := c.Query("gender")
	birthdayParam := c.Query("birthday")
	limitParam := c.Query("limit")
	offsetParam := c.Query("offset")
	orderParam := c.Query("order")
	//格式转换
	idInt, _ := strconv.Atoi(idParam)
	id := uint(idInt)
	genderInt, _ := strconv.Atoi(genderParam)
	gender := uint(genderInt)
	birthday, _ := time.Parse("2006-01-02", birthdayParam)
	limitInt, _ := strconv.Atoi(limitParam)
	offsetInt, _ := strconv.Atoi(offsetParam)
	//创建对象
	var users []User
	queryCondition := User{}
	if idParam != "" {
		queryCondition.ID = id
	}
	if nameParam != "" {
		queryCondition.Name = nameParam
	}
	if genderParam != "" {
		queryCondition.Gender = gender
	}
	if birthdayParam != "" {
		queryCondition.Birthday = birthday
	}
	if limitParam == "" {
		limitInt = -1
	}
	if offsetParam == "" {
		offsetInt = -1
	}
	//数据库操作
	database.Limit(limitInt).Offset(offsetInt).Order(orderParam).Find(&users, queryCondition)
	//封装结果集
	dataList := make(map[string]interface{})
	for k,v := range users {
		dataList[strconv.Itoa(k)] = gin.H{
			"ID": v.ID,
			"name": v.Name,
			"gender": v.Gender,
			"birthday": v.Birthday,
			"created_at": v.CreatedAt,
			"updated_at": v.UpdatedAt,
		}
	}
	//封装返回值
	c.JSON(http.StatusOK,gin.H{
		"op": "R",
		"status": opStatus,
		"data": dataList,
	})
}

//查询指定用户
func getUserByID(c *gin.Context){
	fmt.Println("getUserByID调用")
	opStatus := success
	//取URL参数
	idParam := c.Param("id")
	//格式转换
	idInt, _ := strconv.Atoi(idParam)
	id := uint(idInt)
	//创建对象
	user := User{}
	//数据库操作
	if id > 0 {
		database.First(&user, id)
		fmt.Println("Found ret:", user.ID)
	} else {
		opStatus = paramInvalid
	}
	c.JSON(http.StatusOK, gin.H{
		"op": "R",
		"status": opStatus,
		"ID": user.ID,
		"name": user.Name,
		"gender": user.Gender,
		"birthday": user.Birthday,
		"created_at": user.CreatedAt,
		"updated_at": user.UpdatedAt,
	})
}

//返回API列表
func getAPIList(c *gin.Context){
	fmt.Println("getAPIList调用")
	c.JSON(http.StatusOK, gin.H{
		"user_list_url": "http://localhost:8080/user?{id=queryID}{&name=queryName}{&gender=queryGender}{&birthday=queryBirthday}{&limit=retLimit{&offset=queryOffset}}{&order=retOrder}",
		"specific_user_url": "http://localhost:8080/user/:id",
	})
}

func main() {
	//初始化数据库
	var ok bool
	if database, ok = initDatabase(); ok {
		fmt.Println("数据库连接成功")
	} else {
		return
	}

	//使用默认路由
	router := gin.Default()
	
	//v1版本
	v1 := router.Group("/v1")
	{
		//API for GET user
		v1.GET("/user",getUserList)
		v1.GET("/user/:id", getUserByID)

		//API for POST user
		v1.POST("/user", addUser)

		//API for PUT user
		v1.PUT("/user/:id", updateUserByID)

		//API for DELETE user
		v1.DELETE("/user/:id", deleteUserByID)
	}
	//API List
	router.GET("/", getAPIList)

	//开启服务
	router.Run(":8080")
}
