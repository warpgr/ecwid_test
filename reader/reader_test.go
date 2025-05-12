package reader_test

import (
	"context"
	"os"
	"path"
	"testing"
	"time"

	"github.com/warpgr/ecwid_test/parser"
	"github.com/warpgr/ecwid_test/reader"
	"github.com/warpgr/ecwid_test/store/mocks"

	"github.com/stretchr/testify/suite"
)

type IPReaderEngineTestSuite struct {
	suite.Suite

	testingFilesDir string

	st *mocks.IPStoreMock
	ps parser.IPParser
}

const (
	fileName          = "ip_addresses"
	corruptedFileName = "corrupted"
	count             = 100
	workers           = 4
	chunkSize         = 15
)

func TestRunIPReaderEngineTestSuite(t *testing.T) {
	dir, ok := os.LookupEnv("TESTING_FILES_DIR")
	if !ok {
		return
	}

	suite.Run(t, &IPReaderEngineTestSuite{testingFilesDir: dir})
}

func (s *IPReaderEngineTestSuite) TestIPAddressesParsing() {
	path := path.Join(s.testingFilesDir, fileName)
	file, err := os.OpenFile(path, os.O_RDONLY, os.ModePerm)
	s.Require().NoError(err)
	defer func() {
		err := file.Close()
		s.Require().NoError(err)
	}()

	eg, err := reader.NewSectionedIPReader(file, workers, chunkSize, s.ps)
	s.Require().NoError(err)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	eg.Run(ctx)

	s.Equal(count, len(s.st.IPs))
}

func (s *IPReaderEngineTestSuite) TestIPCorruptedAddressesParsing() {
	path := path.Join(s.testingFilesDir, corruptedFileName)
	file, err := os.OpenFile(path, os.O_RDONLY, os.ModePerm)
	s.Require().NoError(err)
	defer func() {
		err := file.Close()
		s.Require().NoError(err)
	}()

	eg, err := reader.NewSectionedIPReader(file, workers, chunkSize, s.ps)
	s.Require().NoError(err)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	eg.Run(ctx)

	s.Equal(count, len(s.st.IPs))
}

func (s *IPReaderEngineTestSuite) SetupTest() {

	s.st = &mocks.IPStoreMock{}
	s.ps = parser.NewIPParser(s.st)

}

func (s *IPReaderEngineTestSuite) SetupSuite() {}

func (s *IPReaderEngineTestSuite) TearDownSuite() {}
