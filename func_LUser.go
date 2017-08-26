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
r.HandleFunc("/insert_luser",insert_luser)
r.HandleFunc("/delete_luser",delete_luser)
r.HandleFunc("/update_luser",update_luser)
r.HandleFunc("/select_luser",select_luser)
*/

func insert_luser(w http.ResponseWriter,r *http.Request)  {
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

func delete_luser(w http.ResponseWriter,r *http.Request)  {
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

func update_luser(w http.ResponseWriter,r *http.Request)  {
	type testS struct {
		Id int64 `json:"id" nullTag:"1"`
	}
	request := LghRequest{
		w,
		r,
		"update_luser",
		new (testS),
		func(slices []interface{}) bool{
			// 演示使用输入参数的情况
			return false
		},
		func(slices []interface{},inputStruct interface{}) string {
			return "update LUser set u_user_id='444' where id=?"
		}}
	updateDataByStruct(request)
}

func select_luser(w http.ResponseWriter,r *http.Request)  {
	request := LghRequest{
		w,
		r,
		"select_luser",
		nil,
		func(slices []interface{}) bool{
			// 演示使用输入参数的情况
			return false
		},
		func(slices []interface{},inputStruct interface{}) string {
			return "select * from LUser"
		}}
		output := new (LUser)
	selectDataByStruct(request,output)
}

