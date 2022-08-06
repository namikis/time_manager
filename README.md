# time_manager
勤怠管理用SlackBot

## 概要
インターン先で勤怠情報をスプレッドシートに記録する際に、休憩が複数回あると計算に手間がかかる。
そのため　稼働開始・終了・休会開始・再開の四つの時間をSlackBotで管理すれば楽なのではないかと思い作った。

## コマンド
#### 稼働開始
```
@timemanager start
```
#### 終了
```
@timemanager end
```
#### 休憩開始
```
@timemanager break start
```
#### 再開
```
@timemanager break end
```

## 使用技術
Golang, MySQL, Docker-compose
