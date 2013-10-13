package lockserver

import "log"
import "net/http"
import "code.google.com/p/go.net/websocket"
import "time"

type ResultJson struct {
	Duration int
	Tps      int
}

func MonitoringServer(ws *websocket.Conn) {

	log.Println("Connected")

	done := make(chan bool)
	go func() {
		for {
			select {
			case <-done:
				return
			case <-time.After(time.Second / 2):
				m := ResultJson{Duration: 1, Tps: time.Now().Nanosecond() / 1000000}
				err := websocket.JSON.Send(ws, m)
				if err != nil {
					log.Println("Error send", err)
					return
				}
			}
		}
	}()

	var msg string
	for {
		err := websocket.Message.Receive(ws, &msg)
		if err != nil {
			done <- true
			break
		}
	}
	log.Println("Disconnected")
}

func monitoringServer() {
	http.Handle("/monitoring", websocket.Handler(MonitoringServer))
	http.ListenAndServe(":4010", nil)
}
