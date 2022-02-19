package server

import (
	"context"
	"daily-wallpaper/api"
	settings2 "daily-wallpaper/settings"
	"daily-wallpaper/sources"
	"daily-wallpaper/utils"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
)

func handleGetAllSettings(c *gin.Context) {
	settings := settings2.ReadSettings()
	utils.GinJsonResult(c, settings)
}

func handleModifySettings(c *gin.Context) {
	var settings settings2.Settings
	err := c.ShouldBind(&settings)
	if err != nil {
		utils.GinJsonError(c, err.Error())
		return
	}
	settings2.WriteSettings(settings)
	utils.GinJsonResult(c, "数据修改成功")
}

func handleGetSources(c *gin.Context) {
	utils.GinJsonResult(c, utils.GetDescriptions())
}

func handleTodayImage(c *gin.Context) {
	name := settings2.ReadSettings().CurrentSource
	if name == nil || *name == "" {
		*name = "bing"
	}
	source := utils.GetSource(*name)
	res, err := source.GetToday()
	if err != nil {
		utils.GinJsonError(c, err.Error())
		return
	}
	utils.GinJsonResult(c, res)
}

func handleArchiveImages(c *gin.Context) {
	name := settings2.ReadSettings().CurrentSource
	if name == nil || *name == "" {
		*name = "bing"
	}
	source := utils.GetSource(*name)
	var param sources.ArchiveParam
	_ = c.ShouldBind(&param)
	res, err := source.GetArchive(param)
	if err != nil {
		utils.GinJsonError(c, err.Error())
		return
	}
	utils.GinJsonResult(c, res)
}

func handleDownload(c *gin.Context) {
	link := c.Query("src")
	if link == "" {
		utils.GinJsonError(c, "请指定链接地址")
		return
	}
	savedPath, err := api.DownloadFileAndSave(link)
	if err != nil {
		utils.GinJsonError(c, err.Error())
		return
	}
	utils.GinJsonResult(c, savedPath)
}

func serveRoot(urlPrefix, root string) gin.HandlerFunc {
	// https://github.com/golang/go/issues/32350
	builtinMimeTypesLower := map[string]string{
		".css":  "text/css; charset=utf-8",
		".gif":  "image/gif",
		".htm":  "text/html; charset=utf-8",
		".html": "text/html; charset=utf-8",
		".jpg":  "image/jpeg",
		".js":   "application/javascript",
		".wasm": "application/wasm",
		".pdf":  "application/pdf",
		".png":  "image/png",
		".svg":  "image/svg+xml",
		".xml":  "text/xml; charset=utf-8",
	}
	fs := static.LocalFile(root, true)
	fileserver := http.FileServer(fs)
	if urlPrefix != "" {
		fileserver = http.StripPrefix(urlPrefix, fileserver)
	}
	return func(c *gin.Context) {
		if fs.Exists(urlPrefix, c.Request.URL.Path) {
			if v, ok := builtinMimeTypesLower[filepath.Ext(c.Request.URL.Path)]; ok {
				c.Writer.Header().Set("Content-Type", v)
			}
			fileserver.ServeHTTP(c.Writer, c.Request)
			c.Abort()
		}
	}
}

func StartServer() {
	router := gin.Default()
	router.Use(serveRoot("/", utils.GetStaticPath()))
	{
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
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
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
		<-quit
		stopServer()
	}()
}
