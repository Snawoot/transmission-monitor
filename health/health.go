package health

import (
	"errors"
	"fmt"
	"strings"

	"github.com/hekmon/transmissionrpc/v2"
)

var (
	NilTorrent = errors.New("nil torrent")
)

type TorrentError struct {
	code *int64
	msg  *string
}

func (te *TorrentError) Error() string {
	fCode := "<nil>"
	fMsg := "<nil>"
	if te.code != nil {
		fCode = fmt.Sprintf("%d", *te.code)
	}
	if te.msg != nil {
		fMsg = *te.msg
	}
	return fmt.Sprintf("torrent error: code=%s, message=%q", fCode, fMsg)
}

func (te *TorrentError) Code() *int64 {
	return te.code
}

func (te *TorrentError) Msg() *string {
	return te.msg
}

func newTorrentError(code *int64, msg *string) *TorrentError {
	return &TorrentError{
		code: code,
		msg:  msg,
	}
}

type ErrorCheck struct {
	skip_no_response bool
}

func NewErrorCheck(skip_no_response bool) *ErrorCheck {
	return &ErrorCheck{
		skip_no_response: skip_no_response,
	}
}

func (c *ErrorCheck) CheckTorrent(t *transmissionrpc.Torrent) error {
	if t == nil {
		return NilTorrent
	}
	if t.ErrorString != nil && *t.ErrorString != "" || t.Error != nil && *t.Error != 0 {
		if c.skip_no_response &&
			t.ErrorString != nil &&
			(strings.Contains(*t.ErrorString, "No Response") ||
				strings.Contains(*t.ErrorString, "HTTP 0")) {
			return nil
		}
		return newTorrentError(t.Error, t.ErrorString)
	}
	return nil
}
