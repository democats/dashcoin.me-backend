package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/rs/cors"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

//var walletRPC = "http://localhost:9090/json_rpc"
//var daemonRPC = "http://45.63.42.189:29081/"
//var difficultyTarget = 120

var walletRPC = "http://localhost:9080/json_rpc"
var daemonRPC = "http://localhost:7701/"
//var daemonRPC = "http://45.63.42.189:7701/"
var difficultyTarget = 90

type CoinConfig struct {
	Jsonrpc string `json:"jsonrpc"`
	Result  struct {
		BaseCoin struct {
			Git  string `json:"git"`
			Name string `json:"name"`
		} `json:"base_coin"`
		Core struct {
			BYTECOINNETWORK                        string   `json:"BYTECOIN_NETWORK"`
			CHECKPOINTS                            []string `json:"CHECKPOINTS"`
			CRYPTONOTEBLOCKGRANTEDFULLREWARDZONE   int      `json:"CRYPTONOTE_BLOCK_GRANTED_FULL_REWARD_ZONE"`
			CRYPTONOTEBLOCKGRANTEDFULLREWARDZONEV1 int      `json:"CRYPTONOTE_BLOCK_GRANTED_FULL_REWARD_ZONE_V1"`
			CRYPTONOTEBLOCKGRANTEDFULLREWARDZONEV2 int      `json:"CRYPTONOTE_BLOCK_GRANTED_FULL_REWARD_ZONE_V2"`
			CRYPTONOTECOINVERSION                  int      `json:"CRYPTONOTE_COIN_VERSION"`
			CRYPTONOTEDISPLAYDECIMALPOINT          int      `json:"CRYPTONOTE_DISPLAY_DECIMAL_POINT"`
			CRYPTONOTEMINEDMONEYUNLOCKWINDOW       int      `json:"CRYPTONOTE_MINED_MONEY_UNLOCK_WINDOW"`
			CRYPTONOTENAME                         string   `json:"CRYPTONOTE_NAME"`
			CRYPTONOTEPUBLICADDRESSBASE58PREFIX    int      `json:"CRYPTONOTE_PUBLIC_ADDRESS_BASE58_PREFIX"`
			DEFAULTDUSTTHRESHOLD                   int      `json:"DEFAULT_DUST_THRESHOLD"`
			DIFFICULTYCUT                          int      `json:"DIFFICULTY_CUT"`
			DIFFICULTYLAG                          int      `json:"DIFFICULTY_LAG"`
			DIFFICULTYTARGET                       int      `json:"DIFFICULTY_TARGET"`
			EMISSIONSPEEDFACTOR                    int      `json:"EMISSION_SPEED_FACTOR"`
			EXPECTEDNUMBEROFBLOCKSPERDAY           int      `json:"EXPECTED_NUMBER_OF_BLOCKS_PER_DAY"`
			GENESISBLOCKREWARD                     string   `json:"GENESIS_BLOCK_REWARD"`
			GENESISCOINBASETXHEX                   string   `json:"GENESIS_COINBASE_TX_HEX"`
			KILLHEIGHT                             int      `json:"KILL_HEIGHT"`
			MANDATORYTRANSACTION                   int      `json:"MANDATORY_TRANSACTION"`
			MAXBLOCKSIZEINITIAL                    int      `json:"MAX_BLOCK_SIZE_INITIAL"`
			MINIMUMFEE                             int      `json:"MINIMUM_FEE"`
			MONEYSUPPLY                            string   `json:"MONEY_SUPPLY"`
			P2PDEFAULTPORT                         int      `json:"P2P_DEFAULT_PORT"`
			RPCDEFAULTPORT                         int      `json:"RPC_DEFAULT_PORT"`
			SEEDNODES                              []string `json:"SEED_NODES"`
			UPGRADEHEIGHTV2                        int      `json:"UPGRADE_HEIGHT_V2"`
			UPGRADEHEIGHTV3                        int      `json:"UPGRADE_HEIGHT_V3"`
		} `json:"core"`
		Extensions []string `json:"extensions"`
		Status     string   `json:"status"`
	} `json:"result"`
}

