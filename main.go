package main

import (
	"github.com/nspcc-dev/neo-go/pkg/interop/runtime"
	"github.com/nspcc-dev/neo-go/pkg/interop/storage"
)

const prefix = "flags:"

// Main handles method dispatch
func Main(operation string, args []interface{}) interface{} {
	switch operation {
	case "flag":
		if len(args) != 2 {
			return false
		}
		return flagWallet(args[0].(string), args[1].(string))
	case "count":
		if len(args) != 1 {
			return 0
		}
		return getReportCount(args[0].(string))
	case "testEvents":
		return generateTestEvents()
	default:
		return "Unknown method"
	}
}

func flagWallet(wallet string, reporter string) bool {
	ctx := storage.GetContext()
	key := prefix + wallet

	// Get the stored value
	value := storage.Get(ctx, key)
	var count int
	if value != nil {
		// Use single-value type assertion
		count = bytesToInt(value.([]byte))
	}

	// Increment and store the new count
	newCount := count + 1
	storage.Put(ctx, key, intToBytes(newCount))

	// Get current timestamp
	timestamp := runtime.GetTime()

	// Emit WalletFlagged event
	runtime.Notify("WalletFlagged", wallet, reporter, timestamp, newCount)

	runtime.Log("Flagged " + wallet + " by " + reporter)
	return true
}

func getReportCount(wallet string) int {
	ctx := storage.GetContext()
	key := prefix + wallet

	value := storage.Get(ctx, key)
	if value == nil {
		return 0
	}

	// Use single-value type assertion
	return bytesToInt(value.([]byte))
}

// intToBytes converts an int to []byte (little endian)
func intToBytes(n int) []byte {
	return []byte{byte(n)} // Neo only supports small integers in this format
}

// bytesToInt converts []byte to int (little endian)
func bytesToInt(b []byte) int {
	if len(b) == 0 {
		return 0
	}
	return int(b[0]) // assuming single-byte value for simplicity
}

// generateTestEvents creates test events for end-to-end Kafka testing
func generateTestEvents() bool {
	// Test addresses for simulation
	testAddresses := []string{
		"NbTiM6h8r99kpRtb428XcsUk1TzKed2gTc",
		"NfgHwwTi3wHAS8aFAN243C5vGbkYDpqLHP",
		"NgaiKFjurmNmiRzDRQGs44yzByXuiruBej",
		"NhVwRpMi5MdgMNbg8dQgBzPX1SZhTyHp1t",
		"NiNmXL8FjEUEs1nfX9uHFBNaenxDHJtmuB",
	}

	testReporters := []string{
		"SecurityBot1",
		"FraudDetector",
		"CommunityMod",
		"AutoScanner",
		"AdminReview",
	}

	// Generate test flags for each address
	for i := 0; i < len(testAddresses); i++ {
		flagWallet(testAddresses[i], testReporters[i])

		// Add a small delay simulation by incrementing a counter
		for j := 0; j < 100; j++ {
			// Simple delay loop
		}
	}

	runtime.Log("Generated test events for 5 addresses")
	return true
}
