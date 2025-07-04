# Trinetra Flag Contract

A Neo blockchain smart contract for wallet flagging and fraud detection with comprehensive RPC stability features.

## 🚀 Features

- **Wallet Flagging System** - Flag suspicious wallets with reporter tracking
- **Event Emission** - Real-time `WalletFlagged` events for integration
- **RPC Stability** - Enterprise-grade reliability and error handling
- **Health Monitoring** - Built-in contract health checks
- **Test Data Generation** - Automated test events for development

## 📋 Contract Methods

| Method | Parameters | Description |
|--------|------------|-------------|
| `flag` | `wallet`, `reporter` | Flag a wallet address |
| `count` | `wallet` | Get flag count for a wallet |
| `testEvents` | - | Generate 5 test flagging events |
| `healthCheck` | - | Test contract functionality |

## 🛠️ Quick Start

### Prerequisites
- [neo-go](https://github.com/nspcc-dev/neo-go) installed
- Neo wallet with GAS tokens
- Go 1.19+ (for RPC testing)

### 1. Clone Repository
```bash
git clone <repository-url>
cd trinetra-flag
```

### 2. Compile Contract
```bash
neo-go contract compile -i main.go -c main.yml -m main.manifest.json -o main.nef
```

### 3. Deploy to TestNet
```bash
neo-go contract deploy -i main.nef -manifest main.manifest.json \
  --rpc-endpoint https://testnet2.neo.coz.io:443 \
  --wallet my-wallet.json
```

## 📡 Usage Examples

### cURL Examples

#### Flag a Wallet
```bash
curl -X POST https://testnet2.neo.coz.io:443 \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "method": "invokefunction",
    "params": [
      "CONTRACT_HASH",
      "flag",
      [
        {"type": "String", "value": "NbTiM6h8r99kpRtb428XcsUk1TzKed2gTc"},
        {"type": "String", "value": "SecurityBot1"}
      ]
    ],
    "id": 1
  }'
```

#### Get Flag Count
```bash
curl -X POST https://testnet2.neo.coz.io:443 \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "method": "invokefunction",
    "params": [
      "CONTRACT_HASH",
      "count",
      [
        {"type": "String", "value": "NbTiM6h8r99kpRtb428XcsUk1TzKed2gTc"}
      ]
    ],
    "id": 1
  }'
```

#### Health Check
```bash
curl -X POST https://testnet2.neo.coz.io:443 \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "method": "invokefunction",
    "params": [
      "CONTRACT_HASH",
      "healthCheck",
      []
    ],
    "id": 1
  }'

### JavaScript Example

```javascript
import { rpc, wallet, tx } from '@cityofzion/neon-js';

const CONTRACT_HASH = 'YOUR_CONTRACT_HASH';
const RPC_ENDPOINT = 'https://testnet2.neo.coz.io:443';

class TrinetraFlagClient {
  constructor(contractHash, rpcEndpoint) {
    this.contractHash = contractHash;
    this.rpcClient = new rpc.RPCClient(rpcEndpoint);
  }

  // Flag a wallet
  async flagWallet(walletAddress, reporter, account) {
    const script = tx.createScript({
      scriptHash: this.contractHash,
      operation: 'flag',
      args: [
        { type: 'String', value: walletAddress },
        { type: 'String', value: reporter }
      ]
    });

    const transaction = new tx.Transaction({
      script: script,
      account: account,
      networkFee: 0,
      systemFee: 0
    });

    return await this.rpcClient.sendRawTransaction(transaction);
  }

  // Get flag count
  async getFlagCount(walletAddress) {
    const result = await this.rpcClient.invokeFunction(
      this.contractHash,
      'count',
      [{ type: 'String', value: walletAddress }]
    );

    return parseInt(result.stack[0].value);
  }

  // Generate test events
  async generateTestEvents(account) {
    const script = tx.createScript({
      scriptHash: this.contractHash,
      operation: 'testEvents',
      args: []
    });

    const transaction = new tx.Transaction({
      script: script,
      account: account
    });

    return await this.rpcClient.sendRawTransaction(transaction);
  }

  // Health check
  async healthCheck() {
    const result = await this.rpcClient.invokeFunction(
      this.contractHash,
      'healthCheck',
      []
    );

    return result.stack[0].value;
  }

  // Listen for WalletFlagged events
  async subscribeToEvents(callback) {
    const ws = new WebSocket('wss://testnet2.neo.coz.io:443/ws');

    ws.onopen = () => {
      ws.send(JSON.stringify({
        jsonrpc: '2.0',
        method: 'subscribe',
        params: ['notification_from_execution'],
        id: 1
      }));
    };

    ws.onmessage = (event) => {
      const data = JSON.parse(event.data);
      if (data.params && data.params[0].contract === this.contractHash) {
        const notification = data.params[0];
        if (notification.eventname === 'WalletFlagged') {
          callback({
            wallet: notification.state.value[0].value,
            reporter: notification.state.value[1].value,
            timestamp: parseInt(notification.state.value[2].value),
            flagCount: parseInt(notification.state.value[3].value)
          });
        }
      }
    };
  }
}

// Usage
const client = new TrinetraFlagClient(CONTRACT_HASH, RPC_ENDPOINT);

// Flag a wallet
await client.flagWallet('NbTiM6h8r99kpRtb428XcsUk1TzKed2gTc', 'SecurityBot1', myAccount);

// Get flag count
const count = await client.getFlagCount('NbTiM6h8r99kpRtb428XcsUk1TzKed2gTc');
console.log(`Wallet flagged ${count} times`);

// Listen for events
client.subscribeToEvents((event) => {
  console.log('Wallet flagged:', event);
});
```

### Go Example

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/nspcc-dev/neo-go/pkg/rpcclient"
    "github.com/nspcc-dev/neo-go/pkg/util"
    "github.com/nspcc-dev/neo-go/pkg/wallet"
)

const (
    ContractHash = "YOUR_CONTRACT_HASH"
    RPCEndpoint  = "https://testnet2.neo.coz.io:443"
)

type TrinetraClient struct {
    client   *rpcclient.Client
    contract util.Uint160
    wallet   *wallet.Wallet
    account  *wallet.Account
}

func NewTrinetraClient(contractHash, rpcEndpoint, walletPath string) (*TrinetraClient, error) {
    // Connect to RPC
    client, err := rpcclient.New(context.Background(), rpcEndpoint, rpcclient.Options{})
    if err != nil {
        return nil, err
    }

    // Load contract hash
    contract, err := util.Uint160DecodeStringLE(contractHash)
    if err != nil {
        return nil, err
    }

    // Load wallet
    w, err := wallet.NewWalletFromFile(walletPath)
    if err != nil {
        return nil, err
    }

    return &TrinetraClient{
        client:   client,
        contract: contract,
        wallet:   w,
        account:  w.Accounts[0],
    }, nil
}

// Flag a wallet
func (tc *TrinetraClient) FlagWallet(walletAddr, reporter string) error {
    result, err := tc.client.InvokeFunction(
        tc.contract,
        "flag",
        []interface{}{walletAddr, reporter},
        nil,
    )
    if err != nil {
        return err
    }

    fmt.Printf("Flag result: %v\n", result)
    return nil
}

// Get flag count
func (tc *TrinetraClient) GetFlagCount(walletAddr string) (int, error) {
    result, err := tc.client.InvokeFunction(
        tc.contract,
        "count",
        []interface{}{walletAddr},
        nil,
    )
    if err != nil {
        return 0, err
    }

    if len(result.Stack) > 0 {
        return int(result.Stack[0].Value().(int64)), nil
    }
    return 0, fmt.Errorf("no result returned")
}

// Generate test events
func (tc *TrinetraClient) GenerateTestEvents() error {
    result, err := tc.client.InvokeFunction(
        tc.contract,
        "testEvents",
        []interface{}{},
        nil,
    )
    if err != nil {
        return err
    }

    fmt.Printf("Test events result: %v\n", result)
    return nil
}

// Health check
func (tc *TrinetraClient) HealthCheck() (string, error) {
    result, err := tc.client.InvokeFunction(
        tc.contract,
        "healthCheck",
        []interface{}{},
        nil,
    )
    if err != nil {
        return "", err
    }

    if len(result.Stack) > 0 {
        return result.Stack[0].Value().(string), nil
    }
    return "", fmt.Errorf("no result returned")
}

func main() {
    client, err := NewTrinetraClient(ContractHash, RPCEndpoint, "my-wallet.json")
    if err != nil {
        log.Fatal(err)
    }

    // Flag a wallet
    err = client.FlagWallet("NbTiM6h8r99kpRtb428XcsUk1TzKed2gTc", "SecurityBot1")
    if err != nil {
        log.Printf("Error flagging wallet: %v", err)
    }

    // Get flag count
    count, err := client.GetFlagCount("NbTiM6h8r99kpRtb428XcsUk1TzKed2gTc")
    if err != nil {
        log.Printf("Error getting count: %v", err)
    } else {
        fmt.Printf("Wallet flagged %d times\n", count)
    }

    // Health check
    health, err := client.HealthCheck()
    if err != nil {
        log.Printf("Error checking health: %v", err)
    } else {
        fmt.Printf("Contract health: %s\n", health)
    }
}
```

## 🔧 RPC Stability Testing

Test endpoint reliability and performance:

```bash
# Run comprehensive RPC tests
go run rpc-stability-test.go

# Or use batch file (Windows)
run-rpc-tests.bat
```

**Recommended Endpoints:**
- **TestNet**: `https://testnet2.neo.coz.io:443` (561ms)
- **MainNet**: `https://mainnet1.neo.coz.io:443` (614ms)

## 📊 Events

### WalletFlagged Event
```json
{
  "eventname": "WalletFlagged",
  "contract": "CONTRACT_HASH",
  "state": {
    "type": "Array",
    "value": [
      {"type": "ByteString", "value": "wallet_address"},
      {"type": "ByteString", "value": "reporter_name"},
      {"type": "Integer", "value": "timestamp"},
      {"type": "Integer", "value": "flag_count"}
    ]
  }
}
```

### HealthCheck Event
```json
{
  "eventname": "HealthCheck",
  "contract": "CONTRACT_HASH",
  "state": {
    "type": "Array",
    "value": [
      {"type": "ByteString", "value": "contract"},
      {"type": "ByteString", "value": "system"},
      {"type": "Integer", "value": "timestamp"},
      {"type": "ByteString", "value": "OK"}
    ]
  }
}
```

## 🧪 Testing

### Test Data Addresses
The `testEvents()` method flags these addresses:
1. `NbTiM6h8r99kpRtb428XcsUk1TzKed2gTc` - SecurityBot1
2. `NfgHwwTi3wHAS8aFAN243C5vGbkYDpqLHP` - FraudDetector
3. `NgaiKFjurmNmiRzDRQGs44yzByXuiruBej` - CommunityMod
4. `NhVwRpMi5MdgMNbg8dQgBzPX1SZhTyHp1t` - AutoScanner
5. `NiNmXL8FjEUEs1nfX9uHFBNaenxDHJtmuB` - AdminReview

## 📁 Project Structure

```
trinetra-flag/
├── main.go                    # Smart contract source
├── main.yml                   # Contract configuration
├── main.nef                   # Compiled contract
├── main.manifest.json         # Contract manifest
├── my-wallet.json            # Neo wallet
├── rpc-stability-test.go     # RPC testing suite
├── run-rpc-tests.bat         # Test runner script
├── test-events.md            # Event testing guide
├── RPC-STABILITY-GUIDE.md    # RPC implementation guide
└── README.md                 # This file
```

## 🚀 Deployment

### TestNet Deployment
```bash
# Get GAS from faucet first
# https://neowish.ngd.network/

neo-go contract deploy -i main.nef -manifest main.manifest.json \
  --rpc-endpoint https://testnet2.neo.coz.io:443 \
  --wallet my-wallet.json
```

### MainNet Deployment
```bash
# Ensure you have real GAS tokens
neo-go contract deploy -i main.nef -manifest main.manifest.json \
  --rpc-endpoint https://mainnet1.neo.coz.io:443 \
  --wallet my-wallet.json
```

## 🔒 Security Features

- **Input validation** for all parameters
- **Type safety** checks
- **Error handling** with descriptive messages
- **Storage integrity** verification
- **Health monitoring** capabilities

## 📞 Support

For integration issues or questions:
- Check the RPC Stability Guide
- Run health checks on deployed contract
- Test with provided examples
- Review event emission logs

## 📄 License

[Add your license here]

---

**Built with Neo N3 blockchain technology** 🚀
```
