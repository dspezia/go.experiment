package lockserver

import "log"
import "net/http"
import "code.google.com/p/go.net/websocket"
import "time"
import "sync/atomic"

/*****************************************************************************/

type ResultJson struct {
	Tps int64
}

var Counter *int64

/*****************************************************************************/

func MonitoringServer(ws *websocket.Conn) {

	log.Println("Connected")

	done := make(chan bool)
	go func() {
		cnt := atomic.LoadInt64(Counter)
		for {
			select {
			case <-done:
				return
			case <-time.After(time.Second / 2):
				cur := atomic.LoadInt64(Counter)
				delta := 2 * (cur - cnt)
				cnt = cur
				err := websocket.JSON.Send(ws, ResultJson{Tps: delta})
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

/*****************************************************************************/

func monitoringServer(core *Core) {
	Counter = &core.count
	http.Handle("/monitoring", websocket.Handler(MonitoringServer))
	http.ListenAndServe(":4010", nil)
}

/*****************************************************************************/
