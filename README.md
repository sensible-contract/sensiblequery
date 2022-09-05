
# Bitcoin SV blockchain API service

For self-deployment please use [sensiblequery-deploy](https://github.com/sensible-contract/sensiblequery-deploy)。

Api Endpoint: `https://api.sensiblequery.com` for public.

Api V2 Endpoint: `https://api-v2.sensiblequery.com` with authorization.

See the supported API：`https://api.sensiblequery.com/swagger/index.html`

Testnet Api Endpoint: `https://api.sensiblequery.com/test/`

Testnet Supported API：`https://api.sensiblequery.com/test/swagger/index.html`

Block data service includes 2 components，Use Clickhouse as a data computing storage engine，Redis as data storage for UTXO collection of each address。

### 1. Node data synchronization program：Sensibled

Sensibled synchronizes block data by accessing the block folder of the full node (default in`~/.bitcoin/blocks/`)，The synchronized data is saved in Clickhouse，UTXO information is kept at Redis in order to support confirmed block data query。

At the same time, through the listening node zmq, real-time synchronization are done to get tx content and update to redis and clickhouse to support real-time queries of tx, balance and UTXO data.

### 2. Data API server：sensiblequery

Query the data in redis and clickhouse to provide data API services to the outside world.

## sensiblequery: Run dependencies

1. The node is required to provide rpc services. to encapsulate the pushtx interface.
2. You need to use the same redis instance, the same clickhouse instance, as the sensibled service. in order to obtain data.

## Profile

There are multiple profiles in the conf directory that are required for the program to run.

* db.yaml

Clickhouse database configuration, including adses, databases, etc.

* chain.yaml
*
Node configuration, rpc address.

* redis.yaml

Redis configuration, including ads, databases, etc.

Currently compatible with both redis cluster and stringle-node. The addrs configuration of a single address is treated as single-node.

## Run with Docker

It is easier to run sensiblequery with docker-compose. First set up the db/redis/node configuration, and then run:

	$ docker-compose up -d

Stop:

	$ docker-compose stop


## Run with the host

The environment variable LISTEN needs to be set to configure the listening port for the API service. Then start the program directly. The log is then output directly to the terminal.

    $ LISTEN=:5555 ./sensiblequery

You can use nohup or other techniques to place programs in the background to run.

The richquery service can be restarted at any time without any eventual data problems, except for interruptions to user access.

## Deployment resource requirements

| deploy               | DISK(minimum) | DISK(recommended) | MEM(minimum) | MEM(recommended) |
|----------------------|---------------|-------------------|--------------|------------------|
| sensiblequery        | 10 GB         | 20 GB             | 1 GB         | 4 GB             |
| bsv-node + sensibled | 512 GB        | 1000 GB           | 16 GB        | 32 GB            |
| clickhouse           | 1000 GB       | 1500 GB           | 16 GB        | 32 GB            |
| redis x 1            | 20 GB         | 50 GB             | 16 GB        | 32 GB            |
| pika x 1             | 20 GB         | 50 GB             | 4 GB         | 8 GB             |
| cache x 1            | 10 GB         | 20 GB             | 1 GB         | 2 GB             |

Where sensible is used to provide API services to the outside world, multiple instances can be deployed. the mindd is a single-instance run. Redis can deploy single nodes or clusters.
