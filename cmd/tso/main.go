package main

import (
	"log"
	"os"
	"os/signal"
	"runtime"
	"runtime/pprof"

	"gopkg.in/alecthomas/kingpin.v2"

	"hybrid/tso"
)

var (
	cpuProfile = kingpin.Flag("cpuProfile", "write cpu profile to file").Default("").String()
	address    = kingpin.Flag("addr", "listen address").Default(":7070").String()
	batchSize  = kingpin.Flag("batch", "batch size").Default("100000").Int32()
)

func main() {
	kingpin.Parse()

	runtime.GOMAXPROCS(runtime.NumCPU())
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	if *cpuProfile != "" {
		f, err := os.Create(*cpuProfile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		interrupt := make(chan os.Signal, 1)
		signal.Notify(interrupt)
		go catchKill(interrupt)
	}

	log.Println("Timestamp Oracle Started")
	orc := tso.NewOracle(*address, *batchSize)
	orc.Recover()
	orc.WaitForClientConnections()
}

func catchKill(interrupt chan os.Signal) {
	<-interrupt
	if *cpuProfile != "" {
		pprof.StopCPUProfile()
	}
	log.Fatalln("Caught Signal")
}
