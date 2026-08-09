package timers

import (
	"time"

	"github.com/boreq/plum/plum-backend/app"
	"github.com/boreq/plum/plum-backend/logging"
)

type RemoveOldDataHandler interface {
	Execute(cmd app.RemoveOldData) error
}

type RemoveOldData struct {
	removeOldData RemoveOldDataHandler
	interval      time.Duration
	log           logging.Logger
}

func NewRemoveOldData(removeOldData RemoveOldDataHandler, interval time.Duration) *RemoveOldData {
	return &RemoveOldData{
		removeOldData: removeOldData,
		interval:      interval,
		log:           logging.New("entrypoints/timers.RemoveOldData"),
	}
}

func (t *RemoveOldData) Run() error {
	ticker := time.NewTicker(t.interval)
	defer ticker.Stop()

	for range ticker.C {
		if err := t.removeOldData.Execute(app.RemoveOldData{Now: time.Now()}); err != nil {
			t.log.Error("could not remove old data", "err", err)
		}
	}

	return nil
}
