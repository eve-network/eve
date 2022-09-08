import os
import json
from pathlib import Path

# TODO
# staking.params.unbondTime = 1814400s default


LAUNCH_TIME = "2022-09-08T23:00:00Z" # 20 = 3pm CST (4pm EST)
CHAIN_ID = "eve-tn1"
GENESIS_FILE=f"{Path.home()}/.eved/config/genesis.json" # Home Dir of the genesis
FOLDER = "gentx"

CUSTOM_GENESIS_ACCOUNT_VALUES = {
    # 100eve for both
    "eve1hj5fveer5cjtn4wd6wstzugjfdxzl0xpysfwwn": f"{100*1_000_000}ueve # reece, using key from ./test_node.sh",
    # "...": f"{100*1_000_000}ueve # Jacobs eve validator address",
}

def main():
    # outputDetails()
    # resetGenesisFile()
    createGenesisAccountsCommands()
    pass

def resetGenesisFile():    
    if not os.path.exists(GENESIS_FILE):
        print(f"# [!] Genesis file doesn't exist. Run `eved init test --chain-id {CHAIN_ID}`")
        return

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
        genesis["app_state"]['gov']["voting_params"]['voting_period'] = '300s' # 5 min = 172800s, 5 days mainet?

        genesis["app_state"]['slashing']['params']["signed_blocks_window"] = "10000"
        genesis["app_state"]['slashing']['params']["min_signed_per_window"] = '0.050000000000000000' # 5% @ 10,000
        genesis["app_state"]['slashing']['params']["slash_fraction_double_sign"] = '0.050000000000000000' # 5% if you SlashLikeMo
        genesis["app_state"]['slashing']['params']["slash_fraction_downtime"] = '0.00000000000000000' # no downtime slash like Juno

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
    gentx_files = os.listdir(FOLDER)
    # give validators their amounts in the genesis (1uexp & some craft)
    output = "# AUTO GENERATED FROM add-genesis-accounts.py\n"
    for file in gentx_files:
        f = open(FOLDER + "/" + file, 'r')
        data = json.load(f)

        validatorData = data['body']['messages'][0]
        moniker = validatorData['description']['moniker']
        val_addr = validatorData['delegator_address'] # craftxxxxx
        amt = validatorData['value']['amount']

        if val_addr not in CUSTOM_GENESIS_ACCOUNT_VALUES.keys():
            # print()
            output += f"eved add-genesis-account {val_addr} {amt}uexp,10000000000ucraft #{moniker}\n"
            continue # 
                
    for account in CUSTOM_GENESIS_ACCOUNT_VALUES:
        # print(f"craftd add-genesis-account {account} {CUSTOM_GENESIS_ACCOUNT_VALUES[account]}")
        output += f"eved add-genesis-account {account} {CUSTOM_GENESIS_ACCOUNT_VALUES[account]}\n"

    # save output to file in this directory
    current_dir = os.path.dirname(os.path.realpath(__file__))
    with open(os.path.join(current_dir, "run_these_genesis_balances.sh"), 'w') as f:
        f.write(output)

    print(f"# [!] COPY-PASTE-RUN THE \"sh run_these_genesis_balances.sh\" ABOVE TO CREATE THE GENESIS ACCOUNTS")
    print(f"# [!] THEN `eved collect-gentxs --gentx-dir gentx/`")
    print(f"# [!] THEN `eved validate-genesis`")
    print(f"# [!] THEN `code (LOCATION_OF_GENESIS_FILE), AND PUT ON MACHINES`")


if __name__ == "__main__":
    main()