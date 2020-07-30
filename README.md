# cocomoni

![CodeQL](https://github.com/tuki0918/cocomoni/workflows/CodeQL/badge.svg)

新型コロナウイルス接触確認アプリ（COCOA) のダウンロード数を確認するコマンドです。

## Usage

```
$ go run main.go | jq .    
{
  "version": "x.x.x",
  "date": "7月30日17:00",
  "downloads": 950,
  "sentence": "ダウンロード数は、7月30日17:00現在、合計で約950万件です。",
  "link": "https://www.mhlw.go.jp/stf/seisakunitsuite/bunya/cocoa_00138.html"
}
```

## Requirements

厚生労働省に掲載されている画像からテキスト抽出しているため、下記設定が必要です。

+ Google Cloud Platform（GCP）の Cloud Vision API を有効にする
+ Google Cloud API の [認証設定](https://cloud.google.com/docs/authentication/getting-started?hl=ja) を行う
