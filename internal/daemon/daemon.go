package daemon

import (
	"chat-transport/internal/entities"
	"log"
	"sync"
	"time"
)

// Daemon ...
type Daemon struct {
	Srs      []entities.Transport
	Dst      []entities.Transport
	interval time.Duration
}

// NewDaemon ...
func NewDaemon(src, dst []entities.Transport, interval time.Duration) *Daemon {
	return &Daemon{
		Srs:      src,
		Dst:      dst,
		interval: interval,
	}
}

func (d *Daemon) resend(tr entities.Transport) {
	for {
		select {
		case <-time.Tick(d.interval):
			log.Printf("starting consume updates after %s chat: \"%s\"", d.interval, tr.GetName())

			messages, err := tr.GetNewMessages()
			if err != nil {
				log.Printf("fault get new messages chat: \"%s\" err: \"%s\"", tr.GetName(), err)
			}

			for _, dt := range d.Dst {
				for _, msg := range messages {
					if err := dt.SendMessage(msg); err != nil {
						log.Printf("fault send message src_chat: \"%s\" dst_chat: \"%s\" err: \"%s\"", tr.GetName(), dt.GetName(), err)
					}
				}
			}

		}
	}
}

// Run ...
func (d *Daemon) Run() error {
	var wg sync.WaitGroup

	for _, transport := range d.Srs {
		wg.Add(1)
		go func(tr entities.Transport) {
			defer wg.Done()
			d.resend(tr)
		}(transport)
	}

	wg.Wait()

	return nil
}
