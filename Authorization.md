# 1. 接口Token鉴权

在请求时HTTP请求头中添加授权accesskey即可。即 accesskey=9a48a35a8dadd8b54d3859865b292aa6e7624e34e83f073c62065212d7ef0f9b

# 2. 接口签名鉴权

## 请求结构

1. API 的所有接口均通过 `HTTPS` 进行通信，均使用 `UTF-8` 编码
2. 支持的 HTTP 请求方法：GET
4. 注意：请勿在前端直接发起开放接口请求，防止泄露 secret

## 公共参数

### 1. 参数说明

用于标识用户和接口鉴权目的的参数，如非必要，在每个接口单独的接口文档中不再对这些参数进行说明，但每次请求均需要携带这些参数，才能正常发起请求。 公共参数定义如下：

* appid 第三方应用id，如：9Hs3kifu8uHuJkmS9ktk
* ts 当前UTC时间戳，如：1661648098
* sign 签名，如：63edca714f5c4e6ca5421b2ce4530a4d5472e8d3cf8a7b93082522b32117b5ba

### 2.GET 请求结构示例

`https://api.sensiblequery.com/address/17PcYSCLs7rx5BWC1rmA35ZLPNrN5u85dW/history/tx?start=0&end=0&cursor=0&size=5&appid=9Hs3kifu8uHuJkmS9ktk&ts=1661648098&sign=63edca714f5c4e6ca5421b2ce4530a4d5472e8d3cf8a7b93082522b32117b5ba`

## 接口鉴权

### 1.申请应用凭证

* 使用开放接口前须先后台申请 appid 和 secret
* API 会对每个访问请求进行身份验证，即每个请求都需要在公共请求参数中包含签名信息（Signature）以验证请求者身份
* 签名信息由安全凭证生成，安全凭证为 secret， 用于加密签名字符串和服务器端验证签名字符串的密钥（即 sign 参数）

### 2.对请求参数排序

* 首先对所有URL中的请求参数（不包括sign）按参数名的字典序（ ASCII 码）升序排序
* 用户可以借助编程语言中的相关排序函数来实现这一功能，如 PHP 中的 ksort 函数、Golang 中的 sort.Strings(keys)
* 上述GET示例参数的排序结果如下: `appid、cursor、end、size、start、ts`

### 3.拼接 query 字符串

* 此步骤生成URL参数请求字符串，将把上一步排序好的请求参数格式化成“参数名称”=“参数值”的形式
* 注意：“参数值”为原始值而非url编码后的值。然后将格式化后的各个参数用"&"拼接在一起
* 最终生成的 query 字符串为: `appid=9Hs3kifu8uHuJkmS9ktk&cursor=0&end=0&size=5&start=0&ts=1661648098`

### 4.拼接签名前原文字符串

* 此步骤生成签名原文字符串
* 请求签名原文串的拼接规则： 接口路由 + "?" + query字符串
* 实际的接口路由根据接口所属模块的不同而不同，详见各接口说明。比如:  `/address/17PcYSCLs7rx5BWC1rmA35ZLPNrN5u85dW/history/tx`。

`/address/17PcYSCLs7rx5BWC1rmA35ZLPNrN5u85dW/history/tx?appid=9Hs3kifu8uHuJkmS9ktk&cursor=0&end=0&size=5&start=0&ts=1661648098`

### 5.生成签名

首先使用HMAC-SHA256算法对上一步中获得的签名原文字符串进行签名，然后将生成的签名串使用十六进制进行编码，即可获得最终的签名串。

* python示例代码

```
    import hmac
    import hashlib

    API_ID = '9Hs3kifu8uHuJkmS9ktk'
    API_SECRET = 'oiwzUTJ9nSqevJW0VYE7QyikOEh2rRgAXlmFed2w'

    path = '/address/17PcYSCLs7rx5BWC1rmA35ZLPNrN5u85dW/history/tx'

    params = [
        "start=0",
    "end=0",
    "cursor=0",
    "size=5",
    "appid="+API_ID,
    "ts=1661648098",
    ]
    params.sort()

    message = path + "?" + ("&".join(params))
    signature = hmac.new(bytes(API_SECRET, 'utf-8'),
                         msg=bytes(message, 'utf-8'),
                         digestmod=hashlib.sha256).hexdigest()
    print(signature)
```

### 6.拼接签名

请求API时添加来自第5步的sign参数，如 xxx&sign={$signStr}。最终请求地址为：

`https://api.sensiblequery.com/address/17PcYSCLs7rx5BWC1rmA35ZLPNrN5u85dW/history/tx?start=0&end=0&cursor=0&size=5&appid=9Hs3kifu8uHuJkmS9ktk&ts=1661648098&sign=63edca714f5c4e6ca5421b2ce4530a4d5472e8d3cf8a7b93082522b32117b5ba`
