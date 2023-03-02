package tracklists

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"os"
	"strings"

	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

func Read(filepath string) ([][]string, error) {
	f, err := os.Open(filepath)
	defer f.Close()
	if err != nil {
		return nil, err
	}

	var records [][]string

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := utf16(scanner.Bytes())
		parts := strings.Split(line, "\t")

		if parts[0] == "#" || len(parts) <= 1 {
			continue
		}

		record := []string{
			strings.TrimSpace(parts[2]),
			strings.TrimSpace(parts[3]),
		}

		records = append(records, record)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return records, nil
}

func utf16(data []byte) string {
	buf := bytes.NewBuffer(data)
	transformer := unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM)
	decoder := transformer.NewDecoder()
	r := transform.NewReader(buf, unicode.BOMOverride(decoder))
	s, _ := ioutil.ReadAll(r)
	return string(s)
}
