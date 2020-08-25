package MicroserviceOutputParameterRead

import (
	"gw/Parameters"
	"bytes"
	"encoding/json"
	"regexp"
	"fmt"
)

type headerRQ struct {
	//======Header
	If  int    `json:"if"`
	Dri int    `json:"dri"`

}

type bodyRQ struct {
	//======Body
	Mis []int	`json:"mis"`
}

type headerRS struct {
	//======Header
	Dri int `json:"dri"`
	Rsc int `json:"rsc"`
}

type bodyRS struct {
	//======Body
	Op string `json:"op"`
}

//Request
func Request(parameters *Parameters.Parameter) {

	tempH := headerRQ{parameters.InterfaceID(),parameters.DisposableIoTRequestID()}
	h, _ := json.Marshal(tempH)
	//	h, _ := json.MarshalIndent(tempH, "", " ")
	k := regexp.MustCompile(`(\s*?{\s*?|\s*?,\s*?)(['"])?([a-zA-Z0-9]+)(['"])?:`).ReplaceAll(h, []byte("$1$3="))
	k = bytes.ReplaceAll(k, []byte(","), []byte(";"))

	tempB := bodyRQ{parameters.MicroserviceIDs()}
	b, _ := json.Marshal(tempB)
	//	b, _ := json.MarshalIndent(tempB, "", " ")
	j := regexp.MustCompile(`(\s*?{\s*?|\s*?,\s*?)(['"])?([a-zA-Z0-9]+)(['"])?:`).ReplaceAll(b, []byte("$1$3="))

	kj := string(k) + string(j)
	result := bytes.ReplaceAll([]byte(kj), []byte("\""), []byte(""))

	//fmt.Println(string(result))
	fmt.Println("[REQ] Send to Device =>>>" , string(result))
	RQparsing([]byte(result),parameters)

}

//Response
func Response(parameters *Parameters.Parameter) {

	tempH := headerRS{parameters.DisposableIoTRequestID(),parameters.ResponseStatusCode()}
	h, _ := json.Marshal(tempH)
	k := regexp.MustCompile(`(\s*?{\s*?|\s*?,\s*?)(['"])?([a-zA-Z0-9]+)(['"])?:`).ReplaceAll(h, []byte("$1$3="))
	k = bytes.ReplaceAll(k, []byte(","), []byte(";"))

	tempB := bodyRS{parameters.OutputParameter()}
	b, _ := json.Marshal(tempB)
	j := regexp.MustCompile(`(\s*?{\s*?|\s*?,\s*?)(['"])?([a-zA-Z0-9]+)(['"])?:`).ReplaceAll(b, []byte("$1$3="))
	j = bytes.ReplaceAll(j, []byte(`\"`), []byte(""))
	j = bytes.ReplaceAll(j, []byte(`"`), []byte(""))
	j = bytes.ReplaceAll(j, []byte(`\n`), []byte(""))
	j = bytes.ReplaceAll(j, []byte(`:`), []byte("="))
	//fmt.Println(string(j))
	result := string(k) + string(j)
	//fmt.Println(result)
	fmt.Println("[RES] Receive from Servce=>>>" , result)

	RSparsing([]byte(result),parameters)
}




