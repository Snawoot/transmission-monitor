package notifier

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/hekmon/transmissionrpc/v2"
)

type LogNotifier struct{}

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

type CommandExportRecord struct {
	Reason  string                   `json:"reason"`
	Torrent *transmissionrpc.Torrent `json:"torrent"`
}
type CommandNotifier struct {
	command []string
	timeout time.Duration
}

func NewCommandNotifier(timeout time.Duration, command []string) *CommandNotifier {
	return &CommandNotifier{
		command: command,
		timeout: timeout,
	}
}

func (n *CommandNotifier) Notify(t *transmissionrpc.Torrent, reason error) error {
	if len(n.command) == 0 {
		return errors.New("empty command provided to CommandNotifier")
	}

	ctx, cl := context.WithTimeout(context.Background(), n.timeout)
	defer cl()

	subprocess := exec.CommandContext(ctx, n.command[0], n.command[1:]...)
	subprocess.Stdout = os.Stdout
	subprocess.Stderr = os.Stderr
	pipe, err := subprocess.StdinPipe()
	if err != nil {
		return fmt.Errorf("unable to get subprocess stdin pipe: %w", err)
	}
	defer pipe.Close()

	if err := subprocess.Start(); err != nil {
		return fmt.Errorf("unable to start subprocess: %w", err)
	}

	enc := json.NewEncoder(pipe)
	enc.SetIndent("", "\t")
	if err := enc.Encode(&CommandExportRecord{
		Reason:  reason.Error(),
		Torrent: t,
	}); err != nil {
		return fmt.Errorf("unable to export record: %w", err)
	}
	pipe.Close()

	if err := subprocess.Wait(); err != nil {
		return fmt.Errorf("subprocess error: %w", err)
	}

	return nil
}
