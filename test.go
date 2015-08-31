package main

import (
	"fmt"
	"net"
	"time"

	log "github.com/inconshreveable/log15"

	"github.com/aarondl/ultimateq/inet"
	"github.com/aarondl/ultimateq/irc"
	"github.com/aarondl/ultimateq/parse"
)

func main() {

	conn, err := net.Dial("tcp", "irc.gamesurge.net:6667")

	if err != nil {
		fmt.Printf("E: %+V", err)
	}

	srvlog := log.New("module", "zounce")

	//srvlog.Warn("hello world", "worllllllllld", 5, "low", 1, "high", 7)

	iclient := inet.NewIrcClient(conn, srvlog, 0, 10*time.Second, 0, 0, 0)
	iclient.SpawnWorkers(true, true)
	h := irc.Helper{iclient}

	err = h.Sendf("%s %s\n", irc.NICK, "zamn-testing")
	if err != nil {
		fmt.Println(err)
	}

	err = h.Sendf("%s %s 0 * :%s\n", "USER", "zamn", "adam")
	if err != nil {
		fmt.Println(err)
	}

	ch := iclient.ReadChannel()

	disconnect := false
	for err == nil && !disconnect {
		ev, ok := <-ch
		if !ok {
			fmt.Println("disconnected")
			break
		}
		event, err := parse.Parse(ev)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("Event %#v\n", event)
		if event.Name == "PING" {
			err = h.Sendf("%s %s", irc.PONG, event.Args[0])
		}
	}

	fmt.Println("dead")
	// irc.NICK
}
