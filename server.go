package main

import (
	"bytes"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	DeviceRegistration "gw/1.DeviceRegistration"
	MicroserviceCreation "gw/10.MicroserviceCreation"
	MicroserviceRun "gw/11.MicroserviceRun"
	MicroserviceOutputParameterRead "gw/13.MicroserviceOutputParameterRead"
	MicroserviceOutputReport "gw/14.MicroserviceOutputReport"
	MicroserviceStop "gw/15.MicroserviceStop"
	DeviceMicroserviceInformationReport "gw/2.DeviceMicroserviceInformationReport"
	TaskParameterSet "gw/21.TaskParameterSet"
	"gw/Builder"
	"gw/Parameters"
	"gw/ResourceName"
	"io"
	"log"
	"net"
)

var mqttClient mqtt.Client
var deviceInfo map[string]Parameters.Parameter

func MsgHandler(client mqtt.Client, mqtt_msg mqtt.Message) {
	data := mqtt_msg.Payload()
	cParm := make(chan *Parameters.Parameter)
	go getParam(cParm)
	param := <-cParm

	if bytes.Contains(data,[]byte("if")){
		ParseRequestMsg(data,param, client)
	} else {
		ParseResponseMsg(data,param,client)
	}

}


var parameters *Parameters.Parameter

func getParam(pa chan *Parameters.Parameter) {
	pa <- parameters
}

func setParam(p *Parameters.Parameter)  {
	parameters = p
}

func main() {
	parameters := Parameters.NewParameter()
	setParam(parameters)
	parameters.SetMyTopic("A")
	parameters.SetToTopic("A")
	opts := mqtt.NewClientOptions().AddBroker("127.0.0.1:1883")
	opts.ClientID=parameters.MyTopic()
	opts.SetDefaultPublishHandler(MsgHandler)
	deviceInfo = make(map[string]Parameters.Parameter)
	mqttClient = mqtt.NewClient(opts)


	driIf := make(map[int]string)
	parameters.SetDriIf(driIf)

	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	if token := mqttClient.Subscribe(parameters.MyTopic(),0,nil); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		return


	}

	fmt.Println("Interface 14===================MicroserviceOutputReport")
	parameters.SetInterfaceID(14)
	parameters.SetDisposableIoTRequestID(12351)
	parameters.SetDeviceID("4dasc44321")
	parameters.DriIf()[parameters.DisposableIoTRequestID()] = ResourceName.MicroserviceOutputReport
	parameters.SetMicroserviceID(1)
	Builder.Op([]string{"Atemp"},[]string{"29"},parameters)
	MicroserviceOutputParameterRead.Response(parameters)
	MicroserviceCreation.Request(parameters,mqttClient)





/*
	//Interface 10===================Microservice Creation
	// Originator: Server
	fmt.Println("Interface 10===================Microservice Creation")
	parameters.SetInterfaceID(10)
	parameters.SetDisposableIoTRequestID(12347)
	parameters.DriIf()[parameters.DisposableIoTRequestID()] = ResourceName.MicroserviceCreation
	fmt.Println(parameters.DriIf()[parameters.DisposableIoTRequestID()])
	parameters.SetMicroserviceIDs([]int{1,2})

	MicroserviceCreation.Request(parameters,mqttClient)



	//Interface 11===================Microservice Run
	// Originator: Server
	fmt.Println("Interface 11===================Microservice Run")
	parameters.SetInterfaceID(11)
	parameters.SetDisposableIoTRequestID(12349)
	parameters.DriIf()[parameters.DisposableIoTRequestID()] = ResourceName.MicroserviceRun
	parameters.SetMicroserviceIDs([]int{1,2})
	MicroserviceRun.Request(parameters, mqttClient)


	//Interface 21===================TaskParameterSet
	// Originator: Server
	fmt.Println("Interface 21===================TaskParameterSet")
	parameters.SetInterfaceID(21)
	parameters.SetDisposableIoTRequestID(12348)
	parameters.DriIf()[parameters.DisposableIoTRequestID()] = ResourceName.TaskParameterSet
	parameters.SetMicroserviceIDs([]int{1})
	parameters.SetTaskIDs([]int{1})
	fp := ObjectTypeParameters.Fp{}
	fp.Oprd ="C"
	parameters.SetFlexibleTaskParameter(fp)
	TaskParameterSet.Request(parameters,mqttClient)


	//Interface 15===================MS Stop
	// Originator: Server
	fmt.Println("Interface 15===================MS Stop")
	parameters.SetInterfaceID(15)
	parameters.SetDisposableIoTRequestID(12352)
	parameters.DriIf()[parameters.DisposableIoTRequestID()] = ResourceName.MicroserviceStop
	parameters.SetMicroserviceIDs([]int{1,2})
	MicroserviceStop.Request(parameters, mqttClient)
*/
	select{}

}

