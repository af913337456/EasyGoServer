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
r.HandleFunc("/insert_lcomment",insert_lcomment)
r.HandleFunc("/delete_lcomment",delete_lcomment)
r.HandleFunc("/update_lcomment",update_lcomment)
r.HandleFunc("/select_lcomment",select_lcomment)
*/

func insert_lcomment(w http.ResponseWriter,r *http.Request)  {
	request := LghRequest{
		w,
		r,
		"insert_lcomment",
		new (LComment),
		func(slices []interface{}) bool{
			return false
		},
		func(slices []interface{},inputStruct interface{}) string {
			return buildInsertSqlByStruct(new (LComment),"LComment")
		}}
	insertDataByStruct(request)
}

func delete_lcomment(w http.ResponseWriter,r *http.Request)  {
	request := LghRequest{
		w,
		r,
		"delete_lcomment",
		nil,
		func(slices []interface{}) bool{
			return false
		},
		func(slices []interface{},inputStruct interface{}) string {
			return "delete from LComment where id='2'"
		}}
	deleteDataByStruct(request)
}

func update_lcomment(w http.ResponseWriter,r *http.Request)  {
	type testS struct {
		Id int64 `json:"id" nullTag:"1"`
	}
	request := LghRequest{
		w,
		r,
		"update_lcomment",
		new (testS),
		func(slices []interface{}) bool{
			// 演示使用输入参数的情况
			return false
		},
		func(slices []interface{},inputStruct interface{}) string {
			return "update LComment set u_user_id='444' where id=?"
		}}
	updateDataByStruct(request)
}

func select_lcomment(w http.ResponseWriter,r *http.Request)  {
	request := LghRequest{
		w,
		r,
		"select_lcomment",
		nil,
		func(slices []interface{}) bool{
			// 演示使用输入参数的情况
			return false
		},
		func(slices []interface{},inputStruct interface{}) string {
			return "select * from LComment"
		}}
		output := new (LComment)
	selectDataByStruct(request,output)
}

