// -------------------------------------------
// Created By : Onns onns@onns.xyz
// File Name : client.go
// Purpose :
// Creation Date : 2021-04-16 14:41:33
// Last Modified : 2021-04-16 14:42:14
// -------------------------------------------

package main

import (
	"encoding/json"
	"flag"
	"github.com/go-vgo/robotgo"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"time"
)
import qrcode "github.com/skip2/go-qrcode"

type GlobalConfig struct {
	Server string `json:"server"`
}

var OnnsGlobal GlobalConfig

func loadConfig() {
	filename := "config.json"
	if _, err := os.Stat(filename); err == nil {
		b, err := ioutil.ReadFile(filename)
		if err != nil {
			log.Println(err)
		}
		json.Unmarshal(b, &OnnsGlobal)
	}
}

// https://blog.csdn.net/Phoenix_smf/article/details/89278398
func init() {
	loadConfig()
}

func getId() string {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://"+OnnsGlobal.Server+"/v1/session", nil)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("get-id:", err)
	}
	bodyText, _ := ioutil.ReadAll(resp.Body)
	s := string(bodyText)
	// log.Printf("get id: %s %T", s, bodyText)

	type tmpjson struct {
		Id string `json:"selfID"`
	}
	var r tmpjson
	json.Unmarshal([]byte(s), &r)
	log.Printf("get id: %s ", r.Id)

	return r.Id
}

func getPairId(id string) string {
	return id[:4] + "1"
}

func main() {

	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	id := getId()
	u := url.URL{Scheme: "wss", Host: OnnsGlobal.Server, Path: "/ws/", RawQuery: "id=" + id}
	log.Printf("connecting to %s", u.String())
	// qrterminal.Generate("https://onns.xyz/ppt-remoter/?id="+getPairId(id), qrterminal.L, os.Stdout)
	qrcode.WriteFile("https://onns.xyz/ppt-remoter/?id="+getPairId(id)+"&site="+OnnsGlobal.Server, qrcode.Medium, 256, "qr.png")
	// qrcode.WriteFile("http://192.168.1.154/GitHub/onns.xyz/ppt-remoter/?id="+getPairId(id)+"&site="+OnnsGlobal.Server, qrcode.Medium, 256, "qr.png")
	log.Printf("please visit %s", "https://onns.xyz/ppt-remoter/?id="+getPairId(id)+"&site="+OnnsGlobal.Server)
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", message)
			s := string(message)
			type cmdjson struct {
				Cmd string `json:"cmd"`
			}
			var r cmdjson
			json.Unmarshal([]byte(s), &r)
			if r.Cmd == "prev" {
				robotgo.KeyTap("up")
			}
			if r.Cmd == "next" {
				robotgo.KeyTap("down")
			}

		}
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		// case t := <-ticker.C:
		// 	err := c.WriteMessage(websocket.TextMessage, []byte(t.String()))
		// 	if err != nil {
		// 		log.Println("write:", err)
		// 		return
		// 	}
		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}