func ParseRequestMsg(data []byte, parameters *Parameters.Parameter,client mqtt.Client) {
	fmt.Println("[Received RQ] ========== ParseRequestMsg")
	headerTemp := make(map[string]string)
	temp := bytes.SplitAfterN(data, []byte("}"), 2)  //헤더 , 바디 분리 , n은 2개로 "}" 나뉜다
	header := bytes.Split(bytes.ReplaceAll(bytes.ReplaceAll(temp[0], []byte("{"), []byte("")), []byte("}"), []byte("")), []byte(";"))
	for j := range header {
		headerTemp[string(bytes.Split(bytes.ReplaceAll(bytes.TrimSpace(header[j]), []byte(";"), []byte("")), []byte("="))[0])] = string(bytes.Split(bytes.ReplaceAll(bytes.TrimSpace(header[j]), []byte(";"), []byte("")), []byte("="))[1])
	}

	switch headerTemp["if"]{
		case ResourceName.MicroserviceCreation:
			fmt.Println("[Received RQ] - Microservice Creation")
			MicroserviceCreation.RQparsing(data,parameters)
			parameters.SetResponseStatusCode(200)
			MicroserviceCreation.Response(parameters,client)
			parameters.DriIf()[parameters.DisposableIoTRequestID()] = ResourceName.MicroserviceCreation

			break
		case ResourceName.MicroserviceRun:
			fmt.Println("[Received RQ] - Microservice Run")
			MicroserviceRun.RQparsing(data,parameters)
			parameters.SetResponseStatusCode(200)
			MicroserviceRun.Response(parameters,client)

			parameters.DriIf()[parameters.DisposableIoTRequestID()] = ResourceName.MicroserviceRun
			break
		case ResourceName.TaskParameterSet:
			fmt.Println("[Received RQ] - TaskParameterSet")
			TaskParameterSet.RQparsing(data,parameters)
			parameters.SetResponseStatusCode(200)
			TaskParameterSet.Response(parameters,client)
			parameters.DriIf()[parameters.DisposableIoTRequestID()] = ResourceName.TaskParameterSet
			break
		case ResourceName.MicroserviceStop:
			fmt.Println("[Received RQ] - MicroserviceStop")
			MicroserviceStop.RQparsing(data,parameters)
			parameters.SetResponseStatusCode(200)
			MicroserviceStop.Response(parameters,client)
			parameters.DriIf()[parameters.DisposableIoTRequestID()] = ResourceName.MicroserviceStop
			break
		case ResourceName.DeviceRegistration:

			DeviceRegistration.RQparsing(data,parameters)
			parameters.SetResponseStatusCode(200)
			DeviceRegistration.Response(parameters)

			//driIf[parameters.DisposableIoTRequestID()] = ResourceName.DeviceRegistration
			break
		case ResourceName.DeviceMicroserviceInformationReport:

			DeviceMicroserviceInformationReport.RQparsing(data, parameters)
//			deviceInfo[parameters.DeviceID()] = parameters
			DeviceMicroserviceInformationReport.Response(parameters)
		//	driIf[parameters.DisposableIoTRequestID()] = ResourceName.DeviceMicroserviceInformationReport
			break
		case ResourceName.DeviceTaskInformationRequest:


			//driIf[parameters.DisposableIoTRequestID()] = ResourceName.DeviceTaskInformationRequest
			break
		case ResourceName.TaskRun:

			//driIf[parameters.DisposableIoTRequestID()] = ResourceName.TaskRun
			break
		case ResourceName.MicroserviceOutputReport:
			fmt.Println("[Received RQ] - MicroserviceOutputReport")
			MicroserviceOutputReport.RQparsing(data,parameters)
			parameters.SetResponseStatusCode(200)
			MicroserviceOutputReport.Response(parameters,client)
			parameters.DriIf()[parameters.DisposableIoTRequestID()] = ResourceName.MicroserviceOutputReport
			//driIf[parameters.DisposableIoTRequestID()] = ResourceName.MicroserviceOutputReport
			break
		case ResourceName.TaskOff:

			//driIf[parameters.DisposableIoTRequestID()] = ResourceName.TaskOff
			break
		}


	//parameters.SetDriIf(driIf)

}

