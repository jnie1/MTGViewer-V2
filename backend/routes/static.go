package routes

import (
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/gin-gonic/gin"
)

func AddStaticPaths(r *gin.Engine) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	indexPath := path.Join(cwd, "dist", "index.html")
	assetsPath := path.Join(cwd, "dist", "assets")

	_, err = os.Stat(indexPath)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}

	_, err = os.Stat(assetsPath)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}

	r.Static("/", indexPath)
	r.StaticFile("/assets", assetsPath)

	r.GET("/index", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/")
	})
	r.GET("/index.html", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/")
	})

	r.NoRoute(func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/api/") {
			c.Status(http.StatusNotFound)
			return
		}
		c.File(indexPath)
	})

	return nil
}