type WalletdCreateAddressResultSuccess struct {
	Jsonrpc string `json:"jsonrpc"`
	ID      string `json:"id"`
	Result  struct {
		Address string `json:"address"`
	} `json:"result"`
}

type WalletdGetBalanceResult struct {
	Jsonrpc string `json:"jsonrpc"`
	ID      string `json:"id"`
	Result  struct {
		LockedAmount     int `json:"lockedAmount"`
		AvailableBalance int `json:"availableBalance"`
	} `json:"result"`
}

type WalletdGetStatusResult struct {
	ID      string `json:"id"`
	Jsonrpc string `json:"jsonrpc"`
	Result  struct {
		BlockCount      int    `json:"blockCount"`
		KnownBlockCount int    `json:"knownBlockCount"`
		LastBlockHash   string `json:"lastBlockHash"`
		PeerCount       int    `json:"peerCount"`
	} `json:"result"`
}

type WalletdGetTransactionsResult struct {
	ID      string `json:"id"`
	Jsonrpc string `json:"jsonrpc"`
	Result  struct {
		Items []struct {
			BlockHash    string `json:"blockHash"`
			Transactions []struct {
				Amount          int64  `json:"amount"`
				BlockIndex      int    `json:"blockIndex"`
				Extra           string `json:"extra"`
				Fee             int    `json:"fee"`
				IsBase          bool   `json:"isBase"`
				Mixin           int    `json:"mixin"`
				PaymentID       string `json:"paymentId"`
				State           int    `json:"state"`
				Timestamp       int    `json:"timestamp"`
				TransactionHash string `json:"transactionHash"`
				Transfers       []struct {
					Address      string `json:"address"`
					Amount       int64  `json:"amount"`
					SpentOutputs []struct {
						Amount   int64  `json:"amount"`
						KeyImage string `json:"key_image"`
						Mixin    int64  `json:"mixin"`
						OutIndex int64  `json:"out_index"`
						TxPubKey string `json:"tx_pub_key"`
					} `json:"spentOutputs"`
					Type int `json:"type"`
				} `json:"transfers"`
				UnlockTime int `json:"unlockTime"`
			} `json:"transactions"`
		} `json:"items"`
	} `json:"result"`
}

type WalletdGetUnspendOutsResult struct {
	ID      string `json:"id"`
	Jsonrpc string `json:"jsonrpc"`
	Result  struct {
		Outputs []struct {
			Amount               int64  `json:"amount"`
			GlobalOutputIndex    int    `json:"globalOutputIndex"`
			OutputInTransaction  int    `json:"outputInTransaction"`
			OutputKey            string `json:"outputKey"`
			RequiredSignatures   int    `json:"requiredSignatures"`
			TransactionHash      string `json:"transactionHash"`
			TransactionPublicKey string `json:"transactionPublicKey"`
			Type                 int    `json:"type"`
		} `json:"outputs"`
	} `json:"result"`
}

type DaemonGetRandomOutsResult struct {
	Outs []struct {
		Amount int `json:"amount"`
		Outs   []struct {
			GlobalAmountIndex int    `json:"global_amount_index"`
			OutKey            string `json:"out_key"`
		} `json:"outs"`
	} `json:"outs"`
	Status string `json:"status"`
}

type DaemonSendRawTransactionResult struct {
	Status string `json:"status"`
}

func check(e error) {
	if e != nil {
		fmt.Println(e.Error())
		return
	}
}

