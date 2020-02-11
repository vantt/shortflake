package idgenerator

// see https://github.com/sumory/uc/blob/master/src/com/sumory/uc/id/IdWorker.java for java implement
// see https://github.com/sumory/idgen/blob/master/idgen.go  for origin go-lang implementation

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

const (
	idgenEpoch           int64 = 1417942588000 //2014-12-07 16:56:28
	worker_id_bits       uint8 = 10
	Max_worker_id        int32 = -1 ^ (-1 << worker_id_bits) //1023
	sequence_bits        uint8 = 12
	worker_id_shift      uint8 = sequence_bits                  // 12
	timestamp_left_shift uint8 = sequence_bits + worker_id_bits // 22
	sequence_mask        int64 = -1 ^ (-1 << sequence_bits)     // 4095
)

type IdWorker struct {
	workerId      int64
	lastTimestamp int64
	sequence      int64
	lock          sync.Mutex
}

func NewIdWorker(worker_id int32) (*IdWorker, error) {

	if worker_id > Max_worker_id || worker_id < 0 {
		return nil, errors.New(fmt.Sprintf("WorkerId must be between 1 and 1023. [%d] was given.", worker_id))
	}

	id_worker := &IdWorker{workerId: int64(worker_id), lastTimestamp: -1, sequence: 0}

	return id_worker, nil
}

func timeGen() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func tilNextMillis(lastTimestamp int64) int64 {
	timestamp := timeGen()

	for timestamp <= lastTimestamp {
		timestamp = timeGen()
	}

	return timestamp
}

// need synchronized
func (worker *IdWorker) NextId() (int64, error) {
	worker.lock.Lock()
	defer worker.lock.Unlock()

	timestamp := timeGen()
	if timestamp < worker.lastTimestamp {
		return 0, errors.New(fmt.Sprintf("Clock moved backwards. Refusing to generate worker for %d milliseconds", worker.lastTimestamp-timestamp))
	}

	if worker.lastTimestamp == timestamp {
		worker.sequence = (worker.sequence + 1) & sequence_mask
		if worker.sequence == 0 {
			timestamp = tilNextMillis(worker.lastTimestamp)
		}
	} else {
		worker.sequence = 0
	}

	worker.lastTimestamp = timestamp
	return ((timestamp - idgenEpoch) << timestamp_left_shift) | (worker.workerId << worker_id_shift) | worker.sequence, nil
}

func (worker *IdWorker) NextIds(num_ids uint16) ([]int64, error) {
	var new_ids = make([]int64, num_ids)
	var i uint16
	var err error

	for i = 0; i < num_ids; i++ {
		new_ids[i], err = worker.NextId()

		if err != nil {
			return new_ids, err
		}
	}

	return new_ids, nil
}
