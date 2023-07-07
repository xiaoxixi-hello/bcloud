package download

import (
	"fmt"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"strconv"
	"sync"
)

type MultiThread struct {
	accessToken string
	base        string
	count       int
	id          string
	mutex       sync.Mutex
	jobs        chan *DownDetail
	db          *gorm.DB
}

func (m *MultiThread) ProducerDown(f []*DownDetail) {
	go func() {
		for _, v := range f {
			m.jobs <- v
		}
		//	close(m.jobs)
	}()
}

func (m *MultiThread) ConsumerDown(ch <-chan *DownDetail, i int) {
	for c := range ch {
		if err := download(c.Dlink, c.Path, m.base, m.id, m.accessToken, m.db); err != nil {
			if err.Error() == "该任务需要删除" {
				continue
			}
			m.jobs <- c
			zap.L().Info("任务重新入队列成功,", zap.String("进程任务ID", m.id), zap.String("name", c.Name), zap.Error(err))
			zap.L().Info(fmt.Sprintf("剩余任务个数:%d", m.count+1), zap.String("任务进程ID", m.id), zap.Int("当前下载协程ID", i+1))
			continue
		}
		zap.L().Info(fmt.Sprintf("剩余任务个数:%d", m.count+1), zap.String("任务进程ID", m.id), zap.Int("当前下载协程ID", i+1))
		m.mutex.Lock()
		m.count = m.count - 1
		if m.count == 0 {
			close(m.jobs)
		}
		m.mutex.Unlock()
	}
}

func MultiThreadDownRun(f []*DownDetail, base, id, token string, db *gorm.DB, processCount string) error {
	m := MultiThread{
		accessToken: token,
		base:        base,
		count:       len(f),
		id:          id,
		mutex:       sync.Mutex{},
		jobs:        make(chan *DownDetail, 10),
		db:          db,
	}
	m.ProducerDown(f)
	var wg sync.WaitGroup
	atomic, _ := strconv.Atoi(processCount)
	for i := 0; i < atomic; i++ {
		wg.Add(1)
		go func(i int) {
			m.ConsumerDown(m.jobs, i)
			wg.Done()
		}(i)
	}
	wg.Wait()
	//m.ConsumerDown(m.jobs)
	if m.count == 0 {
		zap.L().Info("下载成功", zap.String("进程任务ID", m.id))
		return nil
	}
	return nil
}
