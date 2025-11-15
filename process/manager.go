package process

import (
	"Paprika/publisher"
	"context"
	"fmt"
	"log"
	"sync"
)

type Manager struct {
	p  map[string]IProcess // {"pid": IProcess}
	pb *publisher.Publisher

	stats *GlobalInformation

	wg     sync.WaitGroup
	ctx    context.Context
	cancel context.CancelFunc
}

func NewManager(ctx context.Context, pb *publisher.Publisher) *Manager {
	ctx, cancel := context.WithCancel(ctx)

	return &Manager{
		p:      make(map[string]IProcess),
		pb:     pb,
		stats:  newGlobalInformation(),
		ctx:    ctx,
		cancel: cancel,
	}
}

func (this *Manager) get(processName string) IProcess {
	return this.p[processName]
}

func (this *Manager) set(processName string, process IProcess) {
	this.p[processName] = process
}

func (this *Manager) Spawn(pid string, process IProcess) error {
	if _, exists := this.p[pid]; exists {
		return fmt.Errorf("failed to create new process, because process with the same pid ('%s') already exists!", pid)
	}

	this.set(pid, process)

	// run the process here
	this.wg.Add(1)
	go func() {
		var err error
		defer func() {
			this.stats.end(pid, err)
			this.wg.Done()
		}()

		this.stats.start(pid)
		if err := this.get(pid).Do(this.ctx, this.pb); err != nil {
			err = err
			log.Printf("[ERROR]: process with pid '%s' stopped with error: %v", pid, err)
			return
		}
		log.Printf("[INFO]: process with pid '%s' stopped gracefully", pid)
	}()

	return nil
}

func (this *Manager) GetStats(pid string) (*ProcessInformation, bool) {
	if _, ok := this.p[pid]; !ok {
		return nil, false
	}

	return this.stats.get(pid), true
}

func (this *Manager) StopAll() {
	this.cancel()

	for _, proc := range this.p {
		_ = proc.Stop()
	}

	this.wg.Wait()                     // wait to all process stop
	this.p = make(map[string]IProcess) // clean up
}

func (this *Manager) Stop(pid string) error {
	if _, exists := this.p[pid]; !exists {
		return fmt.Errorf("failed to stop '%s' process, because process with the provided pid does not exist!", pid)
	}

	return this.p[pid].Stop()
}
