package main

import (
	"net/http"
)

/**
type LghRequest struct {
	w http.ResponseWriter
	r *http.Request
	funcName string
	inputStruct  interface{}
	slicesCallBack func(slices *[]interface{}) bool
	getSqlCallBack func(slices *[]interface{},inputStruct interface{}) string
}
r.HandleFunc("/insert_luser",insert_luser_sample)
r.HandleFunc("/delete_luser",delete_luser_sample)
r.HandleFunc("/update_luser",update_luser_sample)
r.HandleFunc("/select_luser",select_luser_sample)
*/

/** 演示不需要参数的形式 */
func update_0(w http.ResponseWriter,r *http.Request)  {
	request := LghRequest{
		w,
		r,
		"update_luser",
		nil, /** nil 表示没输入结构体 */
		func(slices *[]interface{}) bool{
			return false
		},
		func(slices *[]interface{},inputStruct interface{}) string {
			return "update LUser set u_user_id='444' where id='1'"
		}}
	updateDataByStruct(request)
}

/** 演示当有参数输入的时候，参数仅做判断，但是不需要组合到 sql的情况 */
func update_1(w http.ResponseWriter,r *http.Request)  {
	type testS struct {
		Id int64 `json:"id" nullTag:"1"` // nullTag==1 指明 id 必须要求在客户端传入 {"id":123}
	}
	request := LghRequest{
		w,
		r,
		"update_luser",
		new (testS),
		func(slices []interface{}) bool{
			// 在这里对 slices 做你想做的操作，增加或者删除等等
			if slices[0] == -1{
				return true /** 返回 true，终止插入，提示错误或者其它 */
			}
			slices = append(slices[:0], nil) /** 自己做完处理删除掉 */
			return false
		},
		func(slices []interface{},inputStruct interface{}) string {
			// 如果你想根据输入的 json 数据来特定生成 sql，那么就可以在这里使用 slices 来操作
			return "update LUser set u_user_id='444' where id='2'"
		}}
	updateDataByStruct(request)
}

/** 演示使用输入参数的情况 */
func update_luser_sample(w http.ResponseWriter,r *http.Request)  {
	type testS struct {
		Id int64 `json:"id" nullTag:"1"`
	}
	request := LghRequest{
		w,
		r,
		"update_luser",
		new (testS),
		func(slices []interface{}) bool{
			return false
		},
		func(slices []interface{},inputStruct interface{}) string {
			return "update LUser set u_user_id='444' where id=?" /** 对应 id */
		}}
	updateDataByStruct(request)
}

func insert_luser_sample(w http.ResponseWriter,r *http.Request)  {
	request := LghRequest{
		w,
		r,
		"insert_luser",
		new (LUser),
		func(slices []interface{}) bool{
			return false
		},
		func(slices []interface{},inputStruct interface{}) string {
			return buildInsertSqlByStruct(new (LUser),"LUser")
		}}
	insertDataByStruct(request)
}

func delete_luser_sample(w http.ResponseWriter,r *http.Request)  {
	request := LghRequest{
		w,
		r,
		"delete_luser",
		nil,
		func(slices []interface{}) bool{
			return false
		},
		func(slices []interface{},inputStruct interface{}) string {
			return "delete from LUser where id='2'"
		}}
	deleteDataByStruct(request)
}

func select_luser_sample(w http.ResponseWriter,r *http.Request)  {
	request := LghRequest{
		w,
		r,
		"select_luser",
		nil,
		func(slices []interface{}) bool{
			return false
		},
		func(slices []interface{},inputStruct interface{}) string {
			return "select * from LUser"
		}}
	output := new (LUser)
	selectDataByStruct(request,output)
}

