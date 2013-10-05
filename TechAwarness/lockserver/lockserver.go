package main

import "log"
import "net"
import "io"
import "encoding/json"
import "os"
import "strconv"

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
	OP_GET
	OP_SET
	OP_INCR
)

var Service = map[string]Operation{
	"lock":   OP_LOCK,
	"unlock": OP_UNLOCK,
	"get":    OP_GET,
	"set":    OP_SET,
	"incr":   OP_INCR,
}

/*****************************************************************************/

type MessageQuery struct {
	Op     string
	Target string
	Arg    string `json:",omitempty"`
	oper   Operation
	clt    *Client
}

type MessageReply struct {
	Status string
	Error  string `json:",omitempty"`
	Value  string `json:",omitempty"`
	oper   Operation
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
		if !end {
			if err := encoder.Encode(reply); err != nil {
				log.Println("Error ", err)
				end = true
			}
		}
	}
}

/*****************************************************************************/

type Core struct {
	in    chan *MessageQuery
	locks *LockArea
	stats map[string]int64
}

/*****************************************************************************/

func NewCore() *Core {
	return &Core{
		in:    make(chan *MessageQuery, channelSize*128),
		locks: NewLockArea(),
		stats: make(map[string]int64),
	}
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
		case OP_GET:
			core.handleGet(m)
		case OP_SET:
			core.handleSet(m)
		case OP_INCR:
			core.handleIncr(m)
		default:
			m.clt.coreOut <- &MessageReply{Status: "KO", Error: "Unknown operation"}
		}
	}
}

/*****************************************************************************/

func (core *Core) handleOpen(query *MessageQuery) {

	log.Println("Opening connection")
	core.locks.AddClient(query.clt)
}

/*****************************************************************************/

func (core *Core) handleClose(query *MessageQuery) {

	log.Println("Closing connection")
	toBeNotified := core.locks.RemoveClient(query.clt)

	reply := &MessageReply{oper: OP_CLOSE}
	query.clt.coreOut <- reply

	for _, c := range toBeNotified {
		reply := &MessageReply{Status: "OK"}
		c.coreOut <- reply
	}

}

/*****************************************************************************/

func (core *Core) handleLock(query *MessageQuery) {

	log.Println("Locking", query.Target)
	if core.locks.Add(query.clt, query.Target) {
		reply := &MessageReply{Status: "OK"}
		query.clt.coreOut <- reply
	}
}

/*****************************************************************************/

func (core *Core) handleUnlock(query *MessageQuery) {

	log.Println("Unlocking", query.Target)
	c, ok := core.locks.Remove(query.clt, query.Target)
	if ok {
		reply := &MessageReply{Status: "OK"}
		query.clt.coreOut <- reply
		if c != nil {
			c.coreOut <- reply
		}
	} else {
		reply := &MessageReply{Status: "KO", Error: "Cannot find this lock"}
		query.clt.coreOut <- reply
	}
}

/*****************************************************************************/

func (core *Core) handleGet(query *MessageQuery) {

	log.Println("Getting", query.Target)
	val := strconv.FormatInt(core.stats[query.Target], 10)
	reply := &MessageReply{Status: "OK", Value: val}
	query.clt.coreOut <- reply
}

/*****************************************************************************/

func (core *Core) handleSet(query *MessageQuery) {

	log.Println("Setting", query.Target)
	var reply *MessageReply
	if n, err := strconv.ParseInt(query.Arg, 10, 64); err != nil {
		reply = &MessageReply{Status: "KO", Error: "Invalid number"}
	} else {
		core.stats[query.Target] = n
		reply = &MessageReply{Status: "OK"}
	}
	query.clt.coreOut <- reply
}

/*****************************************************************************/

func (core *Core) handleIncr(query *MessageQuery) {

	log.Println("Increment", query.Target)
	var reply *MessageReply
	if n, err := strconv.ParseInt(query.Arg, 10, 64); err != nil {
		reply = &MessageReply{Status: "KO", Error: "Invalid number"}
	} else {
		core.stats[query.Target] += n
		val := strconv.FormatInt(core.stats[query.Target], 10)
		reply = &MessageReply{Status: "OK", Value: val}
	}
	query.clt.coreOut <- reply
}

/*****************************************************************************/

func mainServer() {

	core := NewCore()
	go core.main()

	lis := &Listener{core: core}
	go lis.Listen("tcp", *flagServer)

	channel := make(chan os.Signal)
	signal.Notify(channel, os.Interrupt)
	<-channel
	log.Println("Stop")
}

/*****************************************************************************/
