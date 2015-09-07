package entity

import (
	"database/sql"
)

//	大乐透
type SuperLotto struct {
	ID     int64
	NextID sql.NullInt64
	No     string
	Date   string
	Red1   int
	Red2   int
	Red3   int
	Red4   int
	Red5   int
	Blue1  int
	Blue2  int
}

func (b *SuperLotto) ReadRows(rows *sql.Rows) error {
	return rows.Scan(&b.ID, &b.NextID, &b.No, &b.Date, &b.Red1, &b.Red2, &b.Red3, &b.Red4, &b.Red5, &b.Blue1, &b.Blue2)
}

type SuperLottoSummary struct {
	List []SuperLotto
	Red  []AnalyzeResult
	Blue []AnalyzeResult
}
