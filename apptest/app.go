package apptest

import (
	"bufio"
	"context"
	"flag"
	"io"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"testing"
	"time"
)

type App struct {
	tc             *TestCase
	flags          []string
	cmd            *exec.Cmd
	httpListenAddr string
	BaseURL        string
	Cli            *Client
	dbPath         string
}

var listenRe = regexp.MustCompile(`addr="?([0-9.]+:\d+)"?`)
var binPath = flag.String("bin.path", "./../../bin/app", "path to the app binary")

func StartApp(tc *TestCase, flags ...string) *App {
	tc.T().Helper()

	dbPath := tc.T().Name() + ".db"
	flags = append(flags, "-store.sqlitePath="+dbPath)
	flags = setDefaultFlags(flags)

	cmd := exec.Command(*binPath, flags...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		tc.T().Fatalf("stdout pipe: %v", err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		tc.T().Fatalf("stderr pipe: %v", err)
	}

	if err := cmd.Start(); err != nil {
		tc.T().Fatalf("start app: %v", err)
	}

	addrCh := make(chan string, 1)
	go scanLogs(stdout, os.Stdout, addrCh)
	go scanLogs(stderr, os.Stderr, addrCh)

	addr := waitAddr(tc.T(), addrCh)

	app := &App{
		tc:             tc,
		flags:          flags,
		cmd:            cmd,
		httpListenAddr: addr,
		BaseURL:        "http://" + addr,
		Cli:            NewClient(),
		dbPath:         dbPath,
	}

	app.waitForReady("/api/healthz")

	tc.RegisterCleanup(func() {
		_ = cmd.Process.Kill()
		_ = cmd.Wait()
		app.Cli.CloseConnections()
		_ = os.Remove(dbPath)
		_ = os.Remove(dbPath + "-shm")
		_ = os.Remove(dbPath + "-wal")
	})

	return app
}

func scanLogs(r io.Reader, w io.Writer, addrCh chan<- string) {
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		line := sc.Text()
		_, _ = io.WriteString(w, line+"\n")

		if m := listenRe.FindStringSubmatch(line); len(m) == 2 {
			select {
			case addrCh <- m[1]:
			default:
			}
		}
	}
}

func waitAddr(t *testing.T, ch <-chan string) string {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	select {
	case addr := <-ch:
		return addr
	case <-ctx.Done():
		t.Fatalf("app didn't log listening addr=...")
		return ""
	}
}

func (a *App) waitForReady(path string) {
	t := a.tc.T()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	url := a.BaseURL + path

	for {
		_, statusCode := a.Cli.Get(t, url)
		if statusCode == http.StatusOK {
			return
		}
		if ctx.Err() != nil {
			t.Fatalf("app not ready on %s", url)
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func setDefaultFlags(flags []string) []string {
	defaults := []struct {
		key   string
		value string
	}{
		{"-http.listenAddr=", "127.0.0.1:0"},
	}

	for _, def := range defaults {
		found := false
		for _, f := range flags {
			if strings.HasPrefix(f, def.key) {
				found = true
				break
			}
		}
		if !found {
			flags = append(flags, def.key+def.value)
		}
	}

	return flags
}
