package mixcloud

import (
	"reflect"
	"regexp"
	"testing"
)

func TestNewFilter(t *testing.T) {
	type args struct {
		include []string
		exclude []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"simple positive",
			args{
				[]string{"word"},
				[]string{},
			},
			false,
		},
		{
			"positive / negative mix",
			args{
				[]string{"word", " badword"},
				[]string{},
			},
			true,
		},
		{
			"simple negative",
			args{
				[]string{"wo:rd"},
				[]string{},
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := NewFilter(tt.args.include, tt.args.exclude); (err != nil) != tt.wantErr {
				t.Errorf("NewFilter() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFilter_Filter(t *testing.T) {
	var testMixes = []Mix{
		{
			Key: "/ElectronicBunker/elephant-essential-mix-2011-07-09/",
		},
		{
			Key: "/ElectronicBunker/elephant+tiger-mix/",
		},
		{
			Key: "/ElectronicBunker/tiger-mix/",
		},
	}

	type fields struct {
		Include *regexp.Regexp
		Exclude *regexp.Regexp
	}
	type args struct {
		m []Mix
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []Mix
	}{
		{
			name: "Simple include",
			fields: fields{
				regexp.MustCompile("elephant"),
				regexp.MustCompile(""),
			},
			args: args{
				testMixes,
			},
			want: []Mix{
				{
					Key: "/ElectronicBunker/elephant-essential-mix-2011-07-09/",
				},
				{
					Key: "/ElectronicBunker/elephant+tiger-mix/",
				},
			},
		},
		{
			name: "No match",
			fields: fields{
				regexp.MustCompile("goldfish"),
				regexp.MustCompile(""),
			},
			args: args{
				testMixes,
			},
			want: []Mix{},
		},
		{
			name: "Empty filters",
			fields: fields{
				regexp.MustCompile(""),
				regexp.MustCompile(""),
			},
			args: args{
				testMixes,
			},
			want: testMixes,
		},
		{
			name: "Include + Exclude Combo",
			fields: fields{
				regexp.MustCompile("elephant"),
				regexp.MustCompile("tiger"),
			},
			args: args{
				testMixes,
			},
			want: []Mix{
				{
					Key: "/ElectronicBunker/elephant-essential-mix-2011-07-09/",
				},
			},
		},
		{
			name: "Include either",
			fields: fields{
				regexp.MustCompile("elephant|tiger"),
				regexp.MustCompile(""),
			},
			args: args{
				testMixes,
			},
			want: testMixes,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &Filter{
				Include: tt.fields.Include,
				Exclude: tt.fields.Exclude,
			}
			if got := f.Filter(tt.args.m); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Filter.Filter() = %v, want %v", got, tt.want)
			}
		})
	}
}
