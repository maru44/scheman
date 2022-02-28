package file

import (
	"context"
	"os"

	"github.com/friendsofgo/errors"
)

func (f *File) writeCSV(ctx context.Context) error {
	file, err := os.OpenFile(f.definitionFile, os.O_APPEND|os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600)
	if err != nil {
		return errors.Wrap(err, "failed to open definition file")
	}

	return nil
}
