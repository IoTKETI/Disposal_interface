package EdgeXInterface

import (
	EdgeXURL "gw/EdgeXInterface/URL"
	"encoding/json"
	"net/http"
	"bytes"
	"fmt"
	"strings"
	"time"
)

func SetResource(devName string, devCommand string, resourceValue map[string]interface{}){
	tempBody, _ := json.Marshal(resourceValue)
	tempBodyStr := strings.Replace(string(tempBody),"[","\"[",-1)
	tempBodyStr = strings.Replace(tempBodyStr,"]","]\"",-1)
	body := bytes.NewBufferString(tempBodyStr)
	time.Sleep(time.Second * 5)
	fmt.Println(EdgeXURL.SetResource+"/"+devName+"/command/"+devCommand)
	req, _ := http.NewRequest("PUT", EdgeXURL.SetResource+"/"+devName+"/command/"+devCommand, body)
	req.Header.Add("Content-Type", "application/json")
	client := &http.Client {}
	client.Do(req)
}
