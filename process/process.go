package process

import (
	"Paprika/models"
	"context"
)

type IProcess interface {
	// ctx - a global context from a process manager. You need to create an individual context inside process using parent one.
	Do(ctx context.Context, ch chan *models.Ticker) error
	Stop() error
}

// To embed into a child process
// type Process struct {
// 	ctx    context.Context // create a private context from a global
// 	cancel context.CancelFunc
// }

// common exchange process logic is here
// type BaseProcess struct {
