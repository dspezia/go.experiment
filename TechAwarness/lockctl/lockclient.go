package main

import "fmt"
import "net"
import "bufio"
import "time"

/*****************************************************************************/

var query1 = []byte(`{"Op":"incr", "Target":"counter", "Arg":"1"}`)
var query2 = []byte(`{"Op":"set", "Target":"counter", "Arg":"0"}`)

/*****************************************************************************/

func clientLoop(result *chan int) {

	res := 0
	defer func() { *result <- res }()

	conn, err := net.Dial("tcp", *flagTarget)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	defer conn.Close()

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	for i := 0; i < *flagNbIter; {

		pos := 0
		for ; i < *flagNbIter && pos < *flagPipe; pos += 2 {
			writer.Write(query1)
			writer.Write(query2)
			i += 2
		}
		writer.Flush()

		for j := 0; j < pos; j++ {
			_, err := reader.ReadBytes('\n')
			if err != nil {
				break
			}
			res++
		}
	}

}

/*****************************************************************************/

func mainClient() {

	t := time.Now()

	result := make(chan int)
	for i := 0; i < *flagNbCon; i++ {
		go clientLoop(&result)
	}
	sum := 0
	for i := 0; i < *flagNbCon; i++ {
		sum += <-result
	}

	d := time.Now().Sub(t)

	fmt.Println("Result:", sum, "in", d)
	fmt.Println("Throughput:", float64(sum)/d.Seconds())
}

/*****************************************************************************/
