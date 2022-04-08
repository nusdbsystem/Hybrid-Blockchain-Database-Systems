package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	loadReadPath  = kingpin.Flag("load-read-path", "Path of YCSB initial data").Required().String()
	loadStorePath = kingpin.Flag("load-write-path", "Path of YCSB operation data").Required().String()
	runReadPath   = kingpin.Flag("run-read-path", "Path of YCSB initial data").Required().String()
	runStorePath  = kingpin.Flag("run-write-path", "Path of YCSB operation data").Required().String()
)

func EncodeVal(val string) string {
	runes := []rune(val)
	for i := 0; i < len(runes); i++ {
		if (runes[i] >= 'a' && runes[i] <= 'z') ||
			(runes[i] >= 'A' && runes[i] <= 'Z') ||
			(runes[i] >= '0' && runes[i] <= '9') {
			continue
		}
		runes[i] = '0'
	}
	return string(runes)
}

func LineByLine(m map[string]int64, r io.Reader, w io.Writer) error {
	br := bufio.NewReader(r)
	for {
		line, err := br.ReadString('\n')

		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Printf("%v\n", err)
			return err
		}

		lineStr := string(line)
		if !(strings.HasPrefix(lineStr, "INSERT") ||
			strings.HasPrefix(lineStr, "READ") ||
			strings.HasPrefix(lineStr, "UPDATE")) {
			fmt.Fprintln(w, lineStr)
			continue
		}

		tokens := strings.SplitN(lineStr, " ", 4)
		key := tokens[2]
		val := EncodeVal(tokens[3])
		ver, exists := m[key]

		if strings.HasPrefix(lineStr, "INSERT") || strings.HasPrefix(lineStr, "UPDATE") {
			if exists {
				ver = ver + 1
			} else {
				ver = 0
			}
			m[key] = ver
		}
		if strings.HasPrefix(lineStr, "READ") && !exists {
			ver = -1
		}

		fmt.Fprintf(w, "%s %s %s %s %d\n", tokens[0], tokens[1], key, val, ver)
	}
	return nil
}

func main() {
	kingpin.Parse()

	loadDataReadFile, err := os.Open(*loadReadPath)
	if err != nil {
		panic(err)
	}
	defer loadDataReadFile.Close()

	loadDataWriteFile, err := os.Create(*loadStorePath)
	if err != nil {
		panic(err)
	}
	defer loadDataWriteFile.Close()

	runDataReadFile, err := os.Open(*runReadPath)
	if err != nil {
		panic(err)
	}
	defer runDataReadFile.Close()

	runDataWriteFile, err := os.Create(*runStorePath)
	if err != nil {
		panic(err)
	}
	defer runDataWriteFile.Close()

	m := make(map[string]int64)
	LineByLine(m, loadDataReadFile, loadDataWriteFile)
	LineByLine(m, runDataReadFile, runDataWriteFile)
}
