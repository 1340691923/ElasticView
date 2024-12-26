package gm_operater_log

import (
	"context"
	"github.com/1340691923/ElasticView/pkg/infrastructure/logger"
	"github.com/1340691923/ElasticView/pkg/infrastructure/model"
	"go.uber.org/zap"
	"sync"
	"time"
)

type ReportAcceptStatus struct {
	buffer        []*model.GmOperaterLog
	bufferMutex   *sync.RWMutex
	batchSize     int
	flushInterval int
	log           logger.AppLogger
}

func NewReportAcceptStatus() *ReportAcceptStatus {
	reportAcceptStatus := &ReportAcceptStatus{
		buffer:        make([]*model.GmOperaterLog, 0, 300),
		bufferMutex:   new(sync.RWMutex),
		batchSize:     300,
		flushInterval: 10,
	}

	return reportAcceptStatus
}

func (this *ReportAcceptStatus) Flush() (err error) {

	this.bufferMutex.Lock()
	defer this.bufferMutex.Unlock()
	if len(this.buffer) == 0 {
		return nil
	}

	this.buffer = make([]*model.GmOperaterLog, 0, this.batchSize)
	return nil
}

func (this *ReportAcceptStatus) Add(data *model.GmOperaterLog) (err error) {
	this.bufferMutex.Lock()
	this.buffer = append(this.buffer, data)
	this.bufferMutex.Unlock()

	if this.getBufferLength() >= this.batchSize {
		err := this.Flush()
		return err
	}

	return nil
}

func (this *ReportAcceptStatus) getBufferLength() int {
	this.bufferMutex.RLock()
	defer this.bufferMutex.RUnlock()
	return len(this.buffer)
}

func (this *ReportAcceptStatus) FlushAll() error {
	for this.getBufferLength() > 0 {
		if err := this.Flush(); err != nil {
			return err
		}
	}
	return nil
}

func (this *ReportAcceptStatus) Run(ctx context.Context) error {

	ticker := time.NewTicker(time.Duration(this.flushInterval) * time.Second)
	run := true

	for run {
		select {
		case <-ticker.C:
			if this.getBufferLength() > 300 {
				if err := this.Flush(); err != nil {
					this.log.Error("err", zap.Error(err))
				}
			}
		case <-ctx.Done():
			run = false
		}
	}

	this.FlushAll()

	return ctx.Err()
}
