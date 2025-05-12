package store_test

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/warpgr/ecwid_test/store"
)

type SipleStoreTestSuite struct {
	suite.Suite
}

func TestRunSipleStoreTestSuite(t *testing.T) {
	suite.Run(t, &SipleStoreTestSuite{})
}

func (s *SipleStoreTestSuite) TestStoreIPP() {
	st := store.NewIPStore()
	storeIPs(st)
	count := st.UniqueCount()
	s.Equal(uint32(upperBound-lowerBound), count)
}

func (s *SipleStoreTestSuite) SetupSuite() {}

func (s *SipleStoreTestSuite) TearDownSuite() {}
