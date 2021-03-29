package downloadtaskmgr

import (
	"encoding/binary"
	"encoding/json"
	"github.com/daqnext/meson-common/common/logger"
	"github.com/daqnext/meson-common/common/utils"
	"github.com/syndtr/goleveldb/leveldb"
	"os"
	"sync"
)

const LDBPath = "./downloadldb"
const LDBFile = "./downloadldb/index"

func LevelDBInit() {
	if !utils.Exists(LDBPath) {
		err := os.Mkdir(LDBPath, 0700)
		if err != nil {
			logger.Fatal("file dir create failed, please create dir " + LDBPath + " by manual")
		}
	}
}

var DBLock sync.Mutex

func OpenDB() (*leveldb.DB, error) {
	db, err := leveldb.OpenFile(LDBFile, nil)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func SetTaskToLDB(task *DownloadTask) {
	b := make([]byte, binary.MaxVarintLen64)
	binary.LittleEndian.PutUint64(b, task.Id)
	DBLock.Lock()
	defer DBLock.Unlock()
	db, err := OpenDB()
	if err != nil {
		logger.Error("SetTask open level db error", "err", err)
		return
	}
	defer db.Close()

	taskStr, err := json.Marshal(task)
	if err != nil {
		logger.Error("SetTask json.Marshal error", "taskId", task.Id)
		return
	}
	//fmt.Println(string(taskStr))

	err = db.Put(b, taskStr, nil)
	if err != nil {
		logger.Error("leveldb put data error", "err", err, "taskStr", taskStr)
	}
}

func DelTaskFromLDB(taskId uint64) {
	b := make([]byte, binary.MaxVarintLen64)
	binary.LittleEndian.PutUint64(b, taskId)
	DBLock.Lock()
	defer DBLock.Unlock()
	db, err := OpenDB()
	if err != nil {
		logger.Error("DelTask open level db error", "err", err)
		return
	}
	defer db.Close()

	err = db.Delete(b, nil)
	if err != nil {
		logger.Error("DelTask from level db error", "err", err)
	}
}

func LoopTasksInLDB() []*DownloadTask {
	DBLock.Lock()
	defer DBLock.Unlock()
	db, err := OpenDB()
	if err != nil {
		logger.Error("LoopTasksInDB open level db error", "err", err)
		return nil
	}
	defer db.Close()
	iter := db.NewIterator(nil, nil)
	tasks := []*DownloadTask{}
	for iter.Next() {
		// Remember that the contents of the returned slice should not be modified, and
		// only valid until the next call to Next.
		//key := iter.Key()
		value := iter.Value()
		var task DownloadTask
		err := json.Unmarshal(value, &task)
		if err != nil {
			logger.Error("LoopTasksInDB Unmarshal error", "str", string(value))
			continue
		}
		tasks = append(tasks, &task)
	}
	iter.Release()
	err = iter.Error()
	if err != nil {
		logger.Error("loop iter error", "err", err)
	}
	return tasks
}
