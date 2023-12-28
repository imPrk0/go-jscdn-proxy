package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

var (
	PROXY_URL string
)

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.Use(corsMiddleware)
	r.GET("/*uri", handleProxy)
	return r
}

func handleProxy(c *gin.Context) {
	uri := c.Param("uri")
	targetURL := PROXY_URL + uri

	fmt.Println(targetURL)

	resp, err := http.Get(targetURL)
	if nil != err {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"message": "=-=",
			},
		)
		return
	}
	defer resp.Body.Close()

	//if 200 <= resp.StatusCode && 304 >= resp.StatusCode {
	//	fileName := fmt.Sprintf("cache_%s", uri)
	//	file, err := os.Create(fileName)
	//	if nil != err {
	//		c.JSON(
	//			http.StatusInternalServerError,
	//			gin.H{
	//				"message": "fail cache",
	//			},
	//		)
	//		return
	//	}
	//	defer file.Close()
	//
	//	_, err = io.Copy(file, resp.Body)
	//	if nil != err {
	//		c.JSON(
	//			http.StatusInternalServerError,
	//			gin.H{
	//				"message": "fail cache 2",
	//			},
	//		)
	//		return
	//	}
	//
	//	expTime := time.Now().Add(24 * time.Hour)
	//	err = os.Chtimes(fileName, expTime, expTime)
	//	if nil != err {
	//		c.JSON(
	//			http.StatusInternalServerError,
	//			gin.H{
	//				"message": "fail cache 3",
	//			},
	//		)
	//		return
	//	}
	//}

	c.DataFromReader(
		http.StatusOK,
		resp.ContentLength,
		resp.Header.Get("Content-Type"),
		resp.Body,
		nil,
	)
}

func corsMiddleware(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

	if "OPTIONS" == c.Request.Method {
		c.AbortWithStatus(http.StatusOK)
		return
	}

	c.Next()
}

func init() {
	proxyURL := os.Getenv("PROXY_URL")
	if "" == proxyURL {
		PROXY_URL = "https://cdnjs.cloudflare.com/ajax/libs"
	} else {
		PROXY_URL = proxyURL
	}
	//gin.SetMode(gin.ReleaseMode)
}

func main() {
	r := setupRouter()
	r.Run(
		os.Getenv("SERVER_ADDRESS"),
	)
}
