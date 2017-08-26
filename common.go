package main

import (
	"net/http"
	"io/ioutil"
	"fmt"
	"github.com/bitly/go-simplejson"
	"strconv"
	"database/sql"
	//"encoding/json"
	"encoding/json"
	"strings"
	"crypto/md5"
	"encoding/hex"
	"time"
	"math/rand"
	"reflect"
)

/**
  * 作者：林冠宏
  *
  * author: LinGuanHong,
  *
  * My GitHub : https://github.com/af913337456/
  *
  * My Blog   : http://www.cnblogs.com/linguanh/
  *
  * */

func Log(a ...interface{})  {
	fmt.Println(a)
}

/** 抽离：头部公共处理 */
func commonReqBody(req *http.Request)  *simplejson.Json{
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		//WriteHttpError(400, err.Error(), w)
		fmt.Println("read all body error")
		return nil
	}

	obj, err := simplejson.NewJson(body)
	if err != nil {
		fmt.Println("输入的 jsonObj 不存在 或 格式错误")
		return nil
	}
	return obj
}

/** 根据结构体插入数据库 */
/** 过期 */
//func commonInsertByStruct(
//  	w http.ResponseWriter,
//	r *http.Request,
//	strct interface{},
//	funcName string,
//	table string,
//	unSetCallBack func(*[]interface{}),
//	addIdToSql bool) {
//
//	obj := commonReqBody(r)
//
//	if obj ==nil {
//		fmt.Println("json 是 null")
//		return
//	}
//	slices := []interface{}{}
//	buildSlices(&slices,strct,obj,funcName)
//	unSetCallBack(&slices)
//
//	sqlStr := buildInsertSqlByStruct(new (sticker_info),table,addIdToSql)
//
//	log("sql: "+sqlStr)
//	log(slices)
//
//	_,lastId := insertExe(sqlStr,slices...)
//
//	echoResult(w,lastId)
//}

func exeSqlCommonHandler(
	w http.ResponseWriter,
	obj *simplejson.Json,
	funcName string,
	inputStruct  interface{},
	slicesCallBack func(slices []interface{}) bool,
	getSqlCallBack func(slices []interface{},inputStruct interface{}) string) ([]interface{},string) {

	slices := []interface{}{}
	if inputStruct != nil {
		values := reflect.ValueOf(inputStruct).Elem() // 值
		valueLen := values.NumField()

		keys   := reflect.TypeOf(inputStruct).Elem()  // 键

		for i := 0; i < valueLen; i++ {
			fieldV := values.Field(i)
			fieldK := keys  .Field(i)
			keyName := fieldK.Tag.Get("json")
			nullTag := fieldK.Tag.Get("nullTag")
			switch fieldV.Interface().(type) {
			case int:
			case int64:
				if selectSliceHandler(w,&slices,obj,funcName,keyName,nullTag,false) {
					return nil,""
				}
				break
			case string:
			case *string:
				if selectSliceHandler(w,&slices,obj,funcName,keyName,nullTag,true) {
					return nil,""
				}
				break
			}
		}
		if slicesCallBack(slices) {
			/** 终止 */
			fmt.Fprintf(w,"null");
			Log("callBack: stop")
			return nil,""
		}
	}
	return slices,getSqlCallBack(slices,inputStruct)
}


// sql := "insert LSticker(`f_user_id`,`f_content`,`f_time`) values(?,?,?)"
func buildInsertSqlByStrArr(arr []string,table string) string {

	length := len(arr)

	header := "insert into "+table+"("
	partV  := ") values("

	for i := 0; i < length; i++ {

		keyName := arr[i]
		if i==length-1 {
			header = header + "`"+keyName+"`"
			partV  = partV + "?)"
		}else{
			header = header + "`"+keyName+"`,"
			partV  = partV +"?,"
		}
	}
	return header + partV
}

func buildInsertSqlByStruct(stuct interface{},table string) string {

	k := reflect.TypeOf(stuct).Elem() // 键
	length := k.NumField()

	header := "insert into "+table+"("
	partV  := ") values("

	for i := 0; i < length; i++ {
		//if !addIdToSql && i==0 {
		//	continue
		//}
		fieldK := k.Field(i)
		keyName := fieldK.Tag.Get("json") // 对应 struct 里面的 `json`
		if i==length-1 {
			header = header + "`"+keyName+"`"
			partV  = partV + "?)"
		}else{
			header = header + "`"+keyName+"`,"
			partV  = partV +"?,"
		}
	}
	return header + partV
}

func selectSliceHandler(
	w http.ResponseWriter,
	slice *[]interface{},
	obj *simplejson.Json,
	funcName,keyName,nullTag string,
	isStr bool) bool {
	if !addParamToSlice(slice,obj,keyName,isStr) {
		Log(funcName+" 传入的 jsonObj "+keyName+" 为 null")
		if nullTag=="1" {
			errorRet(w)
			return true
		}
	}
	return false
}

func buildDefaultSelectSql(table string) string {
	return "select * from `"+table+"`"
}

func errorRet(w http.ResponseWriter)  {
	fmt.Fprintf(w,"null");
}

