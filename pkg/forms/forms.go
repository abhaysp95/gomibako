package forms

import (
	"fmt"
	"net/url"
	"strings"
	"unicode/utf8"
)

// holds form data and errMap for any validation message related to form data
type Form struct {
	url.Values
	ErrMap errMap
}

// return new Form struct
func New(data url.Values) *Form {
	return &Form{
		data,
		errMap(map[string][]string{}),
	}
}

// validation for "required". If no value for "field" is found, appropriate error is added to errMap
func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Values.Get(field)
		if strings.TrimSpace(value) == "" {
			f.ErrMap.Add(field, fmt.Sprintf("%s can't be empty", field))
		}
	}
}

// validation for "maximum length". In case of failure, appropriate error is added to errMap
func (f *Form) MaxLength(field string, d int) {
	value := f.Get(field)
	if utf8.RuneCountInString(value) > d {
		f.ErrMap.Add(field, fmt.Sprintf("%s's length can't be more than %d", field, d))
	}
}

// validation for "permitted values". In case of failure, appropriate error is added to errMap
func (f *Form) PermittedValues(field string, opts ...string) {
	value := f.Get(field)
	if value == "" {
		return
	}
	status := false
	for _, opt := range opts {
		if value == opt {
			status = true
		}
	}
	if !status {
		f.ErrMap.Add(field, fmt.Sprintf("%s can have values from %v", field, opts))
	}
}

// returns true if there is no validation error found in Form's data
func (f *Form) Valid() bool {
	return len(f.ErrMap) == 0;
}
