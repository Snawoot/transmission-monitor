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
	log.Printf("torrent error: hash=%s error: %v", *t.HashString, reason)
	return nil
}
