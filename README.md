# EasyGoServer

> 作者：林冠宏 / 指尖下的幽灵

> 掘金：https://juejin.im/user/587f0dfe128fe100570ce2d8

> 博客：http://www.cnblogs.com/linguanh/

> GitHub ： https://github.com/af913337456/

> 联系方式 / Contact：913337456@qq.com

----------

[TOC]

----- 概述

----- 脚本介绍

--------- Linux

--------- Windows

--------- Mac

----- 使用流程

----- 部分代码说明

----- TODO


### 概述
一个能够仅仅依赖你创建好的 sql 文件，就能 ``自动帮你生成基础服务端框架代码`` 的 go server 框架。包含有:

1，基础的 增删改查 
2，拓展性强的API
3，客户端的数据传入 与 服务端的输出 全部依赖 struct
* 例如你的一个输入结构体 ``inputStruct`` 设置为
```go
 type inputStruct struct {
     Id   int64   `json:"id"   nullTag:"1"` // nullTag==1 指明 id 必须要求在客户端传入 {"id":123}
     Name string  `json:"name" nullTag:"0"` // ==0 指明 name 在客户端输入的时候可以不必要
 }
```
对应上例，``客户端输入的 json :`` {"id":666, "name":"lgh"}

* 当你在使用 select 的时候，你的 sql 如果是这样的：`` select User.id , User.age from User``
那么你的对应输出结构体 ``outputStruct`` 应该是:
```go
 type inputStruct struct {
	 Id   int64  `json:"id"`   
     Age  int64  `json:"age"` 
 }
```

4，真正需要你写的代码极少，例如第三点的例子，你要写的就那么多，其中默认的 struct 会自动帮你生成

![](http://images2017.cnblogs.com/blog/690927/201708/690927-20170826145343589-2112195908.png)

---

### 脚本介绍

---
根据 sql 文件，自动生成代码文件，包含有 struct.go，每张表对应生成一个包含有增删改查的基础方法文件
``one_key_create_code ``

根据内置的 makefile 或者 .bat 编译并运行默认的 go server 程序，注意是默认的
``make_server  ``

#### Linux
one_key_create_code.sh  
make_server.sh
Makefile

#### Windows
one_key_create_code.bat
make_server.bat

#### Mac
参照 linux 的

### 使用流程

1，在你的 服务器 安装 mysql 或者 mariadb

2，编写好的你的 sql 文件，可以参照我源码里面的 this.sql 

3，运行步骤2编写好的 sql 文件

4，修改 sql_2_api.go 里面 main 内的 sql 文件名称

5，运行 one_key_create_code 脚本，成功后会在同级目录生成下面文件，记得刷新目录

* ``struct.go``，里面包含注释规范
* 对应你 sql 文件里面的表名称生成的函数文件，格式：`` func_表名称.go``

6，自己写好，main.go 或者 使用我提供的默认 ``LghSampleMain.go``，在里面 添加你自己的路由
```go
	router.HandleFunc("/insert",insert_luser_sample).Methods("POST")
	router.HandleFunc("/select",select_luser_sample).Methods("GET")
	router.HandleFunc("/update",update_luser_sample).Methods("POST")
	router.HandleFunc("/delete",delete_luser_sample).Methods("POST")
```
7，配置好 conf.json 文件，我里面有例子
```json
// Host 是绝对路径
// Port 是要被监听的端口
{
  "Host": "127.0.0.1",
  "Port": ":8884",
  "FilePort":":8885",
  "DbName":"database",
  "DbUser":"root",
  "DbPw":"123456",
  "DbPort":"3306"
}
```

8，现在执行 make_server 脚本，观察控制台的输出，即可。

### 部分代码说明

核心的参数结构体
```go
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

```
例子方法

1，演示不需要参数的形式
```go
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
```

2，演示当有参数输入的时候，参数仅做判断，但是不需要组合到 sql 的情况
```go
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
```

3，演示使用输入参数的情况
```go
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
```

### TODO
寻找自愿者帮忙翻译英文版文档

