package main

/**
  * 作者：林冠宏
  *
  * author: LinGuanHong
  *
  * My GitHub : https://github.com/af913337456/
  *
  * My Blog   : http://www.cnblogs.com/linguanh/
  *
  * */

/**

 	根据 sql 文件，自动生成代码文件，包含有 struct.go，每张表对应生成一个包含有增删改查的基础方法文件

*/

import (
	"fmt"
	"strings"
	"os"
	"bufio"
)

const STR_ARR_SIZE  = 50 /** 控制表的字段个数 */

func main()  {

	file,err := os.Open("this.sql")
	defer file.Close()
	if err!=nil {
		fmt.Println("自动生成代码错误")
		fmt.Println(err)
		return
	}
	reader := bufio.NewReader(file)

	isBeginOneTable  := false
	isFinishOneTable := false
	//result := ""
	var tableName string
	values := [STR_ARR_SIZE]string{}
	valueType := [STR_ARR_SIZE]string{}
	i := 0
	deleteOld()
	for{
		line,err := reader.ReadString('\n')
		if err!=nil {
			isBeginOneTable = false
			isFinishOneTable =false
			fmt.Println("读取完")
			fmt.Println(err)
			break
		}
		if  strings.Contains(line,"CREATE database") ||
			strings.Contains(line,"CREATE DATABASE")    {

		} else if strings.Contains(line,"CREATE TABLE") {
			isBeginOneTable  = true
			isFinishOneTable = false
			tableName = getValue(line)
			fmt.Println("表名字 "+tableName)
			//result = "tableName: "+tableName
		}else if ((strings.Contains(line,"COMMENT") && !strings.Contains(line,"ENGINE")) ||
			strings.Contains(line,"varchar") ||
			strings.Contains(line,"datetime") ||
			strings.Contains(line,"int") ){
			//if !strings.Contains(line,"`id`") { /** 跳过主键 id */
				value := getValue(line)
				value_type := getType(line)
				// fmt.Println("value ： "+value+"  type : "+value_type)
				// fmt.Println("下标 ： "+strconv.FormatInt(int64(i),10))
				values[i] = value
				valueType[i] = value_type
				i++
			//}
		}else if strings.Contains(line,"ENGINE") {
			isBeginOneTable  = false;
			isFinishOneTable = true;
			//fmt.Println("建表 "+tableName)
			//fmt.Println(valueType)
			//fmt.Println(values)
			outputStructFile(tableName,&values,&valueType)
			//outputSqlStrFile(tableName,"id",&values,&valueType)
			outputFuncFile(tableName)

			// 重置
			values = [STR_ARR_SIZE]string{}
			valueType = [STR_ARR_SIZE]string{}
			i = 0
		}
	}
	if(isBeginOneTable || isFinishOneTable){

	}
}

func outputSqlStrFile(tableName string,defaultKey string,values,valueType *[STR_ARR_SIZE]string)  {
	fileName := "sql_str.go";
	file,err := os.Create(fileName)
	defer file.Close()
	if err!=nil {
		fmt.Println(err)
	}
	file.WriteString("package main\n")
	file.WriteString("var "+tableName+"Values = []string{"+createValuesVar(values)+"}\n")
	file.WriteString("var "+tableName+"ValueTypes = []string{"+createValuesVar(valueType)+"}\n")

	file.WriteString(createInsertSql(tableName,values))

	file.WriteString("var select"+tableName+"AllSql string = \"select * from `"+tableName+"`\"\n")
	file.WriteString("var select"+tableName+"LimitSql string = \"select * from `"+tableName+"` where `"+defaultKey+"`<? order by `"+defaultKey+"` desc limit 40\"\n")
	file.WriteString("var select"+tableName+"ByIdSql string = \"select * from `"+tableName+"` where `"+defaultKey+"`=?\"\n")

	file.WriteString("var delete"+tableName+"ByIdSql string = \"delete from `"+tableName+"` where `"+defaultKey+"`=?\"\n")

	file.WriteString(createUpdateSql(tableName,defaultKey,values))
	file.WriteString("\n")
}

