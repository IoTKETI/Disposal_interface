package MicroserviceOutputReport

import (
	"bytes"
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"gw/Parameters"
	"regexp"
	"strconv"

	//"strconv"
)

type headerRQ struct {
	//======Header
	If  int    `json:"if"`
	Dri int    `json:"dri"`
	Di string `json:"di"`
}

type bodyRQ struct {
	//======Body
	Mi int	`json:"mi"`
	Op string `json:"op"`

}

type headerRS struct {
	//======Header
	Dri int `json:"dri"`
	Rsc int `json:"rsc"`
}




//Request
func Request(parameters *Parameters.Parameter, mqttClient mqtt.Client) {

	tempH := headerRQ{parameters.InterfaceID(),parameters.DisposableIoTRequestID(),parameters.DeviceID()}
	h, _ := json.Marshal(tempH)
	//	h, _ := json.MarshalIndent(tempH, "", " ")
	k := regexp.MustCompile(`(\s*?{\s*?|\s*?,\s*?)(['"])?([a-zA-Z0-9]+)(['"])?:`).ReplaceAll(h, []byte("$1$3="))
	k = bytes.ReplaceAll(k, []byte(","), []byte(";"))
	fmt.Println(string(h))

	tempB := bodyRQ{parameters.MicroserviceID(),parameters.OutputParameter()}
	b, _ := json.Marshal(tempB)
	//	b, _ := json.MarshalIndent(tempB, "", " ")
	j := regexp.MustCompile(`(\s*?{\s*?|\s*?,\s*?)(['"])?([a-zA-Z0-9]+)(['"])?:`).ReplaceAll(b, []byte("$1$3="))

	kj := string(k) + string(j)
	result := bytes.ReplaceAll([]byte(kj), []byte("\n"), []byte(""))

	//fmt.Println(string(result))
	fmt.Println("[REQ] Send to Server =>>>" , string(result))

	//if token := mqttClient.Publish(parameters.ToTopic(),0,false,result); token.Wait() && token.Error() != nil {
	//	panic(token.Error())
	//}

	RQparsing([]byte(result),parameters)
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
//	var re = regexp.MustCompile(`(\s*?:\s*?|\s*?\,\s*?)(['"])?([a-zA-Z]+)(['"])?`) //string array

	var all = regexp.MustCompile(`([a-zA-Z0-9-]+):([a-zA-Z0-9-]+)?`) //string array
	allquoted := all.ReplaceAll(htemp,[]byte("\"$1\":\"$2\""))
	allquoted = bytes.ReplaceAll([]byte(allquoted),[]byte("\"\""),[]byte(""))
	fmt.Println(string(allquoted))

	var mapTemp map[string]interface{}
	if err := json.Unmarshal(allquoted, &mapTemp); err != nil {
		panic(err)
	}
	//fmt.Println(mapTemp)

	v,_ := strconv.Atoi(mapTemp["if"].(string))
	parameters.SetInterfaceID(v)
	v,_ = strconv.Atoi(mapTemp["dri"].(string))
	parameters.SetDisposableIoTRequestID(v)
	parameters.SetDeviceID(mapTemp["di"].(string))
	fmt.Printf("\tParsing Completed: 'if'=> %d, 'dri'=> %d, 'di'=> %s \r\n",parameters.InterfaceID(), parameters.DisposableIoTRequestID(),parameters.DeviceID())


	//body
	///fmt.Println(string(body))
	btemp := bytes.TrimSpace(bytes.ReplaceAll(bytes.ReplaceAll(body,[]byte("="),[]byte(":")),[]byte(";"),[]byte(",")))
	btemp = bytes.ReplaceAll(btemp,[]byte("\""),[]byte(""))
	btemp = bytes.ReplaceAll(btemp,[]byte("n"),[]byte(""))
	btemp = bytes.ReplaceAll(btemp,[]byte("\\"),[]byte(""))
	fmt.Println(string(btemp))
	//s := re.ReplaceAll(btemp,[]byte("$1\"$3\""))
	//var stringArray = regexp.MustCompile(`(\s*?\[\s*?|\s*?\]\s*?)(['"])?([a-zA-Z]+)(['"])?`) //string array
	//s = stringArray.ReplaceAll(btemp,[]byte("$1\"$3\""))
	//ab := regexp.MustCompile(`(\s*?{\s*?|\s*?,\s*?)(['"])?([a-zA-Z0-9]+)(['"])?:`).ReplaceAll(s,[]byte("$1\"$3\":"))
	//fmt.Println(string(ab))
	var alla = regexp.MustCompile(`([a-zA-Z0-9-]+):([a-zA-Z0-9-]+)?`) //string array
	allquoted = alla.ReplaceAll(btemp,[]byte("\"$1\":\"$2\""))
	allquoted = bytes.ReplaceAll([]byte(allquoted),[]byte("\"\""),[]byte(""))
//	fmt.Println(string(allquoted))

	var bmapTemp map[string]interface{}
	if err := json.Unmarshal(allquoted, &bmapTemp); err != nil {
		panic(err)
	}

//	fmt.Println(bmapTemp)

	v,_ = strconv.Atoi(bmapTemp["mi"].(string))
	parameters.SetMicroserviceID(v)
	//str := fmt.Sprintf("%v",bmapTemp["op"])
	parameters.SetOpMap(bmapTemp["op"].(map[string]interface{}))
	//parameters.SetOutputParameter(mapTemp["op"].(ObjectTypeParameters.Op))
	fmt.Printf("\tParsing Completed: 'mi'=> %s, 'op'=> %v \r\n",parameters.MicroserviceID(),parameters.OpMap())

	for k,v := range parameters.OpMap(){
		fmt.Println(k,v)
	}
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