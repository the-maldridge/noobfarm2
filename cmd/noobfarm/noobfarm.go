package main

import (
	"os"

	"github.com/hashicorp/go-hclog"

	"github.com/the-maldridge/noobfarm2/internal/web"
)

func main() {
	llevel := os.Getenv("NF_LOGLEVEL")
	if llevel == "" {
		llevel = "INFO"
	}

	appLogger := hclog.New(&hclog.LoggerOptions{
		Name:  "noobfarm2",
		Level: hclog.LevelFromString(llevel),
	})

	w := web.New(appLogger, nil)
	w.Serve(os.Getenv("NF_BIND"))
}
