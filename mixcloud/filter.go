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
		regexp.MustCompile(strings.Join(exclude, "|")),
	}

	return f, nil
}

func (f *Filter) Filter(m []Mix) []Mix {
	included := make([]Mix, 0)
	excluded := make([]Mix, 0)
	if f.Include.String() != "" {
		for _, mix := range m {
			b := []byte(mix.Key)
			if f.Include.Match(b) {
				included = append(included, mix)
			}
		}
	} else {

		included = m
	}

	if f.Exclude.String() != "" {
		for _, mix := range included {
			b := []byte(mix.Key)
			if !f.Exclude.Match(b) {
				excluded = append(excluded, mix)
			}
		}
	} else {

		excluded = included
	}

	return excluded
}
