package lockserver

import "log"
import "net"
import "io"
import "encoding/json"
import "os"
import "strconv"
import "os/signal"
import "sync/atomic"

/*****************************************************************************/

const channelSize = 16
const verbose = false

/*****************************************************************************/

// Operation is an enumerate listing the operation codes
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

// Service is a map to convert an operation name into an enumerate
var Service = map[string]Operation{
	"lock":   OP_LOCK,
	"unlock": OP_UNLOCK,
	"get":    OP_GET,
	"set":    OP_SET,
	"incr":   OP_INCR,
}

/*****************************************************************************/

// MessageQuery is the query message structure.
type MessageQuery struct {
	Op     string
	Target string
	Arg    string `json:",omitempty"`
	oper   Operation
	clt    Replier
}

// MessageReply is the reply message structure.
type MessageReply struct {
	Status string
	Error  string `json:",omitempty"`
	Value  string `json:",omitempty"`
	oper   Operation
}

/*****************************************************************************/

// Listener is the main TCP server, waiting for incoming connections
type Listener struct {
	lis  *net.Listener // TCP listener
	core *Core         // Core goroutine
}

/*****************************************************************************/

// Listen starts the TCP server loop, waiting for incoming connections
func (ln *Listener) Listen(t string, addr string) {

	// Declare TCP listening server
	lis, err := net.Listen(t, addr)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Listening to %s-%s\n", t, addr)

	// Main loop
	for {
		// Accept incoming connection
		c, err := lis.Accept()
		if err != nil {
			log.Println(err)
		} else {
			// Connection accepted, create client and spawn associated goroutines
			clt := NewClient(c, ln.core)
			go clt.jsonIn()
			go clt.jsonOut()
		}
	}
}

/*****************************************************************************/

// Client represents a client connection
type Client struct {
	con     net.Conn           // TCP connection
	core    *Core              // Shortcut to the core goroutine
	coreOut chan *MessageReply // Reply channel (to be used by the core)
}

/*****************************************************************************/

// NewClient construct a Client structure
func NewClient(con net.Conn, core *Core) (clt *Client) {
	channel := make(chan *MessageReply, channelSize)
	return &Client{con: con, core: core, coreOut: channel}
}

/*****************************************************************************/

// Reply is used by the core methods to return a reply to the client
func (clt *Client) Reply(r *MessageReply) {
	clt.coreOut <- r
}

/*****************************************************************************/

// jsonIn processes incoming JSON traffic from the client socket, decode it,
// and send messages to the core.
func (clt *Client) jsonIn() {

	// Be sure the core is notified when connection is closed
	defer func() { clt.core.in <- &MessageQuery{clt: clt, oper: OP_CLOSE} }()

	// Declare a JSON decoder
	decoder := json.NewDecoder(clt.con)
	clt.core.in <- &MessageQuery{clt: clt, oper: OP_OPEN}

	for {

		// Read an incoming message and decode it
		m := &MessageQuery{clt: clt}
		if err := decoder.Decode(m); err == io.EOF {
			break
		} else if err != nil {
			// Decoding error: notify the core, close the connection
			clt.core.in <- &MessageQuery{clt: clt, oper: OP_NONE}
			break
		}

		// Convert operation code and forward to the core
		m.oper = Service[m.Op]
		clt.core.in <- m
	}
}

/*****************************************************************************/

// jsonOut is waiting for outgoing traffic from the core, encode it in JSON
// messages, and write it to the client socket.
func (clt *Client) jsonOut() {

	// Be sure the connection is closed in the end
	defer clt.con.Close()

	// Declare a JSON encoder
	encoder := json.NewEncoder(clt.con)

	// Wait for outgoing messages from the core
	end := false
	for reply := range clt.coreOut {
		// Check closing connection notification
		if reply.oper == OP_CLOSE {
			break
		}
		// Ignore all messages after an encoding error
		if !end {
			// Encode a JSON message, and write it to the socket
			if err := encoder.Encode(reply); err != nil {
				log.Println("Error ", err)
				end = true
			}
		}
	}
}

/*****************************************************************************/

// Core is the structure representing the core goroutine, responsible on the
// logic of the application.
type Core struct {
	in    chan *MessageQuery // Incoming channel
	locks *LockArea          // Lock management data structure
	stats map[string]int64   // Key/value data structure
	count int64              // Command counter
}

/*****************************************************************************/

// NewCore builds a Core object
func NewCore() *Core {
	return &Core{
		in:    make(chan *MessageQuery, channelSize*128),
		locks: NewLockArea(),
		stats: make(map[string]int64),
	}
}

