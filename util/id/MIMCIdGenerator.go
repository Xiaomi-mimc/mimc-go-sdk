package id

import (
	"github.com/Xiaomi-mimc/mimc-go-sdk/common/constant"
	"github.com/Xiaomi-mimc/mimc-go-sdk/util/string"
	"strconv"
	"sync"
)

var idGenerator *IdGenerator
var lock *sync.Mutex = &sync.Mutex{}

func Generate() *string {
	if idGenerator == nil {
		lock.Lock()
		defer lock.Unlock()
		if idGenerator == nil {
			idGenerator = &IdGenerator{}
		}
	}
	id := strutil.RandomStrWithLength(10) + "_" + *(idGenerator.generate())
	return &id
}

type IdGenerator struct {
	counter uint64
}

func (this *IdGenerator) generate() *string {
	this.counter += uint64(cnst.MIMC_COUNTER_VALUE)
	str := strconv.FormatUint(this.counter, 10)
	return &str
}
