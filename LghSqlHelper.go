package main

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

import (
	"log"
	"fmt"
	"github.com/bitly/go-simplejson"
	"reflect"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

/**
	采用了事务操作
	不要使用 sql 拼接，容易引起注入
  */

var (
	myDb    *sql.DB
	//dbUser   = flag.String("dbUser", "", "")
	//dbPw     = flag.String("dbPw"  , "", "")
	//dbName   = flag.String("dbName", "", "")
)

/** 获取数据库连接 */
func setDB(user,pw,dbName,dbPort string) *sql.DB{
	cmd := user+":"+pw+"@tcp(127.0.0.1:"+dbPort+")/"+dbName
	fmt.Println("connect mysql cmd is : "+cmd)
	db,err := sql.Open("mysql",cmd)
	//defer db.Close() 这里不能干掉
	if err!=nil{
		fmt.Println("get mysql failed",err)
		return nil
	}else{
		err = db.Ping()
		if err != nil {
			fmt.Println("get mysql failed",err)
		}else{
			fmt.Println("set database successfull!")
		}
	}
	return db
}

/** 获取事务句柄 */
func getSqlTx() *sql.Tx {
	if myDb == nil {
		fmt.Println("db is null")
		return nil
	}
	tx, err := myDb.Begin() /** 开启事务 */
	if err != nil {
		log.Fatal(err)
	}
	return tx
}

/** sql 参数集合，添加到切片 */
func addParamToSlice(slice *[]interface{},obj *simplejson.Json,name string,isStr bool) bool {
	var value interface{}
	var err error
	if isStr {
		value,err = obj.Get(name).String()
	}else{
		value,err = obj.Get(name).Int64()
	}
	if err != nil {
		fmt.Println(name+" addParamToSlice err")
		fmt.Println(err)
		return false
	}else {
		*slice  = append(*slice,value)
	}
	return true
}

/** select 专用,获取数据，注意顺序 */
/** 并入到 selectExe */
//func getRowData(rows *sql.Rows,struc interface{}) []interface{} {
//	defer rows.Close()
//	result := make([]interface{}, 0)
//	s := reflect.ValueOf(struc).Elem()
//	leng := s.NumField()
//	oneRow := make([]interface{}, leng)
//	/** 绑定各个位置 */
//	for i := 0; i < leng; i++ {
//		oneRow[i] = s.Field(i).Addr().Interface()
//	}
//	for rows.Next() {
//		err := rows.Scan(oneRow...)
//		if err != nil {
//			panic(err)
//		}
//		result = append(result, s.Interface())
//	}
//	return result
//}

func insertExe(sqlStr string,args ...interface{}) (int64,int64)  {
	/** 插入成功必须至少影响一行 */
	return commonExeSqlHandler(sqlStr,args...)
}

/** 允许影响条数是 0 */
func updateExe(sqlStr string,args ...interface{}) (int64,int64)  {
	return commonExeSqlHandler(sqlStr,args...) /** >=0 */
}

func deleteExe(sqlStr string,args ...interface{}) (int64,int64)  {
	return commonExeSqlHandler(sqlStr,args...) /** >=0 */
}

/** 由于 row 不能跨方法传递，会引起资源抢夺，所以废弃这个函数
  * 2017-2-13
  */
//func selectExe(sqlStr string,args ...interface{}) (*sql.Rows,error){
//	thisTx := getSqlTx()
//	defer thisTx.Commit() /** select 直接 commit */
//	stmt, err := thisTx.Prepare(sqlStr)
//	defer stmt.Close()
//	if err != nil {
//		log.Fatal(err)
//		return nil,err
//	}
//
//	rows, err := stmt.Query(args...)
//	// defer rows.Close() /** 不关闭居然会出错，不能跨函数！！！todo find out the reason */
//	if err != nil {
//		log.Fatal(err)
//		return nil,err
//	}
//	return rows,nil
//}

/** 完美 */
func selectExe(
sqlStr string,
getSlicePart func(i interface{}) interface{} ,
sctType interface{},
args ...interface{}) ([]interface{},error){
	thisTx := getSqlTx()
	defer thisTx.Commit()
	stmt, err := thisTx.Prepare(sqlStr)
	defer stmt.Close()
	if err != nil {
		log.Fatal(err)
		return nil,err
	}
	var rows *sql.Rows
	if len(args) <=0 {
		rows, err = stmt.Query()
	}else{
		if args[0]==nil {
			rows, err = stmt.Query()
		}else{
			rows, err = stmt.Query(args...)
		}
	}
	if err != nil {
		log.Fatal(err)
		return nil,err
	}
	result := make([]interface{}, 0)
	s := reflect.ValueOf(sctType).Elem()
	length := s.NumField()    /** 以父结构体总的变量数为界，匿名字段当做一 */
	if length<=0 {
		Log("select 的输出结构体不能为 null ！")
		return nil,nil
	}
	cols,_ := rows.Columns()
	lenCols := len(cols)
	if length!=lenCols {
		Log("select 的输出结构体字段数量 != sql语句select 出来的 ！")
		return nil,nil
	}
	oneRow := make([]interface{},lenCols) /** 以 select 出来的字段数为界 */
	index := 0
	/** 绑定各个位置 */
	for i := 0; i < length; i++ {
		field := s.Field(i)
		switch (field.Interface()).(type)  {
		case int:
			oneRow[index] = field.Addr().Interface()
			index++
			break;
		case string:
			oneRow[index] = field.Addr().Interface()
			index++
			break;
		case int64:
			oneRow[index] = field.Addr().Interface()
			index++
			break;
		case *string:
			/** 字符串指针 */
			oneRow[index] = field.Addr().Interface()
			index++
			break;
		default:
			/** 匿名字段，结构体 */
			/** 多重嵌套是不行的 */
			// todo make it work
			structTypeElem := reflect.ValueOf(field.Addr().Interface()).Elem()

			lenOfStruct := structTypeElem.NumField()
			j:=0
			for ;j<lenOfStruct;j++{
				oneRow[index] = structTypeElem.Field(j).Addr().Interface()
				index++
			}
			break;
		}
	}
	defer rows.Close()
	for rows.Next() {

		err := rows.Scan(oneRow...)
		if err != nil {
			panic(err)
		}
		result = append(result, getSlicePart(s.Interface()))//s.Interface())
	}

	return result,nil
}

/** 一些公共的操作 */
func commonExeSqlHandler(sqlStr string,args ...interface{}) (int64,int64) {
	thisTx := getSqlTx()
	stmt,err := thisTx.Prepare(sqlStr)
	defer stmt.Close()
	if err!=nil {
		defer thisTx.Rollback()
		log.Fatal(err)
		return -1,0
	}
	result,err := stmt.Exec(args...)
	if err!=nil {
		defer thisTx.Rollback()
		log.Fatal(err)
		return -2,0
	}

	effectRows,err := result.RowsAffected()
	if err!=nil {
		defer thisTx.Rollback()
		log.Fatal(err)
		return -3,0
	}
	lastId,err := result.LastInsertId()
	if err!=nil {
		defer thisTx.Rollback()
		log.Fatal(err)
		return -4,0
	}
	defer thisTx.Commit()
	return effectRows,lastId
}