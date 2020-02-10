package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"
)

var resultFmt map[string]string
var nameAry map[string]string

const (
	GAD7 string = "gad7"
	PHQ9 string = "phq9"
	PHQ15 string = "phq15"
	GHQ12 string = "ghq12"
	SASRQ string = "sasrq"
	PCLC string = "pcl-c"
	STAI string = "stai"
)

func init() {
	nameAry = make(map[string]string)
	nameAry[GAD7] = `患者健康问卷焦虑分量表`
	nameAry[PHQ15] = `患者健康问卷抑郁筛查量表`
	nameAry[PHQ15] = `躯体症状问卷`
	nameAry[GHQ12] = `一般健康问卷`
	nameAry[SASRQ] = `斯坦福急性应激反应问卷`
	nameAry[PCLC] = `PTSD17项筛查问卷`
	nameAry[STAI] = `状态一特质焦虑问卷`

	resultFmt = make(map[string]string)
	resultFmt[GAD7] = `
1.GAD-7主要用于筛查评估广泛性焦虑障碍，但是也可用于其他焦虑障碍的筛查。主要评估指标为总分：直接将7条项目分相加，总分为21分。
2.焦虑障碍筛查阳性划界分：≥10分
3.焦虑状态的严重程度划分如下：
	0-4分   无焦虑
	5-9分   轻度焦虑
	10-14  中度焦虑
	15-21  重度焦虑`

	resultFmt[PHQ9] = `
1.PHQ-9总分≥7分，提示存在抑郁发作，需要与医生面诊确认具体诊断，并指导治疗。
2.提示抑郁症状的严重程度：
      0-4 无抑郁症状；       5-9 轻微抑郁症状；
      10-14 中度抑郁症状；   15-19 重度抑郁症状；
      ≥20  极重度抑郁症状；
`

	resultFmt[PHQ15] = `
0-4分: 无躯体症状；
5-9分: 轻度躯体症状；
10-14分: 中度躯体症状；
15-30分: 重度躯体症状；
`

	resultFmt[GHQ12] = `
回答前两项者计0分，回答后两项者计1分，总分范围为0-12分，3分为化界分，也就是说，总分≥3，表示存在压力，需要做进一步的临床评估以帮助明确具体的心理困难（如PHQ-9、GAD-7、PHQ-15或SCL-90）。此时一方面需要积极的自我调整，同时在感到自己难以应对的时候需要积极寻求专业的心理卫生服务。 
`

	resultFmt[SASRQ] = `
SASRQ用于评估创伤幸存者的急性应激反应。该问卷共30个项目，28项用于评估急性应激障碍的症状，2项用于评估社交和职业功能损害（9、26），每一项采用0-5分的6级评估，见指导语。急性应激障碍项目分为4个症状维度，包括：1.分离性症状，由3、4、10、13、16、18、20、24、25、28共10项组成，亚量表分范围0-50分。而分离性症状又包括5个症状，分别为麻木和疏离感（20、28）、环境意识水平降低（4、24）、现实解体体验（3、18）、人格解体体验（10、13）和分离性遗忘（16、25）；2.创伤的再体验症状，由6、7、15、19、23、29共6项组成，亚量表分范围0-30分；3.回避症状，由5、11、14、17、22、30共6项组成，亚量表分范围0-30分；4.显著焦虑或唤起反应增加，由1、2、8、12、21、27共6项组成，亚量表分范围0-30分。
判断每一条症状存在的频率分应至少为3分，即每一条的自评分≥3时说明该条症状存在。
若要做出急性应激障碍的临床诊断（印象），需要同时满足下述条件：
1.5项分离性症状中需同时存在至少3条；
2.至少存在1项创伤的再体验症状；
3.至少存在1项回避症状；
4.至少存在1项显著焦虑或唤起反应增加症状。
`

	resultFmt[PCLC] = `
 1) 简介 
创伤后应激障碍清单（PTSD checklist- Civilian Version, PCL-C）包含有DSM-Ⅳ有关创伤后应激障碍的所有17项症状。是自评量表，由DSM- IV中有关PTSD的诊断标准构成。PCL-C总分范围17-85分，分数越高，代表PTSD发生的可能性越大。闯入症状共5项（1-5），亚总分范围5-25分；回避/麻木症状共7项（6-12），亚总分范围7-35分；警觉症状5项（13-17），亚总分范围5-25分。量表为5级评分制，1=从不，5=极重度。要求受试者对于每一项症状按照在过去的一个月受到烦扰的程度进行评定。在以往的研究中，这一测量已经证明与临床用PTSD诊断量表（Clinician Administered PTSD scale, CAPS）相似的心理测量学特性。PCL-C有很好的重测信度和内部一致性，汇聚效度可以由与密西西比量表（Mississippi Scale for PTSD）的高度相关来证明。如果受试者总分大于等于50分，就很有可能诊断为PTSD，50分的分界有较好的诊断灵敏性（0.82）和特异性（0.83），Kappa系数为0.64按照5分（1-5）等级评分法。此量表测定起来简单易行，适合大样本人群的筛查。
2) 评分方法
问卷中每一条目的评分分为5级: “没有什么反应”评1分； “轻度反应” 评2分；“中度反应”评3分；“重度反应” 评4分；“极重度反应” 评5分。累积记分大于50分即为筛选阳性，进入PTSD的诊断程序。
`

	resultFmt[STAI] = `
STAI的内容、评定与计分方法：
一、内容：由指导语和二个分量表共40项描述题组成。第1-20项为状态焦虑量表(STAI，Form Y-I，以下简称S-AD。其中半数为描述负性情绪的条目，半数为正性情绪条目。主要用于评定即刻的或最近某一特定时间或情景的恐惧、紧张、忧虑和神经质的体验或感受。可用来评价应激情况下的状态焦虑。第21-40题为特质焦虑量表（STAI，Form Y-l，简称T-AI)，用于评定人们经常的情绪体验。其中有11项为描述负性情绪条目，9项为正性情绪条目。可广泛应用于评定内科、外科、心身疾病及精神病人的焦虑情绪；也可用来筛查高校学生、军人，和其他职业人群的有关焦虑问题；以及评价心理治疗、药物治疗的效果。
二、评定方法：该问卷由自我评定或自我报告来完成。受试者根据指导语逐题圈出答案。可用于个人或集体测试，受试者一般需具有初中文化水平。测查无时间限制，一般10-20分钟可完成整个量表条目的回答。
计分法：STAI每一项进行1-4级评分S-AI：1一完全没有，2一有些，3一中等程度，4一非常明显。T-AI：1一几乎没有，2一有些，3一经常，4一几乎总是如此。由受试者根据自己的体验选圈最合适的分值。凡正性情绪项目均为反序计分。分别计算S-AI和T-AI量表的累加分，最小值20，最大值为80，反映状态或特质焦虑的程度。
`
}

