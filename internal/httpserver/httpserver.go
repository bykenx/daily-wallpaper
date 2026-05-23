package httpserver

import (
	"context"
	st "daily-wallpaper/internal/settings"
	"daily-wallpaper/internal/source"
	"daily-wallpaper/internal/source/registry"
	"daily-wallpaper/internal/util"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
)

func jsonError(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusBadRequest,
		"msg":  msg,
		"data": nil,
	})
}

func jsonResult(c *gin.Context, obj interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"data": obj,
	})
}

func handleGetAllSettings(c *gin.Context) {
	settings := st.ReadSettings()
	jsonResult(c, settings)
}

func handleModifySettings(c *gin.Context) {
	var settings st.Settings
	err := c.ShouldBind(&settings)
	if err != nil {
		jsonError(c, err.Error())
		return
	}
	st.WriteSettings(settings)
	jsonResult(c, "数据修改成功")
}

func handleGetSources(c *gin.Context) {
	jsonResult(c, registry.Descriptions())
}

func handleTodayImage(c *gin.Context) {
	name := st.ReadSettings().CurrentSource
	if name == nil || *name == "" {
		*name = "bing"
	}
	provider := registry.Get(*name)
	res, err := provider.GetToday()
	if err != nil {
		jsonError(c, err.Error())
		return
	}
	jsonResult(c, res)
}

func handleArchiveImages(c *gin.Context) {
	name := st.ReadSettings().CurrentSource
	if name == nil || *name == "" {
		*name = "bing"
	}
	provider := registry.Get(*name)
	var param source.ArchiveParam
	_ = c.ShouldBind(&param)
	res, err := provider.GetArchive(param)
	if err != nil {
		jsonError(c, err.Error())
		return
	}
	jsonResult(c, res)
}

func handleGetImage(c *gin.Context) {
	link := c.Query("link")
	if dir, err := util.GetOrDownload(link); err == nil {
		if content, err := os.ReadFile(dir); err == nil {
			c.Render(http.StatusOK, render.Data{
				ContentType: http.DetectContentType(content),
				Data:        content,
			})
			return
		}
	}
	c.String(http.StatusNotFound, http.StatusText(http.StatusNotFound))
}

func handleGetImageList(c *gin.Context) {
	var start int
	var limit int
	if num, err := strconv.Atoi(c.Query("start")); err == nil {
		start = num
	}
	if num, err := strconv.Atoi(c.Query("limit")); err == nil {
		limit = num
	}
	jsonResult(c, util.GetImageListPagination(start, limit))
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

func serveSpaFallback(root string) gin.HandlerFunc {
	indexFile := filepath.Join(root, "index.html")
	return func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/api/") {
			c.String(http.StatusNotFound, http.StatusText(http.StatusNotFound))
			return
		}
		c.File(indexFile)
	}
}

func StartServer() {
	router := gin.Default()
	router.Use(serveRoot("/", util.GetStaticPath()))
	{
		router.GET("api/settings", handleGetAllSettings)
		router.PUT("api/settings", handleModifySettings)
		router.GET("api/image/sources", handleGetSources)
		router.GET("api/image/archive", handleArchiveImages)
		router.GET("api/image/today", handleTodayImage)
		router.GET("api/image/get", handleGetImage)
		router.GET("api/image/list", handleGetImageList)
	}
	router.NoRoute(serveSpaFallback(util.GetStaticPath()))
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
