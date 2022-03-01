package file

import (
	"os"
	"strings"

	"github.com/friendsofgo/errors"
)

func (f *File) writeCSV(showAttrs []string) error {
	file, err := os.OpenFile(f.definitionFile, os.O_APPEND|os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600)
	if err != nil {
		return errors.Wrap(err, "failed to open definition file")
	}

	var out string
	rowsSlice := f.makeRowsSlice(showAttrs)
	for _, tables := range rowsSlice {
		for _, columns := range tables {
			out += strings.Join(columns, ",") + "\n"
		}
		out += "\n"
	}

	if _, err := file.Write([]byte(out)); err != nil {
		return errors.Wrap(err, "failed to write definition file")
	}

	return nil
}
