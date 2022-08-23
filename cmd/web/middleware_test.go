package main

import (
	"fmt"
	"net/http"
	"testing"
)

func TestNoSurf(t *testing.T) {
	var myH myHandler

	h := NoSurf(&myH)

	switch v := h.(type) {
	case http.Handler:
		// do nothing, test passed
	default:
		t.Errorf(fmt.Sprintf("type is not http.Handler, it is %T", v))
	}
}

func TestSessionLoad(t *testing.T) {
	var myH myHandler

	h := SessionLoad(&myH)

	switch v := h.(type) {
	case http.Handler:
		// do nothing, test passed
	default:
		t.Errorf(fmt.Sprintf("type is not http.Handler, it is %T", v))
	}
}