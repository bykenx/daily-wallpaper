package main

import (
	"context"
	"daily-wallpaper/api"
	"daily-wallpaper/sources"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func handleFrontend(c *gin.Context) {

}

func handleGetAllSettings(c *gin.Context) {
	settings := readSettings()
	ginJsonResult(c, settings)
}

func handleModifySettings(c *gin.Context) {
	var settings Settings
	err := c.ShouldBind(&settings)
	if err != nil {
		ginJsonError(c, err.Error())
		return
	}
	writeSettings(settings)
	if settings.CurrentImage != nil && *settings.CurrentImage != "" {
		savedPath, err := downloadFileAndSave(*settings.CurrentImage)
		if err != nil {
			ginJsonError(c, err.Error())
			return
		}
		err = api.SetWallpaper(savedPath)
		if err != nil {
			ginJsonError(c, err.Error())
			return
		}
	}
	ginJsonResult(c, "数据修改成功")
}

func handleGetSources(c *gin.Context) {
	ginJsonResult(c, GetDescriptions())
}

func handleTodayImage(c *gin.Context) {
	name := readSettings().CurrentSource
	if name == nil || *name == "" {
		*name = "bing"
	}
	source := GetSource(*name)
	res, err := source.GetToday()
	if err != nil {
		ginJsonError(c, err.Error())
	}
	ginJsonResult(c, res)
}

func handleArchiveImages(c *gin.Context) {
	name := readSettings().CurrentSource
	if name == nil || *name == "" {
		*name = "bing"
	}
	source := GetSource(*name)
	var param sources.ArchiveParam
	_ = c.ShouldBind(&param)
	res, err := source.GetArchive(param)
	if err != nil {
		ginJsonError(c, err.Error())
	}
	ginJsonResult(c, res)
}

func handleDownload(c *gin.Context) {
	link := c.Query("src")
	if link == "" {
		ginJsonError(c, "请指定链接地址")
		return
	}
	savedPath, err := downloadFileAndSave(link)
	if err != nil {
		ginJsonError(c, err.Error())
		return
	}
	ginJsonResult(c, savedPath)
}

func startServer() chan struct{} {
	router := gin.Default()
	{
		router.GET("", handleFrontend)
		router.GET("api/settings", handleGetAllSettings)
		router.PUT("api/settings", handleModifySettings)
		router.GET("api/image/sources", handleGetSources)
		router.GET("api/image/archive", handleArchiveImages)
		router.GET("api/image/today", handleTodayImage)
		router.POST("api/download", handleDownload)
	}
	server := &http.Server{
		Addr:    ":9001",
		Handler: router,
	}
	go func() {
		openDB()
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
		closeDB()
	}()
	closeChan := make(chan struct{})
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	stopServer := func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err := server.Shutdown(ctx)
		if err != nil {
			log.Fatal("Server forced to shutdown: ", err)
		}
		log.Println("Server stopped.")
	}
	go func() {
		select {
		case <-closeChan:
			stopServer()
		case <-quit:
			stopServer()
		}
	}()
	return closeChan
}
