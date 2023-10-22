package api

import (
	"context"

	"github.com/anonimitycash/anonimitycash-classic/account"
	"github.com/anonimitycash/anonimitycash-classic/asset"
	"github.com/anonimitycash/anonimitycash-classic/blockchain/pseudohsm"
	"github.com/anonimitycash/anonimitycash-classic/blockchain/rpc"
	"github.com/anonimitycash/anonimitycash-classic/blockchain/signers"
	"github.com/anonimitycash/anonimitycash-classic/blockchain/txbuilder"
	"github.com/anonimitycash/anonimitycash-classic/errors"
	"github.com/anonimitycash/anonimitycash-classic/net/http/httperror"
	"github.com/anonimitycash/anonimitycash-classic/net/http/httpjson"
	"github.com/anonimitycash/anonimitycash-classic/protocol/validation"
	"github.com/anonimitycash/anonimitycash-classic/protocol/vm"
)

var (
	// ErrDefault is default Anonimitycash API Error
	ErrDefault = errors.New("Anonimitycash API Error")
)

func isTemporary(info httperror.Info, err error) bool {
	switch info.ChainCode {
	case "MITY000": // internal server error
		return true
	case "MITY001": // request timed out
		return true
	case "MITY761": // outputs currently reserved
		return true
	case "MITY706": // 1 or more action errors
		errs := errors.Data(err)["actions"].([]httperror.Response)
		temp := true
		for _, actionErr := range errs {
			temp = temp && isTemporary(actionErr.Info, nil)
		}
		return temp
	default:
		return false
	}
}