func ParseResponseMsg(data []byte,parameters *Parameters.Parameter,client mqtt.Client)  {
	fmt.Println("[Received RS] - ParseResponseMsg")

	switch parameters.DriIf()[parameters.DisposableIoTRequestID()] {

	case ResourceName.MicroserviceCreation:
		fmt.Println("[Received RS] - MicroserviceCreation")
		MicroserviceCreation.RSparsing(data,parameters)
		delete(parameters.DriIf(),parameters.DisposableIoTRequestID())
		break
	case ResourceName.MicroserviceRun:
		fmt.Println("[Received RS] - MicroserviceRun")
		MicroserviceRun.RSparsing(data,parameters)
		delete(parameters.DriIf(),parameters.DisposableIoTRequestID())
		break
	case ResourceName.TaskParameterSet:
		fmt.Println("[Received RS] - TaskParameterSet")
		TaskParameterSet.RSparsing(data,parameters)
		delete(parameters.DriIf(),parameters.DisposableIoTRequestID())
		break
	case ResourceName.MicroserviceStop:
		fmt.Println("[Received RS] - MicroserviceStop")
		MicroserviceStop.RSparsing(data,parameters)
		delete(parameters.DriIf(),parameters.DisposableIoTRequestID())
		break

	case ResourceName.MicroserviceOutputReport:
		fmt.Println("[Received RS] - MicroserviceOutputReport")
		MicroserviceOutputReport.RSparsing(data,parameters)
		delete(parameters.DriIf(),parameters.DisposableIoTRequestID())
		break
	case ResourceName.DeviceRegistration:

		DeviceRegistration.RSparsing(data,parameters)

		delete(parameters.DriIf(),parameters.DisposableIoTRequestID())
		break
	case ResourceName.DeviceMicroserviceInformationReport:


		delete(parameters.DriIf(),parameters.DisposableIoTRequestID())
		break
	case ResourceName.DeviceTaskInformationRequest:

		delete(parameters.DriIf(),parameters.DisposableIoTRequestID())
		break
	case ResourceName.TaskRun:

		delete(parameters.DriIf(),parameters.DisposableIoTRequestID())
		break
	case ResourceName.TaskOff:

		delete(parameters.DriIf(),parameters.DisposableIoTRequestID())
		break
	}
}

func ConnHandler(conn net.Conn) {
	//parameters := Parameters.NewParameter()


	recvBuf := make([]byte, 4096)
	for {
		n, err := conn.Read(recvBuf)
		if nil != err {
			if io.EOF == err {
				log.Println(err)
				return
			}
			log.Println(err)
			return
		}
		if 0 < n {
			data := recvBuf[:n]
			if bytes.Contains(data,[]byte("if")){
			//	ParseRequestMsg(data,parameters)
			} else { }//ParseResponseMsg(data,parameters) }
			//_, err = conn.Write([]byte(""))
			//if err != nil {
			//	log.Println(err)
			//	return
			//}
		}
	}
}


