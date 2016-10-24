package twocolorball

import (
	"bytes"
	"database/sql"
	"log"

	"github.com/nzai/lottery/conn"
	"github.com/nzai/lottery/entity"
)

//  查询
func Query(topN int) (*entity.TwoColorBallSummary, error) {

	//	连接数据库
	db, err := conn.GetConn()
	if err != nil {
		log.Println("数据库初始化失败: ", err.Error())
		return nil, err
	}
	defer db.Close()

	//  查询列表
	list, err := queryList(topN, db)
	if err != nil {
		return nil, err
	}

	//  查询多少期没有出了
	red, blue, err := queryDisappearCount(db)
	if err != nil {
		return nil, err
	}

	return &entity.TwoColorBallSummary{List: list, Red: red, Blue: blue}, nil
}

//  查询列表
func queryList(topN int, db *sql.DB) ([]entity.TwoColorBall, error) {

	//  构造SQL
	param := make([]interface{}, 0)
	querySQL := bytes.NewBufferString("SELECT T.* FROM (SELECT * FROM TwoColorBall ORDER BY No DESC")
	if topN > 0 {
		querySQL.WriteString(" LIMIT ?")
		param = append(param, topN)
	}
	querySQL.WriteString(") T ORDER BY T.NO ASC")
	//log.Println("sql: ", querySQL.String())

	//	查询
	rows, err := db.Query(querySQL.String(), param...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	//  查询所有保存过的记录
	results := make([]entity.TwoColorBall, 0)
	for rows.Next() {
		var item entity.TwoColorBall
		err = item.ReadRows(rows)
		if err != nil {
			return nil, err
		}

		results = append(results, item)
	}

	return results, nil
}

//  查询多少期没有出了
func queryDisappearCount(db *sql.DB) ([]entity.AnalyzeResult, []entity.AnalyzeResult, error) {

	sql := `
SELECT BL.Ball, BL.BallType, COUNT(TCB.No) DisappearCount
FROM
(
    SELECT BL.Ball, BL.BallType, IFNULL(MAX(BL.No), '') MaxNo
    FROM
    (
        SELECT BL.Ball, BL.BallType, TCB.No
        FROM LotteryBall BL
        JOIN TwoColorBall TCB ON TCB.ID = BL.MainID
        WHERE BL.RecordType = 1
    ) BL
    GROUP BY BL.Ball, BL.BallType
) BL
LEFT JOIN TwoColorBall TCB ON TCB.No > BL.MaxNo
GROUP BY BL.BallType, BL.Ball
ORDER BY BL.BallType, BL.Ball`

	//	查询
	rows, err := db.Query(sql)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	//  查询所有保存过的记录
	red := make([]entity.AnalyzeResult, 0)
	blue := make([]entity.AnalyzeResult, 0)
	var ballNum, ballType, disappearCount int
	for rows.Next() {
		err = rows.Scan(&ballNum, &ballType, &disappearCount)
		if err != nil {
			return nil, nil, err
		}

		item := entity.AnalyzeResult{ballNum, disappearCount}
		if ballType == 1 {
			red = append(red, item)
		} else {
			blue = append(blue, item)
		}
	}

	return red, blue, nil
}
