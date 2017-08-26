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
	模板 main.go
*/

import (
	"github.com/gorilla/mux"
	"net/http"
	"fmt"
	"crypto/x509"
	"io/ioutil"
	"crypto/tls"
)

func test(w http.ResponseWriter,r *http.Request)  {
	fmt.Fprintf(w,"======= hello world! =======")
}

// http 监听
func httpListen(router *mux.Router)  {
	url := serverConfig.Host+serverConfig.Port
	err := http.ListenAndServe(url,router)
	if err !=nil {
		Log("http 监听错误 :")
		Log(err)
		return
	}
}

// https 监听
func httpsListen(router *mux.Router)  {
	basePath := "" // /home/lgh/

	pool := x509.NewCertPool()
	caCertPath := basePath+"" // ca.crt

	caCrt, err := ioutil.ReadFile(caCertPath)
	if err != nil {
		Log("Read ca File err:", err)
		return
	}
	pool.AppendCertsFromPEM(caCrt)
	s := &http.Server{
		Addr:    serverConfig.Host+serverConfig.Port, // :8888
		Handler: router,
		TLSConfig: &tls.Config{
			ClientCAs:  pool,
			ClientAuth: tls.RequireAndVerifyClientCert, /** 开启双向验证 */
		},
	}
	s.ListenAndServeTLS(basePath+"server.crt",basePath+"server.key")
}

func setRouter() *mux.Router {
	router := new (mux.Router)
	router.HandleFunc("/",test).Methods("POST")

	router.HandleFunc("/insert",insert_luser_sample).Methods("POST")
	router.HandleFunc("/select",select_luser_sample).Methods("GET")
	router.HandleFunc("/update",update_luser_sample).Methods("POST")
	router.HandleFunc("/delete",delete_luser_sample).Methods("POST")

	/** 在下面添加你的回调方法 */
	/** add your func below */

	return router
}

func main()  {

	bindServerConfig()

	Log("配置信息:")
	Log(serverConfig)

	myDb = setDB(serverConfig.DbUser,serverConfig.DbPw,serverConfig.DbName,serverConfig.DbPort)

	httpListen(setRouter())

}
