# Trinetra Flag Contract - Event Testing Guide

## Updated Contract Features

### ✅ Added Event Emission
- **Event Name**: `WalletFlagged`
- **Parameters**: 
  - `wallet` (String): The flagged wallet address
  - `reporter` (String): Who reported the wallet
  - `timestamp` (Integer): Unix timestamp when flagged
  - `flagCount` (Integer): Total number of flags for this wallet

### ✅ New Methods Available

1. **flag(wallet, reporter)** - Flags a wallet and emits event
2. **count(wallet)** - Returns flag count for a wallet  
3. **testEvents()** - Generates test events for 5 addresses

## Test Addresses Used in testEvents()

1. `NbTiM6h8r99kpRtb428XcsUk1TzKed2gTc` - Reported by `SecurityBot1`
2. `NfgHwwTi3wHAS8aFAN243C5vGbkYDpqLHP` - Reported by `FraudDetector`
3. `NgaiKFjurmNmiRzDRQGs44yzByXuiruBej` - Reported by `CommunityMod`
4. `NhVwRpMi5MdgMNbg8dQgBzPX1SZhTyHp1t` - Reported by `AutoScanner`
5. `NiNmXL8FjEUEs1nfX9uHFBNaenxDHJtmuB` - Reported by `AdminReview`

## Deployment Commands

### Deploy to TestNet
```bash
neo-go contract deploy -i main.nef -manifest main.manifest.json --rpc-endpoint https://testnet1.neo.coz.io:443 --wallet my-wallet.json
```

### Test Event Generation (After Deployment)
```bash
# Generate test events
neo-go contract invokefunction -r https://testnet1.neo.coz.io:443 <CONTRACT_HASH> testEvents --wallet my-wallet.json

# Test individual flag
neo-go contract invokefunction -r https://testnet1.neo.coz.io:443 <CONTRACT_HASH> flag "NTestAddress123" "TestReporter" --wallet my-wallet.json

# Check flag count
neo-go contract invokefunction -r https://testnet1.neo.coz.io:443 <CONTRACT_HASH> count "NTestAddress123" --wallet my-wallet.json
```

## Event Monitoring

Events will be emitted with the following structure:
```json
{
  "eventname": "WalletFlagged",
  "contract": "<CONTRACT_HASH>",
  "state": {
    "type": "Array",
    "value": [
      {"type": "ByteString", "value": "<wallet_address>"},
      {"type": "ByteString", "value": "<reporter>"},
      {"type": "Integer", "value": "<timestamp>"},
      {"type": "Integer", "value": "<flag_count>"}
    ]
  }
}
```

## Next Steps

1. **Fund Wallet**: Get GAS from TestNet faucet for address `Nih11LWK1PEEBshTvfBNgvn3T8oFBrbnto`
2. **Deploy Contract**: Use the deployment command above
3. **Test Events**: Run `testEvents()` method to generate test data
4. **Monitor Events**: Set up Kafka consumer to capture WalletFlagged events
5. **Verify**: Confirm all 5 test addresses generate events properly
