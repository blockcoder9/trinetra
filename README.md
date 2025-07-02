# Trinetra Flagging Contract

## Functions

- `FlagWallet(wallet)` → Flags a wallet address.
- `GetReportCount(wallet)` → Returns how many times the address has been flagged.

## Deployment (Blackhole Testnet)

- Contract Address: `tbd`
- RPC Endpoint: `https://testnet.blackhole.network:8545`

## Example

### Flag Wallet

```bash
curl -X POST http://localhost:8080/execute \
     -H "Content-Type: application/json" \
     -d '{"flag_wallet":{"wallet":"0xabc123"}}'
