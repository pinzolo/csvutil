# csvutil

[![Build Status](https://travis-ci.org/pinzolo/csvutil.png)](http://travis-ci.org/pinzolo/csvutil)
[![Coverage Status](https://coveralls.io/repos/github/pinzolo/csvutil/badge.svg?branch=master)](https://coveralls.io/github/pinzolo/csvutil?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/pinzolo/csvutil)](https://goreportcard.com/report/github.com/pinzolo/csvutil)
[![Godoc](http://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/github.com/pinzolo/csvutil)
[![license](http://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/pinzolo/csvutil/master/LICENSE)

## Description

CSVに対して様々な処理を行うコマンドラインツールです。  
CSVの構成操作、およびダミーデータ生成が可能です。  
それぞれのサブコマンドはパイプで組み合わせて利用できます。

## Usage

ダミーデータ生成サンプル:

```bash
$ csvutil generate --size 5 --count 10 --header 氏名:郵便番号:住所:建物:メール | \
  csvutil name --name 氏名 |\
  csvutil address --zip-code 郵便番号 --prefecture 住所 --city 住所 --town 住所 --block-number |\
  csvutil building --column 建物 |\
  csvutil email --column メール
```

出力結果:

```
氏名,郵便番号,住所,建物,メール
日高 茂晴,563-3262,常滑市中津川21,ルミエール208,normamorales@topicshots.name
森川 優,861-2458,守谷市太秦安井藤ノ木町17-11,シャトー龍神村西1410,loisfuller@skinder.com
荒川 美華,938-5130,球磨郡山江村大町46-13,エクセルパーク丸島町2018,3henry@kazu.gov
野村 安幸,327-2817,東近江市西日暮里14-30-10,シーサイドレジデンス1923,karenphillips@skivee.info
西 唯七,535-3970,大島郡知名町栄町12,クレストフォレスト奔渡1522,annaprice@twinte.edu
松井 璃亜,680-9813,大阪市北区光陽町12-28,サクラライフ池見108,joycereed@blogxs.name
武田 知佳子,951-1992,青ヶ島村翠ケ丘町57,アネックス緑台1115,donnaadams@vidoo.mil
細川 優在,189-5976,尼崎市八尾町東葛坂45-19,フォレストヴィレッジ東城内1414,brucekennedy@skidoo.com
本間 唯斐,102-0997,最上郡戸沢村南方42,ルミエール河原市213,howardmoore@mybuzz.com
戸田 孝弘,593-6542,広島市安佐北区本島町甲生10-6,ヴィラ渡内102,voluptatem_ut_magni@youbridge.com
```

`csvutil help` で全体のヘルプを確認し、`csvutil help [subcommand]` で各サブコマンドの詳細なヘルプを確認してください。

## Install

[Releases · pinzolo/csvutil](https://github.com/pinzolo/csvutil/releases) から最新の自分の環境にあったバイナリをダウンロードしてお使いください。

また、Go環境がある場合 `go get` でインストールできます。

```bash
$ go get github.com/pinzolo/csvutil/cmd/csvutil
```

## Contribution

1. Fork ([https://github.com/pinzolo/csvutil/fork](https://github.com/pinzolo/csvutil/fork))
1. Create a feature branch
1. Commit your changes
1. Rebase your local changes against the master branch
1. Run test suite with the `go test ./...` command and confirm that it passes
1. Run `gofmt -s`
1. Create a new Pull Request

## Author

[pinzolo](https://github.com/pinzolo)
