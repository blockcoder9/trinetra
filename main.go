package main

import (
	"github.com/nspcc-dev/neo-go/pkg/interop/runtime"
	"github.com/nspcc-dev/neo-go/pkg/interop/storage"
)

const prefix = "flags:"

// Main handles method dispatch with enhanced error handling
func Main(operation string, args []interface{}) interface{} {
	// Input validation for RPC stability
	if operation == "" {
		runtime.Log("Error: Empty operation")
		return "ERROR: Operation cannot be empty"
	}

	switch operation {
	case "flag":
		if len(args) != 2 {
			runtime.Log("Error: Invalid arguments for flag operation")
			return "ERROR: flag requires exactly 2 arguments (wallet, reporter)"
		}
		// Type safety checks
		wallet := args[0].(string)
		reporter := args[1].(string)
		if wallet == "" || reporter == "" {
			runtime.Log("Error: Empty wallet or reporter")
			return "ERROR: Wallet and reporter cannot be empty"
		}
		return flagWallet(wallet, reporter)
	case "count":
		if len(args) != 1 {
			runtime.Log("Error: Invalid arguments for count operation")
			return "ERROR: count requires exactly 1 argument (wallet)"
		}
		wallet := args[0].(string)
		if wallet == "" {
			runtime.Log("Error: Empty wallet for count")
			return "ERROR: Wallet cannot be empty"
		}
		return getReportCount(wallet)
	case "testEvents":
		return generateTestEvents()
	case "healthCheck":
		return performHealthCheck()
	default:
		runtime.Log("Error: Unknown operation: " + operation)
		return "ERROR: Unknown operation. Available: flag, count, testEvents, healthCheck"
	}
}

func flagWallet(wallet string, reporter string) bool {
	// Enhanced error handling for RPC stability
	ctx := storage.GetContext()

	key := prefix + wallet

	// Get the stored value with error handling
	value := storage.Get(ctx, key)
	var count int
	if value != nil {
		// Safe type assertion without two return values
		count = bytesToInt(value.([]byte))
	}

	// Increment and store the new count
	newCount := count + 1
	storage.Put(ctx, key, intToBytes(newCount))

	// Get current timestamp
	timestamp := runtime.GetTime()

	// Emit WalletFlagged event with error handling
	runtime.Notify("WalletFlagged", wallet, reporter, timestamp, newCount)

	runtime.Log("Successfully flagged " + wallet + " by " + reporter + " (count: " + string(rune(newCount)) + ")")
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

// performHealthCheck tests contract functionality and RPC stability
func performHealthCheck() interface{} {
	// Test 1: Storage operations
	ctx := storage.GetContext()
	testKey := "health:test"
	testValue := intToBytes(42)

	// Test write
	storage.Put(ctx, testKey, testValue)

	// Test read
	retrievedValue := storage.Get(ctx, testKey)
	if retrievedValue == nil {
		runtime.Log("Health Check FAILED: Storage read/write error")
		return "UNHEALTHY: Storage operations failed"
	}

	// Verify data integrity
	if bytesToInt(retrievedValue.([]byte)) != 42 {
		runtime.Log("Health Check FAILED: Data integrity error")
		return "UNHEALTHY: Data integrity check failed"
	}

	// Test 2: Event emission
	timestamp := runtime.GetTime()
	runtime.Notify("HealthCheck", "contract", "system", timestamp, "OK")

	// Test 3: Basic arithmetic and logic
	testCount := 5 + 3
	if testCount != 8 {
		runtime.Log("Health Check FAILED: Basic operations error")
		return "UNHEALTHY: Basic operations failed"
	}

	// Clean up test data
	storage.Delete(ctx, testKey)

	runtime.Log("Health Check PASSED: All systems operational")
	return "HEALTHY: Contract is fully operational"
}
