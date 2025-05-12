package configs

import (
	"errors"
	"flag"
	"fmt"
	"runtime"

	"github.com/sirupsen/logrus"
)

func GetConfig() (*AppConfig, error) {
	const defaultChunkSize = 10

	var cfg AppConfig

	flag.StringVar(&cfg.FilePath, "file-path", "",
		"Path to the input file (required)")
	flag.StringVar(&cfg.Strategy, "strategy", StrategyLinear,
		fmt.Sprintf("Processing strategy: %q or %q", StrategyLinear, StrategyConcurrent))
	flag.IntVar(&cfg.ChunkSize, "chunk-size", defaultChunkSize,
		fmt.Sprintf("Chunk size in MiB (default = %d MiB)", defaultChunkSize))
	// There is no need to increase goroutines count than CPUs to avoid unnecessary context switching.
	flag.IntVar(&cfg.WorkerCount, "worker-count", runtime.NumCPU()-2,
		"Number of parallel workers (default = number of CPU cores)")
	flag.StringVar(&cfg.LogLevel, "log-level", "info",
		"Log level (debug, info, warn, error)")
	flag.BoolVar(&cfg.EnableProf, "enable-prof", false,
		"Set true if need to start pprof server.")

	flag.Parse()

	if cfg.FilePath == "" {
		return nil, errors.New("--file-path is required")
	}

	switch cfg.Strategy {
	case StrategyLinear, StrategyConcurrent:
	default:
		return nil, fmt.Errorf("invalid --strategy %q; must be %q or %q",
			cfg.Strategy, StrategyLinear, StrategyConcurrent)
	}

	logLevel, err := logrus.ParseLevel(cfg.LogLevel)
	if err != nil {
		return nil, err
	}
	logrus.SetLevel(logLevel)

	cfg.ChunkSize <<= 20

	return &cfg, nil
}
