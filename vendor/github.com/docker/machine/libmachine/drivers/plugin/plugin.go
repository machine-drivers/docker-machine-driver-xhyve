package plugin

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"time"

	"github.com/docker/machine/libmachine/log"
)

var (
	// Timeout where we will bail if we're not able to properly contact the
	// plugin server.
	defaultTimeout = 10 * time.Second
)

const (
	pluginOutPrefix = "PLUGIN OUT => "
	pluginErrPrefix = "PLUGIN ERR => "
)

type PluginStreamer interface {
	// Return a channel for receiving the output of the stream line by
	// line, and a channel for stopping the stream when we are finished
	// reading from it.
	//
	// It happens to be the case that we do this all inside of the main
	// plugin struct today, but that may not be the case forever.
	AttachStream(*bufio.Scanner) (<-chan string, chan<- bool)
}

type PluginServer interface {
	// Get the address where the plugin server is listening.
	Address() (string, error)

	// Serve kicks off the plugin server.
	Serve() error

	// Close shuts down the initialized server.
	Close() error
}

type McnBinaryExecutor interface {
	// Execute the driver plugin.  Returns scanners for plugin binary
	// stdout and stderr.
	Start() (*bufio.Scanner, *bufio.Scanner, error)

	// Stop reading from the plugins in question.
	Close() error
}

// DriverPlugin interface wraps the underlying mechanics of starting a driver
// plugin server and then figuring out where it can be dialed.
type DriverPlugin interface {
	PluginServer
	PluginStreamer
}

type LocalBinaryPlugin struct {
	Executor McnBinaryExecutor
	Addr     string
	addrCh   chan string
	stopCh   chan bool
}

type LocalBinaryExecutor struct {
	pluginStdout, pluginStderr io.ReadCloser
	DriverName                 string
}

func NewLocalBinaryPlugin(driverName string) *LocalBinaryPlugin {
	return &LocalBinaryPlugin{
		stopCh: make(chan bool, 1),
		addrCh: make(chan string, 1),
		Executor: &LocalBinaryExecutor{
			DriverName: driverName,
		},
	}
}

func (lbe *LocalBinaryExecutor) Start() (*bufio.Scanner, *bufio.Scanner, error) {
	log.Debugf("Launching plugin server for driver %s", lbe.DriverName)

	binaryPath, err := exec.LookPath(fmt.Sprintf("docker-machine-%s", lbe.DriverName))
	if err != nil {
		return nil, nil, fmt.Errorf("Error trying to locate plugin binary: %s", err)
	}

	log.Debugf("Found binary path at %s", binaryPath)

	cmd := exec.Command(binaryPath)

	lbe.pluginStdout, err = cmd.StdoutPipe()
	if err != nil {
		return nil, nil, fmt.Errorf("Error getting cmd stdout pipe: %s", err)
	}

	lbe.pluginStderr, err = cmd.StderrPipe()
	if err != nil {
		return nil, nil, fmt.Errorf("Error getting cmd stderr pipe: %s", err)
	}

	outScanner := bufio.NewScanner(lbe.pluginStdout)
	errScanner := bufio.NewScanner(lbe.pluginStderr)

	if err := cmd.Start(); err != nil {
		return nil, nil, fmt.Errorf("Error starting plugin binary: %s", err)
	}

	return outScanner, errScanner, nil
}

func (lbe *LocalBinaryExecutor) Close() error {
	if err := lbe.pluginStdout.Close(); err != nil {
		return err
	}

	if err := lbe.pluginStderr.Close(); err != nil {
		return err
	}

	return nil
}

func stream(scanner *bufio.Scanner, streamOutCh chan<- string, stopCh <-chan bool) {
	for scanner.Scan() {
		select {
		case <-stopCh:
			close(streamOutCh)
			return
		default:
			streamOutCh <- strings.Trim(scanner.Text(), "\n")
			if err := scanner.Err(); err != nil {
				log.Warnf("Scanning stream: %s", err)
			}
		}
	}
}

func (lbp *LocalBinaryPlugin) AttachStream(scanner *bufio.Scanner) (<-chan string, chan<- bool) {
	streamOutCh := make(chan string)
	stopCh := make(chan bool, 1)
	go stream(scanner, streamOutCh, stopCh)
	return streamOutCh, stopCh
}

func (lbp *LocalBinaryPlugin) execServer() error {
	outScanner, errScanner, err := lbp.Executor.Start()
	if err != nil {
		return fmt.Errorf("Plugin server did not start correctly: %s", err)
	}

	// Scan just one line to get the address, then send it to the relevant
	// channel.
	outScanner.Scan()
	addr := outScanner.Text()
	if err := outScanner.Err(); err != nil {
		return fmt.Errorf("Reading plugin address failed: %s", err)
	}

	lbp.addrCh <- strings.TrimSpace(addr)

	stdOutCh, stopStdoutCh := lbp.AttachStream(outScanner)
	stdErrCh, stopStderrCh := lbp.AttachStream(errScanner)

	for {
		select {
		case out := <-stdOutCh:
			log.Debug(pluginOutPrefix, out)
		case err := <-stdErrCh:
			log.Debug(pluginErrPrefix, err)
		case _ = <-lbp.stopCh:
			stopStdoutCh <- true
			stopStderrCh <- true
			if err := lbp.Executor.Close(); err != nil {
				return fmt.Errorf("Error closing local plugin binary: %s", err)
			}
			return nil
		}
	}

	return nil
}

func (lbp *LocalBinaryPlugin) Serve() error {
	return lbp.execServer()
}

func (lbp *LocalBinaryPlugin) Address() (string, error) {
	if lbp.Addr == "" {
		select {
		case lbp.Addr = <-lbp.addrCh:
			log.Debugf("Plugin server listening at address %s", lbp.Addr)
			close(lbp.addrCh)
			return lbp.Addr, nil
		case <-time.After(defaultTimeout):
			return "", fmt.Errorf("Failed to dial the plugin server in %s", defaultTimeout)
		}
	}
	return lbp.Addr, nil
}

func (lbp *LocalBinaryPlugin) Close() {
	lbp.stopCh <- true
	close(lbp.stopCh)
}
