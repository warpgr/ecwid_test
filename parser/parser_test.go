package parser_test

import (
	"os"
	"path"
	"testing"

	"github.com/warpgr/ecwid_test/store/mocks"

	"github.com/warpgr/ecwid_test/parser"

	"github.com/stretchr/testify/suite"
)

type IPParserTestSuite struct {
	suite.Suite
	testingFilesDir string
}

func TestRunIPParserTestSuite(t *testing.T) {
	dir, ok := os.LookupEnv("TESTING_FILES_DIR")
	if !ok {
		return
	}

	suite.Run(t, &IPParserTestSuite{testingFilesDir: dir})
}

func (s *IPParserTestSuite) TestExtractIPs() {
	const (
		fileName = "ip_addresses"
		count    = 100
	)

	path := path.Join(s.testingFilesDir, fileName)

	buff, err := os.ReadFile(path)
	s.Require().NoError(err)

	st := &mocks.IPStoreMock{}
	ps := parser.NewIPParser(st)

	leftover := ps.ExtractIPs(buff)
	s.Empty(leftover)
	s.Equal(count, len(st.IPs))
}

func (s *IPParserTestSuite) TestExtractIPsWithStepHandling() {
	const (
		fileName         = "ip_addresses"
		count            = 100
		partialStepCount = 50
	)

	path := path.Join(s.testingFilesDir, fileName)

	buff, err := os.ReadFile(path)
	s.Require().NoError(err)

	st := &mocks.IPStoreMock{}
	ps := parser.NewIPParser(st)

	md := len(buff) / 2
	subBuff := buff[:md]

	leftover := ps.ExtractIPs(subBuff)
	s.NotEmpty(leftover)
	s.Equal(partialStepCount, len(st.IPs))

	subBuff = buff[md:]
	subBuff = append(leftover, subBuff...)
	leftover = ps.ExtractIPs(subBuff)
	s.Empty(leftover)
	s.Equal(count, len(st.IPs))
}

func (s *IPParserTestSuite) TestExtractIPsCorruptedFile() {
	const (
		fileName = "ip_addresses"
		count    = 100
	)

	path := path.Join(s.testingFilesDir, fileName)

	buff, err := os.ReadFile(path)
	s.Require().NoError(err)

	st := &mocks.IPStoreMock{}
	ps := parser.NewIPParser(st)

	leftover := ps.ExtractIPs(buff)
	s.Empty(leftover)

	s.Equal(count, len(st.IPs))
}

func (s *IPParserTestSuite) SetupSuite() {}

func (s *IPParserTestSuite) TearDownSuite() {}
