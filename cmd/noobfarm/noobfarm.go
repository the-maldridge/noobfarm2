package main

import (
	"os"

	"github.com/hashicorp/go-hclog"

	"github.com/the-maldridge/noobfarm2/internal/qdb"
	_ "github.com/the-maldridge/noobfarm2/internal/qdb/json"
	"github.com/the-maldridge/noobfarm2/internal/web"
	"github.com/the-maldridge/noobfarm2/internal/web/auth"
	_ "github.com/the-maldridge/noobfarm2/internal/web/auth/file"
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
	auth.SetParentLogger(appLogger)
	qdb.DoCallbacks()
	auth.DoCallbacks()

	db, err := qdb.New(os.Getenv("NF_QDB"))
	if err != nil {
		appLogger.Error("Could not initialize quote source", "error", err)
		os.Exit(1)
	}

	auth, err := auth.Initialize(os.Getenv("NF_AUTH"))
	if err != nil {
		appLogger.Error("Could not initialize authenticator", "error", err)
		os.Exit(1)
	}

	w := web.New(appLogger, db, auth)
	if err := w.Serve(os.Getenv("NF_BIND")); err != nil {
		appLogger.Error("Fatal error starting webserver", "error", err)
	}
}
