package graphb

import (
	"fmt"
	"github.com/pkg/errors"
	"strings"
)

type Fragment struct {
	Name   string
	Target string
	Fields []*Field
}

func (f *Fragment) StringChan() (<-chan string, error) {
	ch := make(chan string)

	for _, field := range f.Fields {
		if field == nil {
			close(ch)
			return ch, errors.WithStack(NilFieldErr{})
		}
		if err := field.check(); err != nil {
			close(ch)
			return ch, errors.WithStack(err)
		}
	}
	return f.stringChan(), nil
}

// StringChan returns a read only channel which is guaranteed to be closed in the future.
func (f *Fragment) stringChan() <-chan string {
	tokenChan := make(chan string)
	go func() {
		// fragment {name} on {target}
		tokenChan <- "fragment"
		tokenChan <- tokenSpace
		tokenChan <- f.Name
		tokenChan <- tokenSpace
		tokenChan <- "on"
		tokenChan <- tokenSpace
		tokenChan <- f.Target

		// emit fields
		tokenChan <- tokenLB
		for i, field := range f.Fields {
			if i != 0 {
				tokenChan <- tokenComma
			}
			strs := field.stringChan()
			for str := range strs {
				tokenChan <- str
			}
		}
		tokenChan <- tokenRB
		close(tokenChan)
	}()
	return tokenChan
}

// TODO: make check
func (f *Fragment) check() error {
	return nil
}

////////////////
// Public API //
////////////////

func MakeFragment(name string, target string) *Fragment {
	return &Fragment{Name: name, Target: target}
}

// JSON returns a json string with "query" field.
func (f *Fragment) JSON() (string, error) {
	strCh, err := f.StringChan()
	if err != nil {
		return "", errors.WithStack(err)
	}
	s := StringFromChan(strCh)
	return fmt.Sprintf(`{%s}`, strings.Replace(s, `"`, `\"`, -1)), nil
}

// GetField return the field identified by the name. Nil if not exist.
func (f *Fragment) GetField(name string) *Field {
	for _, field := range f.Fields {
		if field.Name == name {
			return field
		}
	}
	return nil
}

// SetFields sets the Fields field of this Query.
// If q.Fields already contains data, they will be replaced.
func (f *Fragment) SetFields(fields ...*Field) *Fragment {
	f.Fields = fields
	return f
}

// AddFields adds to the Fields field of this Query.
func (f *Fragment) AddFields(fields ...*Field) *Fragment {
	f.Fields = append(f.Fields, fields...)
	return f
}