func outputFuncFile(tableName string)  {
	fileName := "func_"+tableName+".go"
	_,err := os.Stat(fileName)
	if err == nil {
		fmt.Println("文件："+fileName+" 已经存在，为了避免覆盖你的内容，请删除后再生成该文件！")
		return
	}
	file,err := os.Create(fileName)
	defer file.Close()
	if err!=nil {
		fmt.Println(err)
	}
	file.WriteString("package main\n\n")
	file.WriteString("import (\n")
	/** 默认的包 */
	file.WriteString("	\"net/http\"\n")
	//file.WriteString("	\"fmt\"\n")
	file.WriteString(")\n\n")
	file.WriteString("/**\n" +
		"type LghRequest struct {\n"+
		"	w http.ResponseWriter\n"+
		"	r *http.Request\n"+
		"	funcName string\n"+
		"	inputStruct  interface{}\n"+
		"	slicesCallBack func(slices *[]interface{}) bool\n"+
		"	getSqlCallBack func(slices *[]interface{},inputStruct interface{}) string\n"+
		"}\n")
	file.WriteString("r.HandleFunc(\"/insert_"+strings.ToLower(tableName)+"\",insert_"+strings.ToLower(tableName)+")\n")
	file.WriteString("r.HandleFunc(\"/delete_"+strings.ToLower(tableName)+"\",delete_"+strings.ToLower(tableName)+")\n")
	file.WriteString("r.HandleFunc(\"/update_"+strings.ToLower(tableName)+"\",update_"+strings.ToLower(tableName)+")\n")
	file.WriteString("r.HandleFunc(\"/select_"+strings.ToLower(tableName)+"\",select_"+strings.ToLower(tableName)+")\n")
	file.WriteString("*/\n\n")
	newVersion(file,tableName)
}

func createValuesVar(values *[STR_ARR_SIZE]string) string {
	r := ""
	l := len(values)
	for i:=0 ; i<l ; i++ {
		if values[i] == "" {
			break
		}
		if i==l-1 {
			r = r + "`"+values[i]+"\""
		}else{
			if(values[i+1] == ""){
				r = r + "\""+values[i]+"\""
			}else{
				r = r + "\""+values[i]+"\","
			}
		}
	}
	return r
}

func createUpdateSql(tableName,defaultKey string,values *[STR_ARR_SIZE]string) string {
	sqlInsert := ""
	l := len(values)
	for i:=0 ; i<l ; i++ {
		if values[i] == "" {
			break
		}
		if i==l-1 {
			sqlInsert = sqlInsert + "`"+values[i]+"`=?"
		}else{
			if(values[i+1] == ""){
				sqlInsert = sqlInsert + "`"+values[i]+"`=?"
			}else{
				sqlInsert = sqlInsert + "`"+values[i]+"`=?,"
			}
		}
	}
	sql := "var update"+tableName+"ByIdSql string = \"update `"+tableName+"` set "+sqlInsert+" where `"+defaultKey+"`=?\"\n"
	return sql
}

func createInsertSql(tableName string,values *[STR_ARR_SIZE]string) string {
	sqlStart := "var insert"+tableName+"Sql string = \"insert into `"+tableName+"`("
	l := len(values)
	breakIndex := 0
	for i:=0 ; i<l ; i++ {
		if values[i] == "" {
			sqlStart = sqlStart + ") values("
			breakIndex = i;
			break
		}
		if i==l-1 {
			sqlStart = sqlStart + "`"+values[i]+"`) values("
		}else{
			if(values[i+1] == ""){
				sqlStart = sqlStart + "`"+values[i]+ "`"
			}else{
				sqlStart = sqlStart + "`"+values[i]+"`,"
			}
		}
	}
	for j:=0;j<l;j++ {
		if j == breakIndex {
			sqlStart = sqlStart + ")\"\n"
			break
		}
		if j == l-1 {
			sqlStart = sqlStart + "?)\"\n"
		}else {
			if values[j+1] == "" {
				sqlStart = sqlStart + "?"
			}else{
				sqlStart = sqlStart + "?,"
			}
		}
	}
	return sqlStart
}

func deleteOld(){
	os.Remove("struct.go")
}

