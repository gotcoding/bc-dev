package raft

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"
)

// Persister 持久化保存
type Persister struct {
	mu        sync.Mutex
	raftState []byte
	snapshot  []byte

	raftStateMemSynced  bool
	snapshotInMemSynced bool

	dir string // 目录
}

// 创建持久化目录
func MakePersister(dir string) (ps *Persister, err error) {
	if err = os.Mkdir(dir, 0777); err != nil {
		return
	}
	ps = &Persister{
		dir: dir,
	}
	return
}

// 持久化raft状态到磁盘
func (ps *Persister) syncRaftStateToDisk() {
	flieName := path.Join(ps.dir, "rf.state")
	ps.mustWriteFile(flieName, ps.raftState)
}

// 从磁盘中加载raft状态
func (ps *Persister) loadRaftStateFromDisk() {
	if ps.raftStateMemSynced {
		return
	}
	ps.raftStateMemSynced = true

	fileName := path.Join(ps.dir, "rf.state")
	ps.raftState = ps.mustReadFile(fileName)
}

func (ps *Persister) SaveRaftState(state []byte) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	ps.raftState = state
	ps.raftStateMemSynced = true
	ps.syncRaftStateToDisk()
}

func (ps *Persister) ReadRaftState() []byte {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	ps.loadRaftStateFromDisk()
	return ps.raftState
}

func (ps *Persister) RaftStateSize() int {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	ps.loadRaftStateFromDisk()
	return len(ps.raftState)
}

// 将snapshot同步到磁盘
func (ps *Persister) syncSnapshotToDisk(lastIncludedIndex int) {
	fileName := path.Join(ps.dir, fmt.Sprintf("rf.snapshot-%d", lastIncludedIndex))
	ps.mustWriteFile(fileName, ps.snapshot)
}

// 清理更新index的snapshot文件
func (ps *Persister) cleanOlderSnapshot(lastIncludedIndex int) {
	if fileList, err := ioutil.ReadDir(ps.dir); err != nil {
		log.Fatal(err)
	} else {
		for _, file := range fileList {
			if file.IsDir() {
				continue
			}
			p := file.Name()
			ext := path.Ext(p)
			if !strings.HasPrefix(ext, ".snapshot-") {
				continue
			}
			if idx, err := strconv.Atoi(p[len("rf.snapshot-"):]); err != nil {
				continue
			} else {
				if idx < lastIncludedIndex {
					fullPath := path.Join(ps.dir, p)
					os.Remove(fullPath)
					fmt.Printf("删除过期snapshot: %v\n", fullPath)
				}
			}
		}
	}
}

func (ps *Persister) SaveStateAndSnapshot(state, snapshot []byte, lastIncludedIndex int) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	ps.raftState = state
	ps.snapshot = snapshot
	ps.raftStateMemSynced = true
	ps.snapshotInMemSynced = true

	ps.syncSnapshotToDisk(lastIncludedIndex)
	ps.syncRaftStateToDisk()
	ps.cleanOlderSnapshot(lastIncludedIndex)
}

func (ps *Persister) loadSnapshotFromDisk(lastIncludedIndex int) {
	if ps.snapshotInMemSynced {
		return
	}
	ps.snapshotInMemSynced = true

	fileName := path.Join(ps.dir, fmt.Sprintf("rf.snapshot-%d", lastIncludedIndex))
	ps.snapshot = ps.mustReadFile(fileName)
}

func (ps *Persister) ReadSnapshot(lastIncludedIndex int) []byte {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	ps.loadSnapshotFromDisk(lastIncludedIndex)
	return ps.snapshot
}

func (ps *Persister) SnapshotSize(lastIncludedIndex int) int {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	ps.loadSnapshotFromDisk(lastIncludedIndex)
	return len(ps.snapshot)
}

func (ps *Persister) mustReadFile(fileName string) (data []byte) {
	var err error
	var file *os.File

	// 文件不存在则返回
	if _, err := os.Stat(fileName); err != nil {
		if os.IsNotExist(err) {
			return
		}
		log.Fatal(err)
	}
	if file, err = os.Open(fileName); err != nil {
		log.Fatal(err)
	}
	if data, err = ioutil.ReadAll(file); err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	return
}

func (ps *Persister) mustWriteFile(fileName string, data []byte) {
	var file *os.File
	var err error
	tmpFileName := fileName + "-tmp"
	if file, err = os.Create(tmpFileName); err != nil {
		log.Fatal(err)
	}
	if _, err := file.Write(data); err != nil {
		log.Fatal(err)
	}
	if err = file.Sync(); err != nil {
		log.Fatal(err)
	}
	file.Close()
	if err = os.Rename(tmpFileName, fileName); err != nil {
		log.Fatal(err)
	}
}
