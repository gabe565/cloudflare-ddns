package main

import (
	"log/slog"
	"os"

	"gabe565.com/cloudflare-ddns/cmd"
	"gabe565.com/cloudflare-ddns/internal/config"
	"gabe565.com/utils/slogx"
)

func main() {
	config.InitLog(os.Stderr, slogx.LevelInfo, slogx.FormatAuto)

	root := cmd.New()
	if err := root.Execute(); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}
