package monitor

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/hekmon/transmissionrpc/v2"
	"github.com/hashicorp/go-multierror"
)

type Checker interface {
	CheckTorrent(t *transmissionrpc.Torrent) error
}

type Notifier interface {
	Notify(t *transmissionrpc.Torrent, err error) error
}

type DBQuerier interface {
	Get(key string) ([]byte, error)
	Set(key string, data []byte) error
}

type Monitor struct {
	checker  Checker
	db       DBQuerier
	notifier Notifier
}

type MonitorRecord struct {
	When time.Time `json:"when"`
	Reason string `json:"reason,omitempty"`
	Torrent *transmissionrpc.Torrent `json:"torrent,omitempty"`
}

func NewMonitor(db DBQuerier, checker Checker, notifier Notifier) *Monitor {
	return &Monitor{
		checker:  checker,
		db:       db,
		notifier: notifier,
	}
}

func (m *Monitor) Process(torrents []*transmissionrpc.Torrent) (res error) {
	for _, torrent := range torrents {
		if err := m.checker.CheckTorrent(torrent); err != nil {
			if torrent.HashString == nil {
				res = multierror.Append(
					res,
					errors.New("met torrent w/o hash string, don't know how to process it"),
				)
				continue
			}
			hash := *torrent.HashString

			dbRes, dbErr := m.db.Get(hash)
			if dbErr != nil {
				res = multierror.Append(res, fmt.Errorf("hash %q lookup failed: %w", dbErr))
			} else if dbRes != nil {
				// we already know this torrent is faulty
				continue
			}

			res = multierror.Append(
				res,
				fmt.Errorf("torrent hash=%q, error: %w", hash, err),
			)

			notifyErr := m.notifier.Notify(torrent, err)
			if notifyErr != nil {
				res = multierror.Append(
					res,
					fmt.Errorf("notification about %q failed: %w", hash, notifyErr),
				)
			} else {
				saveErr := m.saveRecord(torrent, err)
				if saveErr != nil {
					res = multierror.Append(
						res,
						fmt.Errorf("can't save record to database for hash %q: %w", hash, err),
					)
				}
			}
		}
	}
	return
}

func (m *Monitor) saveRecord(t *transmissionrpc.Torrent, reason error) error {
	record := &MonitorRecord{
		When: time.Now().UTC().Truncate(0),
		Reason: reason.Error(),
		Torrent: t,
	}
	b, err := json.Marshal(record)
	if err != nil {
		return fmt.Errorf("can't marshal monitor record: %w", err)
	}
	if err := m.db.Set(*t.HashString, b); err != nil {
		return fmt.Errorf("can't save monitor record to database: %w", err)
	}
	return nil
}