func main() {
	r := gin.Default()

	// route
	r.GET("/", getIndex)
	//r.GET("/questions", getQuestions)

	///!!!
	// get /q/x -> post /q/answerx/x
	r.GET("/q/:name", getQx)
	r.POST("/q/answerx/:name", postAnswerX)

	r.GET("/qlist", getList)

	//r.POST("/answers", postAnswers)
	//r.GET("/answers", getAnswers)
	//r.POST("/answers_readonly", postAnswersReadonly)
	//r.GET("/answers_readonly", getAnswersReadonly)

	r.StaticFS("/css", http.Dir("../css"))
	r.StaticFS("/js", http.Dir("../js"))
	//r.StaticFS("/templates", http.Dir("templates"))
	r.LoadHTMLGlob("../templates/*")

	//r.StaticFile("/questions", "./questions.html")
	//r.GET("/upload", onUpload)

	r.Run(":8080")
}

func getList(c *gin.Context) {
	c.HTML(http.StatusOK, "qlist.html", gin.H{})
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
	c.HTML(http.StatusOK, name+".html", gin.H{})
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
			"des":  fmt.Sprintf("Here is your descriptions: scroe %v", data["score"]),
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
			"des":  fmt.Sprintf("Here is your descriptions: scroe %v", data["score"]),
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
			"des":  fmt.Sprintf("Here is your descriptions: scroe %v", data["score"]),
		}
		c.HTML(http.StatusOK, "output.html", d2)
	} else {
		c.String(http.StatusOK, "no answers")
	}
}

