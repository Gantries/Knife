package perf

import (
	"fmt"
	"net/http"
	_ "net/http/pprof" // #nosec G108 - pprof is intentionally imported for profiling
	"runtime"
	"time"

	"github.com/gantries/knife/pkg/easy"
	"github.com/gantries/knife/pkg/lang"
	"github.com/gantries/knife/pkg/log"
)

var defaultPprofPort = 6060

func SetupPerformanceProfile(port *int) {
	logger := log.New("knife/perf")

	runtime.SetBlockProfileRate(1)
	runtime.SetMutexProfileFraction(5)

	serve := fmt.Sprintf("0.0.0.0:%d", *lang.Ternary(port == nil, &defaultPprofPort, port))
	go func() {
		server := &http.Server{
			Addr:              serve,
			ReadHeaderTimeout: 10 * time.Second,
			ReadTimeout:       30 * time.Second,
			WriteTimeout:      30 * time.Second,
			IdleTimeout:       5 * time.Minute,
		}
		if err := server.ListenAndServe(); err != nil {
			easy.Panic(err, nil)
		}
	}()
	logger.Info("Enabled pprof, you may want to access through /debug/pprof/", "serve", serve)
}
