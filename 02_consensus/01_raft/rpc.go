package raft

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
)

type ClientEnd struct {
	addr string
}

func (clientEnd *ClientEnd) Call(serviceMethod string, args, reply interface{}) bool {
	clent, err := rpc.DialHTTP("tcp", clientEnd.addr)
	if err != nil {
		return false
	}
	defer clent.Close()

	if err := clent.Call(serviceMethod, args, reply); err != nil {
		return false
	}
	return true
}

func (rf *Raft) initRpcPeers(addrs []string) {
	peers := make([]*ClientEnd, 0)
	for _, addr := range addrs {
		peers = append(peers, &ClientEnd{addr})
	}
	rf.peers = peers
	return
}

func (rf *Raft) initRpcServer() {
	server := rpc.NewServer()
	server.Register(rf)

	listener, err := net.Listen("tcp", rf.peers[rf.me].addr)
	if err != nil {
		log.Fatal(err)
	}
	if err := http.Serve(listener, server); err != nil {
		log.Fatal(err)
	}
	return
}
