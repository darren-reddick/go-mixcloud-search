package mixcloud

import (
	"net/url"
	"regexp"
	"testing"

	"go.uber.org/zap"
)

var logger *zap.Logger

func init() {
	logger, _ = zap.NewProduction()
}

func dataKeys(m map[string]Mix) []string {
	keys := make([]string, len(m))

	i := 0
	for k := range m {
		keys[i] = k
		i++
	}
	return keys

}

func equal(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func TestInt_Search_Get(t *testing.T) {
	type fields struct {
		Term   string
		Filter Filter
		Client ClientIface
		Url    url.URL
		Store  Store
	}
	type args struct {
		offset int
	}
	tests := []struct {
		name          string
		fields        fields
		args          args
		want          bool
		wantErr       bool
		wantStoreKeys []string
	}{
		{
			name: "Simple",
			fields: fields{
				Term: "akumad",
				Filter: Filter{
					regexp.MustCompile("essaie"),
					regexp.MustCompile(""),
				},
				Client: NewMockClient(sampleResponse1),
				Url:    url.URL{},
				Store:  NewStore(0),
			},
			args:          args{0},
			want:          false,
			wantErr:       false,
			wantStoreKeys: []string{"/AkumaDJ/zouklove-essaie-1/"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Search{
				Term:   tt.fields.Term,
				Filter: tt.fields.Filter,
				Client: tt.fields.Client,
				Url:    tt.fields.Url,
				Store:  tt.fields.Store,
				logger: logger,
			}
			got, err := a.Get(tt.args.offset, DefaultPageLimit)
			if (err != nil) != tt.wantErr {
				t.Errorf("Search.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Search.Get() = %v, want %v", got, tt.want)
			}
			d := dataKeys(a.Data)
			if !equal(tt.wantStoreKeys, d) {
				t.Errorf("Wanted %v but got %v", tt.wantStoreKeys, d)
			}
		})
	}
}

func TestInt_Search_GetAllParallel(t *testing.T) {
	mockclient := NewMockPagingClient(10, 5)
	filter, _ := NewFilter([]string{""}, []string{""})
	store := NewStore(0)
	search, _ := NewMixSearch("a", filter, &mockclient, store, logger)

	_ = search.GetAllParallel()

	datalen := len(search.Data)

	if datalen != 50 {
		t.Errorf("Wanted 50 items but got %d", datalen)
	}
}

func TestInt_Search_GetAllParallelStoreLimit(t *testing.T) {
	mockclient := NewMockPagingClient(10, 5)
	filter, _ := NewFilter([]string{""}, []string{""})
	store := NewStore(5)
	search, _ := NewMixSearch("a", filter, &mockclient, store, logger)

	_ = search.GetAllParallel()

	datalen := len(search.Data)

	if datalen != 5 {
		t.Errorf("Wanted 5 items but got %d", datalen)
	}
}