func getTimeMisecond() int64 {
	t := time.Now() //获取当前时间的结构体
	millis := t.Unix() //毫
	return millis
}

func echoResult(w http.ResponseWriter,r int64)  {
	if r<=0 {
		fmt.Fprintf(w,"-1")
	}else{
		fmt.Fprintf(w,"1")
	}
}

/** 结果用json的格式返回 */
func getResultJson(args ...interface{}) string {
	length := len(args)
	var str   string

	var result string
	keys := []string{"result","data","lastId"}

	for i:=0;i<length;i++ {
		switch v := args[i].(type) {
		case int64:
			var s1 int64 = v
			str = "\""+keys[i]+"\":"+strconv.FormatInt(s1,10)
			break;
		case int:
			var s int = v
			str = "\""+keys[i]+"\":"+strconv.FormatInt(int64(s),10)
			break
		case string:
			str = "\""+keys[i]+"\":\""+v+"\""
			break
		}
		if i==length - 1 {
			result = result + str;
			break
		}
		result = result + str+",";
	}
	return "{"+result+"}"
}

func getResultJsonWithKeys(keys []string,args ...interface{}) string {
	length := len(args)
	var str   string

	var result string

	for i:=0;i<length;i++ {
		switch v := args[i].(type) {
		case int64:
			var s1 int64 = v
			str = "\""+keys[i]+"\":"+strconv.FormatInt(s1,10)
			break;
		case int:
			var s int = v
			str = "\""+keys[i]+"\":"+strconv.FormatInt(int64(s),10)
			break
		case string:
			str = "\""+keys[i]+"\":\""+v+"\""
			break
		}
		if i==length - 1 {
			result = result + str;
			break
		}
		result = result + str+",";
	}
	return "{"+result+"}"
}

/** 32 位 MD5 */
func md532(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	cipherStr := h.Sum(nil)
	return hex.EncodeToString(cipherStr)
}

/** 产生范围随机数 */
func getRandNum(start,end int) int64 {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return int64(r.Intn((end - start)) + start)
}

/** 字符串截取 */
func Substr2(str string,) string {
	if str == "[]"{
		return "[]"
	}
	rs := []rune(str)
	end := len(rs)-1
	start := 1

	return string(rs[start:end])
}

//截取字符串 start 起点下标 length 需要截取的长度
func Substr3(str string, start int, length int) string {
	rs := []rune(str)
	rl := len(rs)
	end := 0

	if start < 0 {
		start = rl - 1 + start
	}
	end = start + length

	if start > end {
		start, end = end, start
	}

	if start < 0 {
		start = 0
	}
	if start > rl {
		start = rl
	}
	if end < 0 {
		end = 0
	}
	if end > rl {
		end = rl
	}

	return string(rs[start:end])
}

func string2int(postId string) int64 {
	value,err := strconv.Atoi(postId)
	if err!=nil {
		return -1
	}
	return int64(value)
}

func int2String(tar int64) string {
	return strconv.FormatInt(tar,10)
}

func byte2string( b []byte ) string {
	s := make([]string,len(b))
	for i := range b {
		s[i] = strconv.Itoa(int(b[i]))
	}
	return strings.Join(s,",")
}

/** 评论和点赞的删除 */
/** 同时删除提醒表的 */
func deleteMoreById(tableName string,id int64) string {

	row,_ := deleteExe(
		"delete `"+tableName+"`,`wake` from `"+tableName+"` left join `wake` on `"+tableName+"`.idOfWake=`wake`.id where `"+tableName+"`.id=?",id)
	if row<0 {
		return "-1";
	}else{
		return strconv.FormatInt(row, 10) /** 删除成功，返回影响的行数 */
	}
}

/** 公共的，传入函数来自定义处理 */
func getJsonWithSql(handler func(rows *sql.Rows) interface{},rows *sql.Rows) (string,error)  {
	jsonData, err := json.Marshal(handler(rows))
	if err != nil {
		fmt.Println(err)
		return "",err
	}
	fmt.Println("返回的 json 是"+string(jsonData))
	return string(jsonData),nil
}

func getJsonData(w http.ResponseWriter,data []interface{})  {
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Fprintf(w,"null")
		return
	}
	stringData := string(jsonData)
	if(stringData=="[]"){
		fmt.Println("空数据集合 data=='[]'")
		fmt.Fprintf(w,"null")
		return
	}
	fmt.Fprintf(w,stringData)
}

/** 集合 sqlHelper 一次 */
func getJsonDataWithSql(
w http.ResponseWriter,
req *http.Request,
sct interface{},
getSqlAndSlice func (obj *simplejson.Json) (string,[]interface{}))  {
	sqlStr,slice := getSqlAndSlice(commonReqBody(req))
	dataSlice,err := selectExe(sqlStr,
		func(i interface{}) interface{} {
			return i
		},sct,slice...)

	if err!=nil {
		fmt.Fprintf(w,"null")
		return
	}
	getJsonData(w,dataSlice)
}

/** 输出 struct */
func structToJsonStr(stuct interface{})  {
	fmt.Println(struct {

	}{})
}

