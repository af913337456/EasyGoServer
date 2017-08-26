package test

import "net/http"

func selectS(w http.ResponseWriter,r *http.Request)  {
	request := LghRequest{

	}
	output := struct {

	}{}
	selectDataByStruct(request,output)
}

func insertS(w http.ResponseWriter,r *http.Request)  {
	request := LghRequest{

	}
	insertDataByStruct(request)
}

func updateS(w http.ResponseWriter,r *http.Request)  {
	request := LghRequest{

	}
	updateDataByStruct(request)
}

func deleteS(w http.ResponseWriter,r *http.Request)  {
	request := LghRequest{

	}
	deleteDataByStruct(request)
}

func main()  {

}
