//go:generate go run ./gen/gen.go -in . -out help_generated.go

package helpdoc

import (
	"bytes"
	"compress/gzip"
	"io"
)

var helpDocs = make(map[string][]byte)

func add(name string, data []byte) {
	helpDocs[name] = data
}

func GetHelpDocument(docName string) (io.Reader, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(helpDocs[docName]))
	if err != nil {
		return nil, err
	}
	defer gz.Close()

	b, err := io.ReadAll(gz)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(b), nil
}
