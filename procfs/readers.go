package procfs

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func readNamedFields(path string, wanted []string) (map[string]uint64, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	want := make(map[string]bool, len(wanted))
	for _, k := range wanted {
		want[k] = true
	}

	found := make(map[string]uint64, len(wanted))
	scanner := bufio.NewScanner(file)
	for scanner.Scan() && len(found) < len(wanted) {
		fields := strings.Fields(scanner.Text())
		if len(fields) < 2 {
			continue
		}
		fieldName, _ := strings.CutSuffix(fields[0], ":")
		if !want[fieldName] {
			continue
		}
		v, err := strconv.ParseUint(fields[1], 10, 64)
		if err != nil {
			return nil, err
		}
		found[fieldName] = v
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	for _, k := range wanted {
		if _, ok := found[k]; !ok {
			return nil, fmt.Errorf("%s not found in %s", k, path)
		}
	}
	return found, nil
}
