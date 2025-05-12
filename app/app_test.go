package app_test

import (
	"context"
	"os"
	"path"
	"testing"
	"time"

	"github.com/warpgr/ecwid_test/app"
	"github.com/warpgr/ecwid_test/configs"

	"github.com/stretchr/testify/suite"
)

const (
	fileName = "ip_addresses"
)

type IPParserAppTestSuite struct {
	suite.Suite

	testingFilesDir string
}

func TestRunIPParserAppTestSuite(t *testing.T) {
	dir, ok := os.LookupEnv("TESTING_FILES_DIR")
	if !ok {
		return
	}

	suite.Run(t, &IPParserAppTestSuite{testingFilesDir: dir})
}

func (s *IPParserAppTestSuite) TestLinearStrategy() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cfg := configs.AppConfig{
		FilePath:    path.Join(s.testingFilesDir, fileName),
		Strategy:    configs.StrategyLinear,
		ChunkSize:   20,
		WorkerCount: 10,
	}

	a := app.NewIPParserApp(cfg)
	err := a.Init()
	s.Require().NoError(err)

	doneChan := a.Run(ctx)

	select {
	case <-ctx.Done():
		// Unreachable.
		s.Require().True(false, "Context deadline exceed.Application did not complete job.")
	case <-doneChan:
		err = a.Stop()
		s.Require().NoError(err)
	}
}

func (s *IPParserAppTestSuite) TestConcurrentStrategy() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cfg := configs.AppConfig{
		FilePath:    path.Join(s.testingFilesDir, fileName),
		Strategy:    configs.StrategyConcurrent,
		ChunkSize:   20,
		WorkerCount: 10,
	}

	a := app.NewIPParserApp(cfg)
	err := a.Init()
	s.Require().NoError(err)

	doneChan := a.Run(ctx)

	select {
	case <-ctx.Done():
		// Unreachable.
		s.Require().True(false, "Context deadline exceed.Application did not complete job.")
	case <-doneChan:
		err = a.Stop()
		s.Require().NoError(err)
	}
}

func (s *IPParserAppTestSuite) SetupSuite() {}

func (s *IPParserAppTestSuite) TearDownSuite() {}
