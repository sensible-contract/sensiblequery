
## Bitcoin SV区块浏览器 API接口

### SatoBlock

我们部署了一个浏览器Demo [BSV Browser](http://bl.ocks.org/jiedo/raw/1e3fa74ca3157a16b7133708212f3193/#/blocks) ，可测试查看Blockchain的数据。

Api Endpoint: `http://120.92.153.221:5555/`

支持的API如下：

## GET /blockchain/info

可获知最新区块位置、同步状态等信息。(未完成)

#### Response

data包括字段为：

- chain: main/testnet
- blocks: 最新区块总数
- headers: 最新区块头总数
- bestBlockHash: 最新blockId
- medianTime: 最新区块时间戳
- github: 项目地址

- Body
```
{
  code: 0,
  msg: "ok",
  data: {
    chain: "main",
    blocks: 674936,
    headers: 674936,
    bestBlockHash: "000000000000000007d4c4d5da6f35a29d4a4ed9eba61e9627634d910f93b2ee",
    difficulty: "",
    medianTime: 1613626720,
    chainwork: ""
  }
}
```


## GET /blocks/{start}/{end} 获取指定高度范围内的区块概述列表

#### Request

`GET /blocks/0/10`

#### Response

返回值 code == 0 为正确，其他都是错误

data包括字段为：

- height: 当前区块高度
- id: 当前区块ID
- prev: 前一个区块ID
- ntx: 当前区块包括的交易条数

- Body
```
{
  code: 0,
  msg: "ok",
  data: [
    {
      height: 0,
      id: "000000000019d6689c085ae165831e934ff763ae46a2a6c172b3f1b60a8ce26f",
      prev: "0000000000000000000000000000000000000000000000000000000000000000",
      ntx: 1
    },
    {
      height: 1,
      id: "00000000839a8e6886ab5951d76f411475428afc90947ee320161bbf18eb6048",
      prev: "000000000019d6689c085ae165831e934ff763ae46a2a6c172b3f1b60a8ce26f",
      ntx: 1
    },
    {
      height: 2,
      id: "000000006a625f06636b8bb6ac7b960a8d03705d1ace08b1a19da3fdcc99ddbd",
      prev: "00000000839a8e6886ab5951d76f411475428afc90947ee320161bbf18eb6048",
      ntx: 1
    },
    {
      height: 3,
      id: "0000000082b5015589a3fdf2d4baff403e6f0be035a5d9742c1cae6295464449",
      prev: "000000006a625f06636b8bb6ac7b960a8d03705d1ace08b1a19da3fdcc99ddbd",
      ntx: 1
    }
  ]
}
```


## GET /block/id/{blkid} 通过区块blkid获取区块概述

#### Request

`GET /block/id/0000000082b5015589a3fdf2d4baff403e6f0be035a5d9742c1cae6295464449`

#### Response

返回值 code == 0 为正确，其他都是错误

data包括字段为：

- height: 当前区块高度
- id: 当前区块ID
- prev: 前一个区块ID
- ntx: 当前区块包括的交易条数

- Body
```
{
  code: 0,
  msg: "ok",
  data: {
    height: 3,
    id: "0000000082b5015589a3fdf2d4baff403e6f0be035a5d9742c1cae6295464449",
    prev: "000000006a625f06636b8bb6ac7b960a8d03705d1ace08b1a19da3fdcc99ddbd",
    ntx: 1
  }
}
```


## GET /block/txs/{blkid} 通过区块blkid获取区块包含的transaction概述列表

#### Request

`GET /block/txs/0000000082b5015589a3fdf2d4baff403e6f0be035a5d9742c1cae6295464449`

#### Response

返回值 code == 0 为正确，其他都是错误

data包括字段为：

- txid: 当前txid
- nIn: 当前交易包括的输入数量
- nOut: 当前交易包括的输出数量
- height: 当前交易被打包的区块高度
- blkid: 当前交易被打包的区块ID
- idx: 当前交易在区块中的序号

- Body
```
{
  code: 0,
  msg: "ok",
  data: [
    {
      txid: "999e1c837c76a1b7fbb7e57baf87b309960f5ffefbf2a9b95dd890602272f644",
      nIn: 1,
      nOut: 1,
      height: 3,
      blkid: "",
      idx: 0
    }
  ]
}
```

## GET /tx/{txid} 通过交易txid获取交易概述

#### Request

`GET /tx/999e1c837c76a1b7fbb7e57baf87b309960f5ffefbf2a9b95dd890602272f644`

#### Response

返回值 code == 0 为正确，其他都是错误

data包括字段为：

- txid: 当前txid
- nIn: 当前交易包括的输入数量
- nOut: 当前交易包括的输出数量
- height: 当前交易被打包的区块高度
- blkid: 当前交易被打包的区块ID
- idx: 当前交易在区块中的序号

- Body
```
{
  code: 0,
  msg: "ok",
  data: {
    txid: "999e1c837c76a1b7fbb7e57baf87b309960f5ffefbf2a9b95dd890602272f644",
    nIn: 1,
    nOut: 1,
    height: 3,
    blkid: "0000000082b5015589a3fdf2d4baff403e6f0be035a5d9742c1cae6295464449",
    idx: 0
  }
}
```

## GET /tx/{txid}/ins 通过交易txid获取交易所有输入信息列表

#### Request

`GET /tx/999e1c837c76a1b7fbb7e57baf87b309960f5ffefbf2a9b95dd890602272f644/ins`

#### Response

返回值 code == 0 为正确，其他都是错误

data包括字段为：

- height: 当前交易被打包的区块高度
- txid: 当前txid
- idx: 当前输入序号
- script_sig: 当前输入解锁脚本类型
- height_txo: 当前输入花费的utxo所属的区块高度，如果为0则未花费
- utxid: 当前输入花费的outpoint的txid
- vout: 当前输入花费的outpoint的index
- address: 当前输入花费的outpoint的address
- genesis: 当前输入花费的outpoint的genesis
- satoshi: 当前输入花费的outpoint的satoshi
- script_type: 当前输入锁定脚本类型

- Body
```
{
  code: 0,
  msg: "ok",
  data: [
    {
      height: 3,
      txid: "999e1c837c76a1b7fbb7e57baf87b309960f5ffefbf2a9b95dd890602272f644",
      idx: 0,
      script_sig: "04ffff001d010e",
      height_txo: 0,
      utxid: "0000000000000000000000000000000000000000000000000000000000000000",
      vout: 4294967295,
      address: "1111111111111111111114oLvT2",
      genesis: "0000000000000000000000000000000000000000",
      satoshi: 0,
      script_type: ""
    }
  ]
}
```

## GET /tx/{txid}/outs 通过交易txid获取交易所有输出信息列表

#### Request

`GET /tx/999e1c837c76a1b7fbb7e57baf87b309960f5ffefbf2a9b95dd890602272f644/outs`

#### Response

返回值 code == 0 为正确，其他都是错误

data包括字段为：

- txid: 当前txid
- vout: 当前输出序号
- address: 当前输出的address
- genesis: 当前输出的genesis
- satoshi: 当前输出的satoshi
- script_type: 当前输出锁定脚本类型
- script: 当前输出锁定脚本
- height: 当前交易被打包的区块高度
- txid_spent: 当前输出被花费的txid
- height_spent: 当前输出被花费的区块高度，如果为0则未花费

- Body
```
{
  code: 0,
  msg: "ok",
  data: [
    {
      txid: "999e1c837c76a1b7fbb7e57baf87b309960f5ffefbf2a9b95dd890602272f644",
      vout: 0,
      address: "1111111111111111111114oLvT2",
      genesis: "0000000000000000000000000000000000000000",
      satoshi: 5000000000,
      script_type: "41ac",
      script: "410494b9d3e76c5b1629ecf97fff95d7a4bbdac87cc26099ada28066c6ff1eb9191223cd897194a08d0c2726c5747f1db49e8cf90e75dc3e3550ae9b30086f3cd5aaac",
      height: 3,
      txid_spent: "0000000000000000000000000000000000000000000000000000000000000000",
      height_spent: 0
    }
  ]
}
```

## GET /tx/{txid}/in/{index} 通过交易txid和index获取指定交易输入信息

#### Request

`GET /tx/999e1c837c76a1b7fbb7e57baf87b309960f5ffefbf2a9b95dd890602272f644/in/0`

#### Response

返回值 code == 0 为正确，其他都是错误

data包括字段为：

- height: 当前交易被打包的区块高度
- txid: 当前txid
- idx: 当前输入序号
- script_sig: 当前输入解锁脚本类型
- height_txo: 当前输入花费的utxo所属的区块高度，如果为0则未花费
- utxid: 当前输入花费的outpoint的txid
- vout: 当前输入花费的outpoint的index
- address: 当前输入花费的outpoint的address
- genesis: 当前输入花费的outpoint的genesis
- satoshi: 当前输入花费的outpoint的satoshi
- script_type: 当前输入锁定脚本类型

- Body
```
{
  code: 0,
  msg: "ok",
  data: {
    height: 3,
    txid: "999e1c837c76a1b7fbb7e57baf87b309960f5ffefbf2a9b95dd890602272f644",
    idx: 0,
    script_sig: "04ffff001d010e",
    height_txo: 0,
    utxid: "0000000000000000000000000000000000000000000000000000000000000000",
    vout: 4294967295,
    address: "1111111111111111111114oLvT2",
    genesis: "0000000000000000000000000000000000000000",
    satoshi: 0,
    script_type: ""
  }
}
```

## GET /tx/{txid}/out/{index} 通过交易txid和index获取指定交易输出信息

#### Request

`GET /tx/999e1c837c76a1b7fbb7e57baf87b309960f5ffefbf2a9b95dd890602272f644/out/0`

#### Response

返回值 code == 0 为正确，其他都是错误

data包括字段为：

- txid: 当前txid
- vout: 当前输出序号
- address: 当前输出的address
- genesis: 当前输出的genesis
- satoshi: 当前输出的satoshi
- script_type: 当前输出锁定脚本类型
- script: 当前输出锁定脚本
- height: 当前交易被打包的区块高度
- txid_spent: 当前输出被花费的txid
- height_spent: 当前输出被花费的区块高度，如果为0则未花费

- Body
```
{
  code: 0,
  msg: "ok",
  data: {
    txid: "999e1c837c76a1b7fbb7e57baf87b309960f5ffefbf2a9b95dd890602272f644",
    vout: 0,
    address: "1111111111111111111114oLvT2",
    genesis: "0000000000000000000000000000000000000000",
    satoshi: 5000000000,
    script_type: "41ac",
    script: "410494b9d3e76c5b1629ecf97fff95d7a4bbdac87cc26099ada28066c6ff1eb9191223cd897194a08d0c2726c5747f1db49e8cf90e75dc3e3550ae9b30086f3cd5aaac",
    height: 3
  }
}
```

## GET /tx/{txid}/out/{index}/spent 通过交易txid和index获取指定交易输出是否被花费状态

#### Request

`GET /tx/0437cd7f8525ceed2324359c2d0ba26006d92d856a9c20fa0241106ee5a597c9/out/0/spent`

#### Response

返回值 code == 0 则已被花费，并返回花费结果。-1 则未花费

data包括字段为：

- height: 输出被花费的区块高度
- txid: 输出被花费的txid
- idx: 输出被花费的txid所在区块内序号
- utxid: 输出txid参数
- vout: 输出index参数

- Body
```
{
  code: 0,
  msg: "ok",
  data: {
    height: 170,
    txid: "f4184fc596403b9d638783cf57adfe4c75c605f6356fbc91338530e9831e9e16",
    idx: 0,
    utxid: "0437cd7f8525ceed2324359c2d0ba26006d92d856a9c20fa0241106ee5a597c9",
    vout: 0
  }
}
```


## GET /address/{address}/utxo 通过地址address获取相关utxo列表

#### Request

`GET /address/17SkEw2md5avVNyYgj6RiXuQKNwkXaxFyQ/utxo`

#### Response

返回值 code == 0 为正确，其他都是错误

data包括字段为：

- txid: 当前txid
- vout: 当前输出序号
- address: 当前输出的address
- genesis: 当前输出的genesis
- satoshi: 当前输出的satoshi
- script_type: 当前输出锁定脚本类型
- script: 当前输出锁定脚本
- height: 当前交易被打包的区块高度

- Body
```
{
  code: 0,
  msg: "ok",
  data: [
    {
      txid: "62104aa084f4a158cb9aa545ee30d68db88bb22d4a66904b78d41e4512c1969a",
      vout: 0,
      address: "17SkEw2md5avVNyYgj6RiXuQKNwkXaxFyQ",
      genesis: "0000000000000000000000000000000000000000",
      satoshi: 100000,
      script_type: "76a91488ac",
      script: "76a91446af3fb481837fadbb421727f9959c2d32a3682988ac",
      height: 474237
    },
    {
      txid: "2c63ac6d71e696dea43ef1ef7fba8c376a6a220383e73a17ba6c3795996db112",
      vout: 0,
      address: "17SkEw2md5avVNyYgj6RiXuQKNwkXaxFyQ",
      genesis: "0000000000000000000000000000000000000000",
      satoshi: 10000,
      script_type: "76a91488ac",
      script: "76a91446af3fb481837fadbb421727f9959c2d32a3682988ac",
      height: 352701
    }
  ]
}
```

## GET /genesis/{genesis}/utxo 通过溯源genesis获取相关utxo列表

#### Request

`GET /genesis/74967a27ce3b46244e2e1fba60844c56bb99afc3/utxo`

#### Response

返回值 code == 0 为正确，其他都是错误

data包括字段为：

- txid: 当前txid
- vout: 当前输出序号
- address: 当前输出的address
- genesis: 当前输出的genesis
- satoshi: 当前输出的satoshi
- script_type: 当前输出锁定脚本类型
- script: 当前输出锁定脚本
- height: 当前交易被打包的区块高度

- Body
```
{
  code: 0,
  msg: "ok",
  data: [
    {
      txid: "2b9e310ddc66214067ff71be9247ff0e25e9cf1a7628e77dafcd0a9836540199",
      vout: 1,
      address: "154vy1TEGYoorGfWXNgqeNWMrux5HoR87k",
      genesis: "74967a27ce3b46244e2e1fba60844c56bb99afc3",
      satoshi: 24000,
      script_type: "510101010101580101580151515a0154580152795279935a7951799300795b79",
      script: "...",
      height: 673050
    },
    {
      txid: "2b9e310ddc66214067ff71be9247ff0e25e9cf1a7628e77dafcd0a9836540199",
      vout: 0,
      address: "17KBdi6KLHorxrxpPtXiTvcN7BgQBasrQn",
      genesis: "74967a27ce3b46244e2e1fba60844c56bb99afc3",
      satoshi: 24000,
      script_type: "510101010101580101580151515a0154580152795279935a7951799300795b79",
      script: "...",
      height: 673050
    },
    {
      txid: "c98fde33eea8b1c283d1ca14b1337818bc6dae7cc3bf3151c22c9fee9b21da80",
      vout: 1,
      address: "1GexnF7rDXS9qJ8UjJr2oddjW7ESxKUH8T",
      genesis: "74967a27ce3b46244e2e1fba60844c56bb99afc3",
      satoshi: 24000,
      script_type: "510101010101580101580151515a0154580152795279935a7951799300795b79",
      script: "...",
      height: 673038
    }
  ]
}

```

## GET /address/{address}/history 通过地址address获取相关tx历史列表

#### Request

`GET /address/17SkEw2md5avVNyYgj6RiXuQKNwkXaxFyQ/history`

#### Response

返回值 code == 0 为正确，其他都是错误

data包括字段为：

- txid: 当前txid
- vout: 当前输出序号
- address: 当前输出的address
- genesis: 当前输出的genesis
- satoshi: 当前输出的satoshi
- script_type: 当前输出锁定脚本类型
- height: 当前交易被打包的区块高度
- io_type: 1为输出包含(即收入)，0为输入包含(即花费)

- Body
```
{
  code: 0,
  msg: "ok",
  data: [
    {
      txid: "62104aa084f4a158cb9aa545ee30d68db88bb22d4a66904b78d41e4512c1969a",
      vout: 0,
      address: "17SkEw2md5avVNyYgj6RiXuQKNwkXaxFyQ",
      genesis: "0000000000000000000000000000000000000000",
      satoshi: 100000,
      script_type: "76a91488ac",
      height: 474237,
      io_type: 1
    },
    {
      txid: "2c63ac6d71e696dea43ef1ef7fba8c376a6a220383e73a17ba6c3795996db112",
      vout: 0,
      address: "17SkEw2md5avVNyYgj6RiXuQKNwkXaxFyQ",
      genesis: "0000000000000000000000000000000000000000",
      satoshi: 10000,
      script_type: "76a91488ac",
      height: 352701,
      io_type: 1
    },
    {
      txid: "cca7507897abc89628f450e8b1e0c6fca4ec3f7b34cccf55f3f531c659ff4d79",
      vout: 0,
      address: "17SkEw2md5avVNyYgj6RiXuQKNwkXaxFyQ",
      genesis: "0000000000000000000000000000000000000000",
      satoshi: 1000000000000,
      script_type: "76a91488ac",
      height: 57044,
      io_type: 0
    },
    {
      txid: "a1075db55d416d3ca199f55b6084e2115b9345e16c5cf302fc80e9d5fbf5d48d",
      vout: 0,
      address: "17SkEw2md5avVNyYgj6RiXuQKNwkXaxFyQ",
      genesis: "0000000000000000000000000000000000000000",
      satoshi: 1000000000000,
      script_type: "76a91488ac",
      height: 57043,
      io_type: 1
    }
  ]
}
```

## GET /genesis/{genesis}/history 通过溯源genesisId获取相关tx历史列表

#### Request

`GET /genesis/74967a27ce3b46244e2e1fba60844c56bb99afc3/history`

#### Response

返回值 code == 0 为正确，其他都是错误

data包括字段为：

- txid: 当前txid
- vout: 当前输出序号
- address: 当前输出的address
- genesis: 当前输出的genesis
- satoshi: 当前输出的satoshi
- script_type: 当前输出锁定脚本类型
- height: 当前交易被打包的区块高度
- io_type: 1为输出包含(即收入)，0为输入包含(即花费)

- Body
```
{
  code: 0,
  msg: "ok",
  data: [
    {
      txid: "2b9e310ddc66214067ff71be9247ff0e25e9cf1a7628e77dafcd0a9836540199",
      vout: 0,
      address: "17KBdi6KLHorxrxpPtXiTvcN7BgQBasrQn",
      genesis: "74967a27ce3b46244e2e1fba60844c56bb99afc3",
      satoshi: 24000,
      script_type: "510101010101580101580151515a0154580152795279935a7951799300795b79",
      height: 673050,
      io_type: 1
    },
    {
      txid: "2b9e310ddc66214067ff71be9247ff0e25e9cf1a7628e77dafcd0a9836540199",
      vout: 1,
      address: "154vy1TEGYoorGfWXNgqeNWMrux5HoR87k",
      genesis: "74967a27ce3b46244e2e1fba60844c56bb99afc3",
      satoshi: 24000,
      script_type: "510101010101580101580151515a0154580152795279935a7951799300795b79",
      height: 673050,
      io_type: 1
    },
    {
      txid: "2b9e310ddc66214067ff71be9247ff0e25e9cf1a7628e77dafcd0a9836540199",
      vout: 0,
      address: "17MkeHkZGvyKR1N57kYKhs4GYU7nNuRtWP",
      genesis: "74967a27ce3b46244e2e1fba60844c56bb99afc3",
      satoshi: 24000,
      script_type: "510101010101580101580151515a0154580152795279935a7951799300795b79",
      height: 673050,
      io_type: 0
    },
    {
      txid: "c98fde33eea8b1c283d1ca14b1337818bc6dae7cc3bf3151c22c9fee9b21da80",
      vout: 0,
      address: "17MkeHkZGvyKR1N57kYKhs4GYU7nNuRtWP",
      genesis: "74967a27ce3b46244e2e1fba60844c56bb99afc3",
      satoshi: 24000,
      script_type: "510101010101580101580151515a0154580152795279935a7951799300795b79",
      height: 673038,
      io_type: 1
    },
    {
      txid: "6b30e971313f1ea366e9d49cb99500d8df941e9037fcc1ed586976baeffada5a",
      vout: 0,
      address: "17MkeHkZGvyKR1N57kYKhs4GYU7nNuRtWP",
      genesis: "74967a27ce3b46244e2e1fba60844c56bb99afc3",
      satoshi: 24000,
      script_type: "510101010101580101580151515a0154580152795279935a7951799300795b79",
      height: 673038,
      io_type: 1
    },
    {
      txid: "c98fde33eea8b1c283d1ca14b1337818bc6dae7cc3bf3151c22c9fee9b21da80",
      vout: 1,
      address: "1GexnF7rDXS9qJ8UjJr2oddjW7ESxKUH8T",
      genesis: "74967a27ce3b46244e2e1fba60844c56bb99afc3",
      satoshi: 24000,
      script_type: "510101010101580101580151515a0154580152795279935a7951799300795b79",
      height: 673038,
      io_type: 1
    },
    {
      txid: "c98fde33eea8b1c283d1ca14b1337818bc6dae7cc3bf3151c22c9fee9b21da80",
      vout: 0,
      address: "17MkeHkZGvyKR1N57kYKhs4GYU7nNuRtWP",
      genesis: "74967a27ce3b46244e2e1fba60844c56bb99afc3",
      satoshi: 24000,
      script_type: "510101010101580101580151515a0154580152795279935a7951799300795b79",
      height: 673038,
      io_type: 0
    }
  ]
}
```
