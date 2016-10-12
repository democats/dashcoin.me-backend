package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "github.com/rs/cors"
    "io/ioutil"
    "log"
    "net/http"
)

var walletRPC = "http://localhost:9090"

type CreateAddressResult struct {
    Jsonrpc string `json:"jsonrpc"`
    ID      string `json:"id"`
    Result  struct {
        Address string `json:"address"`
    } `json:"result"`
}

type GetBalanceResult struct {
    Jsonrpc string `json:"jsonrpc"`
    ID      string `json:"id"`
    Result  struct {
        LockedAmount     int `json:"lockedAmount"`
        AvailableBalance int `json:"availableBalance"`
    } `json:"result"`
}

func check(e error) {
    if e != nil {
        log.Fatal(e)
        return
    }
}

func loginCall(w http.ResponseWriter, r *http.Request) {
    r.ParseForm()

    // Create address
    jsonStr := []byte(`{"method": "createAddress","params": {"spendPublicKey": "` + r.FormValue("spendPublicKey") + `"}}`)
    req, err := http.NewRequest("POST", walletRPC, bytes.NewBuffer(jsonStr))
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    check(err)
    defer resp.Body.Close()
    body, _ := ioutil.ReadAll(resp.Body)

    var car CreateAddressResult
    err = json.Unmarshal(body, &car)
    check(err)

    fmt.Fprintf(w, "Page viewed: %s", r.URL.Path[1:])
    //    fmt.Fprintf(w, "{\"new_address\": true}")
}

func getBalance(w http.ResponseWriter, r *http.Request) {
    r.ParseForm()

    // Get balance
    jsonStr := []byte(`{"method": "getBalance","params": {"address": "` + r.FormValue("address") + `"}}`)
    req, err := http.NewRequest("POST", walletRPC, bytes.NewBuffer(jsonStr))
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    check(err)
    defer resp.Body.Close()
    body, _ := ioutil.ReadAll(resp.Body)

    var gbr GetBalanceResult
    err = json.Unmarshal(body, &gbr)
    check(err)

    fmt.Fprintf(w, "Page viewed: %s", r.URL.Path[1:])
}

func getTransactions(w http.ResponseWriter, r *http.Request) {
    r.ParseForm()
    fmt.Fprintf(w, "Page viewed: %s", r.URL.Path[1:])
}

func getUnconfirmedTransactionHashes(w http.ResponseWriter, r *http.Request) {
    r.ParseForm()
    fmt.Fprintf(w, "Page viewed: %s", r.URL.Path[1:])
}

func main() {
    crs := cors.New(cors.Options{
        AllowedOrigins:   []string{"*"},
        AllowedHeaders:   []string{"*, DNT,X-CustomHeader,Keep-Alive,User-Agent,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type,Set-Cookie"},
        AllowCredentials: true,
        AllowedMethods:   []string{"POST"},
        ExposedHeaders:   []string{"POST"}})

    http.Handle("/login", crs.Handler(http.HandlerFunc(loginCall)))
    http.HandleFunc("/get_address_info", getBalance)
    http.HandleFunc("/get_address_txs", getTransactions)
    http.HandleFunc("/getUnconfirmedTransactionHashes/", getUnconfirmedTransactionHashes)

    http.ListenAndServeTLS(":443", "/etc/letsencrypt/live/api.dashcoin.me/fullchain.pem", "/etc/letsencrypt/live/api.dashcoin.me/privkey.pem", nil)
}
