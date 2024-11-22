package bot

import (
	"context"
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

func generateSeedAddress() (seeds string, addr string) {
	seed := wallet.NewSeed()
	seeds = strings.Join(seed, " ")

	client := liteclient.NewConnectionPool()

	cfg, err := liteclient.GetConfigFromUrl(context.Background(), TonConfig)
	if err != nil {
		loge(err)
	}

	err = client.AddConnectionsFromConfig(context.Background(), cfg)
	if err != nil {
		loge(err)
	}

	api := ton.NewAPIClient(client, ton.ProofCheckPolicyFast).WithRetry()
	api.SetTrustedBlockFromConfig(cfg)

	w, err := wallet.FromSeed(api, seed, wallet.V4R2)
	if err != nil {
		loge(err)
	}

	addr = w.WalletAddress().String()

	return seeds, addr
}

func getBalance(addr string) uint64 {
	balance := uint64(0)

	client := liteclient.NewConnectionPool()

	cfg, err := liteclient.GetConfigFromUrl(context.Background(), TonConfig)
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

	cfg, err := liteclient.GetConfigFromUrl(context.Background(), TonConfig)
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

	cfg, err := liteclient.GetConfigFromUrl(context.Background(), TonConfig)
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

func splitPayment(balance uint64, u *User) {
	half := balance / 2
	send(int64(half), AddressTonAd, u.Seed)

	time.Sleep(time.Minute * 2)

	balance = getBalance(u.AddressDeposit)
	sendall(int64(balance), AddressReward, u.Seed)

	time.Sleep(time.Minute * 2)

	balance = getBalance(u.AddressDeposit)
	u.Balance = balance
	if err := db.Save(u).Error; err != nil {
		loge(err)
	}
}

// func test() {
// 	client := liteclient.NewConnectionPool()

// 	// get config
// 	cfg, err := liteclient.GetConfigFromUrl(context.Background(), TonConfig)
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
