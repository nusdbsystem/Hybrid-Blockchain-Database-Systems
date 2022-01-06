package benchmark

import (
	"bufio"
	"io"
	"math/rand"
	"strings"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

type RequestType int

const (
	SetOp RequestType = iota
	GetOp
)

type Request struct {
	ReqType RequestType
	Key     string
	Val     string
}

func LineByLine(r io.Reader, fn func(line string) error) error {
	br := bufio.NewReader(r)
	// num := 0
	for {
		line, _, err := br.ReadLine()

		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		lineStr := string(line)
		if !(strings.HasPrefix(lineStr, "INSERT") ||
			strings.HasPrefix(lineStr, "READ") ||
			strings.HasPrefix(lineStr, "UPDATE")) {
			continue
		}

		if err := fn(lineStr); err != nil {
			return err
		}
		/*
		num = num + 1
		if num >= 10000 {
			break
		}
		*/
	}
	return nil
}

func GenRandString(n int) string {
	b := make([]byte, n)
	// A rand.Int63() generates 63 random bits, enough for letterIdxMax letters!
	for i, cache, remain := n-1, rand.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = rand.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return string(b)
}