func loginCall(w http.ResponseWriter, r *http.Request) {
	type LoginRequest struct {
		WithCredentials bool   `json:"withCredentials"`
		Address         string `json:"address"`
		ViewKey         string `json:"view_key"`
		CreateAccount   bool   `json:"create_account"`
	}

	type LoginRespond struct {
		NewAddress bool `json:"new_address"`
	}

	var loginRespond LoginRespond

	decoder := json.NewDecoder(r.Body)
	var loginRequest LoginRequest
	err := decoder.Decode(&loginRequest)
	check(err)
	defer r.Body.Close()

	// Create address
	jsonStr := []byte(`{"method": "createAddress","params": {"spendPublicKey": "` + loginRequest.ViewKey + `"}}`)

	req, err := http.NewRequest("POST", walletRPC, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	check(err)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))

	var walletdCreateAddressResultSuccess WalletdCreateAddressResultSuccess
	err = json.Unmarshal(body, &walletdCreateAddressResultSuccess)
	check(err)

	if walletdCreateAddressResultSuccess.Result.Address != "" {
		loginRespond.NewAddress = true
	} else {
		loginRespond.NewAddress = false
	}

	loginRespondJSON, _ := json.Marshal(loginRespond)
	fmt.Fprintln(w, string(loginRespondJSON))
}

