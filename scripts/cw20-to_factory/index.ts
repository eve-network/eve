// Connect to Cosmwasm client as a queryier


import { CosmWasmClient, fromUtf8, GasPrice, Secp256k1HdWallet, SigningCosmWasmClient, toAscii } from "cosmwasm";

import * as fs from 'fs';

// axios
import axios from 'axios';

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
    const addresses: any = await axios.get('https://status.cosmos.directory/juno').then((response) => {
        const v = response.data?.rpc?.current;
        return Object.keys(v);
    });
    // console.log(addresses);

    // loop through all addresses, and save `await CosmWasmClient.connect(address);` to a list
    const clients: CosmWasmClient[] = await Promise.all(addresses.map(async (address: string) => {
        return await CosmWasmClient.connect(address);
    }));    
        

    // const data = await getAccountFromMnemonic(mnemonic, config.prefix);

    // const client = await SigningCosmWasmClient.connectWithSigner(`${process.env.RPC_NODE}`, data.wallet, config);
    // const client = await CosmWasmClient.connect(randomChain);
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

    const addrs = [
        'juno10ecsdy9ey7se73aprhgsuf2zrd6u355zqlh66c',
        'juno10wxn2lv29yqnw2uf4jf439kwy5ef00qdxzlw96',
        'juno12dywmje98gs9jq2xa9m3n9pkdxpyynkv802yzy',
        'juno130mdu9a0etmeuw52qfxk73pn0ga6gawk4k539x',
        'juno13v3cdp77pe24hpek3u886enyugtccqe3k289gk',
        'juno15t6ds4yxyky73kwk7ntymsm50gmef56aa2v9je',
        'juno15zwp7ksfzp9kk5705apu9t4tlky4pwh3xnn9ws',
        // 'juno16arfhunnguzvndwtzas28t8cpe20j7pt0pfjeh',
        // 'juno16mrjtqffn3awme2eczhlpwzj7mnatkeluvhj6c',
        // 'juno175q6smvgnuec5e62rs4chnu5cs8d98q2xgf4rx',
        // 'juno17py8gfneaam64vt9kaec0fseqwxvkq0flmsmhg',
        // 'juno18vqfvzeqcr4fmek864wl6lrumlxs05y4m5gdkv2gdpwautvsz9yq8t86kf',
        // 'juno19manc0mynrgu23zewq4zxdtdxhh0s298lu5wvq',
        // 'juno19wlk8gkfjckqr8d73dyp4n0f0k89q4h7s3j8fw',
        // 'juno19zeuukrexw5c87a828hw40qyd3lst9vlfwfwp8',
        // 'juno1an9w6uvs5el5393j9u6cwewqgs4uq4mus43gkp',
        // 'juno1cry9rwy5c4664685g6y06g9nanm7g28h0njc7f',
        // 'juno1e2e7hf2crkgmwu2c3ln9ycfrst4x54drz5jz3j',
        // 'juno1e3zzwtlvy0wlfef2ufc63j5pe6v8wr9gk6west',
        // 'juno1ejrhx56qy43xelteseq2z4ayhy74azyzyue7j5',
        // 'juno1gywj7qdsga3zv9h6vggxzezg56fxghguyl0qlr',
        // 'juno1gzzws3y22mjqk7pghwh8nasfgnpkdq6t2kuxsj',
        // 'juno1k0xlet79wsakwtwvf8lkjcfuds7ht9lhelsh2y',
        // 'juno1m7a7nva00p82xr0tssye052r8sxsxvcy2v5qz6',
        // 'juno1mk02dyn92xfwaqrfu8rmwgc6tt9xc0r2sh2qhd',
        // 'juno1mkwjmcya6329eyjkswlzeshaqsuc2m5qsx9qwa',
        // 'juno1nnarsp29ntj67fgzglmx48s4utf0ws387ggp4u',
        // 'juno1q2ke5n8rgr8ncjyeeqed3xyhyneluxgej3gt4m',
        // 'juno1q56z6th5umwg9787k4vn0p8f2wffcewask4f23',
        // 'juno1qhaysnsmv3q27f4xwmxhy4e5yjwjv0ax9jxdndtce593tnd7x2ysu56y0d',
        // 'juno1qnqgdcq83kpaa9qvajqgxllheppk4p5qnqgrr8',
        // 'juno1qt8qa27jzlfk8tsehp3fvltdnnn0u04g2ket3f',
        // 'juno1s4ckh9405q0a3jhkwx9wkf9hsjh66nmu769tz5',
        // 'juno1td5dacyexdd6xvzz2sskwd5t5tl7n3hfp7eqty',
        // 'juno1tec2pkmcgh30ljc9eelxxkyh6tet9pkgdvnwd5',
        // 'juno1uurel7z9ztjkruvedzdmlf9qtzy4s4yuk47kq4',
        // 'juno1v03v82syyy54ernxqa6zpmda9c2evgwjmsuav4',
        // 'juno1wgzqfntgygqum63vhkpttfpcrdeqtrg6ruxdj8',
        // 'juno1wmtep3pn8tjughhzpjag78fd9gj3n4z7ucfdy2',
        // 'juno1xn6u2z5qgl5agwee7688h6h7j9jcn5tftujdvj',
        // 'juno1xxpr2t8dmcxzcvwdy0g8n53kzxvwd0wwvmyzqd',
        // 'juno1y9h9q4dpmpssj2twer3s8dmd50nk0fc8huphma'
    ]

    // const data = await clients[0].queryContractSmart("juno14pgsyw2uwsjyvx5z5cwcnc89k8dtmz5xkaxw0cupglj3r6hjvlxsvcpghk", { balance: { address: "juno1wmtep3pn8tjughhzpjag78fd9gj3n4z7ucfdy2" } })
    // console.log(data.balance);

    // loop through addrs, then get the balance
    let balances: any = {}
    let queries: any = []
    for (let i = 0; i < addrs.length; i++) {
        // const data = await clients[0].queryContractSmart("juno14pgsyw2uwsjyvx5z5cwcnc89k8dtmz5xkaxw0cupglj3r6hjvlxsvcpghk", { balance: { address: addrs[i] } })
        // // console.log(data.balance);
        // balances.addrs[i] = data.balance
        // get random client
        const client = clients[Math.floor(Math.random() * clients.length)]
        queries.push({
            c: client.queryContractSmart("juno14pgsyw2uwsjyvx5z5cwcnc89k8dtmz5xkaxw0cupglj3r6hjvlxsvcpghk", { balance: { address: addrs[i] } }), 
            addr: addrs[i]}
        )
    }

    // const data = await a promise.all of every c
    console.log(balances);
    

    // process.exit(0);

    // let contract_addresses_to_query: string[] = []
    // let lastAddr = "";
    // for(let i = 0; i < 2; i++) {
    //     // get a random client
    //     const client = clients[Math.floor(Math.random() * clients.length)];

    //     let query: Record<string, unknown> = { all_accounts: { limit: 25 } }
    //     if(lastAddr.length > 0) {
    //         query = { all_accounts: { limit: 25, start_after: lastAddr } }
    //     }        
        

    //     let all_accounts = await client.queryContractSmart("juno14pgsyw2uwsjyvx5z5cwcnc89k8dtmz5xkaxw0cupglj3r6hjvlxsvcpghk",  query).catch((e) => {
    //         console.log(e);
    //     });                
    //     // console.log(all_accounts);
    //     contract_addresses_to_query.push(...all_accounts.accounts);
    //     lastAddr = all_accounts.accounts[all_accounts.accounts.length - 1];
    // }
    // // print contract_addresses_to_query
    // console.log(contract_addresses_to_query);
    
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