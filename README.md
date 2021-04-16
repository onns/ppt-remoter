# ppt-remoter

PowerPoint Remote Control

# 演示文档遥控器

## 工作原理

1. 首先去服务区询问`id`
2. 建立`socket`链接，接收播放信息
3. 生成命令行二维码

```bash

git clone git@github.com:mdp/qrterminal.git
go get -u github.com/skip2/go-qrcode/...
go get -u github.com/gorilla/websocket
go get rsc.io/qr
```


```bash
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build client.go
```