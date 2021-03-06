package util

// These are misc. utility functions that didn't really fit anywhere else

import "runtime"
import "sync"
import "time"

// A wrapper around runtime.Gosched() so it doesn't need to be imported elsewhere.
func Yield() {
	runtime.Gosched()
}

// A wrapper around runtime.LockOSThread() so it doesn't need to be imported elsewhere.
func LockThread() {
	runtime.LockOSThread()
}

// A wrapper around runtime.UnlockOSThread() so it doesn't need to be imported elsewhere.
func UnlockThread() {
	runtime.UnlockOSThread()
}

// This is used to buffer recently used slices of bytes, to prevent allocations in the hot loops.
var byteStoreMutex sync.Mutex
var byteStore [][]byte

// Gets an empty slice from the byte store.
func GetBytes() []byte {
	byteStoreMutex.Lock()
	defer byteStoreMutex.Unlock()
	if len(byteStore) > 0 {
		var bs []byte
		bs, byteStore = byteStore[len(byteStore)-1][:0], byteStore[:len(byteStore)-1]
		return bs
	} else {
		return nil
	}
}

// Puts a slice in the store.
func PutBytes(bs []byte) {
	byteStoreMutex.Lock()
	defer byteStoreMutex.Unlock()
	byteStore = append(byteStore, bs)
}

// This is a workaround to go's broken timer implementation
func TimerStop(t *time.Timer) bool {
	stopped := t.Stop()
	select {
	case <-t.C:
	default:
	}
	return stopped
}

// Run a blocking function with a timeout.
// Returns true if the function returns.
// Returns false if the timer fires.
// The blocked function remains blocked--the caller is responsible for somehow killing it.
func FuncTimeout(f func(), timeout time.Duration) bool {
	success := make(chan struct{})
	go func() {
		defer close(success)
		f()
	}()
	timer := time.NewTimer(timeout)
	defer TimerStop(timer)
	select {
	case <-success:
		return true
	case <-timer.C:
		return false
	}
}

// This calculates the difference between two arrays and returns items
// that appear in A but not in B - useful somewhat when reconfiguring
// and working out what configuration items changed
func Difference(a, b []string) []string {
	ab := []string{}
	mb := map[string]bool{}
	for _, x := range b {
		mb[x] = true
	}
	for _, x := range a {
		if !mb[x] {
			ab = append(ab, x)
		}
	}
	return ab
}
