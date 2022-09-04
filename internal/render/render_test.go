package render

import (
	"net/http"
	"testing"

	"github.com/magier25/go-bookings/internal/models"
)

func TestAddDefaultData(t *testing.T) {
	var td models.TemplateData

	r, err := getSession()
	if err != nil {
		t.Error(err)
	}

	flashText := "abc"
	session.Put(r.Context(), "flash", flashText)
	result := AddDefaultData(&td, r)
	if result.Flash != flashText {
		t.Errorf("flash value of '%s' not found in session", flashText)
	}
}

func getSession() (*http.Request, error) {
	r, err := http.NewRequest("GET", "/some-url", nil)
	if err != nil {
		return nil, err
	}

	ctx := r.Context()
	ctx, _ = session.Load(ctx, r.Header.Get("X-Session"))
	r = r.WithContext(ctx)

	return r, nil
}
