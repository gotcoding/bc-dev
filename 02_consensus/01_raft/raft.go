package raft

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type ApplyMsg struct {
	CommandValid bool // true为log，false为snapshot
	// 向application层提交日志
	Command      interface{}
	CommandIndex int
	CommandTerm  int

	// 向application层安装快照
	Snapshot          []byte
	LastIncludedIndex int
	LastIncludedTerm  int
}

// 日志项
type LogEntry struct {
	Command interface{}
	Term    int
}

// 角色
const ROLE_LEADER = "Leader"
const ROLE_FOLLOWER = "Follower"
const ROLE_CANDIDATES = "Candidates"

type Raft struct {
	mu        sync.Mutex
	peers     []*ClientEnd
	persister *Persister
	me        int
	dead      int32

	currentTerm       int
	votedFor          int        // 记录再currentTerm任期投票给谁了
	log               []LogEntry //操作日志
	lastIncludedIndex int        // snapshot最后1个logEntry的index，没有snapshot则为0
	lastIncludedTerm  int        // snapshot最后1个logEntry的term，没有snapshot则无意义

	commitIndex int
	lastApplied int //当前应用到状态机的索引

	nextIndex  []int // 每个follower的log同步起点索引
	matchIndex []int // 每个follower的log同步进度

	// 所有服务器，选举状态
	role              string
	leaderId          int
	lastActiveTime    time.Time
	lastBroadcastTime time.Time

	applyCh chan ApplyMsg //应用层提交队列
}

func (rf *Raft) GetState() (int, bool) {
	rf.mu.Lock()
	defer rf.mu.Unlock()

	term := rf.currentTerm
	isLeader := rf.role == ROLE_LEADER
	return term, isLeader
}

func (rf *Raft) Start(command interface{}) (int, int, bool) {
	index := -1
	term := -1
	isLeader := true

	rf.mu.Lock()
	defer rf.mu.Unlock()
	// 只有leader才能写入
	if rf.role != ROLE_LEADER {
		return -1, -1, false
	}

	LogEntry := LogEntry{
		Command: command,
		Term:    rf.currentTerm,
	}
	rf.log = append(rf.log, LogEntry)
	index = rf.lastIndex()
	term = rf.currentTerm
	rf.persist()

	fmt.Printf("RaftNode[%d] Add Command, logIndex[%d] currentTerm[%d]\n", rf.me, index, term)
	return index, term, isLeader
}

func (rf *Raft) Kill() {
	atomic.StoreInt32(&rf.dead, 1)
}

func Make(addrs []string, me int, dataDir string) (rf *Raft, applyCh chan ApplyMsg, err error) {
	rf = &Raft{}
	rf.me = me
	if rf.persister, err = MakePersister(dataDir); err != nil {
		return
	}

	rf.role = ROLE_FOLLOWER
	rf.leaderId = -1
	rf.votedFor = -1
	rf.lastIncludedIndex = 0
	rf.lastIncludedTerm = 0
	rf.lastActiveTime = time.Now()
	applyCh = make(chan ApplyMsg, 1)
	rf.applyCh = applyCh

	// 读取raft持久化状态
	rf.loadPersist()
	// 向application层安装快照
	rf.installSnapshotToApplication()

	// peers
	rf.initRpcPeers(addrs)
	// me
	go rf.initRpcServer()
	// election逻辑
	go rf.electionLoop()
	// leader逻辑
	go rf.appendEntriesLoop()
	// apply逻辑
	go rf.applyLogLoop()

	fmt.Printf("Raftnode[%d]启动", me)
	return
}
