package main

import (
	"flag"
	"log/slog"
	"os"
	"runtime"
	"runtime/pprof"
	"strings"
)

type report struct {
	ID      int
	Payload string
}

func main() {
	var (
		memProfile = flag.String("memprofile", "mem.prof", "write memory profile to file")
		count      = flag.Int("count", 40_000, "number of reports to create")
	)
	flag.Parse()

	reports := generateReports(*count)
	slog.Info("generated reports", "count", len(reports))

	runtime.GC()

	file, err := os.Create(*memProfile)
	if err != nil {
		slog.Error("create memory profile", "error", err)
		os.Exit(1)
	}
	defer file.Close()

	if err := pprof.WriteHeapProfile(file); err != nil {
		slog.Error("write memory profile", "error", err)
		os.Exit(1)
	}

	slog.Info("memory profile written", "profile", *memProfile, "example", reports[0].Payload[:24])
}

func generateReports(count int) []report {
	reports := make([]report, 0, count)
	for i := 0; i < count; i++ {
		reports = append(reports, report{
			ID:      i,
			Payload: strings.Repeat("payload-", 64),
		})
	}
	return reports
}
