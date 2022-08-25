from operator import ge
import os
import json
from pathlib import Path

# cd networks/eve-wip

LAUNCH_TIME = "2022-12-30T01:00:00Z"
CHAIN_ID = "eve-v1"
GENESIS_FILE=f"{Path.home()}/.eve/config/genesis.json" # .eve in future
current_path = os.path.dirname(os.path.realpath(__file__))
FOLDER = "gentx" # local dir where we store the gentx JSON files

CUSTOM_GENESIS_ACCOUNT_VALUES = {
    # test_node account example
    "eve1hj5fveer5cjtn4wd6wstzugjfdxzl0xpysfwwn": "2000000ueve # note here", # useful to auto gen extra for non gentx accounts, or to give them more
}

def main():
    outputDetails()
    resetGenesisFile()
    createGenesisAccountsCommands()
    pass

def resetGenesisFile():
    # load genesis.json & remove all values for accounts & supply
    with open(GENESIS_FILE) as f:
        genesis = json.load(f)
        genesis["genesis_time"] = LAUNCH_TIME
        genesis["chain_id"] = str(CHAIN_ID)

        genesis["app_state"]['auth']["accounts"] = []
        genesis["app_state"]['bank']["balances"] = []
        genesis["app_state"]['bank']["supply"] = []        

        genesis["app_state"]['genutil']["gen_txs"] = []

        genesis["app_state"]['gov']["deposit_params"]['min_deposit'][0]['denom'] = 'ueve'
        genesis["app_state"]['gov']["voting_params"]['voting_period'] = '43200s' # 2 days = 172800s

        genesis["app_state"]['slashing']['params']["signed_blocks_window"] = "10000"
        genesis["app_state"]['slashing']['params']["min_signed_per_window"] = '0.050000000000000000' # 5% * 10,000
        genesis["app_state"]['slashing']['params']["slash_fraction_double_sign"] = '0.050000000000000000' # 5% if you SlashLikeMo
        genesis["app_state"]['slashing']['params']["slash_fraction_downtime"] = '0.01000000000000000' # 0.01% for downtime, like Juno


        genesis["app_state"]['staking']['params']["min_commission_rate"] = '0.10000000000000000' # 10% min commission
        genesis["app_state"]['distribution']['params']["community_tax"] = '0.20000000000000000' # 20% community tax

        genesis["app_state"]['staking']['params']["bond_denom"] = 'ueve' 
        genesis["app_state"]['crisis']['constant_fee']["denom"] = 'ueve' 

        genesis["app_state"]['mint']["minter"]["inflation"] = '0.150000000000000000' # 15% inflation
        genesis["app_state"]['mint']["params"]["mint_denom"] = 'ueve'     

    # save genesis.json
    with open(GENESIS_FILE, 'w') as f:
        json.dump(genesis, f, indent=4)
    print(f"# RESET: {GENESIS_FILE}\n")


def outputDetails() -> str:
    # get the seconds until LAUNCH_TIME
    launch_time = int(os.popen("date -d '" + LAUNCH_TIME + "' +%s").read())
    now = int(os.popen("date +%s").read())
    seconds_until_launch = launch_time - now

    # convert seconds_until_launch to hours, minutes, and seconds
    hours = seconds_until_launch // 3600
    minutes = (seconds_until_launch % 3600) // 60

    print(f"# {LAUNCH_TIME} ({hours}h {minutes}m) from now\n# {CHAIN_ID}\n# GenesisFile: {GENESIS_FILE}")



def createGenesisAccountsCommands():
    os.chdir(current_path)
    os.makedirs(FOLDER, exist_ok=True)
    gentx_files = os.listdir(FOLDER)
    # give validators their amounts in the genesis (1ueve, or more if provided in the custom dict)
    for file in gentx_files:
        f = open(FOLDER + "/" + file, 'r')
        data = json.load(f)

        validatorData = data['body']['messages'][0]
        moniker = validatorData['description']['moniker']
        val_addr = validatorData['delegator_address'] # evexxxxx
        amt = validatorData['value']['amount']

        if val_addr not in CUSTOM_GENESIS_ACCOUNT_VALUES.keys():
            print(f"eve add-genesis-account {val_addr} {amt}ueve #{moniker}")
            continue # 
                
    for account in CUSTOM_GENESIS_ACCOUNT_VALUES:
        print(f"eve add-genesis-account {account} {CUSTOM_GENESIS_ACCOUNT_VALUES[account]}")

    print(f"# [!] COPY-PASTE-RUN THE ABOVE TO CREATE THE GENESIS ACCOUNTS")
    print(f"# [!] THEN `eve collect-gentxs --gentx-dir gentx/`")
    print(f"# [!] THEN `eve validate-genesis`")
    print(f"# [!] THEN `code (LOCATION_OF_GENESIS_FILE), AND PUT ON MACHINES`")


if __name__ == "__main__":
    main()