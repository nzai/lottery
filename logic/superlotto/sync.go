package superlotto

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"strconv"

	"github.com/nzai/lottery/conn"
	"github.com/nzai/lottery/entity"
	"github.com/nzai/lottery/logic/crypto"
	"github.com/nzai/lottery/logic/util"
)

//	同步数据登录验证
func SyncData() error {

	//  抓取开奖结果
	results, err := fetchData()
	if err != nil {
		log.Print("fetchData")
		return err
	}
	//log.Println(results)

	//  保存开奖结果
	err = SaveData(results)
	if err != nil {
		log.Print("SaveData")
		return err
	}

	return nil
}

//  保存开奖结果
func SaveData(fetched []entity.SuperLotto) error {

	//	连接数据库
	db, err := conn.GetConn()
	if err != nil {
		log.Println("数据库初始化失败 : ", err.Error())
		return err
	}
	defer db.Close()

	//	查询
	rows, err := db.Query("SELECT * FROM SuperLotto ORDER BY No ASC")
	if err != nil {
		return err
	}
	defer rows.Close()

	//  查询所有保存过的记录
	saved := make([]entity.SuperLotto, 0)
	for rows.Next() {
		var item entity.SuperLotto
		err = item.ReadRows(rows)
		if err != nil {
			return err
		}

		saved = append(saved, item)
	}

	//  过滤出需要保存的结果
	toSave := make([]entity.SuperLotto, 0)
	for _, item := range fetched {
		exists := false
		for _, savedItem := range saved {
			if item.No == savedItem.No {
				exists = true
				break
			}
		}

		if !exists {
			//  没保存过的结果添加进待保存队列
			toSave = append(toSave, item)
		}
	}
	//log.Println("需要新增的开奖结果:", toSave)

	//  启动事务
	transaction, err := db.Begin()

	//  保存结果的语句
	stmtSL, err := transaction.Prepare("INSERT INTO SuperLotto VALUES(?, NULL, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		//  回滚事务
		transaction.Rollback()
		return err
	}
	defer stmtSL.Close()

	//  保存彩球的语句
	stmtBall, err := transaction.Prepare("INSERT INTO LotteryBall VALUES(?, ?, 2, ?, ?)")
	if err != nil {
		//  回滚事务
		transaction.Rollback()
		return err
	}
	defer stmtBall.Close()

	for _, item := range toSave {

		//  保存结果
		_, err = stmtSL.Exec(item.ID, item.No, item.Date, item.Red1, item.Red2, item.Red3, item.Red4, item.Red5, item.Blue1, item.Blue2)
		if err != nil {
			//  回滚事务
			transaction.Rollback()

			return err
		}

		//  保存彩球
		_, err = stmtBall.Exec(crypto.GetUniqueInt64(), item.ID, item.Red1, 1)
		if err != nil {
			//  回滚事务
			transaction.Rollback()

			return err
		}

		_, err = stmtBall.Exec(crypto.GetUniqueInt64(), item.ID, item.Red2, 1)
		if err != nil {
			//  回滚事务
			transaction.Rollback()

			return err
		}

		_, err = stmtBall.Exec(crypto.GetUniqueInt64(), item.ID, item.Red3, 1)
		if err != nil {
			//  回滚事务
			transaction.Rollback()

			return err
		}

		_, err = stmtBall.Exec(crypto.GetUniqueInt64(), item.ID, item.Red4, 1)
		if err != nil {
			//  回滚事务
			transaction.Rollback()

			return err
		}

		_, err = stmtBall.Exec(crypto.GetUniqueInt64(), item.ID, item.Red5, 1)
		if err != nil {
			//  回滚事务
			transaction.Rollback()

			return err
		}

		_, err = stmtBall.Exec(crypto.GetUniqueInt64(), item.ID, item.Blue1, 2)
		if err != nil {
			//  回滚事务
			transaction.Rollback()

			return err
		}

		_, err = stmtBall.Exec(crypto.GetUniqueInt64(), item.ID, item.Blue2, 2)
		if err != nil {
			//  回滚事务
			transaction.Rollback()

			return err
		}
	}

	//  提交事务
	transaction.Commit()

	sql := `
	UPDATE SuperLotto
	SET NextID = (SELECT B.ID FROM SuperLotto B WHERE B.Date > SuperLotto.Date ORDER BY B.Date ASC LIMIT 1)
	WHERE NextID IS NULL`

	_, err = db.Exec(sql)
	if err != nil {
		return err
	}

	if len(toSave) > 0 {
		log.Println("新增大乐透开奖结果", len(toSave), "组")
	}

	return nil
}

