package raft

import "sync/atomic"

// 最后的index
func (rf *Raft) lastIndex() int {
	return rf.lastIncludedIndex + len(rf.log)
}

// 最后的Term
func (rf *Raft) lastTerm() (lastLogTerm int) {
	if len(rf.log) != 0 {
		lastLogTerm = rf.log[len(rf.log)-1].Term
	} else {
		lastLogTerm = rf.lastIncludedTerm
	}
	return
}

// 日志index转化为log数组下标
func (rf *Raft) index2LogPos(index int) (pos int) {
	return index - rf.lastIncludedIndex - 1
}

func (rf *Raft) killed() bool {
	z := atomic.LoadInt32(&rf.dead)
	return z == 1
}
