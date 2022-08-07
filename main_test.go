package main

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestArgsSplit(t *testing.T) {
	g := NewWithT(t)

	cmd, defaults, args := splitArgs([]string{"cmd", "subcmd1", "subcmd2", "--opt1", "val1", "-f", "--", "test.yaml"})
	g.Expect(cmd).To(Equal("cmd"))
	g.Expect(defaults).To(Equal("test.yaml"))
	g.Expect(args).To(Equal([]string{"subcmd1", "subcmd2", "--opt1", "val1", "-f"}))

	cmd, defaults, args = splitArgs([]string{"cmd", "subcmd1", "subcmd2", "--opt1", "val1", "--"})
	g.Expect(cmd).To(Equal("cmd"))
	g.Expect(defaults).To(Equal(""))
	g.Expect(args).To(Equal([]string{"subcmd1", "subcmd2", "--opt1", "val1", "--"}))
}

func TestOptsToMap(t *testing.T) {
	g := NewWithT(t)
	opts := []string{"subcmd1", "subcmd2", "--opt1", "val1", "-f", "--opt2=2"}
	optsMap, err := optsToMap(getLongOpts(opts))
	g.Expect(err).To(BeNil())
	g.Expect(optsMap).To(Equal(map[string]string{
		"opt1": "val1",
		"opt2": "2",
	}))
}

func TestMakeOpts(t *testing.T) {
	g := NewWithT(t)

	args := []string{"subcmd1", "subcmd2", "--opt1", "val1", "-f", "--opt2=2"}
	defaultsMap := map[string]string{
		"opt":  "0",
		"opt1": "1",
	}
	defaults, err := selectDefaultOpts(defaultsMap, getLongOpts(args))

	g.Expect(err).To(BeNil())
	g.Expect(defaults).To(Equal([]string{"--opt", "0"}))
	args, err = makeArgs(defaultsMap, args)
	g.Expect(err).To(BeNil())
	g.Expect(args).To(Equal([]string{"subcmd1", "subcmd2", "--opt1", "val1", "-f", "--opt2=2", "--opt", "0"}))
}