func outputStructFile(tableName string,values,valuesType *[STR_ARR_SIZE]string)  {
	fileName := "struct.go";
	isExists := checkFileAndCreate(fileName);
	file,err := os.OpenFile(fileName,os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	defer file.Close()
	if err!=nil {
		fmt.Println("outputStructFile 207")
		fmt.Println(err)
		return
	}
	if !isExists {
		file.WriteString("package main\n")
		file.WriteString("\n")
		file.WriteString("// nullTag: 0 代表可以在读取post过来的参数时候为 null. 1 不能为 null \n")
	}

	file.WriteString("type "+tableName+" struct { ")
	file.WriteString("\n")
	length := len(values)
	for i:=0 ; i<length ; i++ {
		if values[i]=="" {
			file.WriteString("}\n")
			break
		}
		rv := ""
		temp := strings.Split(values[i], "_")
		tempL := len(temp)
		for j:=1;j<tempL;j++{
			if j==tempL-1 {
				rv = rv + temp[j]
			}else{
				rv = rv + temp[j] +"_"
			}
		}
		if rv =="" {
			rv = strings.ToUpper(Substr4(temp[0],0,1))+Substr4(temp[0],1,len(temp[0]))
		}else{
			rv = strings.ToUpper(temp[0])+"_"+rv
		}
		file.WriteString("	"+rv+"			"+valuesType[i]+"	`json:\""+values[i]+"\" nullTag:\"1\"`")
		file.WriteString("\n")
		//if valuesType[i] == "int" {
		//	file.WriteString("	"+rv+"			int64	`json:\""+values[i]+"\" nullTag:\"1\"`")
		//	file.WriteString("\n")
		//}else if valuesType[i] == "*string" {
		//	file.WriteString("	"+rv+"			*string	`json:\""+values[i]+"\" nullTag:\"1\"`")
		//	file.WriteString("\n")
		//}
	}
}

func checkFileAndCreate(fileName string) bool {
	_,err := os.Stat(fileName)
	if err!=nil {
		os.Create(fileName)
		return false
	}else{
		fmt.Println("已存在 "+fileName)
		return true
	}
}

func getValue(line string) string {
	firstIndex := strings.Index(line,"`")+1
	value := Substr4(line,firstIndex,strings.LastIndex(line,"`")-firstIndex)
	return value
}

func getType(line string) string {
	if strings.Contains(line,"datetime") || strings.Contains(line,"varchar") {
		return "*string";
	}else if strings.Contains(line,"int"){
		return "int";
	}else {
		return "*string"
	}
}

func Substr4(str string, start int, length int) string {
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

func newVersion(file *os.File,tableName string)  {

	/** 插入 */
	file.WriteString(
		"func insert_"+strings.ToLower(tableName)+"(w http.ResponseWriter,r *http.Request)  {\n"+
		"	request := LghRequest{\n"+
		"		w,\n"+
		"		r,\n"+
		"		\"insert_"+strings.ToLower(tableName)+"\",\n"+
		"		new ("+tableName+"),\n"+
		"		func(slices []interface{}) bool{\n"+
		"			return false\n"+
		"		},\n"+
		"		func(slices []interface{},inputStruct interface{}) string {\n"+
		"			return buildInsertSqlByStruct(new ("+tableName+"),\""+tableName+"\")\n"+
		"		}}\n"+
		"	insertDataByStruct(request)\n"+
		"}\n\n")

	/** 删除 */
	file.WriteString(
		"func delete_"+strings.ToLower(tableName)+"(w http.ResponseWriter,r *http.Request)  {\n"+
		"	request := LghRequest{\n"+
		"		w,\n"+
		"		r,\n"+
		"		\"delete_"+strings.ToLower(tableName)+"\",\n"+
		"		nil,\n"+
		"		func(slices []interface{}) bool{\n"+
		"			return false\n"+
		"		},\n"+
		"		func(slices []interface{},inputStruct interface{}) string {\n"+
		"			return \"delete from "+tableName+" where id='2'\"\n"+
		"		}}\n"+
		"	deleteDataByStruct(request)\n"+
		"}\n\n")

	/** 改 */
	file.WriteString(
		"func update_"+strings.ToLower(tableName)+"(w http.ResponseWriter,r *http.Request)  {\n"+
		"	type testS struct {\n"+
		"		Id int64 `json:\"id\" nullTag:\"1\"`\n"+
		"	}\n"+
		"	request := LghRequest{\n"+
		"		w,\n"+
		"		r,\n"+
		"		\"update_"+strings.ToLower(tableName)+"\",\n"+
		"		new (testS),\n"+
		"		func(slices []interface{}) bool{\n"+
		"			// 演示使用输入参数的情况\n"+
		"			return false\n"+
		"		},\n"+
		"		func(slices []interface{},inputStruct interface{}) string {\n"+
		"			return \"update "+tableName+" set u_user_id='444' where id=?\"\n"+
		"		}}\n"+
		"	updateDataByStruct(request)\n"+
		"}\n\n")
	/** 查 */
	file.WriteString(
		"func select_"+strings.ToLower(tableName)+"(w http.ResponseWriter,r *http.Request)  {\n"+
		"	request := LghRequest{\n"+
		"		w,\n"+
		"		r,\n"+
		"		\"select_"+strings.ToLower(tableName)+"\",\n"+
		"		nil,\n"+
		"		func(slices []interface{}) bool{\n"+
		"			// 演示使用输入参数的情况\n"+
		"			return false\n"+
		"		},\n"+
		"		func(slices []interface{},inputStruct interface{}) string {\n"+
		"			return \"select * from "+tableName+"\"\n"+
		"		}}\n"+
		"		output := new ("+tableName+")\n"+
		"	selectDataByStruct(request,output)\n"+
		"}\n\n")

}

func olderVersion(file *os.File,tableName string)  {
	/** 插入 */
	file.WriteString(
		"func insert_"+strings.ToLower(tableName)+"(w http.ResponseWriter,r *http.Request)  {\n" +
			"	json := commonReqBody(r)\n" +
			"	slice := []interface{}{}\n" +
			"	l := len("+tableName+"Values)\n" +
			"	for i:=0;i<l;i++ {\n" +
			"		if "+tableName+"ValueTypes[i] == \"int\"  {\n" +
			"			addParamToSlice(&slice,json,"+tableName+"Values[i],false)\n" +
			"		}else{\n" +
			"			addParamToSlice(&slice,json,"+tableName+"Values[i],true)\n" +
			"		}\n" +
			"	}\n" +
			"	rowNums,lastId := insertExe(insert"+tableName+"Sql,slice...)\n" +
			"	if rowNums <=0 {\n" +
			"		keys := []string{\"lastId\",\"code\"}\n" +
			"		fmt.Fprintf(w,getResultJsonWithKeys(keys,true,lastId,1))\n" +
			"	}else{\n" +
			"		keys := []string{\"code\"}\n" +
			"		fmt.Fprintf(w,getResultJsonWithKeys(keys,true,-1))\n" +
			"	}\n" +
			"}\n\n")

	/** 删除 */
	file.WriteString(
		"func delete_"+strings.ToLower(tableName)+"(w http.ResponseWriter,r *http.Request) {\n" +
			"	json := commonReqBody(r)\n" +
			"	slice := []interface{}{}\n" +
			"	addParamToSlice(&slice,json,\"id\",false)\n" +
			"	rowNums,lastId := deleteExe(delete"+tableName+"ByIdSql,slice...)\n" +
			"	if rowNums <=0 {\n" +
			"		keys := []string{\"lastId\",\"code\"}\n" +
			"		fmt.Fprintf(w,getResultJsonWithKeys(keys,true,lastId,1))\n" +
			"	}else{\n" +
			"		keys := []string{\"code\"}\n" +
			"		fmt.Fprintf(w,getResultJsonWithKeys(keys,true,-1))\n" +
			"	}\n" +
			"}\n\n")

	/** 改 */
	file.WriteString(
		"func update_"+strings.ToLower(tableName)+"(w http.ResponseWriter,r *http.Request)  {\n" +
			"	json := commonReqBody(r)\n" +
			"	slice := []interface{}{}\n" +
			"	l := len("+tableName+"Values)\n" +
			"	for i:=0;i<l;i++ {\n" +
			"		if "+tableName+"ValueTypes[i] == \"int\"  {\n" +
			"			addParamToSlice(&slice,json,"+tableName+"Values[i],false)\n" +
			"		}else{\n" +
			"			addParamToSlice(&slice,json,"+tableName+"Values[i],true)\n" +
			"		}\n" +
			"	}\n" +
			"	addParamToSlice(&slice,json,\"id\",false)\n" +
			"	rowNums,lastId := updateExe(update"+tableName+"ByIdSql,slice...)\n" +
			"	if rowNums <=0 {\n" +
			"		keys := []string{\"lastId\",\"code\"}\n" +
			"		fmt.Fprintf(w,getResultJsonWithKeys(keys,true,lastId,1))\n" +
			"	}else{\n" +
			"		keys := []string{\"code\"}\n" +
			"		fmt.Fprintf(w,getResultJsonWithKeys(keys,true,-1))\n" +
			"	}\n" +
			"}\n\n")

	/** 查 */
	file.WriteString(
		"func get_"+strings.ToLower(tableName)+"(w http.ResponseWriter,r *http.Request)  {\n" +
			"	json := commonReqBody(r)\n" +
			"	slice := []interface{}{}\n" +
			"	sql := \"\"\n" +
			"	if addParamToSlice(&slice,json,\"lastId\",false) {\n" +
			"		sql = select"+tableName+"LimitSql;\n" +
			"	}else if addParamToSlice(&slice,json,\"id\",false){\n" +
			"		sql = select"+tableName+"ByIdSql\n" +
			"	}else{\n" +
			"		sql = select"+tableName+"AllSql\n" +
			"	}\n" +
			"	dataSlice,err := selectExe(sql,func (i interface{}) interface{} {\n" +
			"		return i\n" +
			"	},new ("+tableName+"),slice)\n" +
			"	if err!=nil {\n" +
			"		keys := []string{\"code\"}\n" +
			"		fmt.Fprintf(w,getResultJsonWithKeys(keys,true,-1))\n" +
			"		return\n" +
			"	}\n" +
			"	getJsonData(w,dataSlice)\n" +
			"}\n\n")
}

