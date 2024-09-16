rd /s /q tmp
md tmp\blocks
md tmp\wallets
md tmp\ref_list
md tmp\utxoset
go build -o main.exe
.\main.exe createwallet
.\main.exe walletslist
.\main.exe createwallet -refname Arno
.\main.exe walletinfo -refname Arno
.\main.exe createwallet -refname Krad
.\main.exe createwallet -refname Exia
.\main.exe createwallet
.\main.exe walletslist
.\main.exe createblockchain -refname Arno
.\main.exe blockchaininfo
.\main.exe initutxoset
.\main.exe balance -refname Arno
.\main.exe sendbyrefname -from Arno -to Krad -amount 100
.\main.exe balance -refname Krad
.\main.exe mine
.\main.exe blockchaininfo
.\main.exe balance -refname Arno
.\main.exe balance -refname Krad
.\main.exe sendbyrefname -from Arno -to Exia -amount 100
.\main.exe sendbyrefname -from Krad -to Exia -amount 30
.\main.exe mine
.\main.exe blockchaininfo
.\main.exe balance -refname Arno
.\main.exe balance -refname Krad
.\main.exe balance -refname Exia
.\main.exe sendbyrefname -from Exia -to Arno -amount 90
.\main.exe sendbyrefname -from Exia -to Krad -amount 90
.\main.exe mine
.\main.exe blockchaininfo
.\main.exe balance -refname Arno
.\main.exe balance -refname Krad
.\main.exe balance -refname Exia