package sync

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/koding/websocketproxy"
	"net/http"
	"net/url"
)

const wsUrl = "ws://localhost:3001/sync"

func SetSyncHandler() (*websocketproxy.WebsocketProxy, error) {
	u, err := url.Parse(wsUrl)
	if err != nil {
		return &websocketproxy.WebsocketProxy{}, err
	}
	w := websocketproxy.NewProxy(u)
	w.Upgrader = &websocket.Upgrader{
		CheckOrigin: checkOrigin,
	}
	return w, nil
}

func checkOrigin(r *http.Request) bool {
	o := r.Header.Get("origin")
	if o == "" {
		return false
	}
	if o == "http://localhost:4200" {
		fmt.Println("working in dev mode")
		return true
	}
	p, err := url.Parse(o)
	if err != nil {
		fmt.Printf("error parsing request origin header in WS origin check method: %s", err)
		return false
	}
	fmt.Printf("the parsed Host from the request origin: %s \n", p.Host)
	fmt.Printf("the request Host: %s \n", r.Host)
	if o == "http://raspberrypi.local:4444" {
		fmt.Println("running on the Pi")
		return true
	}
	return false
}