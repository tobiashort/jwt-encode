package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

func printUsage() {
	fmt.Fprintf(os.Stderr, `Usage: jwt-encode [HEADER#PAYLOAD#SIGNATURE]
Reads from STDIN if HEADER#PAYLOAD#SIGNATURE is not defined as an argument.

Flags:
`)
	flag.PrintDefaults()
	os.Exit(1)
}

func encodeJson(str string) (string, error) {
	var buf bytes.Buffer
	err := json.Compact(&buf, []byte(str))
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(buf.Bytes()), nil
}

func main() {
	help := flag.Bool("h", false, "print help")
	flag.Parse()
	if *help {
		printUsage()
		return
	}
	if flag.NArg() > 1 {
		printUsage()
		return
	}
	input := ""
	if flag.NArg() == 1 {
		input = flag.Arg(0)
	} else {
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			panic(err)
		}
		input = strings.TrimSpace(string(data))
	}
	parts := strings.Split(input, "#")
	if len(parts) != 3 {
		fmt.Fprintln(os.Stderr, "Invalid input. Make sure HEADER, PAYLOAD and SIGNATURE are delimited by a '#'.")
		os.Exit(1)
		return
	}
	header, err := encodeJson(parts[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid header '%s'. Not JSON format.\n", parts[0])
		os.Exit(1)
		return
	}
	payload, err := encodeJson(parts[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid payload '%s'. Not JSON format.\n", parts[1])
		os.Exit(1)
		return
	}
	signature := parts[2]
	fmt.Printf("%s.%s.%s\n", header, payload, signature)
}
