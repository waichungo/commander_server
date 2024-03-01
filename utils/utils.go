package utils

import (
	"bytes"
	"compress/gzip"
	"io"
	"strings"
)

func InSlice(arr []string, search string, strict bool) bool {
	element := ""
	for i := 0; i < len(arr); i++ {
		element = arr[i]
		if strict {
			if strings.EqualFold(element, search) {
				return true
			}
		} else {
			if element == search {
				return true
			}
		}
	}
	return false
}
func RepeatString(value string, n int) []string {
	arr := make([]string, n)
	for i := 0; i < n; i++ {
		arr[i] = value
	}
	return arr
}
func DecompressGzip(reader io.Reader) ([]byte, error) {
	gzReader, err := gzip.NewReader(reader)
	if err != nil {
		return nil, err
	}
	defer gzReader.Close()
	return io.ReadAll(gzReader)
}
func CompressGzip(data []byte) ([]byte, error) {

	var b bytes.Buffer
	w := gzip.NewWriter(&b)
	_, err := w.Write(data)
	if err != nil {
		return nil, err
	}
	w.Close()
	return b.Bytes(), err
}
func RemoveEmptyFromSlice(src []string) []string {
	final := make([]string, 0, 20)
	for _, el := range src {
		el = strings.TrimSpace(el)
		if len(el) > 0 {
			final = append(final, el)
		}

	}
	return final
}
