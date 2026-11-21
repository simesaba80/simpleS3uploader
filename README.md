# 使い方

## 前提
miseが使えること

## 1. バイナリのビルド
```
mise install
mise run build
```

## 2. 環境変数の設定
```
cp .env.copy .env
```
後に適切な環境変数を.envに設定

## 3. 実行
```
./simples3uploader  [<filename>]
```
ファイル名を複数指定すれば同時に複数枚アップロードできます  
出力は `Uploaded <ファイル名> as <URL>` で出力されます
