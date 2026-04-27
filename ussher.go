package main

import (
	"errors"
	"log"
	"os"
)

func initLog() {
	// Attempt to write a log file to a standard location
	file, err1 := os.OpenFile("/var/log/ussher/ussher.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err1 != nil {
		// If we can't write to a standard location, then try to write to the current working directory
		file, err2 := os.OpenFile("ussher.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err2 != nil {
			// If we can't produce real log files, then use stdout and abort
			log.Fatal("Refusing to run without being able to log to /var/log/ussher/ (", err1, ") or current working directory (", err2, ")")
		}

		// At least we can log to the current working directory
		log.SetOutput(file)
		log.Print("Failed to write to /var/log/ussher: ", err1)
	}

	log.SetOutput(file)
}

// validateArgs requires args to look like `ussher <single-arg>` (a username
// or `--version`) and returns the single arg or a usage error. The wrong-arg-
// count check has to live before any os.Args[1] access; otherwise invoking
// `ussher` with no args panics with "index out of range".
func validateArgs(args []string) (string, error) {
	if len(args) != 2 {
		return "", errors.New("usage: ussher <username>")
	}
	return args[1], nil
}

// runDeps wires the side-effecting bits main() depends on (security gates,
// log setup, the actual fetch loop) so tests can substitute fakes. Only the
// fields that have side effects or read host state need to be injected;
// pure helpers like validateArgs are called directly.
type runDeps struct {
	executableWritable func() bool
	runningAsRoot      func() bool
	initLog            func()
	validUser          func(name string) bool
	loadConfig         func(username string, c *Config)
	runFn              func(*Config)
}

// defaultDeps wires the production implementations. main uses this; tests
// supply their own.
func defaultDeps() runDeps {
	return runDeps{
		executableWritable: isExecutableWritable,
		runningAsRoot:      isRunningAsRoot,
		initLog:            initLog,
		validUser:          isValidUser,
		loadConfig:         func(u string, c *Config) { c.LoadConfigByUser(u) },
		runFn:              Run,
	}
}

// runMain is main's testable body: same gate ordering, but every gate's
// outcome is observable as an error return rather than a process-killing
// log.Fatal. main() is the thin wrapper that calls log.Fatal on error.
func runMain(args []string, d runDeps) error {
	arg, err := validateArgs(args)
	if err != nil {
		return err
	}

	// Support `ussher --version` to print versioning info about the executable
	// before proceeding to security-hardening checks.
	if arg == "--version" {
		PrintVersion()
		return nil
	}

	// Security sanity checks
	if d.executableWritable() {
		return errors.New("Refusing to run due to permissions issue on the ussher executable")
	}
	if d.runningAsRoot() {
		return errors.New("Refusing to run as root")
	}

	// Initialize logging AFTER security checks to ensure we're writing logs as
	// a non-root user.
	d.initLog()

	// Check if the input username is valid.
	username := arg
	if !d.validUser(username) {
		return errors.New("User not found")
	}

	// At this point, we know that the input username is valid and safe to use.
	log.Print("Sourcing authorized_keys for ", username)

	var c Config
	d.loadConfig(username, &c)
	d.runFn(&c)
	return nil
}

func main() {
	if err := runMain(os.Args, defaultDeps()); err != nil {
		log.Fatal(err)
	}
}
