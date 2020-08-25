package DeviceRegistration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gw/Parameters"
	"gw/EdgeXInterface"
	"regexp"
	gonanoid "github.com/matoous/go-nanoid"
	"time"
)


type headerRQ struct {
	//======Header
	If  int    `json:"if"`
	Dri int    `json:"dri"`
}

type bodyRQ struct {
	//======Body
	Mis []int `json:"mis"`
}

type headerRS struct {
	//======Header
	Dri int `json:"dri"`
	Rsc int `json:"rsc"`
}

type bodyRS struct {
	//======Body
	Di  string `json:"di"`
}


//Request
func Request(parameters *Parameters.Parameter) {

	tempH := headerRQ{parameters.InterfaceID(),parameters.DisposableIoTRequestID()}
	h, _ := json.Marshal(tempH)
	k := regexp.MustCompile(`(\s*?{\s*?|\s*?,\s*?)(['"])?([a-zA-Z0-9]+)(['"])?:`).ReplaceAll(h, []byte("$1$3="))
	k = bytes.ReplaceAll(k, []byte(","), []byte(";"))


	tempB := bodyRQ{parameters.MicroserviceIDs()}
	b, _ := json.Marshal(tempB)
	j := regexp.MustCompile(`(\s*?{\s*?|\s*?,\s*?)(['"])?([a-zA-Z0-9]+)(['"])?:`).ReplaceAll(b, []byte("$1$3="))
	//fmt.Println(string(j))
	result := string(k) + string(j)
	fmt.Println("[REQ] Send to Server =>>>" , result)


	RQparsing([]byte(result),parameters)
	//var re = regexp.MustCompile(`(\s*?\[\s*?|\s*?\]\s*?)(['"])?([a-zA-Z]+)(['"])?`) //string array
	//s := re.ReplaceAll(test,[]byte("$1\"$3\""))

}
//Response
func Response(parameters *Parameters.Parameter) {

	tempH := headerRS{parameters.DisposableIoTRequestID(),parameters.ResponseStatusCode()}
	h, _ := json.Marshal(tempH)
	k := regexp.MustCompile(`(\s*?{\s*?|\s*?,\s*?)(['"])?([a-zA-Z0-9]+)(['"])?:`).ReplaceAll(h, []byte("$1$3="))
	k = bytes.ReplaceAll(k, []byte(","), []byte(";"))

	tempB := bodyRS{parameters.DeviceID()}
	b, _ := json.Marshal(tempB)
	j := regexp.MustCompile(`(\s*?{\s*?|\s*?,\s*?)(['"])?([a-zA-Z0-9]+)(['"])?:`).ReplaceAll(b, []byte("$1$3="))
	j = bytes.ReplaceAll(j, []byte("\""), []byte(""))
	//fmt.Println(string(j))
	result := string(k) + string(j)
	//fmt.Println(result)
	fmt.Println("[RES] Receive from Device=>>>" , result)

	RSparsing([]byte(result),parameters)

}

func RQparsing(data []byte,parameters *Parameters.Parameter)  {



	temp := bytes.SplitAfterN(data, []byte("}"), 2)  //헤더 , 바디 분리 , n은 2개로 "}" 나뉜다
	header := temp[0]
	body := temp [1]

	//header
	htemp := bytes.TrimSpace(bytes.ReplaceAll(bytes.ReplaceAll(header,[]byte("="),[]byte(":")),[]byte(";"),[]byte(",")))
	var re = regexp.MustCompile(`(\s*?\[\s*?|\s*?\]\s*?)(['"])?([a-zA-Z]+)(['"])?`) //string array
	s := re.ReplaceAll(htemp,[]byte("$1\"$3\""))
	ab := regexp.MustCompile(`(\s*?{\s*?|\s*?,\s*?)(['"])?([a-zA-Z0-9_-]+)(['"])?:`).ReplaceAll(s,[]byte("$1\"$3\":"))
	harray := headerRQ{}
	e := json.Unmarshal(ab, &harray)

	if e != nil {
		panic(e)
	}

	parameters.SetInterfaceID(harray.If)
	parameters.SetDisposableIoTRequestID(harray.Dri)
	fmt.Printf("\tParsing Completed: 'if'=> %d, 'dri'=> %d \r\n",parameters.InterfaceID(), parameters.DisposableIoTRequestID())


	//body
	btemp := bytes.TrimSpace(bytes.ReplaceAll(bytes.ReplaceAll(body,[]byte("="),[]byte(":")),[]byte(";"),[]byte(",")))
	s = re.ReplaceAll(btemp,[]byte("$1\"$3\""))
	ab = regexp.MustCompile(`(\s*?{\s*?|\s*?,\s*?)(['"])?([a-zA-Z0-9_-]+)(['"])?:`).ReplaceAll(s,[]byte("$1\"$3\":"))
	barray := bodyRQ{}
	e = json.Unmarshal(ab, &barray)
	if e != nil {
		panic(e)
	}
	parameters.SetMicroserviceIDs(barray.Mis)
	fmt.Printf("\tParsing Completed: 'mis'=> %d \r\n",parameters.MicroserviceIDs())

	id, _ := gonanoid.Nanoid()
	parameters.SetDeviceID(id)

	devRegParam := EdgeXInterface.NewParameter()
	devRegParam.Name = id
	devRegParam.SetDeviceService("disposable-device")
	devRegParam.SetDeviceProfile("Disposable Temperature Sensing Device")
	EdgeXInterface.DeviceRegistration(devRegParam)
	time.Sleep(5000)

	resourceValue := make(map[string]interface{})
	resourceValue["mis"] = make([]int, len(barray.Mis))
	for i, v := range barray.Mis {
		resourceValue["mis"].([]int)[i] = v
	}
	EdgeXInterface.SetResource(id, "MicroserviceID", resourceValue)
}

func RSparsing(data []byte,parameters *Parameters.Parameter)  {



	temp := bytes.SplitAfterN(data, []byte("}"), 2)  //헤더 , 바디 분리 , n은 2개로 "}" 나뉜다
	header := temp[0]
	body := temp [1]



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


	//body
	btemp := bytes.TrimSpace(bytes.ReplaceAll(bytes.ReplaceAll(body,[]byte("="),[]byte(":")),[]byte(";"),[]byte(",")))
	s = re.ReplaceAll(btemp,[]byte("$1\"$3\""))
	ab = regexp.MustCompile(`(\s*?{\s*?|\s*?,\s*?)(['"])?([a-zA-Z0-9_-]+)(['"])?:`).ReplaceAll(s,[]byte("$1\"$3\":"))
	ab = regexp.MustCompile(`(\s*?:\s*?|\s*?}\s*?)(['"])?([a-zA-Z0-9_-]+)(['"])?}`).ReplaceAll(ab,[]byte("$1\"$3\"}"))
	fmt.Println(string(ab))

	barray := bodyRS{}
	e = json.Unmarshal(ab, &barray)
	if e != nil {
		panic(e)
	}
	parameters.SetDeviceID(barray.Di)
	fmt.Printf("\tParsing Completed: 'di'=> %s \r\n",parameters.DeviceID())

}
