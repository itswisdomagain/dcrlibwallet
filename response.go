package mobilewallet

type UnsignedTransaction struct {
	UnsignedTransaction       []byte
	EstimatedSignedSize       int
	ChangeIndex               int
	TotalOutputAmount         int64
	TotalPreviousOutputAmount int64
}

type Balance struct {
	Total                   int64
	Spendable               int64
	ImmatureReward          int64
	ImmatureStakeGeneration int64
	LockedByTickets         int64
	VotingAuthority         int64
	UnConfirmed             int64
}

type Account struct {
	Number           int32
	Name             string
	Balance          *Balance
	TotalBalance     int64
	ExternalKeyCount int32
	InternalKeyCount int32
	ImportedKeyCount int32
}

type Accounts struct {
	Count              int
	ErrorMessage       string
	ErrorCode          int
	ErrorOccurred      bool
	Acc                *[]Account
	CurrentBlockHash   []byte
	CurrentBlockHeight int32
}

type BlockScanResponse interface {
	OnScan(rescannedThrough int32) bool
	OnEnd(height int32, cancelled bool)
	OnError(code int32, message string)
}

/*
Direction
0: Sent
1: Received
2: Transfered
*/
type Transaction struct {
	Hash        string
	Transaction []byte
	Fee         int64
	Timestamp   int64
	Type        string
	Amount      int64
	Status      string
	Height      int32
	Direction   int32
	Debits      *[]TransactionDebit
	Credits     *[]TransactionCredit
}

type TransactionDebit struct {
	Index           int32
	PreviousAccount int32
	PreviousAmount  int64
	AccountName     string
}

type TransactionCredit struct {
	Index    int32
	Account  int32
	Internal bool
	Amount   int64
	Address  string
}

type getTransactionsResponse struct {
	Transactions  []Transaction
	ErrorOccurred bool
	ErrorMessage  string
}

type GetTransactionsResponse interface {
	OnResult(json string)
}

type TransactionListener interface {
	OnTransaction(transaction string)
	OnTransactionConfirmed(hash string, height int32)
	OnBlockAttached(height int32, timestamp int64)
}

type BlockNotificationError interface {
	OnBlockNotificationError(err error)
}

type DecodedTransaction struct {
	Hash     string
	Type     string
	Version  int32
	LockTime int32
	Expiry   int32
	Inputs   []DecodedInput
	Outputs  []DecodedOutput
}

type DecodedInput struct {
	PreviousTransactionHash  string
	PreviousTransactionIndex int32
	Sequence                 int32
	AmountIn                 int64
	BlockHeight              int32
	BlockIndex               int32
}

type DecodedOutput struct {
	Index     int32
	Value     int64
	Version   int32
	Addresses []string
}

type SpvSyncResponse interface {
	OnPeerConnected(peerCount int32)
	OnPeerDisconnected(peerCount int32)
	OnFetchMissingCFilters(fetchedCFiltersCount int32)
	OnFetchedHeaders(peerInitialHeight, fetchedHeadersCount int32, lastHeaderTime int64)
	OnDiscoveredAddresses(finished bool)
	OnRescanProgress(rescannedThrough int32)
	OnSynced(synced bool)
	/*
	* Handled Error Codes
	* -1 - Unexpected Error
	*  1 - Context Canceled
	*  2 - Deadline Exceeded
	*  3 - Invalid Address
	 */
	OnSyncError(code int, err error)
}
