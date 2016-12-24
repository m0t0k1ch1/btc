# btctx

[![GoDoc](https://godoc.org/github.com/m0t0k1ch1/btctx?status.svg)](https://godoc.org/github.com/m0t0k1ch1/btctx) [![wercker status](https://app.wercker.com/status/18bad0756ae8ad372b622d6ce8b3691c/s/master "wercker status")](https://app.wercker.com/project/byKey/18bad0756ae8ad372b622d6ce8b3691c)

a package for mapping [bitcoin transaction](https://en.bitcoin.it/wiki/Protocol_documentation#tx) and golang struct

``` sh
$ go get github.com/m0t0k1ch1/btctx
```

## Example

``` go
package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/m0t0k1ch1/btctx"
)

func main() {
	btctx.UseTestnet()

	validTxHex := "0100000001ce3cf2e2b334e7e9fa84619469d9edc49368c2f752ea30fb48b080fc794f6d56010000006a473044022065fe1ea4e94a9b44fb62c2b874b63a947504273a60b99b8f7bbf77b4db9331b002205559d8ee93cf341d75866f9eb912af05904fb6eed7372a837308c4e37f3ab58f012103bae5f04799c40862358560e42e441c3080b997a3dec161dd40395e992362bfc9feffffff0200f2052a010000001976a914cbc222711a230ecdd9a5aa65b61ed39c24db2b3488acc08d931a1d0000001976a914426c1ad9fa94f9ea3e6f9248b8bff6768e3ac8c488ac951a1000"

	tx, err := btctx.NewTxFromHex(validTxHex)
	if err != nil {
		log.Fatal(err)
	}

	b, err := json.Marshal(tx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(b))

	txHex, err := tx.ToHex()
	if err != nil {
		log.Fatal(err)
	}

	if txHex == validTxHex {
		log.Println("match!")
	}
}
```

``` json
{
  "version": 1,
  "txIns": [
    {
      "hash": "566d4f79fc80b048fb30ea52f7c26893c4edd969946184fae9e734b3e2f23cce",
      "index": 1,
      "sigScript": {
        "hex": "473044022065fe1ea4e94a9b44fb62c2b874b63a947504273a60b99b8f7bbf77b4db9331b002205559d8ee93cf341d75866f9eb912af05904fb6eed7372a837308c4e37f3ab58f012103bae5f04799c40862358560e42e441c3080b997a3dec161dd40395e992362bfc9",
        "asm": "3044022065fe1ea4e94a9b44fb62c2b874b63a947504273a60b99b8f7bbf77b4db9331b002205559d8ee93cf341d75866f9eb912af05904fb6eed7372a837308c4e37f3ab58f01 03bae5f04799c40862358560e42e441c3080b997a3dec161dd40395e992362bfc9",
        "addresses": null
      },
      "sequence": 4294967294
    }
  ],
  "txOuts": [
    {
      "value": 5000000000,
      "pkScript": {
        "hex": "76a914cbc222711a230ecdd9a5aa65b61ed39c24db2b3488ac",
        "asm": "OP_DUP OP_HASH160 cbc222711a230ecdd9a5aa65b61ed39c24db2b34 OP_EQUALVERIFY OP_CHECKSIG",
        "addresses": [
          "mz6L2hYM8jPR5nhH6kEsc3DQFiSDA1Jqpa"
        ]
      }
    },
    {
      "value": 124999929280,
      "pkScript": {
        "hex": "76a914426c1ad9fa94f9ea3e6f9248b8bff6768e3ac8c488ac",
        "asm": "OP_DUP OP_HASH160 426c1ad9fa94f9ea3e6f9248b8bff6768e3ac8c4 OP_EQUALVERIFY OP_CHECKSIG",
        "addresses": [
          "mmaAPyTMoK3p1K7qCYytuSRRQMPCbrUAMN"
        ]
      }
    }
  ],
  "lockTime": 1055381
}
```