var respErrFormatter = map[error]httperror.Info{
	ErrDefault: {500, "MITY000", "Anonimitycash API Error"},

	// Signers error namespace (2xx)
	signers.ErrBadQuorum: {400, "MITY200", "Quorum must be greater than or equal to 1, and must be less than or equal to the length of xpubs"},
	signers.ErrBadXPub:   {400, "MITY201", "Invalid xpub format"},
	signers.ErrNoXPubs:   {400, "MITY202", "At least one xpub is required"},
	signers.ErrDupeXPub:  {400, "MITY203", "Root XPubs cannot contain the same key more than once"},

	// Contract error namespace (3xx)
	ErrCompileContract: {400, "MITY300", "Compile contract failed"},
	ErrInstContract:    {400, "MITY301", "Instantiate contract failed"},

	// Transaction error namespace (7xx)
	// Build transaction error namespace (70x ~ 72x)
	account.ErrInsufficient:         {400, "MITY700", "Funds of account are insufficient"},
	account.ErrImmature:             {400, "MITY701", "Available funds of account are immature"},
	account.ErrReserved:             {400, "MITY702", "Available UTXOs of account have been reserved"},
	account.ErrMatchUTXO:            {400, "MITY703", "UTXO with given hash not found"},
	ErrBadActionType:                {400, "MITY704", "Invalid action type"},
	ErrBadAction:                    {400, "MITY705", "Invalid action object"},
	ErrBadActionConstruction:        {400, "MITY706", "Invalid action construction"},
	txbuilder.ErrMissingFields:      {400, "MITY707", "One or more fields are missing"},
	txbuilder.ErrBadAmount:          {400, "MITY708", "Invalid asset amount"},
	account.ErrFindAccount:          {400, "MITY709", "Account not found"},
	asset.ErrFindAsset:              {400, "MITY710", "Asset not found"},
	txbuilder.ErrBadContractArgType: {400, "MITY711", "Invalid contract argument type"},
	txbuilder.ErrOrphanTx:           {400, "MITY712", "Transaction input UTXO not found"},
	txbuilder.ErrExtTxFee:           {400, "MITY713", "Transaction fee exceeded max limit"},
	txbuilder.ErrNoGasInput:         {400, "MITY714", "Transaction has no gas input"},

	// Submit transaction error namespace (73x ~ 79x)
	// Validation error (73x ~ 75x)
	validation.ErrTxVersion:                 {400, "MITY730", "Invalid transaction version"},
	validation.ErrWrongTransactionSize:      {400, "MITY731", "Invalid transaction size"},
	validation.ErrBadTimeRange:              {400, "MITY732", "Invalid transaction time range"},
	validation.ErrNotStandardTx:             {400, "MITY733", "Not standard transaction"},
	validation.ErrWrongCoinbaseTransaction:  {400, "MITY734", "Invalid coinbase transaction"},
	validation.ErrWrongCoinbaseAsset:        {400, "MITY735", "Invalid coinbase assetID"},
	validation.ErrCoinbaseArbitraryOversize: {400, "MITY736", "Invalid coinbase arbitrary size"},
	validation.ErrEmptyResults:              {400, "MITY737", "No results in the transaction"},
	validation.ErrMismatchedAssetID:         {400, "MITY738", "Mismatched assetID"},
	validation.ErrMismatchedPosition:        {400, "MITY739", "Mismatched value source/dest position"},
	validation.ErrMismatchedReference:       {400, "MITY740", "Mismatched reference"},
	validation.ErrMismatchedValue:           {400, "MITY741", "Mismatched value"},
	validation.ErrMissingField:              {400, "MITY742", "Missing required field"},
	validation.ErrNoSource:                  {400, "MITY743", "No source for value"},
	validation.ErrOverflow:                  {400, "MITY744", "Arithmetic overflow/underflow"},
	validation.ErrPosition:                  {400, "MITY745", "Invalid source or destination position"},
	validation.ErrUnbalanced:                {400, "MITY746", "Unbalanced asset amount between input and output"},
	validation.ErrOverGasCredit:             {400, "MITY747", "Gas credit has been spent"},
	validation.ErrGasCalculate:              {400, "MITY748", "Gas usage calculate got a math error"},

	// VM error (76x ~ 78x)
	vm.ErrAltStackUnderflow:  {400, "MITY760", "Alt stack underflow"},
	vm.ErrBadValue:           {400, "MITY761", "Bad value"},
	vm.ErrContext:            {400, "MITY762", "Wrong context"},
	vm.ErrDataStackUnderflow: {400, "MITY763", "Data stack underflow"},
	vm.ErrDisallowedOpcode:   {400, "MITY764", "Disallowed opcode"},
	vm.ErrDivZero:            {400, "MITY765", "Division by zero"},
	vm.ErrFalseVMResult:      {400, "MITY766", "False result for executing VM"},
	vm.ErrLongProgram:        {400, "MITY767", "Program size exceeds max int32"},
	vm.ErrRange:              {400, "MITY768", "Arithmetic range error"},
	vm.ErrReturn:             {400, "MITY769", "RETURN executed"},
	vm.ErrRunLimitExceeded:   {400, "MITY770", "Run limit exceeded because the MITY Fee is insufficient"},
	vm.ErrShortProgram:       {400, "MITY771", "Unexpected end of program"},
	vm.ErrToken:              {400, "MITY772", "Unrecognized token"},
	vm.ErrUnexpected:         {400, "MITY773", "Unexpected error"},
	vm.ErrUnsupportedVM:      {400, "MITY774", "Unsupported VM because the version of VM is mismatched"},
	vm.ErrVerifyFailed:       {400, "MITY775", "VERIFY failed"},

	// Mock HSM error namespace (8xx)
	pseudohsm.ErrDuplicateKeyAlias: {400, "MITY800", "Key Alias already exists"},
	pseudohsm.ErrLoadKey:           {400, "MITY801", "Key not found or wrong password"},
	pseudohsm.ErrDecrypt:           {400, "MITY802", "Could not decrypt key with given passphrase"},
}

// Map error values to standard anonimitycash error codes. Missing entries
// will map to internalErrInfo.
//
// TODO(jackson): Share one error table across Chain
// products/services so that errors are consistent.
var errorFormatter = httperror.Formatter{
	Default:     httperror.Info{500, "MITY000", "Anonimitycash API Error"},
	IsTemporary: isTemporary,
	Errors: map[error]httperror.Info{
		// General error namespace (0xx)
		context.DeadlineExceeded: {408, "MITY001", "Request timed out"},
		httpjson.ErrBadRequest:   {400, "MITY002", "Invalid request body"},
		rpc.ErrWrongNetwork:      {502, "MITY103", "A peer core is operating on a different blockchain network"},

		//accesstoken authz err namespace (86x)
		errNotAuthenticated: {401, "MITY860", "Request could not be authenticated"},
	},
}
