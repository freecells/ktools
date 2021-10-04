package pages

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

//PageNum 获取分页 页码
func pageCount(dataCount int, pageDataCount int) (pageCount int) {

	t := dataCount % pageDataCount
	if t > 0 {
		t = 1
	}
	pageCount = t + int(dataCount/pageDataCount)

	return
}

//AllPages 返回分页 html
func Pagination(dataCount int, pageDataCount int, c *gin.Context) (links string) {

	pageNum := 0

	pageNum = pageCount(dataCount, pageDataCount)

	if pageNum <= 1 {
		return ""
	}

	nowPage, _ := strconv.Atoi(c.DefaultQuery("page", "1"))

	params := ""

	t1 := strings.Split(c.Request.RequestURI, "?")

	if len(t1) > 1 {

		parameters := strings.Split(t1[1], "&")

		if len(parameters) > 1 {

			for index, val := range parameters {

				if strings.HasPrefix(val, "page") {

					parameters = append(parameters[:index], parameters[index+1:]...)
				}

			}

			params = "&" + strings.Join(parameters, "&")
		}

	}

	start := `<nav aria-label="Page navigation">
	<ul class="pagination">`

	end := `</ul></nav>`

	pre := `<li class='page-item'>
	<a class='page-link' href='?page=%d%s'> << </a>
	</li>`

	next := `<li class='page-item'>
	<a class='page-link' href='?page=%d%s'> >> </a>
	</li>`

	page := "<li class='page-item %s'><a class='page-link' href='?page=%d%s'>%d</a></li>"

	preStr := ""
	if nowPage <= 1 {
		preStr = ""
	} else {
		preStr = fmt.Sprintf(pre, nowPage-1, params)
	}

	nextStr := ""
	if nowPage >= pageNum {
		nextStr = ""
	} else {
		nextStr = fmt.Sprintf(next, nowPage+1, params)
	}

	pages := ""

	active := ""

	startNum := nowPage - 6
	endNum := nowPage + 6

	if startNum <= 0 {
		startNum = 1
	}

	if endNum > pageNum {
		endNum = pageNum
	}

	for i := startNum; i <= endNum; i++ {
		if nowPage == i {
			active = "active"
		} else {
			active = ""
		}
		pages += fmt.Sprintf(page, active, i, params, i)
	}

	links = start + preStr + pages + nextStr + end

	return
}
