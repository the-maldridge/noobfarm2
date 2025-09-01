package main

import (
	"os"

	"github.com/hashicorp/go-hclog"

	"github.com/the-maldridge/noobfarm2/internal/qdb"
	_ "github.com/the-maldridge/noobfarm2/internal/qdb/json"
	"github.com/the-maldridge/noobfarm2/internal/web"

	_ "github.com/the-maldridge/authware/backend/htpasswd"
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
	qdb.SetParentLogger(appLogger)
	qdb.DoCallbacks()

	db, err := qdb.New(os.Getenv("NF_QDB"))
	if err != nil {
		appLogger.Error("Could not initialize quote source", "error", err)
		os.Exit(1)
	}

	w, err := web.New(appLogger, db)
	if err != nil {
		appLogger.Error("Fatal error initializing server", "error", err)
		os.Exit(1)
	}
	if err := w.Serve(os.Getenv("NF_BIND")); err != nil {
		appLogger.Error("Fatal error starting webserver", "error", err)
	}
}
