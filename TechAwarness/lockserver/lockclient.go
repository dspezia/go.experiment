package main

import "flag"
import "runtime"
import "fmt"
import "net"
import "bufio"

/*****************************************************************************/

var flagTarget = flag.String( "t", "localhost:4002", "Target (host:port)")
var flagNbCon  = flag.Int( "c", 50, "Number of connections")
var flagNbIter = flag.Int( "n", 10000, "Number of iterations")
var flagPipe   = flag.Int( "p", 1, "Pipelining factor")

var query1  = []byte(`{"op":"lock", "obj":[ "111", "222", "333"]}`)
var query2  = []byte(`{"op":"unlock", "obj":[ "111", "222", "333"]}`)

/*****************************************************************************/

func clientLoop( result *chan int ) {

	res := 0
	defer func(){ *result <- res }()

	conn, err := net.Dial( "tcp", *flagTarget )
	if err != nil {
		fmt.Println("Error: ",err)
		return
	}
	defer conn.Close()

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	for i:=0; i < *flagNbIter; {

		pos := 0
		for ; i < *flagNbIter && pos < *flagPipe; pos += 2 {
			writer.Write( query1 )
			writer.Write( query2 )
			i += 2
		}
		writer.Flush()

		for j:=0; j<pos; j++ {
			_, err := reader.ReadBytes('\n')
			if ( err != nil ) {
				break
			}
			//fmt.Println(string(json))
			res++
		}
	}

}

/*****************************************************************************/

func main() {

	runtime.GOMAXPROCS(1)
	flag.Parse()

	result := make(chan int)
	for i := 0; i < *flagNbCon; i++ {
		go clientLoop(&result)
	}
	sum := 0
	for i := 0; i < *flagNbCon; i++ {
		sum += <-result
	}

	fmt.Println("Result: ",sum)
}

/*****************************************************************************/
