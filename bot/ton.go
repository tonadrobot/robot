package bot

import (
	"context"
	"encoding/hex"
	"log"
	"math/big"
	"strings"
	"time"

	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/liteclient"
	"github.com/xssnick/tonutils-go/tlb"
	"github.com/xssnick/tonutils-go/ton"
	"github.com/xssnick/tonutils-go/ton/wallet"
)

func generateSeedAddress() (seeds string, addr string, err error) {
	seed := wallet.NewSeed()
	seeds = strings.Join(seed, " ")

	client := liteclient.NewConnectionPool()

	cfg, err := liteclient.GetConfigFromUrl(context.Background(), getTonConfig())
	if err != nil {
		loge(err)
		return "", "", err
	}

	err = client.AddConnectionsFromConfig(context.Background(), cfg)
	if err != nil {
		loge(err)
		return "", "", err
	}

	api := ton.NewAPIClient(client, ton.ProofCheckPolicyFast).WithRetry()
	api.SetTrustedBlockFromConfig(cfg)

	w, err := wallet.FromSeed(api, seed, wallet.V4R2)
	if err != nil {
		loge(err)
		return "", "", err
	}

	addr = w.WalletAddress().String()

	return seeds, addr, nil
}

func getBalance(addr string) uint64 {
	balance := uint64(0)

	client := liteclient.NewConnectionPool()

	cfg, err := liteclient.GetConfigFromUrl(context.Background(), getTonConfig())
	if err != nil {
		loge(err)
		return 0
	}

	err = client.AddConnectionsFromConfig(context.Background(), cfg)
	if err != nil {
		loge(err)
		return 0
	}

	api := ton.NewAPIClient(client, ton.ProofCheckPolicyFast).WithRetry()
	api.SetTrustedBlockFromConfig(cfg)

	ctx := client.StickyContext(context.Background())

	b, err := api.CurrentMasterchainInfo(ctx)
	if err != nil {
		loge(err)
		return 0
	}

	a := address.MustParseAddr(addr)

	res, err := api.WaitForBlock(b.SeqNo).GetAccount(ctx, b, a)
	if err != nil {
		loge(err)
		return 0
	}

	// list, _ := api.ListTransactions(context.Background(), a, 1000, res.LastTxLT, res.LastTxHash)
	// log.Println(prettyPrint(list[len(list)-1]))

	if res.IsActive {
		balance = res.State.Balance.Nano().Uint64()
	}

	return balance
}

func send(amount int64, to string, seed string) {
	client := liteclient.NewConnectionPool()

	cfg, err := liteclient.GetConfigFromUrl(context.Background(), getTonConfig())
	if err != nil {
		loge(err)
	}

	err = client.AddConnectionsFromConfig(context.Background(), cfg)
	if err != nil {
		loge(err)
	}

	api := ton.NewAPIClient(client, ton.ProofCheckPolicyFast).WithRetry()
	api.SetTrustedBlockFromConfig(cfg)

	ctx := client.StickyContext(context.Background())

	words := strings.Split(seed, " ")

	w, err := wallet.FromSeed(api, words, wallet.V4R2)
	if err != nil {
		loge(err)
	}

	log.Println("wallet address:", w.WalletAddress())

	block, err := api.CurrentMasterchainInfo(context.Background())
	if err != nil {
		loge(err)
	}

	balance, err := w.GetBalance(ctx, block)
	if err != nil {
		loge(err)
	}

	if balance.Nano().Uint64() >= 3000000 {
		addr := address.MustParseAddr(to)

		// bounce := false

		// transfer, err := w.BuildTransfer(addr, tlb.MustFromNano(big.NewInt(amount), 9), bounce, "TON Miners withdraw.")
		// if err != nil {
		// 	loge(err)
		// }

		// _, _, err = w.SendWaitTransaction(ctx, transfer)
		// if err != nil {
		// 	loge(err)
		// }
		err = w.Send(context.Background(), &wallet.Message{
			Mode: 1, // pay fees separately (from balance, not from amount)
			InternalMessage: &tlb.InternalMessage{
				Bounce:  false, // return amount in case of processing error
				DstAddr: addr,  // destination is domain contract
				Amount:  tlb.MustFromNano(big.NewInt(amount), 9),
				Body:    nil,
			},
		}, true)
		if err != nil {
			loge(err)
		}

	}
}

