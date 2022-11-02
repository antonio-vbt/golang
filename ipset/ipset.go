package ipset

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"strings"
)

const (
	IPv6           = "-6"
	SuppressErrors = "-!"
)

type Error struct {
	exec.ExitError
	cmd exec.Cmd
	msg string
}

func (e *Error) Error() string {
	return fmt.Sprintf("running %v: exit status %v: %v", e.cmd.Args, e.ExitCode(), e.msg)
}

type Ipset struct {
	path string // The path to the ipset executable
}

func New() (*Ipset, error) {
	path, err := exec.LookPath("ipset")
	if err != nil {
		return nil, fmt.Errorf("cannot find ipset executable: %w", err)
	}
	return &Ipset{
		path: path,
	}, nil
}

func (ips *Ipset) Create(setname, typename string, opts ...string) error {
	args := append([]string{"create", setname, typename}, opts...)
	return ips.run(args...)
}

func (ips *Ipset) Destroy(setname string) error {
	return ips.run("destroy", setname)
}

func (ips *Ipset) SetExists(setname string) (bool, error) {
	err := ips.run("list", setname)
	eerr, eok := err.(*Error)
	switch {
	case err == nil:
		return true, nil
	case eok && eerr.ExitCode() == 1:
		return false, nil
	default:
		return false, err
	}
}

func (ips *Ipset) Add(setname, addr string, opts ...string) error {
	args := append([]string{"add", setname, addr}, opts...)
	return ips.run(args...)
}

func (ips *Ipset) Del(setname, addr string, opts ...string) error {
	args := append([]string{"del", setname, addr}, opts...)
	return ips.run(args...)
}

func (ips *Ipset) Flush(setname string) error {
	return ips.run("flush", setname)
}

func (ips *Ipset) run(args ...string) error {
	return ips.runWithOutput(args, nil)
}

func (ips *Ipset) runWithOutput(args []string, stdout io.Writer) error {
	args = append([]string{ips.path}, args...)

	var stderr bytes.Buffer
	cmd := exec.Cmd{
		Path:   ips.path,
		Args:   args,
		Stdout: stdout,
		Stderr: &stderr,
	}

	if err := cmd.Run(); err != nil {
		switch e := err.(type) {
		case *exec.ExitError:
			return &Error{
				ExitError: *e,
				cmd:       cmd,
				msg:       strings.TrimRight(stderr.String(), "\r\n"),
			}
		default:
			return err
		}
	}
	return nil
}
