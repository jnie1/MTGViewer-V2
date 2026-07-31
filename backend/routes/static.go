package routes

import (
	"net/http"
	"os"
	"path"

	"github.com/gin-gonic/gin"
)

func AddStaticPaths(r gin.IRouter) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	routeMap := map[string]string{
		"/":       path.Join(cwd, "dist/index.html"),
		"/assets": path.Join(cwd, "dist/assets"),
	}

	for route, filepath := range routeMap {
		stat, err := os.Stat(filepath)
		if os.IsNotExist(err) {
			continue
		}
		if err != nil {
			return err
		}
		if stat.IsDir() {
			r.Static(route, filepath)
		} else {
			r.StaticFile(route, filepath)
		}
	}

	r.GET("/index", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/")
	})

	r.GET("/index.html", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/")
	})

	return nil
}
