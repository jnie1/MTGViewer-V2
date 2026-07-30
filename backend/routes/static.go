package routes

import (
	"os"

	"github.com/gin-gonic/gin"
)

func AddStaticPaths(r gin.IRouter) error {
	routeMap := map[string]string{
		"/":           "./dist/index.html",
		"/index.html": "./dist/index.html",
		"/assets":     "./dist/assets",
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

	return nil
}
