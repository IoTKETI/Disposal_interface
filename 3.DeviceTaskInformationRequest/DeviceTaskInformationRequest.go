package DeviceTaskInformationRequest

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
	//Tif ObjectTypeParameters.Tif `json:"tif"`
	Tif []ObjectTypeParameters.TifObject `json:"tif"`
}

//============================

type RsparsingTif struct {
	Tif []TifObject `json:"tif"`
}

type TifObject struct {   /// 이걸 만들어야함
	Ti int `json:"ti"`
	Sp SpObject `json:"sp,omitempty"`  //들어올때마다 달라짐
	Fp string `json:"fp,omitempty"`
	To bool `json:"to,omitempty"`
}

type SpObject struct {
	value string
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
	fmt.Println("[REQ] Send to Server =>>>" , string(result))

	RQparsing([]byte(result),parameters)
}

//Response
func Response(parameters *Parameters.Parameter) {

	tempH := headerRS{parameters.DisposableIoTRequestID(),parameters.ResponseStatusCode()}
	h, _ := json.Marshal(tempH)
	k := regexp.MustCompile(`(\s*?{\s*?|\s*?,\s*?)(['"])?([a-zA-Z0-9]+)(['"])?:`).ReplaceAll(h, []byte("$1$3="))
	k = bytes.ReplaceAll(k, []byte(","), []byte(";"))

	tempB := bodyRS{parameters.TaskInformation()}
	b, _ := json.Marshal(tempB)
	j := regexp.MustCompile(`(\s*?{\s*?|\s*?,\s*?)(['"])?([a-zA-Z0-9]+)(['"])?:`).ReplaceAll(b, []byte("$1$3="))
	j = bytes.ReplaceAll(j, []byte(`\"`), []byte(""))
	j = bytes.ReplaceAll(j, []byte(`"`), []byte(""))
	j = bytes.ReplaceAll(j, []byte(`\n`), []byte(""))
	j = bytes.ReplaceAll(j, []byte(`:`), []byte("="))
	//fmt.Println(string(j))
	result := string(k) + string(j)
//	fmt.Println(result)

	fmt.Println("[RES] Receive from Device=>>>" , result)
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
	fmt.Println(string(s))
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
	ab := regexp.MustCompile(`(\s*?{\s*?|\s*?,\s*?)(['"])?([a-zA-Z0-9]+)(['"])?:`).ReplaceAll(s,[]byte("$1\"$3\":"))

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


	//body  하드코딩
	btemp := bytes.TrimSpace(bytes.ReplaceAll(bytes.ReplaceAll(body,[]byte("="),[]byte(":")),[]byte(";"),[]byte(",")))
	//fmt.Println(string(btemp))
	btemp = bytes.ReplaceAll(btemp,[]byte(":{"),[]byte(":\""))
	btemp = bytes.ReplaceAll(btemp,[]byte("}}"),[]byte("\"}"))
	btemp = bytes.ReplaceAll(btemp,[]byte("0,"),[]byte("?"))
	//fmt.Println(string(btemp))
	s = re.ReplaceAll(btemp,[]byte("$1\"$3\""))
	ab = regexp.MustCompile(`(\s*?{\s*?|\s*?,\s*?)(['"])?([a-zA-Z0-9]+)(['"])?:`).ReplaceAll(s,[]byte("$1\"$3\":"))
	ab = bytes.ReplaceAll(ab,[]byte("?"),[]byte("0,"))
	//fmt.Println(string(ab))

	barray := bodyRS{}   //파싱 문제
	e = json.Unmarshal(ab, &barray)
	if e != nil {
		panic(e)
	}
	for _, v := range barray.Tif{
		fmt.Printf("\tti: %d, sp: %s, to: %v \r\n", v.Ti , v.Sp, v.To)
	}
	parameters.SetTaskInformation(barray.Tif)
	fmt.Printf("\tParsing Completed: 'tif'=> %v \r\n",parameters.TaskInformation())


}