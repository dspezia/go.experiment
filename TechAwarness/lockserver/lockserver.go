package main

import "log"
import "net"
import "io"
import "runtime"
import "encoding/json"
import "os"
import "flag"
//import "runtime/pprof"
import "os/signal"

/*****************************************************************************/

const channelSize = 16

/*****************************************************************************/

type Operation int

const (
	OP_NONE = iota
	OP_OPEN
	OP_CLOSE
	OP_LOCK
	OP_UNLOCK
	OP_RELEASE
)

var Service = map[string]Operation{
	"lock":    OP_LOCK,
	"unlock":  OP_UNLOCK,
	"release": OP_RELEASE,
}

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

/*****************************************************************************/

type MessageQuery struct {
	oper Operation
	clt  *Client
	Op   string   `json:"op"`
	Obj  []string `json:"obj"`
}

type MessageReply struct {
	oper   Operation
	Status string `json:"status"`
}

/*****************************************************************************/

type Listener struct {
	lis  *net.Listener
	core *Core
}

/*****************************************************************************/

func (ln *Listener) Listen(t string, addr string) {

	lis, err := net.Listen(t, addr)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Listening to %s-%s\n", t, addr)

	for {
		c, err := lis.Accept()
		if err != nil {
			log.Println(err)
		} else {
			clt := NewClient(c, ln.core)
			go clt.jsonIn()
			go clt.jsonOut()
		}
	}
}

/*****************************************************************************/

type Client struct {
	con     net.Conn
	core    *Core
	coreOut chan *MessageReply
}

/*****************************************************************************/

func NewClient(con net.Conn, core *Core) (clt *Client) {
	channel := make(chan *MessageReply, channelSize)
	return &Client{con: con, core: core, coreOut: channel}
}

/*****************************************************************************/

func (clt *Client) jsonIn() {

	defer func() { clt.core.in <- &MessageQuery{clt: clt, oper: OP_CLOSE} }()

	decoder := json.NewDecoder(clt.con)
	clt.core.in <- &MessageQuery{clt: clt, oper: OP_OPEN}

	for {

		m := &MessageQuery{clt: clt}
		if err := decoder.Decode(m); err == io.EOF {
			break
		} else if err != nil {
			clt.core.in <- &MessageQuery{clt: clt, oper: OP_NONE}
			break
		}

		m.oper = Service[m.Op]
		clt.core.in <- m
	}
}

/*****************************************************************************/

func (clt *Client) jsonOut() {

	defer clt.con.Close()
	encoder := json.NewEncoder(clt.con)

	end := false
	for reply := range clt.coreOut {
		if reply.oper == OP_CLOSE {
			break
		}
		if ( !end ) {
			if err := encoder.Encode(reply); err != nil {
				log.Println("Error ",err)
				end = true
			}
		}
	}
}

/*****************************************************************************/

type Core struct {
	in chan *MessageQuery
	//locks map[string]
}

/*****************************************************************************/

func (core *Core) main() {

	for m := range core.in {
		switch m.oper {
		case OP_OPEN:
			core.handleOpen(m)
		case OP_CLOSE:
			core.handleClose(m)
		case OP_LOCK:
			core.handleLock(m)
		case OP_UNLOCK:
			core.handleUnlock(m)
		case OP_RELEASE:
			core.handleRelease(m)
		default:
			m.clt.coreOut <- &MessageReply{Status: "Error"}
		}
	}
}

/*****************************************************************************/

func (core *Core) handleOpen(query *MessageQuery) {
	log.Println("Opening connection")	
}

/*****************************************************************************/

func (core *Core) handleClose(query *MessageQuery) {
	log.Println("Closing connection")
	reply := &MessageReply{oper: OP_CLOSE}
	query.clt.coreOut <- reply
//	close(query.clt.coreOut)
}

/*****************************************************************************/

func (core *Core) handleLock(query *MessageQuery) {
	//log.Println("Lock")
	reply := &MessageReply{Status: "Ok"}
	query.clt.coreOut <- reply
}

/*****************************************************************************/

func (core *Core) handleUnlock(query *MessageQuery) {
	//log.Println("unlock")
	reply := &MessageReply{Status: "Ok"}
	query.clt.coreOut <- reply
}

/*****************************************************************************/

func (core *Core) handleRelease(query *MessageQuery) {
	//log.Println("release")
	reply := &MessageReply{Status: "Ok"}
	query.clt.coreOut <- reply
}

/*****************************************************************************/

func main() {

	runtime.GOMAXPROCS(16)
/*
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}
*/
	core := &Core{in: make(chan *MessageQuery, channelSize*128)}
	go core.main()

	lis := &Listener{core: core}
	go lis.Listen("tcp", ":4002")

	channel := make(chan os.Signal)
	signal.Notify( channel, os.Interrupt )
	<-channel
	log.Println("Stop")
}

/*****************************************************************************/
