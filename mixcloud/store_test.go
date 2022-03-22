//go:build unit

package mixcloud

import (
	"sync"
	"testing"
)

func TestStore_Put(t *testing.T) {
	type fields struct {
		Data    map[string]Mix
		RWMutex *sync.RWMutex
	}
	type args struct {
		m Mix
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Store{
				Data:    tt.fields.Data,
				RWMutex: tt.fields.RWMutex,
			}
			s.Put(tt.args.m)
		})
	}
}
