package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"time"

	"github.com/warpgr/ecwid_test/app"
	"github.com/warpgr/ecwid_test/configs"

	_ "net/http/pprof"

	"github.com/sirupsen/logrus"
)

func main() {
	start := time.Now()
	defer dumpSysUsage(start)

	// Initializing configs.
	cfg, err := configs.GetConfig()
	if err != nil {
		logrus.WithError(err).Fatal("Broken configs.")
	}
	ctx := context.Background()

	a := app.NewIPParserApp(*cfg)
	if err := a.Init(); err != nil {
		logrus.WithError(err).Fatal("Error occurs in app init.")
	}

	doneChan := a.Run(ctx)

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt)

		sig := <-sigChan
		logrus.WithField("signal", sig).Warn("Signal received.")

		if err := a.Stop(); err != nil {
			logrus.WithError(err).Error("Error occurs on app interruption.")
		}
	}()

	if cfg.EnableProf {
		go func() {
			const addr = ":6060"

			logrus.
				WithField("endpoint", fmt.Sprintf("http://localhost%s/debug/pprof/", addr)).
				Info("Starting pprof server.")

			if err := http.ListenAndServe(addr, nil); err != nil {
				logrus.WithError(err).Fatal("Can't start pprof server.")
			}
		}()
	}

	<-doneChan

	if err := a.Stop(); err != nil {
		logrus.WithError(err).Error("Error occurs on app shutdown.")
	}
}

func dumpSysUsage(start time.Time) {
	elapsed := time.Since(start)

	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	logrus.WithFields(logrus.Fields{
		"elapsed":        elapsed,
		"total_alloc_mb": m.TotalAlloc >> 20,
		"heap_in_use_mb": m.HeapInuse >> 20,
	}).Info("Sys usage.")
}
