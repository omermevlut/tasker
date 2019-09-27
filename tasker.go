package tasker

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/omermevlut/tasker/utils"
)

// Tasker ...
type Tasker struct {
	RedisUtil utils.RedisUtilInterface
	OnRunChan chan *Task
	wg        *sync.WaitGroup
}

// New tasker
func New(redisAddr string, workers int) *Tasker {
	t := &Tasker{
		RedisUtil: utils.NewRedisUtil(redisAddr),
		wg:        &sync.WaitGroup{},
	}

	for i := 0; i < workers; i++ {
		t.wg.Add(2)

		go t.delayedQueueWorker()
		go t.defaultQueueWorker()
	}

	return t
}

// Delayed task
func (t *Tasker) Delayed(task *Task) *Tasker {
	jTask, _ := json.Marshal(task)
	t.RedisUtil.SetDelayed(jTask, float64(task.RunAt))

	return t
}

// SetOnRunChan accepts a task channel to be notified
func (t *Tasker) SetOnRunChan(task chan *Task) {
	t.OnRunChan = task
}

func (t *Tasker) delayedQueueWorker() {
	for {
		t.RedisUtil.MoveExpiredItems(time.Now().Unix())
		time.Sleep(time.Second)
	}
}

func (t *Tasker) defaultQueueWorker() {
	for {
		time.Sleep(time.Second)

		var task Task
		res := t.RedisUtil.PopFromActiveQueue()

		if res == "" {
			continue
		}

		json.Unmarshal([]byte(res), &task)
		t.OnRunChan <- &task

		if task.isExpired() {
			continue
		}

		if task.IsRepeating {
			task.setNextRun()
			t.Delayed(&task)
		}
	}
}
