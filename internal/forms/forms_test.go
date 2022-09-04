package forms

import (
	"net/url"
	"testing"
)

func TestForm_Valid(t *testing.T) {
	postedData := url.Values{}
	form := New(postedData)

	isValid := form.Valid()
	if !isValid {
		t.Error("got invalid when should have been valid")
	}
}

func TestForm_GetError(t *testing.T) {
	postedData := url.Values{}
	form := New(postedData)

	isError := form.Errors.Get("a")
	if isError != "" {
		t.Error("should not have an error, but got one")
	}

	form.Required("a")
	isError = form.Errors.Get("a")
	if isError == "" {
		t.Error("should have an error but did not get one")
	}

}

func TestForm_Required(t *testing.T) {
	postedData := url.Values{}
	form := New(postedData)

	form.Required("a", "b", "c")
	if form.Valid() {
		t.Error("form shows valid when required fields missing")
	}

	postedData = url.Values{}
	postedData.Add("a", "1")
	postedData.Add("b", "2")
	postedData.Add("c", "3")

	form = New(postedData)
	form.Required("a", "b", "c")
	if !form.Valid() {
		t.Error("shows does not have required fields when it does")
	}
}

func TestForm_IsEmail(t *testing.T) {
	postedData := url.Values{}
	postedData.Add("email", "abc")

	form := New(postedData)
	form.IsEmail("email")
	if form.Valid() {
		t.Error("shows invalid email as valid")
	}

	postedData = url.Values{}
	postedData.Add("email", "john@example.com")

	form = New(postedData)
	form.IsEmail("email")
	if !form.Valid() {
		t.Error("shows valid email as invalid")
	}
}

func TestForm_Has(t *testing.T) {
	postedData := url.Values{}
	form := New(postedData)

	has := form.Has("whatever")
	if has {
		t.Error("form shows has filed when it does not")
	}

	postedData.Add("a", "1")
	form = New(postedData)
	has = form.Has("a")
	if !has {
		t.Error("shows form does not have field when it should")
	}
}

func TestForm_MinLength(t *testing.T) {
	postedData := url.Values{}
	form := New(postedData)
	form.MinLength("a", 1)
	if form.Valid() {
		t.Error("shows min length for non-existent field")
	}

	postedData.Add("a", "a")
	form = New(postedData)
	form.MinLength("a", 100)
	if form.Valid() {
		t.Error("shows min length of 100 met data is shorter")
	}

	form = New(postedData)
	form.MinLength("a", 1)
	if !form.Valid() {
		t.Error("shows min length of 1 is not when it is")
	}
}
