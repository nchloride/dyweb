package main 

import (
	"fmt"
	"flag"
	"os"
	"encoding/json"
	"strings"
	"net/http"
	"io/ioutil"
	"bytes"
	"io"
)	


	

type EndpointHeaders struct {
	Header []string 
}

type Endpoint struct{
	Method string `json:"method"`
	Param []string `json:"params"`
	Name string `json:"name"`
	Status int `json:"status"`
	Content string `json:"content"`
	Headers []EndpointHeaders `json:"headers"`

}	




func main(){

	var (
		jsonPath string
		port int
	)
	
	flag.StringVar(&jsonPath, "json","", "no json file declared")
	flag.IntVar(&port,"p",8000, "default port")
	
	flag.Parse()
	

	
	fmt.Println(fmt.Sprintf("running on port:%v",port))
	file,err := os.ReadFile(jsonPath)

	fmt.Println(string(file))
	if err != nil {
		fmt.Println(err)
	}
	var ep []Endpoint
	json.NewDecoder(strings.NewReader(string(file))).Decode(&ep)
	//fmt.Println(ep)

	
	for i:=0; i<len(ep); i++{
		fmt.Println(ep[i])
		pathName := ep[i].Name
		method := ep[i].Method
		status := ep[i].Status
		content := ep[i].Content
		headers := ep[i].Headers
		http.HandleFunc(pathName,func(w http.ResponseWriter, r *http.Request){
			fmt.Println(fmt.Sprintf("Endpoint '%v' invoked \n Method used: %v",pathName,method))
			if method != r.Method {
				message := fmt.Sprintf("Not a %v endpoint", r.Method)
				
				w.WriteHeader(404)
				w.Write([]byte(message))
				return
			}
			if method == "POST" {
				if r.Body != nil{
					data,err:=ioutil.ReadAll(r.Body)
					if err != nil {
						fmt.Println(err)
					}
					body := strings.NewReader(string(data))
					
					var buffer bytes.Buffer
					var postBody map[string]interface{}
					io.Copy(&buffer,body)
					convertedString := buffer.String()
					json.Unmarshal([]byte(convertedString),&postBody)
					fmt.Println(postBody["test2"])


				}
			}
			if len(headers) > 0  {
				Add := w.Header().Add
				AddHeaders(headers,&Add)
//				for index:=0; index<len(headers); index++{
//					key := headers[index].Header[0]
//					value := headers[index].Header[1]
//					w.Header().Add(key,value)
//				}
			}
			
			w.WriteHeader(status)
			w.Write([]byte(content))

		})
	}
	http.ListenAndServe(fmt.Sprintf(":%v",port),nil)
}

//type AddHeadersFunc func(string,string)

func AddHeaders(headers []EndpointHeaders, Add *func(key, value string)){
	 for index:=0; index<len(headers); index++{
	 	key := headers[index].Header[0]
	 	value := headers[index].Header[1]
	 	(*Add)(key,value)
	 }
}
