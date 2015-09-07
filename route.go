package main

import (
	"fmt"
	"log"
	"lottery/app/entity"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nzai/lottery/logic/superlotto"
	"github.com/nzai/lottery/logic/twocolorball"
)

//	注册路由
func RegisterRoute(route *gin.Engine) {

	//	静态文件目录
	route.Static("static", "./static")

	//	注册模板
	registerTemplates("static/html/layout.html", "static/html/twocolorball/index.html")
	registerTemplates("static/html/layout.html", "static/html/twocolorball/analyze1.html")
	registerTemplates("static/html/layout.html", "static/html/twocolorball/analyze2.html")
	registerTemplates("static/html/layout.html", "static/html/superlotto/index.html")
	registerTemplates("static/html/layout.html", "static/html/superlotto/analyze1.html")
	registerTemplates("static/html/layout.html", "static/html/superlotto/analyze2.html")

	//	首页默认转到双色球列表
	route.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/twocolorball")
	})

	//	双色球
	t := route.Group("/twocolorball")
	{
		//	列表页面
		t.GET("/", func(c *gin.Context) {
			err := processTemplate(c.Writer, "static/html/twocolorball/index.html", nil)
			if err != nil {
				log.Print(err)
			}
		})

		//	列表查询
		t.GET("/list", func(c *gin.Context) {

			paramN := c.Query("n")
			n, err := strconv.Atoi(paramN)
			if err != nil {
				c.JSON(http.StatusOK, entity.ResultSM(false, fmt.Sprintf("错误的查询参数N:%s", paramN)))
				return
			}

			//  查询
			results, err := twocolorball.Query(n)
			if err != nil {
				c.JSON(http.StatusOK, entity.ResultSM(false, fmt.Sprintf("查询双色球数据失败:%v", err)))
				return
			}

			c.JSON(http.StatusOK, entity.ResultD(results))
		})

		//	分析1页面
		t.GET("/analyze1", func(c *gin.Context) {
			err := processTemplate(c.Writer, "static/html/twocolorball/analyze1.html", nil)
			if err != nil {
				log.Print(err)
			}
		})

		//	分析1查询
		t.GET("/doanalyze1", func(c *gin.Context) {

			reds := c.Query("reds")
			blues := c.Query("blues")

			//  红球
			redNums := make([]int, 0)

			if len(reds) > 0 {
				parts := strings.Split(reds, ",")
				for _, value := range parts {
					i, err := strconv.Atoi(value)
					if err != nil {
						c.JSON(http.StatusOK, entity.ResultSM(false, fmt.Sprintf("输入的参数有误:%s", value)))
						return
					}

					redNums = append(redNums, i)
				}
			}

			//  蓝球
			blueNum, err := strconv.Atoi(blues)

			//  查询
			results, err := twocolorball.Analyze1(redNums, blueNum)
			if err != nil {
				c.JSON(http.StatusOK, entity.ResultSM(false, fmt.Sprintf("分析双色球数据失败:%v", err)))
				return
			}

			c.JSON(http.StatusOK, entity.ResultD(results))
		})

		//	分析2页面
		t.GET("/analyze2", func(c *gin.Context) {
			err := processTemplate(c.Writer, "static/html/twocolorball/analyze2.html", nil)
			if err != nil {
				log.Print(err)
			}
		})

		//	分析2查询
		t.GET("/doanalyze2", func(c *gin.Context) {

			reds := c.Query("reds")
			blues := c.Query("blues")

			//  红球
			redNums := make([]int, 0)

			if len(reds) > 0 {
				parts := strings.Split(reds, ",")
				for _, value := range parts {
					i, err := strconv.Atoi(value)
					if err != nil {
						c.JSON(http.StatusOK, entity.ResultSM(false, fmt.Sprintf("输入的参数有误:%s", value)))
						return
					}

					redNums = append(redNums, i)
				}
			}

			//  蓝球
			blueNum, err := strconv.Atoi(blues)

			//  查询
			results, err := twocolorball.Analyze2(redNums, blueNum)
			if err != nil {
				c.JSON(http.StatusOK, entity.ResultSM(false, fmt.Sprintf("分析双色球数据失败:%v", err)))
				return
			}

			c.JSON(http.StatusOK, entity.ResultD(results))
		})
	}

	//	大乐透
	s := route.Group("/superlotto")
	{
		//	列表页面
		s.GET("/", func(c *gin.Context) {
			err := processTemplate(c.Writer, "static/html/superlotto/index.html", nil)
			if err != nil {
				log.Print(err)
			}
		})

		//	列表查询
		s.GET("/list", func(c *gin.Context) {

			paramN := c.Query("n")
			n, err := strconv.Atoi(paramN)
			if err != nil {
				c.JSON(http.StatusOK, entity.ResultSM(false, fmt.Sprintf("错误的查询参数N:%s", paramN)))
				return
			}

			//  查询
			results, err := superlotto.Query(n)
			if err != nil {
				c.JSON(http.StatusOK, entity.ResultSM(false, fmt.Sprintf("查询大乐透数据失败:%v", err)))
				return
			}

			c.JSON(http.StatusOK, entity.ResultD(results))
		})

		//	分析1页面
		s.GET("/analyze1", func(c *gin.Context) {
			err := processTemplate(c.Writer, "static/html/superlotto/analyze1.html", nil)
			if err != nil {
				log.Print(err)
			}
		})

		//	分析1查询
		s.GET("/doanalyze1", func(c *gin.Context) {

			reds := c.Query("reds")
			blues := c.Query("blues")

			//  红球
			redNums := make([]int, 0)

			if len(reds) > 0 {
				parts := strings.Split(reds, ",")
				for _, value := range parts {
					i, err := strconv.Atoi(value)
					if err != nil {
						c.JSON(http.StatusOK, entity.ResultSM(false, fmt.Sprintf("输入的参数有误:%s", value)))
						return
					}

					redNums = append(redNums, i)
				}
			}

			//  蓝球
			blueNums := make([]int, 0)
			if len(blues) > 0 {
				parts := strings.Split(blues, ",")
				for _, value := range parts {
					i, err := strconv.Atoi(value)
					if err != nil {
						c.JSON(http.StatusOK, entity.ResultSM(false, fmt.Sprintf("输入的参数有误:%s", value)))
						return
					}

					blueNums = append(blueNums, i)
				}
			}

			//  查询
			results, err := superlotto.Analyze1(redNums, blueNums)
			if err != nil {
				c.JSON(http.StatusOK, entity.ResultSM(false, fmt.Sprintf("分析大乐透数据失败:%v", err)))
				return
			}

			c.JSON(http.StatusOK, entity.ResultD(results))
		})

		//	分析2页面
		s.GET("/analyze2", func(c *gin.Context) {
			err := processTemplate(c.Writer, "static/html/superlotto/analyze2.html", nil)
			if err != nil {
				log.Print(err)
			}
		})

		//	分析2查询
		s.GET("/doanalyze2", func(c *gin.Context) {

			reds := c.Query("reds")
			blues := c.Query("blues")

			//  红球
			redNums := make([]int, 0)

			if len(reds) > 0 {
				parts := strings.Split(reds, ",")
				for _, value := range parts {
					i, err := strconv.Atoi(value)
					if err != nil {
						c.JSON(http.StatusOK, entity.ResultSM(false, fmt.Sprintf("输入的参数有误:%s", value)))
						return
					}

					redNums = append(redNums, i)
				}
			}

			//  蓝球
			blueNums := make([]int, 0)
			if len(blues) > 0 {
				parts := strings.Split(blues, ",")
				for _, value := range parts {
					i, err := strconv.Atoi(value)
					if err != nil {
						c.JSON(http.StatusOK, entity.ResultSM(false, fmt.Sprintf("输入的参数有误:%s", value)))
						return
					}

					blueNums = append(blueNums, i)
				}
			}

			//  查询
			results, err := superlotto.Analyze2(redNums, blueNums)
			if err != nil {
				c.JSON(http.StatusOK, entity.ResultSM(false, fmt.Sprintf("分析大乐透数据失败:%v", err)))
				return
			}

			c.JSON(http.StatusOK, entity.ResultD(results))
		})
	}
}
