package test

import (
	"context"
	"time"

	"github.com/anonimitycash/anonimitycash-classic/account"
	"github.com/anonimitycash/anonimitycash-classic/blockchain/pseudohsm"
	"github.com/anonimitycash/anonimitycash-classic/blockchain/txbuilder"
	"github.com/anonimitycash/anonimitycash-classic/consensus"
	"github.com/anonimitycash/anonimitycash-classic/crypto/ed25519/chainkd"
	"github.com/anonimitycash/anonimitycash-classic/database"
	dbm "github.com/anonimitycash/anonimitycash-classic/database/leveldb"
	"github.com/anonimitycash/anonimitycash-classic/event"
	"github.com/anonimitycash/anonimitycash-classic/protocol"
	"github.com/anonimitycash/anonimitycash-classic/protocol/bc"
	"github.com/anonimitycash/anonimitycash-classic/protocol/bc/types"
	"github.com/anonimitycash/anonimitycash-classic/protocol/vm"
)

const (
	vmVersion    = 1
	assetVersion = 1
)

// MockChain mock chain with genesis block
func MockChain(testDB dbm.DB) (*protocol.Chain, *database.Store, *protocol.TxPool, error) {
	store := database.NewStore(testDB)
	dispatcher := event.NewDispatcher()
	txPool := protocol.NewTxPool(store, dispatcher)
	chain, err := protocol.NewChain(store, txPool)
	return chain, store, txPool, err
}

// MockUTXO mock a utxo
func MockUTXO(controlProg *account.CtrlProgram) *account.UTXO {
	utxo := &account.UTXO{}
	utxo.OutputID = bc.Hash{V0: 1}
	utxo.SourceID = bc.Hash{V0: 2}
	utxo.AssetID = *consensus.MITYAssetID
	utxo.Amount = 1000000000
	utxo.SourcePos = 0
	utxo.ControlProgram = controlProg.ControlProgram
	utxo.AccountID = controlProg.AccountID
	utxo.Address = controlProg.Address
	utxo.ControlProgramIndex = controlProg.KeyIndex
	utxo.Change = controlProg.Change
	return utxo
}

// MockTx mock a tx
func MockTx(utxo *account.UTXO, testAccount *account.Account) (*txbuilder.Template, *types.TxData, error) {
	txInput, sigInst, err := account.UtxoToInputs(testAccount.Signer, utxo)
	if err != nil {
		return nil, nil, err
	}

	b := txbuilder.NewBuilder(time.Now())
	if err := b.AddInput(txInput, sigInst); err != nil {
		return nil, nil, err
	}
	out := types.NewTxOutput(*consensus.MITYAssetID, 100, []byte{byte(vm.OP_FAIL)})
	if err := b.AddOutput(out); err != nil {
		return nil, nil, err
	}
	return b.Build()
}

// MockSign sign a tx
func MockSign(tpl *txbuilder.Template, hsm *pseudohsm.HSM, password string) (bool, error) {
	err := txbuilder.Sign(nil, tpl, password, func(_ context.Context, xpub chainkd.XPub, path [][]byte, data [32]byte, password string) ([]byte, error) {
		return hsm.XSign(xpub, path, data[:], password)
	})
	if err != nil {
		return false, err
	}
	return txbuilder.SignProgress(tpl), nil
}

// MockBlock mock a block
func MockBlock() *bc.Block {
	return &bc.Block{
		BlockHeader: &bc.BlockHeader{Height: 1},
	}
}
