package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

	"gopkg.in/yaml.v3"
)

func main() {
	flag.Parse()

	var err error
	var defaultsMap map[string]string
	var defaultsBytes []byte

	if len(flag.Args()) < 3 {
		printUsage()
		log.Fatal("expecting an opts file and a cmd.")
	}

	cmd, defaultsFile, args := splitArgs(flag.Args())

	if defaultsBytes, err = ioutil.ReadFile(defaultsFile); err != nil {
		log.Fatalf("failed to read %s: %q", defaultsFile, err)
	}

	switch {
	case strings.HasSuffix(defaultsFile, "yaml"):
		if err = yaml.Unmarshal(defaultsBytes, &defaultsMap); err != nil {
			log.Fatalf("failed to unmarshal %s: %q", defaultsFile, err)
		}
	case strings.HasSuffix(defaultsFile, "json"):
		if err = json.Unmarshal(defaultsBytes, &defaultsMap); err != nil {
			log.Fatalf("failed to unmarshal %s: %q", defaultsFile, err)
		}
	default:
		log.Fatalf("unexpected input file %s, file name must end with yaml/json", defaultsFile)
	}

	if args, err = makeArgs(defaultsMap, args); err != nil {
		log.Fatal(err)
	}

	var command = exec.Command(cmd, args...)
	command.Stderr = os.Stderr
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout

	if err = command.Run(); err != nil {
		log.Fatal(err)
	}
}

func printUsage() {
	fmt.Println(`usage:
	optifiy <yaml/json> -- cmd [--opt1=o1] [--opt2 o2]
decompose the key-values in yaml/json as long opts to invoke cmd.`)
}

func splitArgs(args []string) (cmd string, defaultsFile string, opts []string) {
	i := -1
	for j := range args {
		if args[j] == "--" {
			i = j
			break
		}
	}

	if i == -1 {
		printUsage()
		panic(fmt.Errorf("no defaults provided"))
	}

	if i == len(args)-1 {
		printUsage()
		panic(fmt.Errorf("no cmd provided"))
	}

	if i == len(args)-2 {
		return args[i+1], args[i-1], nil
	}

	return args[i+1], args[i-1], args[i+2:]
}

func getLongOpts(args []string) []string {
	var longOpts []string
	for i := 0; i < len(args); i++ {
		if strings.HasPrefix(args[i], "--") {
			longOpts = append(longOpts, args[i])
			if !strings.Contains(args[i], "=") {
				longOpts = append(longOpts, args[i+1])
				i++
			}
		}
	}

	return longOpts
}

func selectDefaultOpts(defaults map[string]string, opts []string) ([]string, error) {
	var err error
	var overrides map[string]string
	if overrides, err = optsToMap(opts); err != nil {
		return nil, fmt.Errorf("failed to parse options %v: %q", opts, err)
	}

	var out []string
	for k, v := range defaults {
		if _, exist := overrides[k]; !exist {
			out = append(out, []string{"--" + k, v}...)
		}
	}

	return out, nil
}

func makeArgs(defaultsMap map[string]string, args []string) ([]string, error) {
	var err error
	var defaultOpts []string
	var finalArgs []string
	var longOpts = getLongOpts(args)
	if defaultOpts, err = selectDefaultOpts(defaultsMap, longOpts); err != nil {
		return nil, err
	}

	finalArgs = append(finalArgs, args...)
	finalArgs = append(finalArgs, defaultOpts...)

	return finalArgs, nil
}

func optsToMap(opts []string) (map[string]string, error) {
	res := make(map[string]string)

	i := 0
	for i < len(opts) {
		if !strings.HasPrefix(opts[i], "--") {
			return nil, fmt.Errorf("expecting long form opts like --opt, but was: %q", opts[i])
		}

		o := strings.SplitN(opts[i], "=", 2)
		switch len(o) {
		case 1:
			res[o[0][2:]] = opts[i+1]
			i += 2
		case 2:
			res[o[0][2:]] = o[1]
			i += 1
		default:
			panic(fmt.Errorf("unexpected opt split: %d", len(o)))
		}
	}

	return res, nil
}
