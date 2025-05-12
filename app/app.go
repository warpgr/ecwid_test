package app

import (
	"context"
	"os"

	"github.com/warpgr/ecwid_test/configs"
	"github.com/warpgr/ecwid_test/parser"
	"github.com/warpgr/ecwid_test/reader"
	"github.com/warpgr/ecwid_test/store"

	"github.com/sirupsen/logrus"
)

func NewIPParserApp(cfg configs.AppConfig) IPParserApp {
	return &ipParserApp{
		cfg: cfg,
	}
}

type ipParserApp struct {
	cfg      configs.AppConfig
	cancel   context.CancelFunc
	doneChan chan struct{}
	file     *os.File

	ipReader reader.IPReaderEngine
	ipStore  store.IPStore
}

func (a *ipParserApp) Init() (err error) {
	logrus.Info("Initializing app.")
	// Preparing file.
	a.file, err = os.OpenFile(a.cfg.FilePath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		logrus.WithError(err).Fatal("Can't open file.")
	}

	var (
		workerCount int
	)

	// Creating processor instances depending on provided strategy.
	switch a.cfg.Strategy {
	case configs.StrategyConcurrent:
		a.ipStore = store.NewShardedAtomicStore()
		workerCount = a.cfg.WorkerCount
		a.ipReader, err = reader.NewSectionedIPReader(a.file, workerCount, a.cfg.ChunkSize, parser.NewIPParser(a.ipStore))
		if err != nil {
			return err
		}

	case configs.StrategyLinear:
		fallthrough

	default:
		a.ipStore = store.NewIPStore()
		workerCount = 1
		a.ipReader, err = reader.NewSectionedIPReader(a.file, workerCount, a.cfg.ChunkSize, parser.NewIPParser(a.ipStore))
		if err != nil {
			return err
		}
	}

	if err != nil {
		if err := a.file.Close(); err != nil {
			logrus.WithError(err).Error("Error occurs on file closing.")
		}
		return err
	}

	return nil
}

func (a *ipParserApp) Run(ctx context.Context) <-chan struct{} {
	logrus.Info("Running app.")
	child, cancel := context.WithCancel(ctx)
	a.cancel = cancel

	a.doneChan = make(chan struct{}, 1)
	go func() {
		defer close(a.doneChan)
		a.ipReader.Run(child)
	}()

	logrus.WithFields(logrus.Fields{
		"strategy":  a.cfg.Strategy,
		"chunkSize": a.cfg.ChunkSize,
		"filePath":  a.cfg.FilePath,
	}).Info("Application flow started.")

	return a.doneChan
}

func (a *ipParserApp) Stop() error {
	logrus.Info("Stopping app.")
	a.cancel()

	<-a.doneChan
	// Closing resources.
	_ = a.file.Close()

	logrus.WithFields(logrus.Fields{
		"path":      a.cfg.FilePath,
		"uniqueIPs": a.ipStore.UniqueCount(),
	}).Info("Processed.")

	return nil
}
