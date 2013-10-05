// lockserver project main.go
package main

import "flag"
import "runtime"
import "fmt"

/*****************************************************************************/

var flagListen = flag.Bool("l", false, "Listen (server mode)")
var flagServer = flag.String("s", ":4002", "(host:port)")

var flagTarget = flag.String("t", "localhost:4002", "Target (host:port)")
var flagNbCon = flag.Int("c", 50, "Number of connections")
var flagNbIter = flag.Int("n", 10000, "Number of iterations")
var flagPipe = flag.Int("p", 1, "Pipelining factor")

/*****************************************************************************/

func main() {

	runtime.GOMAXPROCS(4)
	flag.Parse()

	if *flagListen {
		fmt.Println("Server starting ...")
		mainServer()
	} else {
		fmt.Println("Client starting ...")
		mainClient()
	}
}

/*****************************************************************************/
