import os
import json

'''
Ensures validators data is correct (% & the token amount)
'''

FOLDER="gentx"
if not os.path.exists(FOLDER):
    print('gentx folder not found')
    exit(1)

# get all files within the gentx folder
gentx_files = os.listdir(FOLDER)

invalids = ""
output = ""

for file in gentx_files:
    f = open('gentx/' + file, 'r')
    data = json.load(f)

    validatorData = data['body']['messages'][0]
    moniker = validatorData['description']['moniker']
    rate = float(validatorData['commission']['rate']) * 100
    valop = validatorData['validator_address']
    bond_value = validatorData['value']
    amt = int(bond_value["amount"])/1_000_000

    if bond_value['denom'] != 'ueve':
        invalids += f'[!] Invalid denomination for validator: {moniker} {bond_value["denom"]} \n'    
    elif amt > 1.0: # TODO is there a limit we are setting for validators gentxs?
        invalids += f'[!] Invalid ueve amount for validator: {moniker} {amt}\n'
    else:
        output += (f"{valop} {rate}%\texp: {amt}, {moniker.strip()}\n")    

print(output)
print(f"{invalids}")