package forms

import (
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/asaskevich/govalidator"
)

func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}

type Form struct {
	url.Values
	Errors errors
}

//New initialises a form struct
func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "This field cannot be blank")
		}
	}
}

func (f *Form) Has(field string) bool {
	x := f.Get(field)
	if x == "" {
		return false
	}
	return true
}

func (f *Form) MinLength(field string, length int) bool {
	x := f.Get(field)
	if len(x) < length {
		f.Errors.Add(field, fmt.Sprintf("This field must be at least %d characters long", length))
		return false
	}
	return true
}

func (f *Form) IsEmail(field string) {
	if !govalidator.IsEmail(f.Get(field)) {
		f.Errors.Add(field, "Invalid email address")
	}
}

func (f *Form) IsPhoneNumber(field string) {
	if !govalidator.IsNumeric(f.Get(field)) {
		f.Errors.Add(field, "Invalid phone number")
	}
}

func (f *Form) IsSame(field, field2 string) {
	p := f.Get(field)
	p2 := f.Get(field2)

	fmt.Println(p)
	fmt.Println(p2)

	if !(p == p2) {
		log.Println("password not same")
		f.Errors.Add(field2, "passwords doesnt match")
	}
}
