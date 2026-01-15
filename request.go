package graphb

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

type Element interface {
	stringChan() <-chan string
	check() error
}

type Request struct {
	Elements []Element
}

func (r *Request) StringChan() (<-chan string, error) {
	ch := make(chan string)

	// TODO: make check
	//if err := r.check(); err != nil {
	//	close(ch)
	//	return ch, errors.WithStack(err)
	//}

	for _, f := range r.Elements {
		if f == nil {
			close(ch)
			return ch, errors.WithStack(NilFieldErr{})
		}
		if err := f.check(); err != nil {
			close(ch)
			return ch, errors.WithStack(err)
		}
	}
	return r.stringChan(), nil
}

func (r *Request) stringChan() <-chan string {
	tokenChan := make(chan string)
	go func() {
		for i, el := range r.Elements {
			if i != 0 {
				tokenChan <- tokenSpace
			}
			strs := el.stringChan()
			for str := range strs {
				tokenChan <- str
			}
		}
		close(tokenChan)
	}()
	return tokenChan
}

// TODO: make check
func (r *Request) check() error {
	return nil
}

////////////////
// Public API //
////////////////

// MakeRequest constructs a Request and returns a pointer of it.
func MakeRequest() *Request { return &Request{} }

func (r *Request) JSON() (string, error) {
	strCh, err := r.StringChan()
	if err != nil {
		return "", errors.WithStack(err)
	}
	s := StringFromChan(strCh)
	return fmt.Sprintf(`%s`, strings.Replace(s, `"`, `\"`, -1)), nil
}

func (r *Request) SetElements(elements ...Element) *Request {
	r.Elements = elements
	return r
}
