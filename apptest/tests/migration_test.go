package tests

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/AltSoyuz/adequate/apptest"
)

func TestMigrationVersion(t *testing.T) {
	tc := apptest.NewTestCase(t)
	defer tc.Stop()

	app := apptest.StartApp(tc)

	res, statusCode := app.Cli.Get(t, app.BaseURL+"/api/migrations/version")
	var resp struct {
		Version int64 `json:"version"`
	}

	if statusCode != http.StatusOK {
		t.Fatalf("unexpected status code: got %d, want %d, resp body: %s", statusCode, http.StatusOK, res)
	}
	err := json.Unmarshal([]byte(res), &resp)
	if err != nil {
		t.Fatalf("could not unmarshal response: %v", err)
	}

	if resp.Version != 1 {
		t.Fatalf("unexpected migration version: got %d, want %d", resp.Version, 1)
	}
}
