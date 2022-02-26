package core

import (
	"fmt"

	"github.com/maru44/scheman/definition"
)

type (
	rel struct {
		TableName    string
		RelTableName string
		Nullable     bool
	}
)

func (s *SchemanState) genMermaid(isIgnoreView bool) string {
	var (
		fKeys, ones, manys []rel
		cols               string
		rels               string
	)
	for _, t := range s.Tables {
		if t.IsView && isIgnoreView {
			continue
		}

		for _, r := range t.FKeys {
			fKeys = append(fKeys, rel{
				TableName:    r.Table,
				RelTableName: r.ForeignTable,
				Nullable:     r.Nullable,
			})
		}
		for _, r := range t.ToOneRelationships {
			ones = append(ones, rel{
				TableName:    r.Table,
				RelTableName: r.ForeignTable,
				Nullable:     r.Nullable,
			})
		}
		for _, r := range t.ToManyRelationships {
			manys = append(manys, rel{
				TableName:    r.Table,
				RelTableName: r.ForeignTable,
				Nullable:     r.Nullable,
			})
		}

		cols += fmt.Sprintf("%s {\n", t.Name)
		for _, col := range t.Columns {
			c := definition.ConvertCol(col, t.PKey, s.Config.DriverName)
			cols += fmt.Sprintf("  %s %s\n", c.DBType, c.Name)
		}
		cols += "}\n\n"
	}

	for _, f := range fKeys {
		isOne := false
		for _, r := range ones {
			if f.TableName == r.RelTableName && f.RelTableName == r.TableName {
				switch f.Nullable {
				case true:
					if r.Nullable {
						rels += fmt.Sprintf("%s |o--o| %s : own\n", r.TableName, f.TableName)
					} else {
						rels += fmt.Sprintf("%s |o--|| %s : own\n", r.TableName, f.TableName)
					}
				case false:
					if r.Nullable {
						rels += fmt.Sprintf("%s ||--o| %s : own\n", r.TableName, f.TableName)
					} else {
						rels += fmt.Sprintf("%s ||--|| %s : own\n", r.TableName, f.TableName)
					}
				}
				isOne = true
				break
			}
		}
		if isOne {
			continue
		}
		for _, r := range manys {
			if f.TableName == r.RelTableName && f.RelTableName == r.TableName {
				switch f.Nullable {
				case true:
					if r.Nullable {
						rels += fmt.Sprintf("%s |o--o{ %s : has \n", r.TableName, f.TableName)
					} else {
						rels += fmt.Sprintf("%s |o--|{ %s : has \n", r.TableName, f.TableName)
					}
				case false:
					if r.Nullable {
						rels += fmt.Sprintf("%s ||--o{ %s : has \n", r.TableName, f.TableName)
					} else {
						rels += fmt.Sprintf("%s ||--|{ %s : has \n", r.TableName, f.TableName)
					}
				}
				isOne = true
				break
			}
		}
	}
	return "erDiagram\n\n" + rels + "\n\n" + cols
}
