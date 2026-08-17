// Package store 提供 SQLite 持久化。
package store

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

// Draw 一期开奖结果（JSON 结构即 /api/draws 返回结构）。
type Draw struct {
	Issue string `json:"issue"` // 期号，如 "2026094"
	Date  string `json:"date"`  // 开奖日期 "YYYY-MM-DD"
	Red   [6]int `json:"red"`   // 6 个红球，升序
	Blue  int    `json:"blue"`  // 蓝球
}

// Store 封装 SQLite 读写。
type Store struct {
	db *sql.DB
}

const upsertSQL = `INSERT INTO draws (issue, date, red1, red2, red3, red4, red5, red6, blue)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	ON CONFLICT(issue) DO UPDATE SET
		date = excluded.date,
		red1 = excluded.red1, red2 = excluded.red2, red3 = excluded.red3,
		red4 = excluded.red4, red5 = excluded.red5, red6 = excluded.red6,
		blue = excluded.blue`

func upsertArgs(d Draw) []any {
	return []any{d.Issue, d.Date, d.Red[0], d.Red[1], d.Red[2], d.Red[3], d.Red[4], d.Red[5], d.Blue}
}

// Open 打开（不存在则创建）数据库并建表。
func Open(path string) (*Store, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	const schema = `CREATE TABLE IF NOT EXISTS draws (
		issue TEXT PRIMARY KEY,
		date  TEXT NOT NULL,
		red1  INTEGER NOT NULL,
		red2  INTEGER NOT NULL,
		red3  INTEGER NOT NULL,
		red4  INTEGER NOT NULL,
		red5  INTEGER NOT NULL,
		red6  INTEGER NOT NULL,
		blue  INTEGER NOT NULL
	)`
	if _, err := db.Exec(schema); err != nil {
		db.Close()
		return nil, err
	}
	return &Store{db: db}, nil
}

// Close 关闭数据库。
func (s *Store) Close() error { return s.db.Close() }

// Upsert 幂等写入一期：期号已存在则覆盖。
func (s *Store) Upsert(d Draw) error {
	_, err := s.db.Exec(upsertSQL, upsertArgs(d)...)
	return err
}

// UpsertMany 批量幂等写入（单事务）。
func (s *Store) UpsertMany(draws []Draw) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	for _, d := range draws {
		if _, err := tx.Exec(upsertSQL, upsertArgs(d)...); err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

// LatestIssue 返回库中最新期号；空库返回 ""。
func (s *Store) LatestIssue() (string, error) {
	var issue string
	err := s.db.QueryRow(`SELECT issue FROM draws ORDER BY issue DESC LIMIT 1`).Scan(&issue)
	if err == sql.ErrNoRows {
		return "", nil
	}
	return issue, err
}

// List 返回最新在前的最多 limit 期；before 非空时仅返回期号小于 before 的期。
func (s *Store) List(limit int, before string) ([]Draw, error) {
	query := `SELECT issue, date, red1, red2, red3, red4, red5, red6, blue FROM draws`
	var args []any
	if before != "" {
		query += ` WHERE issue < ?`
		args = append(args, before)
	}
	query += ` ORDER BY issue DESC LIMIT ?`
	args = append(args, limit)

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var draws []Draw
	for rows.Next() {
		var d Draw
		if err := rows.Scan(&d.Issue, &d.Date, &d.Red[0], &d.Red[1], &d.Red[2], &d.Red[3], &d.Red[4], &d.Red[5], &d.Blue); err != nil {
			return nil, err
		}
		draws = append(draws, d)
	}
	return draws, rows.Err()
}

// Count 返回库中总期数。
func (s *Store) Count() (int, error) {
	var n int
	err := s.db.QueryRow(`SELECT COUNT(*) FROM draws`).Scan(&n)
	return n, err
}
