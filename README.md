# ppt-remoter

PowerPoint Remote Control

## Usage

1. Go to the [release page](https://github.com/onns/ppt-remoter/releases/latest), download [cilent.exe](https://github.com/onns/ppt-remoter/releases/download/v0.0.1/client.exe) and [config.json](https://github.com/onns/ppt-remoter/releases/download/v0.0.1/config.json), put them in the same folder.
2. Click the **client.exe**.
3. Copy the `https://...` to your phone or just open the **qr.png** in the same folder and scan the qrcode.
4. Open your PowerPoint and **play it** (Full Screen Mode).
5. Use your phone to control it!

![Demo Page](https://user-images.githubusercontent.com/16622934/115057501-b9c93400-9f16-11eb-9d29-61290221d20d.jpg)

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