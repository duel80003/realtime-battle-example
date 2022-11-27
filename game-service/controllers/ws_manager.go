package controllers

import (
	"encoding/json"
	"game-service/tools"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
)

var (
	upgrader  = websocket.Upgrader{}
	WSConnMap sync.Map
)

func EntryHandler(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Errorf("upgrade: %s", err)
		return
	}
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		tools.Logger.Debugf("recv: %s", message)
		var request *Request
		err = json.Unmarshal(message, &request)
		if err != nil {
			tools.Logger.Errorf("Unmarshal error: %s", err)
			break
		}
		setConnInfo(request.PlayerID, c)
		MatchHandler(request)
	}
}

func setConnInfo(playerID string, c *websocket.Conn) {
	WSConnMap.Store(playerID, c)
}

func getConnInfo(playerID string) *websocket.Conn {
	connI, ok := WSConnMap.Load(playerID)
	if !ok {
		return nil
	}
	return connI.(*websocket.Conn)
}

func deleteConn(playerID string) {
	WSConnMap.Delete(playerID)
}
