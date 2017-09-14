package main

import (
	"net/http"
	"os/exec"
	"fmt"
	"flag"
	"os"
)

type my struct {
	fuck func()
}

func (h *my) ServeHTTP(w http.ResponseWriter,r *http.Request) {
	h.fuck()
}

var cmd *exec.Cmd

func main() {

	arg := flag.String(
		"arg",
		"",
		"更新书使用 --arg=\"update\" \n 第一次上传 --arg=\"upload\"")

	flag.Parse()

	if arg==nil{
		fmt.Println("命令参数不能为 null,输入: -help 查看帮助")
		return
	}
	if *arg==""{
		fmt.Println("命令参数不能为 null,输入: -help 查看帮助")
		return
	}

	create_gitignore()

	if *arg=="upload" {
		fmt.Println("第一次上传")
		// 删除一次 .git
		cmd = exec.Command("rm","-rf", ".git")
		err := cmd.Run()
		if err != nil {
			fmt.Println(err)
		}
		if git_init() {
			if git_add_all() {
				if git_commit("第一次提交") {
					if git_remote() {
						if git_push() {
							echoJson("success","成功")
						}
					}
				}
			}
		}
	}

	if *arg=="update" {
		fmt.Println("更新书")
		if git_add_all() {
			if git_commit("更新提交") {
				if git_push() {
					echoJson("success","成功")
				}
			}
		}
	}
}

func git_init() bool {
	cmd = exec.Command("git", "init")
	return exeCmd(func(err error) {
		fmt.Println(err)
		echoFailed("git_init "+err.Error())
	}, func() {
		fmt.Println("git init 成功")
	})
}

func git_add_all() bool {
	cmd = exec.Command("git", "add",".")
	return exeCmd(func(err error) {
		fmt.Println(err)
		echoFailed("git_add_all "+err.Error())
	}, func() {
		fmt.Println("git add . 成功")
	})
}

func git_commit(str string) bool {
	cmd = exec.Command("git", "commit","-m",str)

	return exeCmd(func(err error) {
		fmt.Println(err)
		echoFailed("git_commit "+err.Error())
	}, func() {
		fmt.Println("git commit -m \"...\" 成功")
	})
}

const (
	user   = "lgh"
	pw	   = "123aaa"
	name   = "lgh"
	git	   = "test3"
	footer = ".git"
	host   = "localhost"
	port   = "3000"
	url    = "http://"+user+":"+pw+"@"+host+":"+port+"/"+name+"/"+git+footer
)

func git_remote() bool {
	cmd = exec.Command("git", "remote","add","origin",url)
	return exeCmd(func(err error) {
		fmt.Println(err.Error()+" ---> "+url)
		echoFailed("git_remote "+err.Error())
	}, func() {
		fmt.Println("git remote add origin \""+url+"\" 成功")
	})
}

func git_push() bool {
	cmd = exec.Command("git", "push","-u","origin","master")
	return exeCmd(func(err error) {
		fmt.Println(err)
		echoFailed("git_push "+err.Error())
	}, func() {
		fmt.Println("git push -u origin master \"...\" 成功")
	})
}

func exeCmd(failed func(err error),success func()) bool {
	err := cmd.Run()
	if err != nil {
		failed(err)
		return false
	}
	success()
	return true
}

func echoFailed(msg string)  {
	echoJson("failed",msg)
}

func echoJson(status,msg string)  {
	fmt.Println("{\"ret\":\""+status+"\",\"msg\":\""+msg+"\"}")
}

func create_gitignore()  {
	fileName := ".gitignore"
	_,err := os.Stat(fileName)
	if err == nil {
		fmt.Println("文件：.gitignore 已经存在")
		return
	}
	file,err := os.Create(fileName)
	defer file.Close()
	if err!=nil {
		fmt.Println(err)
		return
	}
	file.WriteString("lgh_git_util")
}