func getAddressInfoCall(w http.ResponseWriter, r *http.Request) {
	type AddressInfoRequest struct {
		Address string `json:"address"`
		ViewKey string `json:"view_key"`
	}

	type AddressInfoRespondSpentOutput struct {
		Amount   string `json:"amount"`
		KeyImage string `json:"key_image"`
		TxPubKey string `json:"tx_pub_key"`
		OutIndex int    `json:"out_index"`
		Mixin    int    `json:"mixin"`
	}

	type AddressInfoRespond struct {
		LockedFunds        string                          `json:"locked_funds"`
		TotalReceived      string                          `json:"total_received"`
		TotalSent          string                          `json:"total_sent"`
		ScannedHeight      int                             `json:"scanned_height"`
		ScannedBlockHeight int                             `json:"scanned_block_height"`
		StartHeight        int                             `json:"start_height"`
		TransactionHeight  int                             `json:"transaction_height"`
		BlockchainHeight   int                             `json:"blockchain_height"`
		SpentOutputs       []AddressInfoRespondSpentOutput `json:"spent_outputs,omitempty"`
	}

	var addressInfoRespond AddressInfoRespond

	decoder := json.NewDecoder(r.Body)
	var addressInfoRequest AddressInfoRequest
	err := decoder.Decode(&addressInfoRequest)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer r.Body.Close()

	// Get status
	jsonStr := []byte(`{"method": "getStatus","params": {}}`)
	req, err := http.NewRequest("POST", walletRPC, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	check(err)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	var walletdGetStatusResult WalletdGetStatusResult
	err = json.Unmarshal(body, &walletdGetStatusResult)
	check(err)

	addressInfoRespond.BlockchainHeight = walletdGetStatusResult.Result.KnownBlockCount
	addressInfoRespond.ScannedBlockHeight = walletdGetStatusResult.Result.BlockCount

	// Get balance
	jsonStr = []byte(`{"method": "getBalance","params": {"address": "` + addressInfoRequest.Address + `"}}`)
	req, err = http.NewRequest("POST", walletRPC, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client = &http.Client{}
	resp, err = client.Do(req)
	check(err)
	defer resp.Body.Close()
	body, _ = ioutil.ReadAll(resp.Body)

	var walletdGetBalanceResult WalletdGetBalanceResult
	err = json.Unmarshal(body, &walletdGetBalanceResult)
	check(err)

	addressInfoRespond.LockedFunds = strconv.Itoa(walletdGetBalanceResult.Result.LockedAmount)

	// Get SpentOutputs
	jsonStr = []byte(`{"method": "getTransactions","params": {"addresses": ["` + addressInfoRequest.Address + `"],"firstBlockIndex":0` + `,"blockCount": ` + strconv.Itoa(addressInfoRespond.BlockchainHeight) + `}}`)
	req, err = http.NewRequest("POST", walletRPC, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client = &http.Client{}
	resp, err = client.Do(req)
	check(err)
	defer resp.Body.Close()
	body, _ = ioutil.ReadAll(resp.Body)

	var walletdGetTransactionsResult WalletdGetTransactionsResult
	err = json.Unmarshal(body, &walletdGetTransactionsResult)
	check(err)

	var totalReceived int64 = 0
	var totalSent int64 = 0
	for _, v1 := range walletdGetTransactionsResult.Result.Items[:] {
		for _, v2 := range v1.Transactions[:] {
			for _, v3 := range v2.Transfers[:] {
				if v3.Address == addressInfoRequest.Address {
					totalReceived += v3.Amount
				}

				if v3.Amount < 0 {
					for _, v4 := range v3.SpentOutputs[:] {
						// todo: add real tx
						var spendOutput AddressInfoRespondSpentOutput
						spendOutput.Amount = strconv.FormatInt(v4.Amount, 10)
						spendOutput.KeyImage = v4.KeyImage
						spendOutput.TxPubKey = v4.TxPubKey
						spendOutput.OutIndex = int(v4.OutIndex)
						spendOutput.Mixin = int(v4.Mixin)

						addressInfoRespond.SpentOutputs = append(addressInfoRespond.SpentOutputs, spendOutput)
						totalSent += v4.Amount
					}
				}
			}
		}
	}
	addressInfoRespond.TotalReceived = strconv.FormatInt(totalReceived, 10)
	addressInfoRespond.TotalSent = strconv.FormatInt(totalSent, 10)

	addressInfoRespondJSON, _ := json.Marshal(addressInfoRespond)
	fmt.Fprintln(w, string(addressInfoRespondJSON))
}

func getAddressTxsCall(w http.ResponseWriter, r *http.Request) {
	type AddressTxRequest struct {
		Address string `json:"address"`
		ViewKey string `json:"view_key"`
	}

	type AddressTxRespondSpentOutput struct {
		Amount   string `json:"amount"`
		KeyImage string `json:"key_image"`
		TxPubKey string `json:"tx_pub_key"`
		OutIndex int    `json:"out_index"`
		Mixin    int    `json:"mixin"`
	}

	type AddressTxRespondTransaction struct {
		ID            int                           `json:"id"`
		Hash          string                        `json:"hash"`
		Timestamp     int                           `json:"timestamp"`
		TotalReceived string                        `json:"total_received"`
		TotalSent     string                        `json:"total_sent"`
		UnlockTime    int                           `json:"unlock_time"`
		Height        int                           `json:"height"`
		Coinbase      bool                          `json:"coinbase"`
		Mempool       bool                          `json:"mempool"`
		Mixin         int                           `json:"mixin"`
		SpentOutputs  []AddressTxRespondSpentOutput `json:"spent_outputs,omitempty"`
		PaymentID     string                        `json:"payment_id,omitempty"`
	}

	type AddressTxRespond struct {
		TotalReceived      string                        `json:"total_received"`
		ScannedHeight      int                           `json:"scanned_height"`
		ScannedBlockHeight int                           `json:"scanned_block_height"`
		StartHeight        int                           `json:"start_height"`
		TransactionHeight  int                           `json:"transaction_height"`
		Transactions       []AddressTxRespondTransaction `json:"transactions"`
		BlockchainHeight   int                           `json:"blockchain_height"`
	}

	var addressTxRespond AddressTxRespond

	decoder := json.NewDecoder(r.Body)
	var addressTxRequest AddressTxRequest
	err := decoder.Decode(&addressTxRequest)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer r.Body.Close()

	// Get status
	jsonStr := []byte(`{"method": "getStatus","params": {}}`)
	req, err := http.NewRequest("POST", walletRPC, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	check(err)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	var walletdGetStatusResult WalletdGetStatusResult
	err = json.Unmarshal(body, &walletdGetStatusResult)
	check(err)

	blockchainHeight := walletdGetStatusResult.Result.KnownBlockCount

	// Get transactions
	jsonStr = []byte(`{"method": "getTransactions","params": {"addresses": ["` + addressTxRequest.Address + `"],"firstBlockIndex":0` + `,"blockCount": ` + strconv.Itoa(blockchainHeight) + `}}`)
	req, err = http.NewRequest("POST", walletRPC, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client = &http.Client{}
	resp, err = client.Do(req)
	check(err)
	defer resp.Body.Close()
	body, _ = ioutil.ReadAll(resp.Body)

	var walletdGetTransactionsResult WalletdGetTransactionsResult
	err = json.Unmarshal(body, &walletdGetTransactionsResult)
	check(err)

	addressTxRespond.ScannedBlockHeight = walletdGetStatusResult.Result.BlockCount
	addressTxRespond.BlockchainHeight = walletdGetStatusResult.Result.KnownBlockCount
	addressTxRespond.StartHeight = 0

	var totalReceived int64

	for _, v1 := range walletdGetTransactionsResult.Result.Items[:] {
		for _, v2 := range v1.Transactions[:] {
			// todo: add real tx
			var txTotalReceived int64
			var txTotalSent int64
			var transaction AddressTxRespondTransaction

			transaction.Hash = v2.TransactionHash
			transaction.Timestamp = v2.Timestamp
			transaction.Height = v2.BlockIndex
			transaction.UnlockTime = v2.UnlockTime
			transaction.PaymentID = v2.PaymentID
			transaction.Coinbase = v2.IsBase
			transaction.Mixin = v2.Mixin
			transaction.Mempool = false
			if v2.Timestamp == 0 {
				transaction.Timestamp = 2147483647
				transaction.Mempool = true
			}

			for _, v3 := range v2.Transfers[:] {
				if v3.Address == addressTxRequest.Address {
					totalReceived += v3.Amount
					txTotalReceived += v3.Amount
				}

				if v3.Amount < 0 {
					for _, v4 := range v3.SpentOutputs[:] {
						var spendOutput AddressTxRespondSpentOutput
						spendOutput.Amount = strconv.FormatInt(v4.Amount, 10)
						spendOutput.KeyImage = v4.KeyImage
						spendOutput.TxPubKey = v4.TxPubKey
						spendOutput.OutIndex = int(v4.OutIndex)
						spendOutput.Mixin = int(v4.Mixin)

						transaction.SpentOutputs = append(transaction.SpentOutputs, spendOutput)
						txTotalSent += v4.Amount
					}
				}
			}

			transaction.TotalReceived = strconv.FormatInt(txTotalReceived, 10)
			transaction.TotalSent = strconv.FormatInt(txTotalSent, 10)
			addressTxRespond.Transactions = append(addressTxRespond.Transactions, transaction)
		}
	}

	addressTxRespond.TotalReceived = strconv.FormatInt(totalReceived, 10)

	addressTxRespondJSON, _ := json.Marshal(addressTxRespond)
	//REMOVE
	//fmt.Println(string(addressTxRespondJSON))

	fmt.Fprintln(w, string(addressTxRespondJSON))
}

func getUnspentOuts(w http.ResponseWriter, r *http.Request) {
	type UnspentOutsRequest struct {
		Address       string `json:"address"`
		ViewKey       string `json:"view_key"`
		Amount        string `json:"amount"`
		Mixin         int    `json:"mixin"`
		UseDust       bool   `json:"use_dust"`
		DustThreshold string `json:"dust_threshold"`
	}

	type UnspentOutsRespondOutput struct {
		Amount         string   `json:"amount"`
		PublicKey      string   `json:"public_key"`
		Index          int      `json:"index"`
		GlobalIndex    int      `json:"global_index"`
//		TxID           int      `json:"tx_id"`
		TxHash         string   `json:"tx_hash"`
		TxPubKey       string   `json:"tx_pub_key"`
//		TxPrefixHash   string   `json:"tx_prefix_hash"`
		SpendKeyImages []string `json:"spend_key_images"`
	}

	type UnspentOutsRespond struct {
		Amount  string `json:"amount"`
		Outputs       []UnspentOutsRespondOutput `json:"outputs"`
	}

	var unspentOutsRespond UnspentOutsRespond

	decoder := json.NewDecoder(r.Body)
	var unspentOutsRequest UnspentOutsRequest
	err := decoder.Decode(&unspentOutsRequest)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer r.Body.Close()

	m := make(map[string][]string)

	// Get status
	jsonStr := []byte(`{"method": "getStatus","params": {}}`)
	req, err := http.NewRequest("POST", walletRPC, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	check(err)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	var walletdGetStatusResult WalletdGetStatusResult
	err = json.Unmarshal(body, &walletdGetStatusResult)
	check(err)

	blockchainHeight := walletdGetStatusResult.Result.KnownBlockCount

	// Spent outputs
	jsonStr = []byte(`{"method": "getTransactions","params": {"addresses": ["` + unspentOutsRequest.Address + `"],"firstBlockIndex":0` + `,"blockCount": ` + strconv.Itoa(blockchainHeight) + `}}`)
	req, err = http.NewRequest("POST", walletRPC, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client = &http.Client{}
	resp, err = client.Do(req)
	check(err)
	defer resp.Body.Close()
	body, _ = ioutil.ReadAll(resp.Body)

	var walletdGetTransactionsResult WalletdGetTransactionsResult
	err = json.Unmarshal(body, &walletdGetTransactionsResult)
	check(err)

	for _, v1 := range walletdGetTransactionsResult.Result.Items[:] {
		for _, v2 := range v1.Transactions[:] {
			for _, v3 := range v2.Transfers[:] {
				if v3.Amount < 0 {
					for _, v4 := range v3.SpentOutputs[:] {
						m[v4.TxPubKey] = append(m[v4.TxPubKey], v4.KeyImage)
					}
				}
			}
		}
	}

	// Get UnspendOuts
	jsonStr = []byte(`{"method": "getUnspendOuts","params": {"address": "` + unspentOutsRequest.Address + `"}}`)
	req, err = http.NewRequest("POST", walletRPC, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client = &http.Client{}
	resp, err = client.Do(req)
	check(err)
	defer resp.Body.Close()
	body, _ = ioutil.ReadAll(resp.Body)

	var walletdGetUnspendOutsResult WalletdGetUnspendOutsResult
	err = json.Unmarshal(body, &walletdGetUnspendOutsResult)
	check(err)

	var amount int64 = 0
	for _, v := range walletdGetUnspendOutsResult.Result.Outputs[:] {
	  var Output UnspentOutsRespondOutput
	  Output.Amount = strconv.FormatInt(v.Amount, 10)

	  Output.PublicKey = v.OutputKey
	  Output.Index = v.OutputInTransaction
	  Output.GlobalIndex = v.GlobalOutputIndex
	  // Output.TxID = ???
		Output.TxHash = v.TransactionHash
    Output.TxPubKey = v.TransactionPublicKey
	  // Output.TxPrefixHash = ???

		if _, exists := m[v.TransactionPublicKey]; exists {
    	Output.SpendKeyImages = m[v.TransactionPublicKey]
		} else {
			Output.SpendKeyImages = []string{}
		}

    unspentOutsRespond.Outputs = append(unspentOutsRespond.Outputs, Output)

		amount = amount + v.Amount
  }

	unspentOutsRespond.Amount = strconv.FormatInt(amount, 10)
	unspentOutsRespondJSON, _ := json.Marshal(unspentOutsRespond)

	fmt.Fprintln(w, string(unspentOutsRespondJSON))
}

func getRandomOuts(w http.ResponseWriter, r *http.Request) {
	type RandomOutsRequest struct {
		Amounts []string `json:"amounts"`
		Count   int      `json:"count"`
	}

	type RandomOutsRespondAmountOutOutput struct {
		GlobalIndex string `json:"global_index"`
		PublicKey   string `json:"public_key"`
	}

	type RandomOutsRespondAmountOut struct {
		Amount  string `json:"amount"`
		Outputs []RandomOutsRespondAmountOutOutput `json:"outputs"`
	}

	type RandomOutsRespond struct {
		AmountOuts []RandomOutsRespondAmountOut `json:"amount_outs"`
	}

	var randomOutsRespond RandomOutsRespond

	decoder := json.NewDecoder(r.Body)
	var randomOutsRequest RandomOutsRequest
	err := decoder.Decode(&randomOutsRequest)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer r.Body.Close()

	// Get random outs
	amountsString := strings.Join(randomOutsRequest.Amounts, `, `)
	jsonStr := []byte(`{"amounts": [` + amountsString + `],"outs_count":` + strconv.FormatInt(int64(randomOutsRequest.Count), 10) + `}`)
	req, err := http.NewRequest("POST", daemonRPC + "getrandom_outs", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	check(err)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	var daemonGetRandomOutsResult DaemonGetRandomOutsResult
	err = json.Unmarshal(body, &daemonGetRandomOutsResult)
	check(err)

	for _, v1 := range daemonGetRandomOutsResult.Outs[:] {
	  var randomOutsRespondAmountOut RandomOutsRespondAmountOut
	  randomOutsRespondAmountOut.Amount = strconv.FormatInt(int64(v1.Amount), 10)

		for _, v2 := range v1.Outs[:] {
			var randomOutsRespondAmountOutOutput RandomOutsRespondAmountOutOutput
			randomOutsRespondAmountOutOutput.GlobalIndex = strconv.FormatInt(int64(v2.GlobalAmountIndex), 10)
			randomOutsRespondAmountOutOutput.PublicKey = v2.OutKey

			randomOutsRespondAmountOut.Outputs = append(randomOutsRespondAmountOut.Outputs, randomOutsRespondAmountOutOutput)
		}

		randomOutsRespond.AmountOuts = append(randomOutsRespond.AmountOuts, randomOutsRespondAmountOut)
  }

	randomOutsRespondJSON, _ := json.Marshal(randomOutsRespond)

	fmt.Fprintln(w, string(randomOutsRespondJSON))
}

func submitRawTx(w http.ResponseWriter, r *http.Request) {
	type SubmitRawTxRequest struct {
		Address       string `json:"address"`
		Tx    		    string `json:"tx"`
		ViewKey       string `json:"view_key"`
	}

	type SubmitRawTxRespond struct {
		Status     	  string `json:"status"`
	}

	var submitRawTxRespond SubmitRawTxRespond

	decoder := json.NewDecoder(r.Body)
	var submitRawTxRequest SubmitRawTxRequest
	err := decoder.Decode(&submitRawTxRequest)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer r.Body.Close()

	// Send transaction
	jsonStr := []byte(`{"tx_as_hex": "` + submitRawTxRequest.Tx + `"}`)
	req, err := http.NewRequest("POST", daemonRPC + "sendrawtransaction", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	check(err)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	var daemonSendRawTransactionResult DaemonSendRawTransactionResult
	err = json.Unmarshal(body, &daemonSendRawTransactionResult)
	check(err)

	submitRawTxRespond.Status = daemonSendRawTransactionResult.Status

	submitRawTxRespondJSON, _ := json.Marshal(submitRawTxRespond)

	fmt.Fprintln(w, string(submitRawTxRespondJSON))
}

func main() {
	crs := cors.New(cors.Options{
		AllowedHeaders:   []string{"*", "DNT", "X-CustomHeader", "Keep-Alive", "User-Agent", "X-Requested-With", "If-Modified-Since", "Cache-Control", "Content-Type", "Set-Cookie"},
		AllowCredentials: true,
		ExposedHeaders:   []string{"*", "DNT", "X-CustomHeader", "Keep-Alive", "User-Agent", "X-Requested-With", "If-Modified-Since", "Cache-Control", "Content-Type", "Set-Cookie"},
		MaxAge:           86400,
		Debug:            true,
	})

	http.Handle("/.well-known/acme-challenge/", http.FileServer(http.FileSystem(http.Dir("/var/tmp/letsencrypt/"))))

	http.Handle("/login", crs.Handler(http.HandlerFunc(loginCall)))
	http.Handle("/get_address_info", crs.Handler(http.HandlerFunc(getAddressInfoCall)))
	http.Handle("/get_address_txs", crs.Handler(http.HandlerFunc(getAddressTxsCall)))
	http.Handle("/get_unspent_outs", crs.Handler(http.HandlerFunc(getUnspentOuts)))
	http.Handle("/get_random_outs", crs.Handler(http.HandlerFunc(getRandomOuts)))
	http.Handle("/submit_raw_tx", crs.Handler(http.HandlerFunc(submitRawTx)))

	http.ListenAndServeTLS(":443", "/etc/letsencrypt/live/api.dashcoin.me/fullchain.pem", "/etc/letsencrypt/live/api.dashcoin.me/privkey.pem", nil)
}
