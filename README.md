# RPCHub Aggregator

RPCHub is an RPC aggregator that offers you the fastest and most robust RPC services by integrating your owned nodes, private and public endpoints.

RPCHub is an open-source software that allows you to customize your own strategy configurations.

With customized configurations, RPCHub is able to achieve to make the RPC service to be scalable, exclusive, stable, low-cost, and high-performance.

Configurations are stored and utilized locally to best protect your privacy.


## Installation
Download the last release [here](https://github.com/BlockPILabs/aggregator/releases)

## Building from source
```shell
    git clone https://github.com/BlockPILabs/aggregator.git
    cd aggregator
    make
```
To start the aggregator, the following command can be run:
```shell
    build/aggregator
```
## Configuration
Default password is `123456`.

Visit https://ag-cfg.rpchub.io/ to configure the aggregator. See the [documents](https://docs.rpchub.io/) for more details.

Get current configuration by using command, replace `<password>` to what you set or using the default password:
```shell
curl -u rpchub:<password> 'http://localhost:8012/config'
```
To update the configuration, the following command can be run:
```shell
curl -u rpchub:<password> -X POST 'http://localhost:8012/config' --header 'Content-Type: application/json' --data-raw '{"password":"123456","request_timeout":30,"max_retries":3,"nodes":{"arbitrum":[{"name":"blockpi-public-arbitrum","endpoint":"https://arbitrum.blockpi.network/v1/rpc/public","weight":90,"read_only":false,"disabled":false},{"name":"arbitrum-official","endpoint":"https://arb1.arbitrum.io/rpc","weight":10,"read_only":false,"disabled":false}],"bsc":[{"name":"blockpi-public-bsc","endpoint":"https://bsc.blockpi.network/v1/rpc/public","weight":100,"read_only":false,"disabled":false}]},"phishing_db":["https://cfg.rpchub.io/agg/scam-addresses.json"],"phishing_db_update_interval":3600}'
```

## Reset configuration
1. Stop the aggregator. 
2. Delete the configuration directory `rm -rf $HOME/.rpchub/aggregator/`.
3. Restart the aggregator