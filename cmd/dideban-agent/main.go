package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/dideban-project/dideban/internal/collector"
)

func main() {
	mem, err := collector.ReadMemory()
	if err != nil {
		fmt.Fprintln(os.Stderr, "dideban-agent:", err)
		os.Exit(1)
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(mem); err != nil {
		fmt.Fprintln(os.Stderr, "dideban-agent:", err)
		os.Exit(1)
	}
}
