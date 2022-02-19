package db

import (
	"fmt"

	"github.com/volatiletech/sqlboiler/v4/boilingcore"
)

type SchemanState struct {
	*boilingcore.State
}

func (s *SchemanState) Run() error {
	fmt.Println(s.Tables)
	return nil
}
