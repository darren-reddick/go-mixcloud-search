package mixcloud

import (
	"fmt"
	"regexp"
	"strings"
)

type Filter struct {
	Include *regexp.Regexp
	Exclude *regexp.Regexp
}

type invalidFilterTermError struct {
	term string
	msg  string
}

func (i *invalidFilterTermError) Error() string {
	return fmt.Sprintf("%s: %s", i.term, i.msg)
}

func validateFilterTerm(s string) error {
	matched, _ := regexp.MatchString(`^\w+\z`, s)
	if !matched {
		return &invalidFilterTermError{s, "Invalid search term"}
	}

	return nil
}

func NewFilter(include []string, exclude []string) (Filter, error) {
	for _, s := range include {
		err := validateFilterTerm(s)
		if err != nil {
			return Filter{}, err
		}
	}
	f := Filter{
		regexp.MustCompile(strings.Join(include, "|")),
		regexp.MustCompile(""),
	}

	return f, nil
}

func (f *Filter) Filter(m []Mix) []Mix {
	replace := make([]Mix, 0)
	if f.Include.String() != "" {
		for _, mix := range m {
			b := []byte(mix.Key)
			if f.Include.Match(b) {
				replace = append(replace, mix)
			}
		}
	} else {

		replace = m
	}

	return replace
}
