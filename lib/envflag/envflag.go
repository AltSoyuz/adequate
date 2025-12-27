package envflag

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

var (
	prefix = flag.String("envflag.prefix", "", "Prefix for environment variables")
)

// Parse parses environment vars and command-line flags.
//
// Flags set via command-line override flags set via environment vars.
//
// This function must be called instead of flag.Parse() before using any flags in the program.
func Parse() {
	ParseFlagSet(flag.CommandLine, os.Args[1:])
}

// ParseFlagSet parses the given args into the given fs.
func ParseFlagSet(fs *flag.FlagSet, args []string) {
	// Keep existing behavior: fatal on error for backward compatibility.
	if err := ParseFlagSetErr(fs, args); err != nil {
		// Do not use lib/logger here, since it is uninitialized yet.
		log.Fatalf("%s", err)
	}
}

// ParseFlagSetErr behaves like ParseFlagSet but returns an error instead of calling log.Fatalf.
// Use this when the caller prefers to handle parse errors instead of exiting the process.
func ParseFlagSetErr(fs *flag.FlagSet, args []string) error {
	if err := fs.Parse(args); err != nil {
		return fmt.Errorf("cannot parse flags %q: %w", args, err)
	}
	if fs.NArg() > 0 {
		return fmt.Errorf("unprocessed command-line args left: %s; the most likely reason is missing `=` between boolean flag name and value; see https://pkg.go.dev/flag#hdr-Command_line_flag_syntax", fs.Args())
	}
	// Remember explicitly set command-line flags.
	flagsSet := make(map[string]bool)
	fs.Visit(func(f *flag.Flag) {
		flagsSet[f.Name] = true
	})

	// Obtain the remaining flag values from environment vars.
	fs.VisitAll(func(f *flag.Flag) {
		if flagsSet[f.Name] {
			// The flag is explicitly set via command-line.
			return
		}
		// Get flag value from environment var.
		fname := getEnvFlagName(f.Name)
		if v, ok := os.LookupEnv(fname); ok {
			if err := fs.Set(f.Name, v); err != nil {
				// Return an error instead of exiting; the caller can decide what to do.
				// Do not use lib/logger here, since it is uninitialized yet.
				// Preserve original context in the error message.
				// Example: cannot set flag number to "not-a-number", which is read from env var "number": parse error
				// We wrap the underlying error for callers to inspect if needed.
				panicErr := fmt.Errorf("cannot set flag %s to %q, which is read from env var %q: %w", f.Name, v, fname, err)
				// Since VisitAll runs in a closure, we can't return directly here; use panic to bubble the error.
				panic(panicErr)
			}
		}
	})

	// If a panic was used to bubble an fs.Set error, recover and return it as an error.
	// This pattern avoids duplicating VisitAll logic while still returning an error.
	// (The panic is limited to the closure above and immediately recovered here.)
	if r := recover(); r != nil {
		if err, ok := r.(error); ok {
			return err
		}
		return fmt.Errorf("unexpected error: %v", r)
	}

	return nil
}

func getEnvFlagName(s string) string {
	// Substitute dots with underscores, since env var names cannot contain dots.
	s = strings.ReplaceAll(s, ".", "_")
	return *prefix + s
}
