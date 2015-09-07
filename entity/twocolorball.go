package entity

import (
	"database/sql"
)

//	双色球
type TwoColorBall struct {
	ID     int64
	NextID sql.NullInt64
	No     string
	Date   string
	Red1   int
	Red2   int
	Red3   int
	Red4   int
	Red5   int
	Red6   int
	Blue1  int
}

func (b *TwoColorBall) ReadRows(rows *sql.Rows) error {
	return rows.Scan(&b.ID, &b.NextID, &b.No, &b.Date, &b.Red1, &b.Red2, &b.Red3, &b.Red4, &b.Red5, &b.Red6, &b.Blue1)
}

type TwoColorBallSummary struct {
	List []TwoColorBall
	Red  []AnalyzeResult
	Blue []AnalyzeResult
}
