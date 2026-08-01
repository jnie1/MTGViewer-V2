package routes

import (
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/gin-gonic/gin"
)

func AddStaticRoutes(r *gin.Engine, knownRoutes ...string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	dist := path.Join(cwd, "dist")
	index := path.Join(dist, "index.html")
	assets := path.Join(dist, "assets")

	_, err = os.Stat(index)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}

	_, err = os.Stat(assets)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}

	r.StaticFile("/", index)
	r.Static("/assets", assets)

	r.GET("/index", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/")
	})
	r.GET("/index.html", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/")
	})

	r.NoRoute(func(c *gin.Context) {
		for _, rt := range knownRoutes {
			if strings.HasPrefix(c.Request.URL.Path, rt) {
				c.Status(http.StatusNotFound)
				return
			}
		}
		c.File(index)
	})

	return nil
}
