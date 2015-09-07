package superlotto

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/nzai/lottery/conn"
	"github.com/nzai/lottery/entity"
)

//  分析1
func Analyze1(reds []int, blues []int) ([]entity.AnalyzeResult, error) {

	if len(reds) == 0 && len(blues) == 0 {
		return nil, errors.New("参数为空")
	}

	//    log.Println("len:", reds)
	condition := bytes.NewBufferString("")
	params := make([]interface{}, 0)
	sql := `
        SELECT (LB.BallType - 1) * 100 + LB.Ball Num, COUNT(LB.Ball) Count
        FROM 
        (
            SELECT SL.ID
            FROM SuperLotto SL
            WHERE 1=1 %s
        ) SL
        LEFT JOIN LotteryBall LB ON LB.MainID = SL.ID
        GROUP BY (LB.BallType - 1) * 100 + LB.Ball`

	//  查询红球
	if len(reds) != 0 {
		redsString := bytes.NewBufferString("")
		for _, value := range reds {
			if redsString.Len() > 0 {
				redsString.WriteString(",")
			}

			_, err := redsString.WriteString(strconv.Itoa(value))
			if err != nil {
				return nil, err
			}
		}

		condition.WriteString(fmt.Sprintf(`
            AND EXISTS
            (
                SELECT 1
                FROM LotteryBall LB
                WHERE LB.MainID = SL.ID
                AND LB.Ball IN (%s)
                AND LB.RecordType = 2
                AND LB.BallType = 1
                GROUP BY LB.MainID
                HAVING COUNT(LB.MainID) >= ?
            ) `, redsString.String()))

		params = append(params, len(reds))
	}

	//  查询蓝球
	if len(blues) != 0 {
		bluesString := bytes.NewBufferString("")
		for _, value := range blues {
			if bluesString.Len() > 0 {
				bluesString.WriteString(",")
			}

			_, err := bluesString.WriteString(strconv.Itoa(value))
			if err != nil {
				return nil, err
			}
		}

		condition.WriteString(fmt.Sprintf(`
            AND EXISTS
            (
                SELECT 1
                FROM LotteryBall LB
                WHERE LB.MainID = SL.ID
                AND LB.Ball IN (%s)
                AND LB.RecordType = 2
                AND LB.BallType = 2
                GROUP BY LB.MainID
                HAVING COUNT(LB.MainID) >= ?
            ) `, bluesString.String()))

		params = append(params, len(blues))
	}

	//	连接数据库
	db, err := conn.GetConn()
	if err != nil {
		log.Println("数据库初始化失败 : ", err.Error())
		return nil, err
	}
	defer db.Close()

	s := fmt.Sprintf(sql, condition.String())
	//log.Println("SQL:", s)
	//log.Println("Parameter:", params)

	//	查询
	rows, err := db.Query(s, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	//  查询所有保存过的记录
	results := make([]entity.AnalyzeResult, 0)
	for rows.Next() {
		var item entity.AnalyzeResult
		err = item.ReadRows(rows)
		if err != nil {
			return nil, err
		}

		results = append(results, item)
	}

	return results, nil
}

//  分析2
func Analyze2(reds []int, blues []int) ([]entity.AnalyzeResult, error) {

	if len(reds) == 0 && len(blues) == 0 {
		return nil, errors.New("参数为空")
	}

	//    log.Println("len:", reds)
	condition := bytes.NewBufferString("")
	params := make([]interface{}, 0)
	sql := `
        SELECT (LB.BallType - 1) * 100 + LB.Ball Num, COUNT(LB.Ball) Count
        FROM 
        (
            SELECT SL.NextID
            FROM SuperLotto SL
            WHERE 1=1 %s
        ) SL
        JOIN SuperLotto SLN ON SLN.ID = SL.NextID
        LEFT JOIN LotteryBall LB ON LB.MainID = SLN.ID
        GROUP BY (LB.BallType - 1) * 100 + LB.Ball`

	//  查询红球
	if len(reds) != 0 {
		redsString := bytes.NewBufferString("")
		for _, value := range reds {
			if redsString.Len() > 0 {
				redsString.WriteString(",")
			}

			_, err := redsString.WriteString(strconv.Itoa(value))
			if err != nil {
				return nil, err
			}
		}

		condition.WriteString(fmt.Sprintf(`
            AND EXISTS
            (
                SELECT 1
                FROM LotteryBall LB
                WHERE LB.MainID = SL.ID
                AND LB.Ball IN (%s)
                AND LB.RecordType = 2
                AND LB.BallType = 1
                GROUP BY LB.MainID
                HAVING COUNT(LB.MainID) >= ?
            ) `, redsString.String()))

		params = append(params, len(reds))
	}

	//  查询蓝球
	if len(blues) != 0 {
		bluesString := bytes.NewBufferString("")
		for _, value := range blues {
			if bluesString.Len() > 0 {
				bluesString.WriteString(",")
			}

			_, err := bluesString.WriteString(strconv.Itoa(value))
			if err != nil {
				return nil, err
			}
		}

		condition.WriteString(fmt.Sprintf(`
            AND EXISTS
            (
                SELECT 1
                FROM LotteryBall LB
                WHERE LB.MainID = SL.ID
                AND LB.Ball IN (%s)
                AND LB.RecordType = 2
                AND LB.BallType = 2
                GROUP BY LB.MainID
                HAVING COUNT(LB.MainID) >= ?
            ) `, bluesString.String()))

		params = append(params, len(blues))
	}

	//	连接数据库
	db, err := conn.GetConn()
	if err != nil {
		log.Println("数据库初始化失败 : ", err.Error())
		return nil, err
	}
	defer db.Close()

	s := fmt.Sprintf(sql, condition.String())
	//log.Println("SQL:", s)
	//log.Println("Parameter:", params)

	//	查询
	rows, err := db.Query(s, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	//  查询所有保存过的记录
	results := make([]entity.AnalyzeResult, 0)
	for rows.Next() {
		var item entity.AnalyzeResult
		err = item.ReadRows(rows)
		if err != nil {
			return nil, err
		}

		results = append(results, item)
	}

	return results, nil
}
