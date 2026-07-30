//go:build prod

package main

import "github.com/jnie1/MTGViewer-V2/config"

func main() {
	cfg := config.Load()
	RegisterRouter(cfg)
}
