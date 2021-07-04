package agin

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func getModelsList(c *gin.Context) {
	var modelList []string
	for k, _ := range _register.Models {
		modelList = append(modelList, k)
	}
	c.JSON(200, gin.H{
		"modelList": modelList,
	})
}

func getRecords(c *gin.Context) {

	table := getAdminModel(c)
	data, err := table.GetRecord(c)
	if err != nil {
		c.String(400, "err get")
		return
	}

	c.JSON(200, data)
}

func updateRecord(c *gin.Context) {
	table := getAdminModel(c)
	table.Updates(c)
}

func createRecord(c *gin.Context) {
	table := getAdminModel(c)
	table.Create(c)
}

func deleteRecord(c *gin.Context) {
	table := getAdminModel(c)
	table.Delete(c)
}

func getAdminModel(c *gin.Context) IAdminModel {
	tableName := c.Param("table")
	table := _register.Models[tableName]
	return table
}

func uploadFile(c *gin.Context) {
	file, _ := c.FormFile("file")
	log.Println("file.Filename:", file.Filename)

	//// Upload the file to specific dst.
	err := c.SaveUploadedFile(file, file.Filename)
	if err != nil {
		c.String(500, "get img err")
		return
	}

	url, err := qiNiuUpload.Upload(file.Filename)
	if err != nil {
		fmt.Println("url: ", url, "err: ", err)
		c.String(400, "upload to cloud err")
		return
	}

	c.String(http.StatusOK, url)
}
