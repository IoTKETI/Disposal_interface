package TaskParameterSet

import (
	"bytes"
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"gw/ObjectTypeParameters"
	"gw/Parameters"
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
	Tis []int 	`json:"tis"`
	fp ObjectTypeParameters.Fp
}

type headerRS struct {
	//======Header
	Dri int `json:"dri"`
	Rsc int `json:"rsc"`
}


//Request
func Request(parameters *Parameters.Parameter, mqttClient mqtt.Client) {

	tempH := headerRQ{parameters.InterfaceID(),parameters.DisposableIoTRequestID()}
	h, _ := json.Marshal(tempH)
	//	h, _ := json.MarshalIndent(tempH, "", " ")
	k := regexp.MustCompile(`(\s*?{\s*?|\s*?,\s*?)(['"])?([a-zA-Z0-9]+)(['"])?:`).ReplaceAll(h, []byte("$1$3="))
	k = bytes.ReplaceAll(k, []byte(","), []byte(";"))

	tempB := bodyRQ{parameters.MicroserviceIDs(), parameters.TaskIDs(),parameters.FlexibleTaskParameter()}
	b, _ := json.Marshal(tempB)
	//	b, _ := json.MarshalIndent(tempB, "", " ")
	j := regexp.MustCompile(`(\s*?{\s*?|\s*?,\s*?)(['"])?([a-zA-Z0-9]+)(['"])?:`).ReplaceAll(b, []byte("$1$3="))

	kj := string(k) + string(j)
	result := bytes.ReplaceAll([]byte(kj), []byte("\""), []byte(""))

	//fmt.Println(string(result))
	fmt.Println("[REQ] Send to Device =>>>" , string(result))

	token := mqttClient.Publish(parameters.ToTopic(), 0, false, string(result))
	token.Wait()
	//RQparsing([]byte(result),parameters)

}

//Response
func Response(parameters *Parameters.Parameter,mqttClient mqtt.Client) {

	tempH := headerRS{parameters.DisposableIoTRequestID(),parameters.ResponseStatusCode()}
	h, _ := json.Marshal(tempH)
	k := regexp.MustCompile(`(\s*?{\s*?|\s*?,\s*?)(['"])?([a-zA-Z0-9]+)(['"])?:`).ReplaceAll(h, []byte("$1$3="))
	k = bytes.ReplaceAll(k, []byte(","), []byte(";"))

	//fmt.Println(string(k))
	fmt.Println("[RES] Receive from Servce=>>>" , string(k))
	token := mqttClient.Publish(parameters.ToTopic(), 0, false, string(k))
	token.Wait()
	//RSparsing(k,parameters)
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
	parameters.SetTaskIDs(barray.Tis)
	parameters.SetFlexibleTaskParameter(barray.fp)
	fmt.Printf("\tParsing Completed: 'mis'=> %v 'tis'=> %v 'tis'=> %v \r\n",parameters.MicroserviceIDs(),parameters.TaskIDs(),parameters.FlexibleTaskParameter())

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


}