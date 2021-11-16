package transactionMgr

import "sync"

type TransactionMgr struct {
	orderingMap map[int64]string
	lock        sync.RWMutex
}

func NewTransactionMgr() *TransactionMgr {

	return &TransactionMgr{orderingMap: make(map[int64]string), lock: sync.RWMutex{}}
}

func (txMgr *TransactionMgr) ReadNounce(nounce int64) string {
	txMgr.lock.RLock()
	defer txMgr.lock.RUnlock()
	key := txMgr.orderingMap[nounce]
	return key
}

func (txMgr *TransactionMgr) WriteNounce(nounce int64, key string) bool {
	txMgr.lock.Lock()
	defer txMgr.lock.Unlock()
	if _, ok := txMgr.orderingMap[nounce]; !ok {
		txMgr.orderingMap[nounce] = key
		return true
	} else {
		return false
	}
}