/*****************************************************************************/

// main is the main event loop of the Core goroutine
func (core *Core) main() {

	// Dequeue incoming events
	for m := range core.in {

		// Dispatch event to related function
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
			m.clt.Reply(&MessageReply{Status: "KO", Error: "Unknown operation"})
		}
		atomic.AddInt64(&core.count, 1)
	}
}

/*****************************************************************************/

// handleOpen handles open connection notifications
func (core *Core) handleOpen(query *MessageQuery) {

	if verbose {
		log.Println("Opening connection")
	}
	core.locks.AddClient(query.clt)
}

/*****************************************************************************/

// handleClose handles close connection notification
func (core *Core) handleClose(query *MessageQuery) {

	if verbose {
		log.Println("Closing connection")
	}
	// Remove client from all data structures.
	// All locks will be released.
	toBeNotified := core.locks.RemoveClient(query.clt)

	// Send reply
	query.clt.Reply(&MessageReply{oper: OP_CLOSE})

	// Forward replies to any clients for which the locks have been regranted
	for _, c := range toBeNotified {
		c.Reply(&MessageReply{Status: "OK"})
	}
}

/*****************************************************************************/

// handleLock manages locking operations
func (core *Core) handleLock(query *MessageQuery) {

	if verbose {
		log.Println("Locking", query.Target)
	}

	// Try to add the lock
	if core.locks.Add(query.clt, query.Target) {
		// Only reply if the lock has been granted
		query.clt.Reply(&MessageReply{Status: "OK"})
	}
}

/*****************************************************************************/

// handleUnlock manages any unlocking operation
func (core *Core) handleUnlock(query *MessageQuery) {

	if verbose {
		log.Println("Unlocking", query.Target)
	}

	// Try to remove the lock
	c, ok := core.locks.Remove(query.clt, query.Target)
	if ok {
		// Send reply to the client
		reply := &MessageReply{Status: "OK"}
		query.clt.Reply(reply)
		if c != nil {
			// Forward a reply to another client if the lock has been regranted
			c.Reply(reply)
		}
	} else {
		// Error: could not release the lock
		reply := &MessageReply{Status: "KO", Error: "Cannot find this lock"}
		query.clt.Reply(reply)
	}
}

/*****************************************************************************/

// handleGet implements the GET integer value operation
func (core *Core) handleGet(query *MessageQuery) {

	if verbose {
		log.Println("Getting", query.Target)
	}

	// Retrieve corresponding statistic, and format the value
	val := strconv.FormatInt(core.stats[query.Target], 10)
	query.clt.Reply(&MessageReply{Status: "OK", Value: val})
}

/*****************************************************************************/

// handleSet implements the SET integer value operation
func (core *Core) handleSet(query *MessageQuery) {

	if verbose {
		log.Println("Setting", query.Target)
	}
	var reply *MessageReply

	// Parse integer
	if n, err := strconv.ParseInt(query.Arg, 10, 64); err != nil {
		reply = &MessageReply{Status: "KO", Error: "Invalid number"}
	} else {
		// Update corresponding statistic
		core.stats[query.Target] = n
		reply = &MessageReply{Status: "OK"}
	}
	query.clt.Reply(reply)
}

/*****************************************************************************/

// handleIncr implements the INCR integer value operation
func (core *Core) handleIncr(query *MessageQuery) {

	if verbose {
		log.Println("Increment", query.Target)
	}
	var reply *MessageReply

	// Try to parse the increment
	if n, err := strconv.ParseInt(query.Arg, 10, 64); err != nil {
		reply = &MessageReply{Status: "KO", Error: "Invalid number"}
	} else {
		// Update the corresponding statistic
		core.stats[query.Target] += n
		val := strconv.FormatInt(core.stats[query.Target], 10)
		reply = &MessageReply{Status: "OK", Value: val}
	}
	query.clt.Reply(reply)
}

/*****************************************************************************/

// MainServer is the main entry point of this package. It spawns TCP listener
// and core goroutines
func MainServer(server string) {

	// Build core, and start goroutine
	core := NewCore()
	go core.main()

	// Build TCP listener and start goroutine
	lis := &Listener{core: core}
	go lis.Listen("tcp", server)

	// Register monitoring server
	go monitoringServer(core)

	// Setup SIGINT signal handler, and wait
	channel := make(chan os.Signal)
	signal.Notify(channel, os.Interrupt)
	<-channel
	log.Println("Stop")
}

/*****************************************************************************/
