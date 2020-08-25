package DeviceMicroserviceInformationReport

import (
	"gw/ObjectTypeParameters"
	"gw/Parameters"
	"bytes"
	"encoding/json"
	"fmt"
	"regexp"
)

type headerRQ struct {
	//======Header
	If  int    `json:"if"`
	Dri int    `json:"dri"`
	Di  string `json:"di"`
}

type bodyRQ struct {
	//======Body
	Mif []ObjectTypeParameters.MifObject `json:"mif"`
}

type headerRS struct {
	//======Header
	Dri int `json:"dri"`
	Rsc int `json:"rsc"`
}
//Request
func Request(parameters *Parameters.Parameter) {

	tempH := headerRQ{parameters.InterfaceID(), parameters.DisposableIoTRequestID(), parameters.DeviceID()}
	h, _ := json.Marshal(tempH)
	//	h, _ := json.MarshalIndent(tempH, "", " ")
	k := regexp.MustCompile(`(\s*?{\s*?|\s*?,\s*?)(['"])?([a-zA-Z0-9_-]+)(['"])?:`).ReplaceAll(h, []byte("$1$3="))
	k = bytes.ReplaceAll(k, []byte(","), []byte(";"))

	tempB := bodyRQ{parameters.MicroserviceInformation()}
	b, _ := json.Marshal(tempB)
	//	b, _ := json.MarshalIndent(tempB, "", " ")
	j := regexp.MustCompile(`(\s*?{\s*?|\s*?,\s*?)(['"])?([a-zA-Z0-9_-]+)(['"])?:`).ReplaceAll(b, []byte("$1$3="))

	kj := string(k) + string(j)
	result := bytes.ReplaceAll([]byte(kj), []byte("\""), []byte(""))

	//fmt.Println(string(result))
	fmt.Println("[REQ] Send to Server =>>>" , string(result))

	RQparsing([]byte(result),parameters)
}

//Response
func Response(parameters *Parameters.Parameter) {

	tempH := headerRS{parameters.DisposableIoTRequestID(), parameters.ResponseStatusCode()}
	h, _ := json.Marshal(tempH)
	k := regexp.MustCompile(`(\s*?{\s*?|\s*?,\s*?)(['"])?([a-zA-Z0-9_-]+)(['"])?:`).ReplaceAll(h, []byte("$1$3="))
	result := bytes.ReplaceAll(k, []byte(","), []byte(";"))
	//fmt.Println(string(result))

	fmt.Println("[RES] Receive from Device=>>>" , string(result))

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
	ab := regexp.MustCompile(`(\s*?:\s*?|\s*?}\s*?)(['"])?([a-zA-Z0-9_-]+)(['"])?}`).ReplaceAll(s,[]byte("$1\"$3\"}"))
	//fmt.Println(string(ab))
	harray := headerRQ{}
	e := json.Unmarshal(ab, &harray)

	if e != nil {
		panic(e)
	}

	parameters.SetInterfaceID(harray.If)
	parameters.SetDisposableIoTRequestID(harray.Dri)
	parameters.SetDeviceID(harray.Di)
	fmt.Printf("\tParsing Completed: 'if'=> %d, 'dri'=> %d, 'di'=> %s \r\n",parameters.InterfaceID(), parameters.DisposableIoTRequestID(),parameters.DeviceID())


	//body
	btemp := bytes.TrimSpace(bytes.ReplaceAll(bytes.ReplaceAll(body,[]byte("="),[]byte(":")),[]byte(";"),[]byte(",")))
	s = re.ReplaceAll(btemp,[]byte("$1\"$3\""))

	var stringArray = regexp.MustCompile(`(\s*?\[\s*?|\s*?\]\s*?)(['"])?([a-zA-Z]+)(['"])?`) //string array
	s = stringArray.ReplaceAll(btemp,[]byte("$1\"$3\""))
	ab = regexp.MustCompile(`(\s*?{\s*?|\s*?,\s*?)(['"])?([a-zA-Z0-9_-]+)(['"])?:`).ReplaceAll(s,[]byte("$1\"$3\":"))
//	fmt.Println(string(ab))

	barray := bodyRQ{}
	e = json.Unmarshal(ab, &barray)
	if e != nil {
		panic(e)
	}

	for _, v := range barray.Mif{
		fmt.Printf("\tmi: %d, ops: %s \r\n", v.Mi , v.Ops)
	}

	parameters.SetMicroserviceInformation(barray.Mif)
	fmt.Printf("\tParsing Completed: 'mif'=> %v \r\n",parameters.MicroserviceInformation())

}


func RSparsing(data []byte,parameters *Parameters.Parameter)  {



	temp := bytes.SplitAfterN(data, []byte("}"), 2)  //헤더 , 바디 분리 , n은 2개로 "}" 나뉜다
	header := temp[0]


	//header
	htemp := bytes.TrimSpace(bytes.ReplaceAll(bytes.ReplaceAll(header,[]byte("="),[]byte(":")),[]byte(";"),[]byte(",")))
	var re = regexp.MustCompile(`(\s*?\[\s*?|\s*?\]\s*?)(['"])?([a-zA-Z]+)(['"])?`) //string array
	s := re.ReplaceAll(htemp,[]byte("$1\"$3\""))
	ab := regexp.MustCompile(`(\s*?{\s*?|\s*?,\s*?)(['"])?([a-zA-Z0-9_-]+)(['"])?:`).ReplaceAll(s,[]byte("$1\"$3\":"))
	harray := headerRS{}
	e := json.Unmarshal(ab, &harray)

	if e != nil {
		panic(e)
	}

	parameters.SetDisposableIoTRequestID(harray.Dri)
	parameters.SetResponseStatusCode(harray.Rsc)
	fmt.Printf("\tParsing Completed: 'dri'=> %d, 'rsc'=> %d \r\n",parameters.DisposableIoTRequestID(), parameters.ResponseStatusCode())




}
