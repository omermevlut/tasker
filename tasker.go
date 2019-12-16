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
	wg        *sync.WaitGroup
}

// New tasker
func New(redisAddr string, workers int) *Tasker {
	t := &Tasker{
		RedisUtil: utils.NewRedisUtil(redisAddr),
		wg:        &sync.WaitGroup{},
	}

	for i := 0; i < workers; i++ {
		t.wg.Add(1)

		go t.delayedQueueWorker()
	}

	return t
}

// Delayed task
func (t *Tasker) Delayed(task *Task) *Tasker {
	jTask, _ := json.Marshal(task)
	t.RedisUtil.SetDelayed(jTask, float64(task.RunAt))

	return t
}

// OnRun ...
func (t *Tasker) OnRun(callback func(t Task)) {
	go func() {
		for {
			time.Sleep(time.Second)

			var task Task
			res := t.RedisUtil.PopFromActiveQueue()

			if res == "" {
				continue
			}

			json.Unmarshal([]byte(res), &task)
			callback(task)

			if task.isExpired() || task.isReachedMaxRunCount() {
				continue
			}

			if task.IsRepeating {
				task.RunCount++

				task.setNextRun()
				t.Delayed(&task)
			}
		}
	}()
}

func (t *Tasker) delayedQueueWorker() {
	for {
		t.RedisUtil.MoveExpiredItems(time.Now().Unix())
		time.Sleep(time.Second)
	}
}