//  通过列表抓取开奖结果
func fetchData() ([]entity.SuperLotto, error) {
	//  抓取网页
	frontHtml, err := util.DownloadHtml("http://www.lottery.gov.cn/historykj/history.jspx?_ltype=dlt")
	if err != nil {
		log.Println("下载网页失败: ", err.Error())
		return nil, err
	}
	//log.Printf("%s", html)

	//  获取一共有几页
	regex := regexp.MustCompile(`共\d+条记录 \d+\/(\d+)页`)
	group := regex.FindSubmatch(frontHtml)
	if len(group) < 1 {
		return nil, errors.New("分析结果页数失败")
	}

	pageCount, _ := strconv.Atoi(string(group[1]))
	// log.Println("共", pageCount, "页")

	var html string
	list := make([]entity.SuperLotto, 0)
	for index := 1; index <= pageCount; index++ {
		if index == 1 {
			//  首页已经抓取过了
			html = string(frontHtml)
		} else {
			//  抓取网页
			content, err := util.DownloadHtml(fmt.Sprintf("http://www.lottery.gov.cn/historykj/history_%d.jspx?_ltype=dlt", index))
			if err != nil {
				log.Println("下载网页失败: ", err.Error())
				return nil, err
			}

			html = string(content)
		}

		//  分析Html,抓取信息
		results, err := analyzeHtml(html)
		if err != nil {
			log.Println("下载网页失败: ", err.Error())
			return nil, err
		}

		//  把每页的分析结果添加进结果集
		list = append(list, results...)
	}

	log.Printf("已经获取了%d组大乐透数据\n", len(list))
	return list, nil
}

//  分析网页抓取开奖结果
func analyzeHtml(html string) ([]entity.SuperLotto, error) {
	//  使用正则分析网页
	regex := regexp.MustCompile(`<td height="23"[^>]*?>(\d+)<\/td>\s+<td[^>]*?>(\d+)<\/td>\s+<td[^>]*?>(\d+)<\/td>\s+<td[^>]*?>(\d+)<\/td>\s+<td[^>]*?>(\d+)<\/td>\s+<td[^>]*?>(\d+)<\/td>\s+<td[^>]*?>(\d+)<\/td>\s+<td[^>]*?>(\d+)<\/td>\s+<td[^>]*?>\S+<\/td>\s+<td[^>]*?>\S+<\/td>\s+<td[^>]*?>\S+<\/td>\s+<td[^>]*?>\S+<\/td>\s+<td[^>]*?>\S+<\/td>\s+<td[^>]*?>\S+<\/td>\s+<td[^>]*?>\S+<\/td>\s+<td[^>]*?>\S+<\/td>\s+<td[^>]*?>[\s\S]+?<\/td>\s+<td[^>]*?>\S+<\/td>\s+<td[^>]*?>\S+<\/td>\s+<td[^>]*?>(\S+)<\/td>`)

	group := regex.FindAllStringSubmatch(html, -1)
	//log.Println(group)
	//log.Println("共计:", len(group))

	results := make([]entity.SuperLotto, 0)
	for _, section := range group {

		item := entity.SuperLotto{}
		item.ID = crypto.GetUniqueInt64()
		item.No = section[1]
		item.Date = section[9]
		item.Red1, _ = strconv.Atoi(section[2])
		item.Red2, _ = strconv.Atoi(section[3])
		item.Red3, _ = strconv.Atoi(section[4])
		item.Red4, _ = strconv.Atoi(section[5])
		item.Red5, _ = strconv.Atoi(section[6])
		item.Blue1, _ = strconv.Atoi(section[7])
		item.Blue2, _ = strconv.Atoi(section[8])

		results = append(results, item)
	}

	return results, nil
}
