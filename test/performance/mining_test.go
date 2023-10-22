package performance

import (
	"os"
	"testing"

	"github.com/anonimitycash/anonimitycash-classic/account"
	dbm "github.com/anonimitycash/anonimitycash-classic/database/leveldb"
	"github.com/anonimitycash/anonimitycash-classic/mining"
	"github.com/anonimitycash/anonimitycash-classic/test"
)

// Function NewBlockTemplate's benchmark - 0.05s
func BenchmarkNewBlockTpl(b *testing.B) {
	testDB := dbm.NewDB("testdb", "leveldb", "temp")
	defer os.RemoveAll("temp")

	chain, _, txPool, err := test.MockChain(testDB)
	if err != nil {
		b.Fatal(err)
	}
	accountManager := account.NewManager(testDB, chain)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mining.NewBlockTemplate(chain, txPool, accountManager)
	}
}
