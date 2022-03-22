package mixcloud

import (
	"sync"
	"testing"
)

var testMixes = map[string]Mix{
	"mix1": {
		Key:  "mix1",
		Url:  "mix1url",
		Name: "mix1name",
	},
}

func TestUnit_Store_Put(t *testing.T) {
	type fields struct {
		Data    map[string]Mix
		limit   int
		RWMutex *sync.RWMutex
	}
	type args struct {
		m Mix
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Simple put",
			fields: fields{
				testMixes,
				0,
				&sync.RWMutex{},
			},
			args:    args{},
			wantErr: false,
		},
		{
			name: "Store full",
			fields: fields{
				testMixes,
				1,
				&sync.RWMutex{},
			},
			args:    args{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Store{
				Data:    tt.fields.Data,
				limit:   tt.fields.limit,
				RWMutex: tt.fields.RWMutex,
			}
			if err := s.Put(tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("Store.Put() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
