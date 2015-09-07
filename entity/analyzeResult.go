package entity

import (
	"database/sql"
)

//	分析结果
type AnalyzeResult struct {
	Num   int
	Count int
}

func (r *AnalyzeResult) ReadRows(rows *sql.Rows) error {
	return rows.Scan(&r.Num, &r.Count)
}
