package main

import (
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func main() {
	r := gin.Default()
	config := cors.DefaultConfig()
	config.AllowCredentials = true
	config.AllowOrigins = []string{"http://localhost:8080", "http://suncheng.xyz:7777"}
	r.Use(cors.New(config))
	r.StaticFS("/static", http.Dir("./static"))
	r.POST("/uploadfile", uploadFile)
	r.Run(":81")
}

func uploadFile(c *gin.Context) {
	id := c.PostForm("id")
	fileName := c.PostForm("fileName")
	file, err := c.FormFile("files")
	if err != nil {
		c.JSON(200, gin.H{
			"result": false,
			"msg":    "c.FormFile",
		})
		return
	}
	src, err := file.Open()
	if err != nil {
		c.JSON(200, gin.H{
			"result": false,
			"msg":    "file.Open",
		})
		return
	}
	defer src.Close()
	path := "static"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, os.ModePerm)
	}
	if len(id) > 0 {
		path += "/" + id
		if _, err := os.Stat(path); os.IsNotExist(err) {
			os.Mkdir(path, os.ModePerm)
		}
	}
	path += "/" + fileName + "." + strings.Split(file.Filename, ".")[1]
	dst, err := os.Create(path)
	if err != nil {
		c.JSON(200, gin.H{
			"result": false,
			"msg":    "os.Create",
		})
		return
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)

	c.JSON(200, gin.H{
		"result": err == nil,
		"msg":    "io.Copy",
		"path":   "http://118.24.27.231:81/" + path,
	})
}

func bindExtend(c *gin.Context, obj interface{}) error {
	err := c.ShouldBindWith(obj, binding.JSON)
	return err
}

type fileRequset struct {
	Path     string `json:"path"`     // 路径
	FileName string `json:"fileName"` // 文件名
}
