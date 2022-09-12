### Run 2 Eve validators easily

Steps
```bash
# From the root of this project folder,
- sh docker/init.sh # performs: gentxs, collection, add account balances, genesis file & manupulation, 
- docker-compose up # mounts to the docker's testnet folder to use validators 1 and 2 unique dirs. Also uses ports +10
```