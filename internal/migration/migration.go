package migration

import (
	"net/http"

	"github.com/AltSoyuz/adequate/internal/store"
	"github.com/AltSoyuz/adequate/lib/httpserver"
	"github.com/AltSoyuz/adequate/lib/logger"
)

func MigrationHandler(s *store.Store) http.HandlerFunc {
	type MigrationHandlerResp struct {
		Version int64 `json:"version"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		version, err := s.Queries.GetLastMigrationVersion(ctx)
		if err != nil {
			httpserver.WriteError(w, r, http.StatusInternalServerError, err)
			return
		}

		logger.Info("migration version fetched", "version", version)

		httpserver.WriteJSON(w, r, http.StatusOK, MigrationHandlerResp{
			Version: version,
		})
	}
}
