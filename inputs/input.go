package inputs

import (
	"io"
	"net/http"
	"os"
	"strings"
)

func OpenAll(paths []string) ([]io.ReadCloser, func(), error) {
	files := make([]io.ReadCloser, 0, len(paths))

	closer := func() {
		for _, file := range files {
			file.Close()
		}
	}

	for _, path := range paths {
		if file, err := readCloserFor(path); err != nil {
			closer()
			return nil, nil, err
		} else {
			files = append(files, file)
		}
	}

	return files, closer, nil
}

func readCloserFor(path string) (io.ReadCloser, error) {
	if path == "-" {
		return os.Stdin, nil
	} else if !strings.HasPrefix(path, "https://") && !strings.HasPrefix(path, "http://") {
		return os.Open(path)
	} else if response, err := http.Get(path); err != nil {
		return nil, err
	} else {
		return response.Body, nil
	}
}
