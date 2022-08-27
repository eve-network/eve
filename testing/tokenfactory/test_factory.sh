eved tx tokenfactory create-denom reece "my denom" reece 6 https://reece.sh 100000 false --from eve1 --chain-id eve-t1 --yes -b sync

eved q tokenfactory list-denom

eved q tokenfactory show-denom reece


eved tx tokenfactory mint-and-send-tokens reece 9 eve1hj5fveer5cjtn4wd6wstzugjfdxzl0xpysfwwn --from eve1 --chain-id eve-t1 --yes -b sync
eved q bank balances eve1hj5fveer5cjtn4wd6wstzugjfdxzl0xpysfwwn --denom=reece