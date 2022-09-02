Go and Ethereum exploration. It starts with https://goethereumbook.org/en/client-setup/

## Setup
If running locally, we assume a ganache server is running:
```bash
ganache -m "hole symptom crater bring army industry link later fabric hotel asthma pupil"
```

You will also need to use an infura API key, create a `secrets.json` file that looks like the following:
```
{
  "infura": "your API key here",
}
```


## Running the scripts
To run the remote client, simply compile and run: `go run .`. If you want to run against a local ganache client, run `go run . -local`.

You can also deploy the contract (if running locally, by passing `-deploy`). The full usage is below:
```
Usage of goreth:
  -deploy
    	a flag that, if passed, deploys the Storage contract
  -local
    	a flag that, if passed, runs the script locally
```


## Appendix: compiling solidity contracts
We have solidity contracts in the `contracts/` folder.
To work with contracts, you'll need to install `abigen` (a devtool packaged in `go-ethereum`), and solidity:
```
brew install ethereum
brew install solidity
```

Then, you can generate the ABI of a contract, compile it, and generate go bindings with the respective 3 commands:
```
solc --abi Store.sol
solc --bin Store.sol
abigen --bin=Store_sol_Store.bin --abi=Store_sol_Store.abi --pkg=store --out=Store.go
```
