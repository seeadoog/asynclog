package asynclog

import (
	"reflect"
	"sync"
	"sync/atomic"

	"go.uber.org/zap"
)

type LogWatch[T zap.Logger | zap.SugaredLogger] struct {
	conf *LogConf
	inst atomic.Pointer[T]

	lock sync.Mutex

	typ reflect.Type

	factory func(c *LogConf) (*T, error)
}

func NewLogWatch[T zap.Logger | zap.SugaredLogger](conf *LogConf, factory func(c *LogConf) (*T, error)) (*LogWatch[T], error) {
	lw := &LogWatch[T]{
		factory: factory,
		conf:    conf,
	}
	log, err := factory(conf)
	if err != nil {
		return nil, err
	}
	lw.inst.Store(log)
	return lw, nil
}

func (l *LogWatch[T]) Get() *T {
	return l.inst.Load()
}

func (l *LogWatch[T]) Update(conf *LogConf) (error, bool) {
	l.lock.Lock()
	defer l.lock.Unlock()

	if reflect.DeepEqual(l.conf, conf) {
		return nil, false
	}

	log, err := l.factory(conf)
	if err != nil {
		return err, false
	}
	l.conf = conf
	l.inst.Store(log)
	return nil, true
}