func sendall(amount int64, to string, seed string) {
	client := liteclient.NewConnectionPool()

	cfg, err := liteclient.GetConfigFromUrl(context.Background(), getTonConfig())
	if err != nil {
		loge(err)
	}

	err = client.AddConnectionsFromConfig(context.Background(), cfg)
	if err != nil {
		loge(err)
	}

	api := ton.NewAPIClient(client, ton.ProofCheckPolicyFast).WithRetry()
	api.SetTrustedBlockFromConfig(cfg)

	ctx := client.StickyContext(context.Background())

	words := strings.Split(seed, " ")

	w, err := wallet.FromSeed(api, words, wallet.V4R2)
	if err != nil {
		loge(err)
	}

	log.Println("wallet address:", w.WalletAddress())

	block, err := api.CurrentMasterchainInfo(context.Background())
	if err != nil {
		loge(err)
	}

	balance, err := w.GetBalance(ctx, block)
	if err != nil {
		loge(err)
	}

	if balance.Nano().Uint64() >= 3000000 {
		addr := address.MustParseAddr(to)

		// bounce := false

		// transfer, err := w.BuildTransfer(addr, tlb.MustFromNano(big.NewInt(amount), 9), bounce, "TON Miners withdraw.")
		// if err != nil {
		// 	loge(err)
		// }

		// _, _, err = w.SendWaitTransaction(ctx, transfer)
		// if err != nil {
		// 	loge(err)
		// }
		err = w.Send(context.Background(), &wallet.Message{
			Mode: wallet.CarryAllRemainingBalance, // pay fees separately (from balance, not from amount)
			InternalMessage: &tlb.InternalMessage{
				Bounce:  false, // return amount in case of processing error
				DstAddr: addr,  // destination is domain contract
				Amount:  tlb.MustFromNano(big.NewInt(amount), 9),
				Body:    nil,
			},
		}, true)
		if err != nil {
			loge(err)
		}

	}
}

func splitPayment(u *User) {
	time.Sleep(time.Minute * 2)

	balance := getBalance(u.AddressDeposit)

	half := balance / 2
	send(int64(half), AddressTonAd, u.Seed)

	time.Sleep(time.Minute * 2)

	balance = getBalance(u.AddressDeposit)
	sendall(int64(balance), AddressReward, u.Seed)
}

// func test() {
// 	client := liteclient.NewConnectionPool()

// 	// get config
// 	cfg, err := liteclient.GetConfigFromUrl(context.Background(), getTonConfig())
// 	if err != nil {
// 		log.Fatalln("get config err: ", err.Error())
// 	}

// 	// connect to mainnet lite servers
// 	err = client.AddConnectionsFromConfig(context.Background(), cfg)
// 	if err != nil {
// 		log.Fatalln("connection err: ", err.Error())
// 	}

// 	// api client with full proof checks
// 	api := ton.NewAPIClient(client, ton.ProofCheckPolicyFast).WithRetry()
// 	api.SetTrustedBlockFromConfig(cfg)

// 	// bound all requests to single ton node
// 	ctx := client.StickyContext(context.Background())

// 	// seed words of account, you can generate them with any wallet or using wallet.NewSeed() method
// 	words := strings.Split("earn arrive short number regular area foam combine brick marine fluid labor tiger warm lady novel casino ski increase apart kiss couple frost renew", " ")

// 	w, err := wallet.FromSeed(api, words, wallet.V4R2)
// 	if err != nil {
// 		log.Fatalln("FromSeed err:", err.Error())
// 	}

// 	log.Println("wallet address:", w.WalletAddress())

// 	log.Println("fetching and checking proofs since config init block, it may take near a minute...")
// 	block, err := api.CurrentMasterchainInfo(context.Background())
// 	if err != nil {
// 		log.Fatalln("get masterchain info err: ", err.Error())
// 	}
// 	log.Println("master proof checks are completed successfully, now communication is 100% safe!")

// 	balance, err := w.GetBalance(ctx, block)
// 	if err != nil {
// 		log.Fatalln("GetBalance err:", err.Error())
// 	}

// 	if balance.Nano().Uint64() >= 3000000 {
// 		addr := address.MustParseAddr("0QD1jykb42qB1yE2hOYAGaJWcOReKHYbCCQUjzrOHSmluNGF")

// 		log.Println("sending transaction and waiting for confirmation...")

// 		// if destination wallet is not initialized (or you don't care)
// 		// you should set bounce to false to not get money back.
// 		bounce := false

// 		transfer, err := w.BuildTransfer(addr, tlb.MustFromTON("0.003"), bounce, "Hello from tonutils-go!")
// 		if err != nil {
// 			log.Fatalln("Transfer err:", err.Error())
// 		}

// 		tx, block, err := w.SendWaitTransaction(ctx, transfer)
// 		if err != nil {
// 			log.Fatalln("SendWaitTransaction err:", err.Error())
// 		}

// 		balance, err = w.GetBalance(ctx, block)
// 		if err != nil {
// 			log.Fatalln("GetBalance err:", err.Error())
// 		}

// 		log.Printf("transaction confirmed at block %d, hash: %s balance left: %s", block.SeqNo,
// 			base64.StdEncoding.EncodeToString(tx.Hash), balance.String())

// 	}

// 	log.Println("not enough balance:", balance.String())
// }

