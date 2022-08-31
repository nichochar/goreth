Go and Ethereum exploration. It starts with https://goethereumbook.org/en/client-setup/

If running against local, we assume a ganache server is running:
```bash
ganache -m "hole symptom crater bring army industry link later fabric hotel asthma pupil"
```

To run the remote client, simply compile and run: `go run .`. If you want to run against a local ganache client, run `go run . -local`.


To work with contracts, the order of things, listed in more detail [here](https://www.notion.so/privy-io/Wallet-API-7269112327424b498bb51926ec92932c), requires you to:
```
solc --bin Storage.sol -o Storage.bin
solc --abi Storage.sol -o build
abigen --abi contracts/Storage.abi --pkg main --type Storage --out storage.go --bin contracts/Storage.bin

```
