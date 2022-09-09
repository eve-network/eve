// Connect to Cosmwasm client as a queryier


import { CosmWasmClient, fromUtf8, GasPrice, Secp256k1HdWallet, SigningCosmWasmClient, toAscii } from "cosmwasm";

import * as fs from 'fs';

// load dotenv
import * as dotenv from 'dotenv';
dotenv.config();

const config = {
    chainId: "juno-1",
    rpcEndpoint: `${process.env.RPC_NODE}`,
    prefix: "juno",
    gasPrice: GasPrice.fromString("0.03ucraft"),
};

// THIS IS A TEST MNUMONIC ONLY FOR TESTING PURPOSES WITH CRAFT, DO NOT USE WITH ACTUAL FUNDS. USED IN test_script.sh 
// EVER. FOR ANY REASON. DO ANYTHING FOR ANYONE, FOR ANY REASON, EVER, NO MATTER WHAT. NO MATTER WHERE.
// OR WHO, OR WHO YOU ARE WITH, OR WHERE YOU ARE GOING, OR WHERE YOU'VE BEEN, EVER, FOR ANY REASON WHATSOEVER.
const mnemonic = "decorate bright ozone fork gallery riot bus exhaust worth way bone indoor calm squirrel merry zero scheme cotton until shop any excess stage laundry";

// create async main function
async function main() {
    // GET A RANDOM ADDRESS TO QUERY FROM 

    // const data = await getAccountFromMnemonic(mnemonic, config.prefix);

    // const client = await SigningCosmWasmClient.connectWithSigner(`${process.env.RPC_NODE}`, data.wallet, config);
    const client = await CosmWasmClient.connect("https://rpc.juno.chaintools.tech");
    // let raw = await client.queryContractRaw(
    //     "juno14pgsyw2uwsjyvx5z5cwcnc89k8dtmz5xkaxw0cupglj3r6hjvlxsvcpghk",
    //     new Uint8Array([0, 10, ...toAscii("balance")])
    // )
    // JSON.parse(fromUtf8(raw))   

    // if(raw != null) {
    //     console.log(fromUtf8(raw));
    // } else {
    //     console.log("No data found");
    // }


    let all_accounts = (await client.queryContractSmart("juno14pgsyw2uwsjyvx5z5cwcnc89k8dtmz5xkaxw0cupglj3r6hjvlxsvcpghk", { all_accounts: { limit: 50 } })).accounts;
    console.log(all_accounts);
    

    // console.log("d");

    
    // save v to a file
    // get current dir that this file is in
    // const currentDir = process.cwd();
    // fs.writeFileSync(currentDir + '/balance.json', JSON.stringify(v));
}


// const getAccountFromMnemonic = async (mnemonic: any, prefix: string = "cosmos") => {
//     let wallet = await Secp256k1HdWallet.fromMnemonic(mnemonic, { prefix: prefix });
//     const [account] = await wallet.getAccounts();
//     return {
//         wallet: wallet,
//         account: account,
//     }
// }

main()