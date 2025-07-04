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

	storage.Put(ctx, key, intToBytes(count+1))
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