func checkNewTmu(u *User) uint64 {
	lastProcessedLT := uint64(0)
	var lastProcessedHash []byte
	new := uint64(0)

	client := liteclient.NewConnectionPool()

	cfg, err := liteclient.GetConfigFromUrl(context.Background(), getTonConfig())
	if err != nil {
		loge(err)
		return 0
	}

	// connect to mainnet lite servers
	err = client.AddConnectionsFromConfig(context.Background(), cfg)
	if err != nil {
		loge(err)
		return 0
	}

	// initialize ton api lite connection wrapper with full proof checks
	api := ton.NewAPIClient(client, ton.ProofCheckPolicyFast).WithRetry()
	api.SetTrustedBlockFromConfig(cfg)

	master, err := api.CurrentMasterchainInfo(context.Background()) // we fetch block just to trigger chain proof check
	if err != nil {
		loge(err)
		return 0
	}

	// address on which we are accepting payments
	treasuryAddress, err := address.ParseAddr(u.AddressDeposit)
	if err != nil {
		log.Println(err)
		return 0
	}

	acc, err := api.GetAccount(context.Background(), master, treasuryAddress)
	if err != nil {
		loge(err)
		return 0
	}

	// Cursor of processed transaction, save it to your db
	// We start from last transaction, will not process transactions older than we started from.
	// After each processed transaction, save lt to your db, to continue after restart
	// if u.LastTxLT == 0 {
	// 	lastProcessedLT = acc.LastTxLT
	// 	lastProcessedHash = acc.LastTxHash
	// } else {
	// 	lastProcessedLT = u.LastTxLT
	// 	lastProcessedHash, err = hex.DecodeString(u.LastTxHash)
	// 	if err != nil {
	// 		loge(err)
	// 		return 0
	// 	}
	// }

	lastProcessedLT = acc.LastTxLT
	lastProcessedHash = acc.LastTxHash

	// lastProcessedLT := uint64(28970514000003)
	// channel with new transactions
	// transactions := make(chan *tlb.Transaction)

	// it is a blocking call, so we start it asynchronously
	// go api.SubscribeOnTransactions(context.Background(), treasuryAddress, lastProcessedLT, transactions)
	// hash, _ := hex.DecodeString("655abbdaba882076a649fae19a351cee53bfcbd22d79a908a8365ec7fe9e93ee")
	transactions, err := api.ListTransactions(context.Background(), treasuryAddress, 300, lastProcessedLT, lastProcessedHash)
	if err != nil {
		log.Println(err)
		return 0
	}

	// log.Println("waiting for transfers...")

	// USDT master contract addr, but can be any jetton
	// usdt := jetton.NewJettonMasterClient(api, address.MustParseAddr("EQCxE6mUtQJKFnGfaROTKOt1lZbDiiX1kCixRv7Nw2Id_sDs"))
	// get our jetton wallet address
	// treasuryJettonWallet, err := usdt.GetJettonWalletAtBlock(context.Background(), treasuryAddress, master)
	// if err != nil {
	// 	log.Fatalln("get jetton wallet address err: ", err.Error())
	// 	return
	// }

	// tonWallet, err := wallet.from

	// listen for new transactions from channel
	for _, tx := range transactions {
		// only internal messages can increase the balance
		if tx.IO.In != nil && tx.IO.In.MsgType == tlb.MsgTypeInternal {
			ti := tx.IO.In.AsInternal()
			src := ti.SrcAddr

			// verify that event sender is our jetton wallet
			// if ti.SrcAddr.Equals(treasuryJettonWallet.Address()) {
			// 	var transfer jetton.TransferNotification
			// 	if err = tlb.LoadFromCell(&transfer, ti.Body.BeginParse()); err == nil {
			// 		// convert decimals to 6 for USDT (it can be fetched from jetton details too), default is 9
			// 		amt := tlb.MustFromNano(transfer.Amount.Nano(), 6)

			// 		// reassign sender to real jetton sender instead of its jetton wallet contract
			// 		src = transfer.Sender
			// 		log.Println("received", amt.String(), "USDT from", src.String())
			// 	}
			// }

			// show received ton amount
			log.Println("received", ti.Amount.String(), "TON from", src.String())

			if !isTxProcessed(hex.EncodeToString(tx.Hash), tx.LT) {
				new += ti.Amount.Nano().Uint64()

				lastProcessedLT = tx.LT
				lastProcessedHash = tx.Hash

				log.Printf("Hash: %s", hex.EncodeToString(lastProcessedHash))
				log.Printf("LT: %d", lastProcessedLT)

				processTx(hex.EncodeToString(lastProcessedHash), tx.LT)
			}
		}
	}

	u.LastTxLT = lastProcessedLT
	u.LastTxHash = hex.EncodeToString(lastProcessedHash)

	if err := db.Save(u).Error; err != nil {
		loge(err)
	}

	log.Printf("New TON: %d", new)

	return new
}

func getTonConfig() string {
	if conf.Dev {
		return DevTonConfig
	} else {
		return TonConfig
	}
}
