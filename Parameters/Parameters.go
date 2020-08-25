package Parameters

import (
	"gw/ObjectTypeParameters"
)

type Parameter struct {
	interfaceID             int //if
	disposableIoTRequestID  int //dri
	responseStatusCode      int //rsc
	microserviceID          int //mi
	microserviceIDs         []int //mis
	microserviceInformation []ObjectTypeParameters.MifObject
	inputParameter          string
	inputParameters         []string
	outputParameter         string // op
	outputParameters        []string // ops
	changedOutputParameter  []string //cop
	microserviceConfigure   string
	taskID                  int     //ti
	taskIDs                 []int  //tis
	taskInformation         []ObjectTypeParameters.TifObject
	flexibleTaskParameter   ObjectTypeParameters.Fp
	flexibleTaskParameters  []string   //
	staticTaskParameter     []byte		//sp
	spArray					[]string	//spArray for Tif
	staticTaskParameters    []string   //
	taskOrchestration       bool  //to
	taskConfigure           string
	deviceID                string
	driIf					map[int]string
	myTopic					string
	toTopic					string
}

func (p *Parameter) ToTopic() string {
	return p.toTopic
}

func (p *Parameter) SetToTopic(toTopic string) {
	p.toTopic = toTopic
}

func (p *Parameter) MyTopic() string {
	return p.myTopic
}

func (p *Parameter) SetMyTopic(myTopic string) {
	p.myTopic = myTopic
}


func (p *Parameter) DriIf() map[int]string {
	return p.driIf
}

func (p *Parameter) SetDriIf(dirIf map[int]string) {
	p.driIf = dirIf
}

func (p *Parameter) OutputParameter() string {
	return p.outputParameter
}

func (p *Parameter) SetOutputParameter(outputParameter string) {
	p.outputParameter = outputParameter
}

func (p *Parameter) SpArray() []string {
	return p.spArray
}

func (p *Parameter) SetSpArray(spArray []string) {
	p.spArray = spArray
}


func (p *Parameter) MicroserviceInformation() []ObjectTypeParameters.MifObject {
	return p.microserviceInformation
}

func (p *Parameter) SetMicroserviceInformation(microserviceInformation []ObjectTypeParameters.MifObject) {
	p.microserviceInformation = microserviceInformation
}

func (p *Parameter) TaskInformation() []ObjectTypeParameters.TifObject {
	return p.taskInformation
}

func (p *Parameter) SetTaskInformation(taskInformation []ObjectTypeParameters.TifObject) {
	p.taskInformation = taskInformation
}




//di
func (p *Parameter) DeviceID() string {
	return p.deviceID
}

//di
func (p *Parameter) SetDeviceID(deviceID string) {
	p.deviceID = deviceID
}

func (p *Parameter) TaskConfigure() string {
	return p.taskConfigure
}

func (p *Parameter) SetTaskConfigure(taskConfigure string) {
	p.taskConfigure = taskConfigure
}

func (p *Parameter) TaskOrchestration() bool {
	return p.taskOrchestration
}

func (p *Parameter) SetTaskOrchestration(taskOrchestration bool) {
	p.taskOrchestration = taskOrchestration
}

func (p *Parameter) StaticTaskParameters() []string {
	return p.staticTaskParameters
}

func (p *Parameter) SetStaticTaskParameters(staticTaskParameters []string) {
	p.staticTaskParameters = staticTaskParameters
}

func (p *Parameter) StaticTaskParameter() []byte {
	return p.staticTaskParameter
}

func (p *Parameter) SetStaticTaskParameter(staticTaskParameter []byte) {
	p.staticTaskParameter = staticTaskParameter
}

func (p *Parameter) FlexibleTaskParameters() []string {
	return p.flexibleTaskParameters
}

func (p *Parameter) SetFlexibleTaskParameters(flexibleTaskParameters []string) {
	p.flexibleTaskParameters = flexibleTaskParameters
}

func (p *Parameter) FlexibleTaskParameter() ObjectTypeParameters.Fp {
	return p.flexibleTaskParameter
}

func (p *Parameter) SetFlexibleTaskParameter(flexibleTaskParameter ObjectTypeParameters.Fp) {
	p.flexibleTaskParameter = flexibleTaskParameter
}



func (p *Parameter) TaskIDs() []int {
	return p.taskIDs
}

func (p *Parameter) SetTaskIDs(taskIDs []int) {
	p.taskIDs = taskIDs
}

func (p *Parameter) TaskID() int {
	return p.taskID
}

func (p *Parameter) SetTaskID(taskID int) {
	p.taskID = taskID
}

func (p *Parameter) MicroserviceConfigure() string {
	return p.microserviceConfigure
}

func (p *Parameter) SetMicroserviceConfigure(microserviceConfigure string) {
	p.microserviceConfigure = microserviceConfigure
}

func (p *Parameter) ChangedOutputParameter() []string {
	return p.changedOutputParameter
}

func (p *Parameter) SetChangedOutputParameter(changedOutputParameter []string) {
	p.changedOutputParameter = changedOutputParameter
}

func (p *Parameter) OutputParameters() []string {
	return p.outputParameters
}
//ops string array
func (p *Parameter) SetOutputParameters(outputParameters []string) {
	p.outputParameters = outputParameters
}



func (p *Parameter) InputParameters() []string {
	return p.inputParameters
}

func (p *Parameter) SetInputParameters(inputParameters []string) {
	p.inputParameters = inputParameters
}

func (p *Parameter) InputParameter() string {
	return p.inputParameter
}

func (p *Parameter) SetInputParameter(inputParameter string) {
	p.inputParameter = inputParameter
}


//mis
func (p *Parameter) MicroserviceIDs() []int {
	return p.microserviceIDs
}
//mis
func (p *Parameter) SetMicroserviceIDs(microserviceIDs []int) {
	p.microserviceIDs = microserviceIDs
}

func (p *Parameter) MicroserviceID() int {
	return p.microserviceID
}

func (p *Parameter) SetMicroserviceID(microserviceID int) {
	p.microserviceID = microserviceID
}
//rsc
func (p *Parameter) ResponseStatusCode() int {
	return p.responseStatusCode
}

//rsc
func (p *Parameter) SetResponseStatusCode(responseStatusCode int) {
	p.responseStatusCode = responseStatusCode
}
//dri
func (p *Parameter) DisposableIoTRequestID() int {
	return p.disposableIoTRequestID
}

//dri
func (p *Parameter) SetDisposableIoTRequestID(disposableIoTRequestID int) {
	p.disposableIoTRequestID = disposableIoTRequestID
}
//if
func (p *Parameter) InterfaceID() int {
	return p.interfaceID
}

//Set "if"
func (p *Parameter) SetInterfaceID(interfaceID int) {
	p.interfaceID = interfaceID
}

func NewParameter() *Parameter {

	return &Parameter{}
}
