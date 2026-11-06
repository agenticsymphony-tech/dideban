package collector

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// Memory holds memory statistics in bytes.
type Memory struct {
	TotalBytes     uint64 `json:"total_bytes"`
	AvailableBytes uint64 `json:"available_bytes"`
	UsedBytes      uint64 `json:"used_bytes"`
}

// ParseMeminfo parses the contents of /proc/meminfo.
// It takes an io.Reader (not a file path) so it can be tested with plain strings.
func ParseMeminfo(r io.Reader) (Memory, error) {
	var totalKiB, availableKiB uint64
	var haveTotal, haveAvailable bool

	sc := bufio.NewScanner(r)
	for sc.Scan() {
		key, rest, ok := strings.Cut(sc.Text(), ":")
		if !ok || (key != "MemTotal" && key != "MemAvailable") {
			continue
		}

		fields := strings.Fields(rest) // e.g. ["16061232", "kB"]
		if len(fields) == 0 {
			return Memory{}, fmt.Errorf("meminfo: no value for %s", key)
		}
		v, err := strconv.ParseUint(fields[0], 10, 64)
		if err != nil {
			return Memory{}, fmt.Errorf("meminfo: parsing %s: %w", key, err)
		}

		switch key {
		case "MemTotal":
			totalKiB, haveTotal = v, true
		case "MemAvailable":
			availableKiB, haveAvailable = v, true
		}
	}
	if err := sc.Err(); err != nil {
		return Memory{}, fmt.Errorf("meminfo: reading: %w", err)
	}
	if !haveTotal || !haveAvailable {
		return Memory{}, fmt.Errorf("meminfo: MemTotal or MemAvailable missing")
	}
	if availableKiB > totalKiB {
		return Memory{}, fmt.Errorf("meminfo: MemAvailable (%d) > MemTotal (%d)", availableKiB, totalKiB)
	}

	return Memory{
		TotalBytes:     totalKiB * 1024,
		AvailableBytes: availableKiB * 1024,
		UsedBytes:      (totalKiB - availableKiB) * 1024,
	}, nil
}

// ReadMemory reads the current memory statistics from /proc/meminfo.
func ReadMemory() (Memory, error) {
	f, err := os.Open("/proc/meminfo")
	if err != nil {
		return Memory{}, err
	}
	defer f.Close()
	return ParseMeminfo(f)
}
