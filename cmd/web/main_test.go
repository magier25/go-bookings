package main

import (
	"testing"

	"github.com/magier25/go-bookings/internal/driver"
)

func TestRun(t *testing.T) {
	var db driver.DB

	err := run(false, false, &db)
	if err != nil {
		t.Error("failed run()")
	}
}