func postAnswersReadonly(c *gin.Context) {
	getAnswersReadonly(c)
}

func getScore(name string, qAry []int) int {
	result := 0
	sum := 0
	for i := 0; i < len(qAry); i++ {
		sum += qAry[i]
	}

	switch name {

	default:
		result = sum
	}

	return result
}

func getDesByScore(scoreAry []int , desAry []string, curScore int) string {
	for i := 0; i < len(scoreAry); i++  {
		if curScore <= scoreAry[i] {
			return desAry[i]
		}
	}

	return desAry[len(desAry) - 1]
}

func getDes(name string, totalScore int) string {
	log.Printf("getDes by %v score %v \n", name, totalScore)

	switch name {
	case GAD7:
		{
			//0-4分   无焦虑
			//5-9分   轻度焦虑
			//10-14  中度焦虑
			//15-21  重度焦虑
			scoreAry := []int{4,9,14,21}
			desAry := []string {"无焦虑", "轻度焦虑", "中度焦虑", "重度焦虑"}
			//desAry := []string {"", "", "", ""}
			return getDesByScore(scoreAry, desAry, totalScore)
		}

	case PHQ9:
		{
			scoreAry := []int{4,9,14,19, 1000}
			desAry := []string {"无抑郁症状", "轻微抑郁症状", "中度抑郁症状", "重度抑郁症状", "极重度抑郁症状"}
			//desAry := []string {"", "", "", ""}
			return getDesByScore(scoreAry, desAry, totalScore)
		}

	case PHQ15:
		{
			scoreAry := []int{4,9,14,30}
			desAry := []string {"无躯体症状", "轻度躯体症状", "中度躯体症状", "重度躯体症状"}
			//desAry := []string {"", "", "", ""}
			return getDesByScore(scoreAry, desAry, totalScore)
		}

	case GHQ12:
		{
			scoreAry := []int{3,100}
			desAry := []string {"一切OK", "存在压力"}
			//desAry := []string {"", "", "", ""}
			return getDesByScore(scoreAry, desAry, totalScore)
		}

	case SASRQ:
		{

		}

	case PCLC:
		{

		}

	case STAI:
		{

		}
	}

	return "not found " + name
}

func postAnswerX(c *gin.Context) {
	// 问卷类型
	name := c.Param("name")
	fmt.Printf("name: %v\n", name)

	// 问卷答案
	qAry := make([]int, 0)
	for i := 1; i < 100; i++ {
		// li name
		qx := "q" + strconv.Itoa(i)
		v, ok := c.GetPostForm(qx)
		if !ok {
			fmt.Printf("break: %v\n", qx)
			break
		}

		vi, err := strconv.Atoi(v)
		if err != nil {
			fmt.Errorf("err:%v\n", err)
			break
		}

		qAry = append(qAry, vi)
	}

	log.Printf("qAry:%v\n", qAry)

	score := getScore(name, qAry)
	des := getDes(name, score)
	des_ref := resultFmt[name]
	nameLang := nameAry[name]

	data := gin.H{
		"name":    nameLang,
		"score":   score,
		"des":     des,
		"des_ref": des_ref,
	}

	c.HTML(http.StatusOK, "answerx.html", data)

}
