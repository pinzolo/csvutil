package csvutil

import (
	"math/rand"
	"strconv"

	gimei "github.com/mattn/go-gimei"
)

var (
	firstApartmentNames = []string{
		"アーク",
		"アーバン",
		"アパ",
		"アルファ",
		"ウィング",
		"エクセル",
		"エバー",
		"エル",
		"ガーデン",
		"グラン",
		"グリーン",
		"クレスト",
		"グレース",
		"クローバー",
		"コスモ",
		"サクラ",
		"ザ・",
		"サニー",
		"サン",
		"シーサイド",
		"スカイ",
		"スイート",
		"スター",
		"テラス",
		"ドリーム",
		"ノーブル",
		"パーク",
		"パール",
		"ファイン",
		"ファミリー",
		"フォレスト",
		"ブライト",
		"ベル",
		"ライジング",
		"ライフ",
		"リアル",
		"レインボー",
		"レイク",
		"レジェンド",
		"ロイヤル",
		"ローズ",
	}
	lastApartmentNames = []string{
		"ヴィレッジ",
		"エステート",
		"ガーデン",
		"グリーン",
		"コート",
		"サイド",
		"シティ",
		"シャトレ",
		"スクエア",
		"ステージ",
		"タワー",
		"ハイツ",
		"ハイム",
		"ハウス",
		"パレス",
		"パーク",
		"ヒルズ",
		"プラザ",
		"フォレスト",
		"メゾン",
		"ライフ",
		"レジデンス",
	}
	singleApartmentNames = []string{
		"アクティ",
		"アネックス",
		"ヴィラ",
		"ウィステリア",
		"ヴェルディ",
		"エスポワール",
		"シーサイド",
		"シャーメゾン",
		"シャトー",
		"シャトレー",
		"ジョイフル",
		"ソレイユ",
		"プレステージ",
		"フローラル",
		"プロムナード",
		"ポートアイランド",
		"リバーサイド",
		"リヴィエール",
		"ルミエール",
	}
	campanySuffixes = []string{
		"銀行",
		"工務店",
		"電機",
		"不動産",
		"システムズ",
		"電算",
		"製作所",
		"生命保険",
		"工業",
		"損害保険",
		"商事",
		"製菓",
		"新聞社",
		"出版",
		"印刷",
		"製薬",
		"通信",
	}
)

func fakeApartment(full bool) string {
	name := fakeApartmentName()
	room := strconv.Itoa((rand.Intn(20)+1)*100 + rand.Intn(30) + 1)
	if full {
		return name + toFullWidthNum(room)
	}
	return name + room
}

func fakeOffice(full bool) string {
	name := gimei.NewName().First.Kanji() + "ビル"
	floor := strconv.Itoa(rand.Intn(20)+1) + "F"
	if full {
		return name + toFullWidthNum(floor)
	}
	return name + floor
}

func fakeApartmentName() string {
	var name string
	if lot(70) {
		f, l := sampleString(firstApartmentNames), sampleString(lastApartmentNames)
		for f == l {
			f, l = sampleString(firstApartmentNames), sampleString(lastApartmentNames)
		}
		name = f + l
	} else {
		name = sampleString(singleApartmentNames)
	}
	if lot(60) {
		name += gimei.NewAddress().Town.Kanji()
	}
	return name
}
