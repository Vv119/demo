package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	. "github.com/yuanpinda/modifiedDemo/dbResource"
)


//获取指定ID的资源
func getResourceByID(c *gin.Context){
	//资源名
	resourceName := c.Param("resourceName")
	//唯一ID
	idStr := c.Param("id")
	id := StringToUint(idStr)
	fmt.Println("getResourceByID", resourceName, id)
	//根据资源名获取资源对象
	resObject := GetDBResourceObjectByName(resourceName)
	if resObject != nil {
		status,retInfo,_ := resObject.RetrieveDBDataByID(id)
		c.JSON(status, retInfo)
	}else {
		c.JSON(404, gin.H{
				
		})
		fmt.Println("Invalid Res Name:", resourceName)
	}
}
//更新指定ID的资源
func updateResourceByID(c *gin.Context){
	//资源名
	resourceName := c.Param("resourceName")
	//唯一ID
	idStr := c.Param("id")
	id := StringToUint(idStr)
	fmt.Println("updateResourceByID", resourceName, id)
	//根据资源名获取资源对象
	resObject := GetDBResourceObjectByName(resourceName)
	if resObject != nil {
		attrNameList := resObject.GetAttrNameList()
		formParams := make(map[string]string)
		for _, attrName := range attrNameList {
			value := c.PostForm(attrName)
			formParams[attrName] = value
		}
		status, retInfo, _ := resObject.UpdateDBDataByID(id, formParams)
		c.JSON(status, retInfo)
	} else {
		c.JSON(404, gin.H{

		})
		fmt.Println("Invalid Res Name:", resourceName)
	}
	
}
//删除指定ID的资源
func deleteResourceByID(c *gin.Context){
	//资源名
	resourceName := c.Param("resourceName")
	//唯一ID
	idStr := c.Param("id")
	id := StringToUint(idStr)
	fmt.Println("deleteResourceByID", resourceName, id)
	//根据资源名获取资源对象
	resObject := GetDBResourceObjectByName(resourceName)
	if resObject != nil {
		status, retInfo,_ :=  resObject.DeleteDBDataByID(id)
		c.JSON(status, retInfo)
	} else {
		c.JSON(404, gin.H{

		})
		fmt.Println("Invalid Res Name:", resourceName)
	}
}
//添加资源
func addResource(c *gin.Context){
	//资源名
	resourceName := c.Param("resourceName")
	fmt.Println("addResource", resourceName)
	//根据资源名获取资源对象
	resObject := GetDBResourceObjectByName(resourceName)
	if resObject != nil {
		attrNameList := resObject.GetAttrNameList()
		formParams := make(map[string]string)
		for _, attrName := range attrNameList {
			value := c.PostForm(attrName)
			fmt.Println("attr name:",attrName, value)
			formParams[attrName] = value
		}
		//修改对象属性
		resObject.ModifyAttrWithParams(formParams)
		//操作数据库	
		status,retInfo,_ :=  resObject.CreateDBData()
		c.JSON(status, retInfo)
	} else {
		c.JSON(404, gin.H{

		})
		fmt.Println("Invalid Res Name:", resourceName)
	}
}
func main(){
	//使用默认路由
	router := gin.Default()

	//v1版本
	v1 := router.Group("/v1")
	{
		//API for GET 
		//v1.GET("/:ResourceName", getResourceList)
		v1.GET("/:resourceName/:id", getResourceByID)
		//API for POST
		v1.POST("/:resourceName", addResource)
		//API for PUT
		v1.PUT("/:resourceName/:id", updateResourceByID)
		//API for DELETE
		v1.DELETE("/:resourceName/:id", deleteResourceByID)
	}

	router.Run(":8081")
}
