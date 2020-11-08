package tools

import (
	"errors"
	"sync"
	"time"
)

const (
	baseTimeStamp int64 = 1580572800
	workerBits    uint8 = 10
	numberBits    uint8 = 12

	workerMax int64 = -1 ^ (-1 << workerBits)
	numberMax int64 = -1 ^ (-1 << numberBits)

	timeShift   uint8 = workerBits + numberBits
	workerShift uint8 = numberBits
)

type Worker struct {
	mu         sync.Mutex
	timestamps int64
	workerId   int64
	number     int64
}

func NewWorker(workerId int64) (*Worker, error) {
	if workerId < 0 || workerId > workerMax {
		return nil, errors.New("workerId error")
	}

	return &Worker{
		timestamps: 0,
		workerId:   workerId,
		number:     0,
	}, nil
}

func (w *Worker) GetId() int64 {
	w.mu.Lock()
	defer w.mu.Unlock()

	currentTimeStamp := time.Now().UnixNano() / 1000000
	if w.timestamps == currentTimeStamp {
		w.number++
		if w.number > numberMax {
			for currentTimeStamp < w.timestamps {
				currentTimeStamp = time.Now().UnixNano() / 1000000
			}
		}
	} else {
		w.number = 0
		w.timestamps = currentTimeStamp
	}
	ID := int64((currentTimeStamp-baseTimeStamp)<<timeShift | (w.workerId << workerShift) | (w.number))

	return ID
}
