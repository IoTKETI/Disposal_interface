package MicroserviceOutputReport

import (
	"bytes"
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"gw/Parameters"
	"regexp"
	"strconv"
)

type headerRQ struct {
	//======Header
	If  string    `json:"if"`
	Dri string    `json:"dri"`
	Di string `json:"di"`
}

type bodyRQ struct {
	//======Body
	Mi string	`json:"mi"`
	Op string `json:"op"`

}

type headerRS struct {
	//======Header
	Dri int `json:"dri"`
	Rsc int `json:"rsc"`
}




//Request
func Request(parameters *Parameters.Parameter, mqttClient mqtt.Client) {

	tempH := headerRQ{string(parameters.InterfaceID()),string(parameters.DisposableIoTRequestID()),parameters.DeviceID()}
	h, _ := json.Marshal(tempH)
	//	h, _ := json.MarshalIndent(tempH, "", " ")
	k := regexp.MustCompile(`(\s*?{\s*?|\s*?,\s*?)(['"])?([a-zA-Z0-9]+)(['"])?:`).ReplaceAll(h, []byte("$1$3="))
	k = bytes.ReplaceAll(k, []byte(","), []byte(";"))

	tempB := bodyRQ{string(parameters.MicroserviceID()),parameters.OutputParameter()}
	b, _ := json.Marshal(tempB)
	//	b, _ := json.MarshalIndent(tempB, "", " ")
	j := regexp.MustCompile(`(\s*?{\s*?|\s*?,\s*?)(['"])?([a-zA-Z0-9]+)(['"])?:`).ReplaceAll(b, []byte("$1$3="))

	kj := string(k) + string(j)
	result := bytes.ReplaceAll([]byte(kj), []byte("\""), []byte(""))

	//fmt.Println(string(result))
	fmt.Println("[REQ] Send to Server =>>>" , string(result))

	//if token := mqttClient.Publish(parameters.ToTopic(),0,false,result); token.Wait() && token.Error() != nil {
	//	panic(token.Error())
	//}

	//RQparsing([]byte(result),parameters)
}

//Response
func Response(parameters *Parameters.Parameter, mqttClient mqtt.Client) {

	tempH := headerRS{parameters.DisposableIoTRequestID(),parameters.ResponseStatusCode()}
	h, _ := json.Marshal(tempH)
	k := regexp.MustCompile(`(\s*?{\s*?|\s*?,\s*?)(['"])?([a-zA-Z0-9]+)(['"])?:`).ReplaceAll(h, []byte("$1$3="))
	k = bytes.ReplaceAll(k, []byte(","), []byte(";"))

	//tempB := bodyRS{parameters.()}
	//b, _ := json.Marshal(tempB)
	//j := regexp.MustCompile(`(\s*?{\s*?|\s*?,\s*?)(['"])?([a-zA-Z0-9]+)(['"])?:`).ReplaceAll(b, []byte("$1$3="))
	//j = bytes.ReplaceAll(j, []byte(`\"`), []byte(""))
	//j = bytes.ReplaceAll(j, []byte(`"`), []byte(""))
	//j = bytes.ReplaceAll(j, []byte(`\n`), []byte(""))
	//j = bytes.ReplaceAll(j, []byte(`:`), []byte("="))
	////fmt.Println(string(j))
	//result := string(k) + string(j)

	fmt.Println("[RES] Receive from Device=>>>" , string(k))
	//RSparsing(k,parameters)

	//if token := mqttClient.Publishsjlka;fj(parameters.ToTopic(),0,false,string(k)); token.Wait() && token.Error() != nil {
	//	panic(token.Error())
	//}
	//token := mqttClient.Publish(parameters.ToTopic(), 0, false, string(result))
	//token.Wait()
	//time.Sleep(3 * time.Second)
}



func RQparsing(data []byte,parameters *Parameters.Parameter)  {



	temp := bytes.SplitAfterN(data, []byte("}"), 2)  //헤더 , 바디 분리 , n은 2개로 "}" 나뉜다
	header := temp[0]
	body := temp [1]


	//header
	htemp := bytes.TrimSpace(bytes.ReplaceAll(bytes.ReplaceAll(header,[]byte("="),[]byte(":")),[]byte(";"),[]byte(",")))
	//var re = regexp.MustCompile(`(\s*?:\s*?|\s*?\,\s*?)(['"])?([a-zA-Z]+)(['"])?`) //string array
	//
	var re = regexp.MustCompile(`([a-zA-Z0-9-]+):([a-zA-Z0-9-]+)`) //string array
	s := re.ReplaceAll(htemp,[]byte("\"$1\":\"$2\""))
	fmt.Println(string(s))
	harray := headerRQ{}
	e := json.Unmarshal(s, &harray)

	if e != nil {
		panic(e)
	}

	v,_ := strconv.Atoi(harray.If)
	parameters.SetInterfaceID(v)
	v,_ = strconv.Atoi(harray.Dri)
	parameters.SetDisposableIoTRequestID(v)
	parameters.SetDeviceID(harray.Di)
	fmt.Printf("\tParsing Completed: 'if'=> %d, 'dri'=> %d, 'di'=> %s \r\n",parameters.InterfaceID(), parameters.DisposableIoTRequestID(),parameters.DeviceID())


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

	v,_ = strconv.Atoi(barray.Mi)
	parameters.SetMicroserviceID(v)
	parameters.SetOutputParameter(barray.Op)
	fmt.Printf("\tParsing Completed: 'mi'=> %s, 'op'=> %v \r\n",parameters.MicroserviceID(),parameters.OutputParameter())

}


func RSparsing(data []byte,parameters *Parameters.Parameter)  {



	temp := bytes.SplitAfterN(data, []byte("}"), 2)  //헤더 , 바디 분리 , n은 2개로 "}" 나뉜다
	header := temp[0]
	//body := temp[1]



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
/*	btemp := bytes.TrimSpace(bytes.ReplaceAll(bytes.ReplaceAll(body,[]byte("="),[]byte(":")),[]byte(";"),[]byte(",")))
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
	fmt.Printf("\tParsing Completed: 'op'=> %s \r\n",parameters.OutputParameter())*/

}