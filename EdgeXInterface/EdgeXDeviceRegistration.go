package EdgeXInterface

import (
	EdgeXURL "gw/EdgeXInterface/URL"
	"encoding/json"
	"net/http"
	"bytes"
	"fmt"
	"io/ioutil"
)

type Address struct {
	Address string
}

type ProtocolOther struct {
	Other Address `json:"other"`
}

type ServiceInfo struct {
	Name string `json:"name"`
}

type ProfileInfo struct {
	Name string `json:"name"`
}

type DevRegParams struct {
	Name string `json:"name"`
	Protocols ProtocolOther `json:"protocols"`
	AdminState string
	OperatingState string
	Service ServiceInfo `json:"service"`
	Profile ProfileInfo `json:"profile"`
}

func NewParameter() *DevRegParams{
	res := DevRegParams{}
	addr := Address{"NA"}
	prtcOther := ProtocolOther{addr}
	res.Protocols = prtcOther
	res.AdminState = "unlocked"
	res.OperatingState = "enabled"
	return &res
}

func DeviceRegistration(param *DevRegParams){
	tempBody, _ := json.Marshal(param)
	body := bytes.NewBuffer(tempBody)
	fmt.Println(EdgeXURL.DeviceRegistration)
	res,_ := http.Post(EdgeXURL.DeviceRegistration, "application/json", body)
	defer res.Body.Close()
	ioutil.ReadAll(res.Body)
}

func (p *DevRegParams) SetDeviceService(name string){
	sInfo := ServiceInfo{name}
	p.Service = sInfo
}

func (p *DevRegParams) SetDeviceProfile(name string){
	pInfo := ProfileInfo{name}
	p.Profile = pInfo
}
