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
	"net/http"
	"fmt"
	"github.com/bitly/go-simplejson"
)

type LghRequest struct {
	w http.ResponseWriter
	r *http.Request

	// 标记使用，当前的方法名称
	funcName string

	// 输入的结构体，与客户端输入的 json 成对应关系
	inputStruct  interface{}

	// 自定义 slices 的回调，方便你做参数处理，返回 true 意味着此次操作终止，例如 update
	slicesCallBack func(slices []interface{}) bool

	// 根据传入的 jsonObj 生成的 slices 来回调，方法生成自定义 sql
	getSqlCallBack func(slices []interface{},inputStruct interface{}) string
}

func getJsonObjByRequest(request LghRequest) (*simplejson.Json,int64){
	if request.r == nil {
		Log(request.funcName+" 输入 http.req 不能为 null")
		return nil,-1
	}
	obj := commonReqBody(request.r)
	if obj==nil && request.inputStruct != nil{
		Log(request.funcName+" 输入 jsonObj 为 null，不能与 inputStruct 对应")
		errorRet(request.w)
		return nil,-2
	}else if obj==nil && request.inputStruct == nil {
		return nil,1
	}else if obj!=nil && request.inputStruct == nil{
		Log(request.funcName+" 输入 inputStruct 为 null，不能与 jsonObj 对应")
		errorRet(request.w)
		return nil,-3
	}
	return obj,2
}

func insertDataByStruct(request LghRequest){

	obj,ret := getJsonObjByRequest(request);
	if ret<=0 {
		return
	}
	slices,sqlStr := exeSqlCommonHandler(
		request.w,
		obj,
		request.funcName,
		request.inputStruct,
		request.slicesCallBack,
		request.getSqlCallBack);

	if slices == nil {
		return
	}
	var lastId int64
	if len(slices) <=0 {
		_,lastId = insertExe(sqlStr)
	}else{
		_,lastId = insertExe(sqlStr,slices...)
	}
	echoResult(request.w,lastId)
}

func selectDataByStruct(request LghRequest,outputStruct interface{}){

	obj,ret := getJsonObjByRequest(request);
	if ret <= 0 {
		return
	}
	slices,sqlStr := exeSqlCommonHandler(
		request.w,
		obj,
		request.funcName,
		request.inputStruct,
		request.slicesCallBack,
		request.getSqlCallBack);

	if slices == nil {
		return
	}
	var dataSlice []interface{}
	var err error
	if len(slices) <=0 {
		dataSlice,err = selectExe(sqlStr,func (i interface{}) interface{} {
			return i /** 要处理的情况，再独立出这个函数 */
		},outputStruct,nil)
	}else{
		dataSlice,err = selectExe(sqlStr,func (i interface{}) interface{} {
			return i
		},outputStruct,slices...)
	}
	if dataSlice == nil {
		return
	}
	if err!=nil {
		fmt.Fprintf(request.w,"null");
		Log(err)
		return
	}
	getJsonData(request.w,dataSlice)
}

func updateDataByStruct(request LghRequest)  {

	obj,ret := getJsonObjByRequest(request);
	if ret<=0 {
		Log("stop on line 136")
		return
	}
	slices,sqlStr := exeSqlCommonHandler(
		request.w,
		obj,
		request.funcName,
		request.inputStruct,
		request.slicesCallBack,
		request.getSqlCallBack);

	if slices == nil {
		return
	}
	var rowNum int64
	if len(slices) <=0 {
		rowNum,_ = updateExe(sqlStr)
	}else{
		Log(slices)
		rowNum,_ = updateExe(sqlStr,slices...)
	}
	echoResult(request.w,rowNum)
}

func deleteDataByStruct(request LghRequest)  {

	obj,ret := getJsonObjByRequest(request);
	if ret<=0 {
		return
	}
	slices,sqlStr := exeSqlCommonHandler(
		request.w,
		obj,
		request.funcName,
		request.inputStruct,
		request.slicesCallBack,
		request.getSqlCallBack);
	if slices == nil {
		return
	}
	var rowNum int64
	if len(slices) <=0 {
		rowNum,_ = deleteExe(sqlStr)
	}else{
		rowNum,_ = deleteExe(sqlStr,slices...)
	}
	echoResult(request.w,rowNum)
}








