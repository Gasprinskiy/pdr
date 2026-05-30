package render_pool

import (
	"time"

	"github.com/klippa-app/go-pdfium"
	"github.com/klippa-app/go-pdfium/multi_threaded"
)

type PoolConfig struct {
	WaitDuration  time.Duration
	WorkerBinPath string
	MinIdle       int
	MaxIdle       int
	MaxTotal      int
}

type Pool struct {
	pool         pdfium.Pool
	waitDuration time.Duration
}

func NewPool(conf PoolConfig) *Pool {
	return &Pool{
		pool: multi_threaded.Init(multi_threaded.Config{
			MinIdle:  conf.MinIdle,
			MaxIdle:  conf.MaxIdle,
			MaxTotal: conf.MaxTotal,
			Command: multi_threaded.Command{
				BinPath: conf.WorkerBinPath,
				// Args:    []string{"run", "pdfium/worker/main.go"},
			},
		}),
		waitDuration: conf.WaitDuration,
	}
}

func (p *Pool) Instance() (pdfium.Pdfium, error) {
	return p.pool.GetInstance(p.waitDuration)
}

func (p *Pool) CloseInstances() error {
	return p.pool.Close()
}
