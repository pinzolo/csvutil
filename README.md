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
  csvutil name --name 氏名 | \
  csvutil address --zip-code 郵便番号 --prefecture 住所 --city 住所 --town 住所 --block-number | \
  csvutil building --column 建物 | \
  csvutil email --column メール
```

出力結果:

```
氏名,郵便番号,住所,建物,メール
土田 慶太郎,061-8035,京都府和光市西春別駅前西町42-18,プレステージ前河原1926,tempora_doloribus_inventore@twitterbeat.info
古川 篤人,468-2410,長崎県東根市御蔵島村一円16,テラスライフ1421,charlesstevens@eire.mil
横田 姫葵,197-3356,愛知県志木市神田佐久間河岸44-12-6,ソレイユ麻生1229,qjordan@roomm.edu
荻野 吉文,819-3023,東京都大阪市平野区亀尾町44-24-2,グリーンエステート1727,tempore_numquam_consequatur@midel.org
森山 沙江,289-9318,岡山県大阪市東淀川区桜作42-6,グレースヴィレッジ宇和町杢所1118,ucooper@voonder.net
大野 沙祈,093-0401,新潟県河北郡内灘町木興町46-10,パールレジデンス梨ケ原806,victorpierce@dynabox.edu
山崎 玲菜,065-9621,兵庫県大阪市浪速区八幡町39-18,スイートコート光町506,qui_aut_nihil@meeveo.biz
金沢 優二,201-6394,徳島県岩手郡雫石町日置野田36-10,プレステージ八森樋長421,ibell@demizz.com
松原 誠吾,877-1329,石川県吉野郡野迫川村灘38-5-10,ガーデンタワー藤河内1120,ab_cum@voomm.info
坂井 力也,967-7366,愛媛県大阪市西区白鳥町中津屋41,レイクパーク609,perferendis_earum@quamba.edu
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
