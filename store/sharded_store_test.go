package store_test

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/warpgr/ecwid_test/store"
)

const (
	lowerBound = 10000
	upperBound = 100000
	iteration  = 100
)

type ShardedStoreTestSuite struct {
	suite.Suite
}

func TestRunShardedStoreTestSuite(t *testing.T) {
	suite.Run(t, &ShardedStoreTestSuite{})
}

func (s *ShardedStoreTestSuite) TestStore() {
	st := store.NewShardedAtomicStore()
	storeIPs(st)
	count := st.UniqueCount()
	s.Equal(uint32(upperBound-lowerBound), count)
}

func (s *ShardedStoreTestSuite) SetupSuite() {}

func (s *ShardedStoreTestSuite) TearDownSuite() {}

func storeIPs(st store.IPStore) {
	for ip := uint32(lowerBound); ip < upperBound; ip++ {
		for i := 0; i < iteration; i++ {
			st.StoreIP(ip)
		}
	}
}
