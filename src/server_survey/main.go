package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

func main() {
	r := gin.Default()

	// route
	r.GET("/", getIndex)
	r.GET("/questions", getQuestions)

	r.GET("/q/:name", getQx)

	r.POST("/answers", postAnswers)
	r.GET("/answers", getAnswers)
	r.POST("/answers_readonly", postAnswersReadonly)
	r.GET("/answers_readonly", getAnswersReadonly)

	r.StaticFS("/css", http.Dir("css"))
	r.StaticFS("/js", http.Dir("js"))
	//r.StaticFS("/templates", http.Dir("templates"))
	r.LoadHTMLGlob("templates/*")

	//r.StaticFile("/questions", "./questions.html")
	//r.GET("/upload", onUpload)

	r.Run(":8080")
}

func onUpload(c *gin.Context) {
	c.HTML(http.StatusOK, "upload.html", gin.H{})
}

func tryResetCookie(c *gin.Context) {
	value, err := c.Cookie("session_id")
	if err != nil {
		fmt.Printf("trySetCookie err %v\n", err)
	}
	//if len(value) != 0 {
	//	return
	//}

	value = fmt.Sprintf("%v", time.Now().UnixNano()) //TODO ..UDID...
	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    value,
		Path:     "/",
		HttpOnly: true,
	}
	http.SetCookie(c.Writer, cookie)

	fmt.Printf("cookie %v\n", cookie)
}

func getIndex(c *gin.Context) {
	//c.String(http.StatusOK, "It's root")

	c.Redirect(http.StatusTemporaryRedirect, "/questions")

	// dosent work
	//c.Request.URL.Path = "/questions"
	//gin.Default().HandleContext(c)
}

func getQuestions(c *gin.Context) {
	//c.String(http.StatusOK, "It's questions")
	//panic!
	tryResetCookie(c)
	c.HTML(http.StatusOK, "questions.html", gin.H{})
}

func getQx(c *gin.Context) {
	//c.String(http.StatusOK, "It's questions")
	//panic!
	tryResetCookie(c)

	//x := c.Param("x")
	name := c.Param("name")
	c.HTML(http.StatusOK, name + ".html", gin.H{})
}

var res = make(map[string]gin.H)
var resMx sync.RWMutex

func getAnswer(c *gin.Context) gin.H {
	resMx.RLock()
	defer resMx.RUnlock()

	v, err := c.Cookie("session_id")
	if err != nil {
		//c.String(http.StatusOK, "You have no answers")
		return nil
	}

	if ret, ok := res[v]; ok {
		return ret
	}

	return nil
}

func setAnswer(c *gin.Context, data gin.H) {
	resMx.Lock()
	defer resMx.Unlock()

	v, err := c.Cookie("session_id")
	if err != nil {
		//c.String(http.StatusOK, "You have no answers")
		return
	}
	res[v] = data
}

func postAnswers(c *gin.Context) {
	n := c.PostForm("name")
	q1 := c.PostForm("q1")

	//qAry := c.PostFormArray("q")
	//log.Printf("qAry %v\n", qAry)

	//c.Request.ParseForm()
	//for k, v := range c.Request.PostForm {
	//	fmt.Printf("k:%v\n", k)
	//	fmt.Printf("v:%v\n", v)
	//}

	v, err := c.Cookie("session_id")

	//TODO js层加重复提交的判断
	//TODO 客户端通过重定向到get页面，来防止页面刷新导致的重提交
	// redirect to get, but failed,, still post?

	data := gin.H{
		"name":       n,
		"q1":         q1,
		"score":      fmt.Sprintf("%d", rand.Intn(100)),
		"session_id": v,
	}

	// 不支持cookie
	if err != nil || len(v) == 0 {
		d2 := gin.H{
			"name": data["name"],
			"des": fmt.Sprintf("Here is your descriptions: scroe %v", data["score"]),
		}
		c.HTML(http.StatusOK, "output.html", d2)
	} else {
		setAnswer(c, data)

		// RPG redirect post get
		c.Redirect(http.StatusTemporaryRedirect, "answers_readonly")
	}
}

func getAnswers(c *gin.Context) {
	data := getAnswer(c)
	if data != nil {
		//c.JSON(http.StatusOK, data)
		d2 := gin.H{
			"name": data["name"],
			"des": fmt.Sprintf("Here is your descriptions: scroe %v", data["score"]),
		}
		c.HTML(http.StatusOK, "output.html", d2)
	} else {
		c.String(http.StatusOK, "no answers")
	}
}

func getAnswersReadonly(c *gin.Context) {
	data := getAnswer(c)
	if data != nil {
		//c.JSON(http.StatusOK, data)
		d2 := gin.H{
			"name": data["name"],
			"des": fmt.Sprintf("Here is your descriptions: scroe %v", data["score"]),
		}
		c.HTML(http.StatusOK, "output.html", d2)
	} else {
		c.String(http.StatusOK, "no answers")
	}
}

func postAnswersReadonly(c *gin.Context) {
	getAnswersReadonly(c)
}
