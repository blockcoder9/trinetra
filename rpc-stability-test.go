package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

// RPCEndpoint represents a Neo RPC endpoint
type RPCEndpoint struct {
	Name string
	URL  string
}

// RPCRequest represents a JSON-RPC request
type RPCRequest struct {
	JSONRPC string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	ID      int           `json:"id"`
}

// RPCResponse represents a JSON-RPC response
type RPCResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	Result  interface{} `json:"result,omitempty"`
	Error   *RPCError   `json:"error,omitempty"`
	ID      int         `json:"id"`
}

// RPCError represents an RPC error
type RPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// TestResult holds the results of stability tests
type TestResult struct {
	Endpoint     string
	Online       bool
	ResponseTime time.Duration
	RateLimit    bool
	ErrorHandled bool
	BlockHeight  int
	LastError    string
}

// Neo TestNet and MainNet RPC endpoints
var rpcEndpoints = []RPCEndpoint{
	{"TestNet-1", "https://testnet1.neo.coz.io:443"},
	{"TestNet-2", "https://testnet2.neo.coz.io:443"},
	{"TestNet-3", "https://testnet3.neo.coz.io:443"},
	{"TestNet-4", "https://testnet4.neo.coz.io:443"},
	{"TestNet-5", "https://testnet5.neo.coz.io:443"},
	{"MainNet-1", "https://mainnet1.neo.coz.io:443"},
	{"MainNet-2", "https://mainnet2.neo.coz.io:443"},
	{"MainNet-3", "https://mainnet3.neo.coz.io:443"},
}

func main() {
	fmt.Println("🔍 Neo RPC Stability Test Suite")
	fmt.Println("================================")
	fmt.Println()

	var results []TestResult
	var wg sync.WaitGroup

	// Test all endpoints concurrently
	resultsChan := make(chan TestResult, len(rpcEndpoints))

	for _, endpoint := range rpcEndpoints {
		wg.Add(1)
		go func(ep RPCEndpoint) {
			defer wg.Done()
			result := testEndpoint(ep)
			resultsChan <- result
		}(endpoint)
	}

	// Wait for all tests to complete
	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	// Collect results
	for result := range resultsChan {
		results = append(results, result)
	}

	// Display results
	displayResults(results)

	// Provide recommendations
	provideRecommendations(results)
}

func testEndpoint(endpoint RPCEndpoint) TestResult {
	fmt.Printf("Testing %s (%s)...\n", endpoint.Name, endpoint.URL)

	result := TestResult{
		Endpoint: endpoint.Name,
	}

	// Test 1: Basic connectivity and response time
	start := time.Now()
	blockHeight, err := getBlockCount(endpoint.URL)
	responseTime := time.Since(start)

	if err != nil {
		result.Online = false
		result.LastError = err.Error()
		fmt.Printf("❌ %s - Offline: %s\n", endpoint.Name, err.Error())
		return result
	}

	result.Online = true
	result.ResponseTime = responseTime
	result.BlockHeight = blockHeight
	fmt.Printf("✅ %s - Online (%dms, Block: %d)\n", endpoint.Name, responseTime.Milliseconds(), blockHeight)

	// Test 2: Rate limit handling
	result.RateLimit = testRateLimit(endpoint.URL)

	// Test 3: Error handling
	result.ErrorHandled = testErrorHandling(endpoint.URL)

	return result
}

func getBlockCount(rpcURL string) (int, error) {
	request := RPCRequest{
		JSONRPC: "2.0",
		Method:  "getblockcount",
		Params:  []interface{}{},
		ID:      1,
	}

	response, err := makeRPCCall(rpcURL, request)
	if err != nil {
		return 0, err
	}

	if response.Error != nil {
		return 0, fmt.Errorf("RPC error: %s", response.Error.Message)
	}

	blockCount, ok := response.Result.(float64)
	if !ok {
		return 0, fmt.Errorf("invalid response format")
	}

	return int(blockCount), nil
}

func testRateLimit(rpcURL string) bool {
	fmt.Printf("  🔄 Testing rate limits...\n")

	// Send 20 requests rapidly
	var wg sync.WaitGroup
	successCount := 0
	var mu sync.Mutex

	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := getBlockCount(rpcURL)
			mu.Lock()
			if err == nil {
				successCount++
			}
			mu.Unlock()
		}()
	}

	wg.Wait()

	// If we got responses to most requests, rate limiting is handled well
	rateLimitOK := successCount >= 15
	fmt.Printf("  📊 Rate limit test: %d/20 requests succeeded\n", successCount)

	return rateLimitOK
}