func RQparsing(data []byte,parameters *Parameters.Parameter)  {



	temp := bytes.SplitAfterN(data, []byte("}"), 2)  //헤더 , 바디 분리 , n은 2개로 "}" 나뉜다
	header := temp[0]
	body := temp [1]

	//header
	htemp := bytes.TrimSpace(bytes.ReplaceAll(bytes.ReplaceAll(header,[]byte("="),[]byte(":")),[]byte(";"),[]byte(",")))
	var re = regexp.MustCompile(`(\s*?{\s*?|\s*?\,\s*?)(['"])?([a-zA-Z]+)(['"])?`) //string array
	s := re.ReplaceAll(htemp,[]byte("$1\"$3\""))
	//fmt.Println(string(s))
	ab := regexp.MustCompile(`(\s*?:\s*?|\s*?}\s*?)(['"])?([a-zA-Z0-9]+)(['"])?}`).ReplaceAll(s,[]byte("$1\"$3\"}"))
	//fmt.Println(string(ab))
	harray := headerRQ{}
	e := json.Unmarshal(s, &harray)

	if e != nil {
		panic(e)
	}

	parameters.SetInterfaceID(harray.If)
	parameters.SetDisposableIoTRequestID(harray.Dri)
	fmt.Printf("\tParsing Completed: 'if'=> %d, 'dri'=> %d \r\n",parameters.InterfaceID(), parameters.DisposableIoTRequestID())


	//body
	btemp := bytes.TrimSpace(bytes.ReplaceAll(bytes.ReplaceAll(body,[]byte("="),[]byte(":")),[]byte(";"),[]byte(",")))
	s = re.ReplaceAll(btemp,[]byte("$1\"$3\""))

	var stringArray = regexp.MustCompile(`(\s*?\[\s*?|\s*?\]\s*?)(['"])?([a-zA-Z]+)(['"])?`) //string array
	s = stringArray.ReplaceAll(btemp,[]byte("$1\"$3\""))
	ab = regexp.MustCompile(`(\s*?{\s*?|\s*?,\s*?)(['"])?([a-zA-Z0-9]+)(['"])?:`).ReplaceAll(s,[]byte("$1\"$3\":"))
	//	fmt.Println(string(ab))

	barray := bodyRQ{}
	e = json.Unmarshal(ab, &barray)
	if e != nil {
		panic(e)
	}

	parameters.SetMicroserviceIDs(barray.Mis)
	fmt.Printf("\tParsing Completed: 'mis'=> %v \r\n",parameters.MicroserviceIDs())

}


func RSparsing(data []byte,parameters *Parameters.Parameter)  {



	temp := bytes.SplitAfterN(data, []byte("}"), 2)  //헤더 , 바디 분리 , n은 2개로 "}" 나뉜다
	header := temp[0]
	body := temp[1]

	//header
	htemp := bytes.TrimSpace(bytes.ReplaceAll(bytes.ReplaceAll(header,[]byte("="),[]byte(":")),[]byte(";"),[]byte(",")))
	var re = regexp.MustCompile(`(\s*?\[\s*?|\s*?\]\s*?)(['"])?([a-zA-Z]+)(['"])?`) //string array
	s := re.ReplaceAll(htemp,[]byte("$1\"$3\""))
	ab := regexp.MustCompile(`(\s*?{\s*?|\s*?,\s*?)(['"])?([a-zA-Z0-9]+)(['"])?:`).ReplaceAll(s,[]byte("$1\"$3\":"))
	harray := headerRS{}
	e := json.Unmarshal(ab, &harray)

	if e != nil {
		panic(e)
	}

	parameters.SetDisposableIoTRequestID(harray.Dri)
	parameters.SetResponseStatusCode(harray.Rsc)
	fmt.Printf("\tParsing Completed: 'dri'=> %d, 'rsc'=> %d \r\n",parameters.DisposableIoTRequestID(), parameters.ResponseStatusCode())


	//body 하드코딩
	btemp := bytes.TrimSpace(bytes.ReplaceAll(bytes.ReplaceAll(body,[]byte("="),[]byte(":")),[]byte(";"),[]byte(",")))
	btemp = bytes.ReplaceAll(btemp,[]byte(":{"),[]byte(":\""))
	btemp = bytes.ReplaceAll(btemp,[]byte("}}"),[]byte("\"}"))
	btemp = bytes.ReplaceAll(btemp,[]byte("0,"),[]byte("?"))
	s = re.ReplaceAll(btemp,[]byte("$1\"$3\""))
	ab = regexp.MustCompile(`(\s*?{\s*?|\s*?,\s*?)(['"])?([a-zA-Z0-9]+)(['"])?:`).ReplaceAll(s,[]byte("$1\"$3\":"))
	ab = bytes.ReplaceAll(ab,[]byte("?"),[]byte("0,"))
	fmt.Println(string(ab))

	barray := bodyRS{}
	e = json.Unmarshal(ab, &barray)
	if e != nil {
		panic(e)
	}
	parameters.SetOutputParameter(barray.Op)
	fmt.Printf("\tParsing Completed: 'op'=> %s \r\n",parameters.OutputParameter())

}