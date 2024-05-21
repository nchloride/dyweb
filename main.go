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
	"net/url"
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
	UrlDecode bool `json:"urlDecode"`
}	



var (
	jsonPath string
	port int
)


func main(){

	

	
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
		urlDecode := ep[i].UrlDecode
		http.HandleFunc(pathName,func(w http.ResponseWriter, r *http.Request){

			fmt.Println("+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++")
			fmt.Println(fmt.Sprintf("\n \n Endpoint '%v' invoked \n Method used: %v",pathName,r.Method))

			if r.Method == "OPTIONS"{
				w.Header().Add("Access-Control-Allow-Origin", "*")
				w.Header().Add("Access-Control-Allow-Methods", "GET, POST")
				w.Header().Add("Access-Control-Allow-Headers", "*")

			}

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
					fmt.Println(postBody)

				}
			}
			if len(headers) > 0  {
				Add := w.Header().Add
				AddHeaders(headers,&Add)
			}
			fmt.Println(fmt.Sprintf("Request from: %v", r.RemoteAddr))
			
			if urlDecode {
				unescapedContent,err := url.QueryUnescape(content)
				if err != nil {
					fmt.Println(err)
				}
				content = unescapedContent
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


func init(){
	
	flag.Usage = func() {
		
		h:= "Creates webserver with just a command  \n\n"

		h+="Usage: \n\n\n"

		h+="dyweb -p <port_number.  -json <json_config_path>\n\n"

		h+= "JSON Format: \n \n"

		h+= "'pathName' endpoint name <string> \n"
		h+= "'method' endpoint method <string> \n"
		h+= "'content' response <string>\n"
		h+= "'headers' contains an array of 'header' \n"
		h+= "'header' contains an array that should only have 2 string values, key and value"

		fmt.Fprintf(os.Stderr, h)
		
	}

	
	flag.StringVar(&jsonPath, "json","", "no json file declared")
	flag.IntVar(&port,"p",8000, "default port")
	
	flag.Parse()
}