func testErrorHandling(rpcURL string) bool {
	fmt.Printf("  🧪 Testing error handling...\n")

	// Test 1: Invalid method
	invalidRequest := RPCRequest{
		JSONRPC: "2.0",
		Method:  "invalidmethod",
		Params:  []interface{}{},
		ID:      1,
	}

	response, err := makeRPCCall(rpcURL, invalidRequest)
	if err != nil {
		return false
	}

	// Should return an error response, not crash
	if response.Error == nil {
		return false
	}

	// Test 2: Invalid parameters
	invalidParamsRequest := RPCRequest{
		JSONRPC: "2.0",
		Method:  "getblock",
		Params:  []interface{}{"invalid_hash"},
		ID:      2,
	}

	response2, err := makeRPCCall(rpcURL, invalidParamsRequest)
	if err != nil {
		return false
	}

	// Should handle invalid parameters gracefully
	errorHandled := response2.Error != nil
	fmt.Printf("  ✅ Error handling: Proper error responses received\n")

	return errorHandled
}

func makeRPCCall(rpcURL string, request RPCRequest) (*RPCResponse, error) {
	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Post(rpcURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var rpcResponse RPCResponse
	err = json.Unmarshal(body, &rpcResponse)
	if err != nil {
		return nil, err
	}

	return &rpcResponse, nil
}

func displayResults(results []TestResult) {
	fmt.Println("\n📊 RPC Stability Test Results")
	fmt.Println("==============================")

	onlineCount := 0
	var fastestEndpoint TestResult
	var slowestEndpoint TestResult
	fastestTime := time.Hour
	slowestTime := time.Duration(0)

	for _, result := range results {
		status := "❌ Offline"
		if result.Online {
			status = "✅ Online"
			onlineCount++

			if result.ResponseTime < fastestTime {
				fastestTime = result.ResponseTime
				fastestEndpoint = result
			}
			if result.ResponseTime > slowestTime {
				slowestTime = result.ResponseTime
				slowestEndpoint = result
			}
		}

		rateLimitStatus := "❌"
		if result.RateLimit {
			rateLimitStatus = "✅"
		}

		errorHandlingStatus := "❌"
		if result.ErrorHandled {
			errorHandlingStatus = "✅"
		}

		fmt.Printf("%-12s | %s | %4dms | Rate: %s | Errors: %s | Block: %d\n",
			result.Endpoint,
			status,
			result.ResponseTime.Milliseconds(),
			rateLimitStatus,
			errorHandlingStatus,
			result.BlockHeight)

		if !result.Online && result.LastError != "" {
			fmt.Printf("             Error: %s\n", result.LastError)
		}
	}

	fmt.Printf("\n📈 Summary: %d/%d endpoints online (%.1f%%)\n",
		onlineCount, len(results), float64(onlineCount)/float64(len(results))*100)

	if onlineCount > 0 {
		fmt.Printf("⚡ Fastest: %s (%dms)\n", fastestEndpoint.Endpoint, fastestTime.Milliseconds())
		fmt.Printf("🐌 Slowest: %s (%dms)\n", slowestEndpoint.Endpoint, slowestTime.Milliseconds())
	}
}

func provideRecommendations(results []TestResult) {
	fmt.Println("\n💡 Recommendations")
	fmt.Println("==================")

	// Find best TestNet endpoint
	var bestTestNet TestResult
	bestTestNetTime := time.Hour

	// Find best MainNet endpoint
	var bestMainNet TestResult
	bestMainNetTime := time.Hour

	for _, result := range results {
		if !result.Online {
			continue
		}

		if result.Endpoint[:7] == "TestNet" && result.ResponseTime < bestTestNetTime {
			bestTestNetTime = result.ResponseTime
			bestTestNet = result
		}

		if result.Endpoint[:7] == "MainNet" && result.ResponseTime < bestMainNetTime {
			bestMainNetTime = result.ResponseTime
			bestMainNet = result
		}
	}

	if bestTestNet.Online {
		fmt.Printf("🧪 Best TestNet RPC: %s (%dms response time)\n",
			bestTestNet.Endpoint, bestTestNet.ResponseTime.Milliseconds())
		
		// Get the actual URL for the best endpoint
		for _, ep := range rpcEndpoints {
			if ep.Name == bestTestNet.Endpoint {
				fmt.Printf("   URL: %s\n", ep.URL)
				break
			}
		}
	}

	if bestMainNet.Online {
		fmt.Printf("🚀 Best MainNet RPC: %s (%dms response time)\n",
			bestMainNet.Endpoint, bestMainNet.ResponseTime.Milliseconds())
		
		// Get the actual URL for the best endpoint
		for _, ep := range rpcEndpoints {
			if ep.Name == bestMainNet.Endpoint {
				fmt.Printf("   URL: %s\n", ep.URL)
				break
			}
		}
	}

	fmt.Println("\n🔧 Implementation Tips:")
	fmt.Println("• Use the fastest endpoint for your primary RPC")
	fmt.Println("• Implement fallback to secondary endpoints")
	fmt.Println("• Add retry logic with exponential backoff")
	fmt.Println("• Monitor endpoint health in production")
	fmt.Println("• Cache responses when possible to reduce RPC calls")
}
