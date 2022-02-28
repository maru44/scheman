package file

import (
	"context"
	"os"
	"path/filepath"

	"github.com/friendsofgo/errors"
	"github.com/maru44/scheman/definition"
)

type (
	File struct {
		*definition.CommonInfo
		rawMermaid string

		definitionFile string
		erdFile        string
	}
)

func NewFile(definitionFile, erdFile string, info *definition.CommonInfo) definition.Definition {
	return &File{
		definitionFile: definitionFile,
		erdFile:        erdFile,
		CommonInfo:     info,
	}
}

func (f *File) SetMermaid(m string) {
	f.rawMermaid = m
}

func (f *File) Upsert(ctx context.Context) error {
	if f.definitionFile == "" {
		return nil
	}

	ext := filepath.Ext(f.definitionFile)
	switch ext {
	case "csv":
		return f.writeCSV(ctx)
	case "json":
	case "tsv":
	}
	return nil
}

func (f *File) Mermaid(ctx context.Context) error {
	if f.erdFile == "" {
		return nil
	}

	file, err := os.OpenFile(f.erdFile, os.O_APPEND|os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return errors.Wrap(err, "failed to open mermaid file")
	}

	if _, err := file.Write([]byte(f.rawMermaid)); err != nil {
		return errors.Wrap(err, "failed to write mermaid file")
	}
	return nil
}
