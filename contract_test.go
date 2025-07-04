package main

// import (
// 	"testing"
// )

// func TestFlagWalletAndGetReportCount(t *testing.T) {
// 	contract := NewFlagContract()

// 	// Test initial count
// 	wallet := "0xabc123"
// 	if count := contract.GetReportCount(wallet); count != 0 {
// 		t.Errorf("Expected 0 reports initially, got %d", count)
// 	}

// 	// Flag wallet once
// 	err := contract.FlagWallet(wallet, "reporter1")
// 	if err != nil {
// 		t.Fatalf("FlagWallet failed: %v", err)
// 	}

// 	// Check count after 1 flag
// 	if count := contract.GetReportCount(wallet); count != 1 {
// 		t.Errorf("Expected 1 report, got %d", count)
// 	}

// 	// Flag wallet multiple times
// 	_ = contract.FlagWallet(wallet, "reporter2")
// 	_ = contract.FlagWallet(wallet, "reporter3")

// 	// Final count check
// 	if count := contract.GetReportCount(wallet); count != 3 {
// 		t.Errorf("Expected 3 reports, got %d", count)
// 	}
// }

// func TestMultipleWallets(t *testing.T) {
// 	contract := NewFlagContract()

// 	_ = contract.FlagWallet("0xwallet1", "userA")
// 	_ = contract.FlagWallet("0xwallet2", "userB")
// 	_ = contract.FlagWallet("0xwallet1", "userC")

// 	if count := contract.GetReportCount("0xwallet1"); count != 2 {
// 		t.Errorf("Expected 2 reports for 0xwallet1, got %d", count)
// 	}
// 	if count := contract.GetReportCount("0xwallet2"); count != 1 {
// 		t.Errorf("Expected 1 report for 0xwallet2, got %d", count)
// 	}
// 	if count := contract.GetReportCount("0xwallet3"); count != 0 {
// 		t.Errorf("Expected 0 reports for 0xwallet3, got %d", count)
// 	}
// }
