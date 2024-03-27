# CosmWasm contracts e2e tests
## Guideline

Run local node

```bash
bash scripts/run-node.sh
```

Download the cw20_base wasm file

```bash
bash scripts/wasm/get_contract.sh
```

Then test upload code

```bash
bash scripts/wasm/upload_code.sh
```

After upload, we can test instantiate the cw20_code

```bash
bash scripts/wasm/instantiate_cw20.sh
```
