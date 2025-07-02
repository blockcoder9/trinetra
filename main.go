package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"
)

type Report struct {
	Reporter  string    `json:"reporter"`
	Timestamp time.Time `json:"timestamp"`
}

type FlagContract struct {
	reports map[string][]Report
	mu      sync.Mutex
}

func NewFlagContract() *FlagContract {
	return &FlagContract{
		reports: make(map[string][]Report),
	}
}

func (fc *FlagContract) FlagWallet(wallet string, reporter string) error {
	fc.mu.Lock()
	defer fc.mu.Unlock()

	report := Report{
		Reporter:  reporter,
		Timestamp: time.Now(),
	}
	fc.reports[wallet] = append(fc.reports[wallet], report)
	return nil
}

func (fc *FlagContract) GetReportCount(wallet string) int {
	fc.mu.Lock()
	defer fc.mu.Unlock()

	return len(fc.reports[wallet])
}

var contract = NewFlagContract()

func flagWalletHandler(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Wallet   string `json:"wallet"`
		Reporter string `json:"reporter"`
	}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	if data.Wallet == "" || data.Reporter == "" {
		http.Error(w, "Wallet and Reporter are required", http.StatusBadRequest)
		return
	}
	_ = contract.FlagWallet(data.Wallet, data.Reporter)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Wallet flagged successfully"))
}

func getReportCountHandler(w http.ResponseWriter, r *http.Request) {
	wallet := r.URL.Query().Get("wallet")
	if wallet == "" {
		http.Error(w, "Wallet address required", http.StatusBadRequest)
		return
	}
	count := contract.GetReportCount(wallet)
	resp := map[string]interface{}{
		"wallet":       wallet,
		"report_count": count,
	}
	json.NewEncoder(w).Encode(resp)
}

func main() {
	http.HandleFunc("/flag", flagWalletHandler)
	http.HandleFunc("/count", getReportCountHandler)
	log.Println("[Trinetra Mock Contract] Listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
