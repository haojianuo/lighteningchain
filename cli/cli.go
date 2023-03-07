package cli

import (
	"bytes"
	"flag"
	"fmt"
	"lighteningchain/blockchain"
	"lighteningchain/utils"
	"os"
	"runtime"
	"strconv"
)

type CommandLine struct{}

func (cli *CommandLine) printUsage() {
	fmt.Println("Welcome to Arno's LighteningChain System!")
	fmt.Println("===================================================")
	fmt.Println("Usage:")
	fmt.Println("---------------------------------------------------")
	fmt.Println("To get started, create a new blockchain and declare the owner's address. " +
		"You can then make transactions and mine blocks to add to the blockchain.")
	fmt.Println("---------------------------------------------------")
	fmt.Println("Commands:")
	fmt.Println("---------------------------------------------------")
	fmt.Println("createblockchain -address ADDRESS  -> " +
		"Create a new blockchain with the specified owner's address")
	fmt.Println("balance -address ADDRESS           -> " +
		"Check the balance of a specific address in the blockchain")
	fmt.Println("blockchaininfo                      -> " +
		"View information about all blocks in the blockchain")
	fmt.Println("send -from FROMADDRESS -to TOADDRESS -amount AMOUNT  -> " +
		"Create a new transaction and add it to the candidate block for mining")
	fmt.Println("mine                                -> " +
		"Mine a block and add it to the blockchain")
	fmt.Println("===================================================")
}
func (cli *CommandLine) createBlockChain(address string) {
	newChain := blockchain.InitBlockChain([]byte(address))
	newChain.Database.Close()
	fmt.Println("Finished creating blockchain, and the owner is: ", address)
}
func (cli *CommandLine) balance(address string) {
	chain := blockchain.LoadBlockChain()
	defer chain.Database.Close()

	balance, _ := chain.FindUTXOs([]byte(address))
	fmt.Printf("Address:%s, Balance:%d \n", address, balance)
}

// getblockchaininfo命令需要使用我们之前设计的迭代器遍历区块链。
func (cli *CommandLine) getBlockChainInfo() {
	chain := blockchain.LoadBlockChain()
	defer chain.Database.Close()
	iterator := chain.Iterator()
	ogprevhash := chain.BackOgPrevHash()
	for {
		block := iterator.Next()
		fmt.Println("--------------------------------------------------------------------------------------------------------------")
		fmt.Printf("Timestamp:%d\n", block.Timestamp)
		fmt.Printf("Previous hash:%x\n", block.PrevHash)
		fmt.Printf("Transactions:%v\n", block.Transactions)
		fmt.Printf("hash:%x\n", block.Hash)
		fmt.Printf("Pow: %s\n", strconv.FormatBool(block.ValidatePoW()))
		fmt.Println("--------------------------------------------------------------------------------------------------------------")
		fmt.Println()
		if bytes.Equal(block.PrevHash, ogprevhash) {
			break
		}
	}
}

// send命令将会调用CreateTransaction函数，并将创建的交易信息保存到交易信息池中。
func (cli *CommandLine) send(from, to string, amount int) {
	chain := blockchain.LoadBlockChain()
	defer chain.Database.Close()
	tx, ok := chain.CreateTransaction([]byte(from), []byte(to), amount)
	if !ok {
		fmt.Println("Failed to create transaction")
		return
	}
	tp := blockchain.CreateTransactionPool()
	tp.AddTransaction(tx)
	tp.SaveFile()
	fmt.Println("Success!")
}

// mine命令调用RunMine即可
func (cli *CommandLine) mine() {
	chain := blockchain.LoadBlockChain()
	defer chain.Database.Close()
	chain.RunMine()
	fmt.Println("Finish Mining")
}

// 使用flag库将各命令注册即可
func (cli *CommandLine) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		runtime.Goexit()
	}
}

func (cli *CommandLine) Run() {
	cli.validateArgs()

	createBlockChainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	balanceCmd := flag.NewFlagSet("balance", flag.ExitOnError)
	getBlockChainInfoCmd := flag.NewFlagSet("blockchaininfo", flag.ExitOnError)
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)
	mineCmd := flag.NewFlagSet("mine", flag.ExitOnError)

	createBlockChainOwner := createBlockChainCmd.String("address", "", "The address refer to the owner of blockchain")
	balanceAddress := balanceCmd.String("address", "", "Who need to get balance amount")
	sendFromAddress := sendCmd.String("from", "", "Source address")
	sendToAddress := sendCmd.String("to", "", "Destination address")
	sendAmount := sendCmd.Int("amount", 0, "Amount to send")

	switch os.Args[1] {
	case "createblockchain":
		err := createBlockChainCmd.Parse(os.Args[2:])
		utils.Handle(err)

	case "balance":
		err := balanceCmd.Parse(os.Args[2:])
		utils.Handle(err)

	case "blockchaininfo":
		err := getBlockChainInfoCmd.Parse(os.Args[2:])
		utils.Handle(err)

	case "send":
		err := sendCmd.Parse(os.Args[2:])
		utils.Handle(err)

	case "mine":
		err := mineCmd.Parse(os.Args[2:])
		utils.Handle(err)

	default:
		cli.printUsage()
		runtime.Goexit()
	}

	if createBlockChainCmd.Parsed() {
		if *createBlockChainOwner == "" {
			createBlockChainCmd.Usage()
			runtime.Goexit()
		}
		cli.createBlockChain(*createBlockChainOwner)
	}

	if balanceCmd.Parsed() {
		if *balanceAddress == "" {
			balanceCmd.Usage()
			runtime.Goexit()
		}
		cli.balance(*balanceAddress)
	}

	if sendCmd.Parsed() {
		if *sendFromAddress == "" || *sendToAddress == "" || *sendAmount <= 0 {
			sendCmd.Usage()
			runtime.Goexit()
		}
		cli.send(*sendFromAddress, *sendToAddress, *sendAmount)
	}

	if getBlockChainInfoCmd.Parsed() {
		cli.getBlockChainInfo()
	}

	if mineCmd.Parsed() {
		cli.mine()
	}
}
