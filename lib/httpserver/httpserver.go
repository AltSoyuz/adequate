package httpserver

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/AltSoyuz/adequate/lib/buildinfo"
	"github.com/AltSoyuz/adequate/lib/logger"
)

// Serve starts HTTP servers on the given addresses with the provided handler.
// It listens for context cancellation to initiate a graceful shutdown.
// It returns an error if any server fails to start or if shutdown is problematic.
func Serve(ctx context.Context, addr string, handler http.Handler) error {
	// Listener avec TCP keep-alive configuré simplement.
	lc := net.ListenConfig{
		KeepAlive: 3 * time.Minute,
	}
	ln, err := lc.Listen(ctx, "tcp", addr)
	if err != nil {
		return err
	}
	return ServeWithListener(ctx, ln, handler)
}

// ServeWithListener démarre un serveur HTTP en utilisant un net.Listener fourni.
// Utile pour les tests : on peut créer un listener pour récupérer l'adresse et
// contrôler le cycle de vie du serveur depuis le test.
func ServeWithListener(ctx context.Context, ln net.Listener, handler http.Handler) error {
	logger.InfoSkipframes(2, "listening", "addr", ln.Addr().String())

	srv := &http.Server{
		Handler:           wrapHandlerWithBuiltins(handler),
		ReadHeaderTimeout: 5 * time.Second,
		IdleTimeout:       60 * time.Second,
		ErrorLog:          logger.StdErrorLogger(),
	}

	return serveWithShutdown(ctx, srv, ln)
}

// serveWithShutdown gère le cycle de vie d'un serveur HTTP avec shutdown gracieux.
func serveWithShutdown(ctx context.Context, srv *http.Server, ln net.Listener) error {
	errCh := make(chan error, 1)

	go func() {
		if err := srv.Serve(ln); err != nil && err != http.ErrServerClosed {
			errCh <- err
			return
		}
		close(errCh)
	}()

	select {
	case <-ctx.Done():
		// arrêt gracieux borné
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		shutdownErr := srv.Shutdown(shutdownCtx) // capture l'erreur

		// vide l'erreur éventuelle de Serve
		if err := <-errCh; err != nil {
			return err // vraie erreur serveur
		}
		// si le shutdown a dépassé le délai, signale-le
		if shutdownErr == context.DeadlineExceeded {
			return shutdownErr
		}
		// sinon, arrêt normal → pas d'erreur
		return nil

	case err := <-errCh:
		// échec non prévu de Serve
		return err
	}
}

type ErrResponse struct {
	Error string `json:"error"`
}

func WriteJSON(w http.ResponseWriter, r *http.Request, status int, v any) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(v); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	h := w.Header()
	h.Set("Content-Type", "application/json; charset=utf-8")
	h.Set("X-Content-Type-Options", "nosniff")

	rid := r.Header.Get("X-Request-Id")
	if rid != "" {
		h.Set("X-Request-Id", rid)
	}

	w.WriteHeader(status)
	_, _ = w.Write(buf.Bytes())
}

func WriteError(w http.ResponseWriter, r *http.Request, status int, err error) {
	msg := http.StatusText(status)
	if err != nil {
		msg = err.Error()
		args := []any{
			"status", status,
			"method", r.Method,
			"path", r.URL.Path,
			"err", err.Error(), // flatten error
		}
		if rid := r.Header.Get("X-Request-Id"); rid != "" {
			args = append(args, "rid", rid) // non-empty
		}
		logger.Error("http error", args...)
	}

	WriteJSON(w, r, status, ErrResponse{Error: msg})
}

func DecodeJSON[T any](r *http.Request) (T, error) {
	var v T
	defer func() {
		if err := r.Body.Close(); err != nil {
			logger.Error("error closing request body", "err", err)
		}
	}()
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&v); err != nil {
		return v, err
	}
	if err := dec.Decode(&struct{}{}); err != io.EOF {
		return v, io.ErrUnexpectedEOF
	}

	return v, nil
}

func SPAFileServer(staticDir string) http.HandlerFunc {
	fs := http.FileServer(http.Dir(staticDir))

	return func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {

		case "/app", "/app/":
			http.ServeFile(w, r, filepath.Join(staticDir, "200.html"))
			return
		default:
			if strings.HasPrefix(r.URL.Path, "/app/") {
				http.ServeFile(w, r, filepath.Join(staticDir, "200.html"))
				return
			}
		}

		switch r.URL.Path {
		case "/":
			http.ServeFile(w, r, filepath.Join(staticDir, "index.html"))
			return
		case "/about":
			http.ServeFile(w, r, filepath.Join(staticDir, "about.html"))
			return
		}

		fs.ServeHTTP(w, r)
	}
}

func wrapHandlerWithBuiltins(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/api/healthz":
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("OK"))
			return
		case r.Method == http.MethodGet && r.URL.Path == "/api/version":
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(buildinfo.Version))
			return
		case r.Method == http.MethodGet && r.URL.Path == "/api/metrics":
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			return
		}

		next.ServeHTTP(w, r)
	})
}
