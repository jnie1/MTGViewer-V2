//go:build prod

package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	RegisterRouter()
}
