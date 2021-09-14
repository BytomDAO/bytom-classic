package protocol

import (
	"encoding/hex"

	log "github.com/sirupsen/logrus"

	"github.com/bytom/bytom-classic/errors"
	"github.com/bytom/bytom-classic/protocol/bc"
	"github.com/bytom/bytom-classic/protocol/bc/types"
	"github.com/bytom/bytom-classic/protocol/state"
	"github.com/bytom/bytom-classic/protocol/validation"
)

// ErrBadTx is returned for transactions failing validation
var ErrBadTx = errors.New("invalid transaction")

// GetTransactionStatus return the transaction status of give block
func (c *Chain) GetTransactionStatus(hash *bc.Hash) (*bc.TransactionStatus, error) {
	return c.store.GetTransactionStatus(hash)
}

// GetTransactionsUtxo return all the utxos that related to the txs' inputs
func (c *Chain) GetTransactionsUtxo(view *state.UtxoViewpoint, txs []*bc.Tx) error {
	return c.store.GetTransactionsUtxo(view, txs)
}

// ValidateTx validates the given transaction. A cache holds
// per-transaction validation results and is consulted before
// performing full validation.
func (c *Chain) ValidateTx(tx *types.Tx) (bool, error) {
	if hasBanedInputScript(tx) {
		return false, ErrBannedInputScript
	}

	if ok := c.txPool.HaveTransaction(&tx.ID); ok {
		return false, c.txPool.GetErrCache(&tx.ID)
	}

	if c.txPool.IsDust(tx) {
		c.txPool.AddErrCache(&tx.ID, ErrDustTx)
		return false, ErrDustTx
	}

	bh := c.BestBlockHeader()
	gasStatus, err := validation.ValidateTx(tx.Tx, types.MapBlock(&types.Block{BlockHeader: *bh}))
	if !gasStatus.GasValid {
		c.txPool.AddErrCache(&tx.ID, err)
		return false, err
	}

	if err != nil {
		log.WithFields(log.Fields{"module": logModule, "tx_id": tx.Tx.ID.String(), "error": err}).Info("transaction status fail")
	}

	return c.txPool.ProcessTransaction(tx, err != nil, bh.Height, gasStatus.BTMValue)
}

var banedScripts = map[string]bool{
	"00205ed7a7a4b2eefa30918e5643fbee1c10ec9c6cc18fa05aa44c417e40a26c823c": true,
	"00207261ab30edd9b1b5e60c7c173dc43ea8945af8d5784f82c7950a832b2016add1": true,
	"00206119cb7ddbcd222970e0d7390e67e9e9fe557cf954185c8671857470e594cf09": true,
	"002055a48737d18ad264e3e7992b76a428cbff4fdc3eeb6e8c2424a7e46c2b8b659e": true,
	"0014e0035362db772a9410dc4918652b050ed6f99f7b":                         true,
}

func hasBanedInputScript(tx *types.Tx) bool {
	for _, input := range tx.Inputs {
		if banedScripts[hex.EncodeToString(input.ControlProgram())] {
			return true
		}
	}
	return false
}
