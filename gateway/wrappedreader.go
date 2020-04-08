package gateway

import (
	"errors"
	"github.com/tatsuworks/czlib"
	"io"
	"io/ioutil"
	"sync"
)

type wrappedReader struct {
	io.ReadCloser
	sync.RWMutex
	closeChan chan struct{}
	isClosed bool
}

func (r *wrappedReader) Close() error {
	r.Lock()
	r.isClosed = true
	r.Unlock()

	err := r.ReadCloser.Close()

	select {
	case r.closeChan <- struct{}{}:
	default:
	}

	return err
}

func (r *wrappedReader) Read() ([]byte, error) {
	r.RLock()
	closed := r.isClosed
	r.RUnlock()

	if closed {
		return nil, errors.New("reader was closed")
	}

	select {
	case <-r.closeChan:
		return nil, errors.New("reader was closed")
	default:
		return ioutil.ReadAll(r.ReadCloser)
	}
}

func (r *wrappedReader) Reset(reader io.Reader) {
	r.ReadCloser.(czlib.Resetter).Reset(reader)
}
