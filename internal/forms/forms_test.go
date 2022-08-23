package forms

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestForm_Valid(t *testing.T) {
	// r := httptest.NewRequest("POST", "/some-url", nil)
	// form := New(r.PostForm)
	postData := url.Values{}
	form := New(postData)

	isValid := form.Valid()
	if !isValid {
		t.Error("got invalid when should have been valid")
	}
}

func TestForm_Required(t *testing.T) {
	r := httptest.NewRequest("POST", "/some-url", nil)
	form := New(r.PostForm)

	form.Required("a", "b", "c")
	if form.Valid() {
		t.Error("form shows valid when required fields missing")
	}

	postedData := url.Values{}
	postedData.Add("a", "a")
	postedData.Add("b", "a")
	postedData.Add("c", "a")

	r, _ = http.NewRequest("POST", "/some-url", nil)

	r.PostForm = postedData
	form = New(r.PostForm)
	form.Required("a", "b", "c")
	if !form.Valid() {
		t.Error("shows doesnt have required fields when it does")
	}
}

func TestForm_Has(t *testing.T) {
	// r := httptest.NewRequest("POST", "/some-url", nil)
	// form := New(r.PostForm)
	postData := url.Values{}
	form := New(postData)

	has := form.Has("some-field")
	if has {
		t.Error("form shows has field when it does not")
	}

	postData = url.Values{}
	postData.Add("a", "a")
	form = New(postData)

	has = form.Has("a")
	if !has {
		t.Error("shows form doesnt have field when it should")
	}
}

func TestForm_MinLength(t *testing.T) {
	// r := httptest.NewRequest("POST", "/some-url", nil)
	// form := New(r.PostForm)
	postData := url.Values{}
	form := New(postData)

	form.MinLength("x", 10)
	if form.Valid() {
		t.Error("form shows min length for non-existent field")
	}

	isError := form.Errors.Get("x")
	if isError == "" {
		t.Error("should have an error, but did not get one")
	}

	postData = url.Values{}
	postData.Add("field", "value")
	form = New(postData)

	form.MinLength("field", 10)
	if form.Valid() {
		t.Error("shows min length of 10 met when data is shorter")
	}

	postData = url.Values{}
	postData.Add("field", "value")
	form = New(postData)

	form.MinLength("field", 3)
	if !form.Valid() {
		t.Error("shows min length of 3 is not met when it is")
	}

	isError = form.Errors.Get("field")
	if isError != "" {
		t.Error("should not have an error, but it did")
	}

}

func TestForm_IsEmail(t *testing.T) {
	postData := url.Values{}
	form := New(postData)

	form.IsEmail("x")
	if form.Valid() {
		t.Error("form shows valid email for non-existent field")
	}

	postData = url.Values{}
	postData.Add("email", "jdoe.com")
	form = New(postData)

	form.IsEmail("email")
	if form.Valid() {
		t.Error("got valid for an invalid email")
	}

	postData = url.Values{}
	postData.Add("email", "jdoe@gmail.com")
	form = New(postData)

	form.IsEmail("email")
	if !form.Valid() {
		t.Error("got an invalid email when we should not have")
	}

}