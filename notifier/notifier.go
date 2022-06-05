package notifier

import (
	"log"

	"github.com/hekmon/transmissionrpc/v2"
)

type LogNotifier struct {}

func NewLogNotifier() LogNotifier {
	return LogNotifier{}
}

func (_ LogNotifier) Notify(t *transmissionrpc.Torrent, reason error) error {
	hash := ""
	if t.HashString != nil {
		hash = *t.HashString
	}
	log.Printf("torrent error: hash=%s error: %v", hash, reason)
	return nil
}
