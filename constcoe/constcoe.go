package constcoe

const (
	Difficulty          = 12
	InitCoin            = 1000
	TransactionPoolFile = "./tmp/transaction_pool.data"
	BCPath              = "./tmp/blocks"
	BCFile              = "./tmp/blocks/MANIFEST"
	CheckSumLength      = 4
	NetworkVersion      = byte(0x00)
	Wallets             = "./tmp/wallets/"
	WalletsRefList      = "./tmp/ref_list/"
	UTXOSet             = "./tmp/utxoset/" //This is new
)
