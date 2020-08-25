package Builder

import (
	"gw/ObjectTypeParameters"
	"gw/Parameters"
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)



func MifObjectArray(length int, misArray []int, opsArray [][]string,parameters *Parameters.Parameter) {
//length, mis array , ops array

	tempMifArray := make([]ObjectTypeParameters.MifObject, length )

	for i := range misArray {
		tempMifArray[i].Mi = misArray[i]
	}

	for i := range opsArray {
		tempMifArray[i].Ops = opsArray[i]
	}

	parameters.SetMicroserviceInformation(tempMifArray)
}


func TifObjectArray(length int, tiArray []int,  parameters *Parameters.Parameter) {
//length, ti array, sp obj array,  fp obj array , to array

	tempTifArray := make([]ObjectTypeParameters.TifObject, length )

	for i := range tiArray {
		tempTifArray[i].Ti = tiArray[i]
	}

	for i := range parameters.SpArray() {
		tempTifArray[i].Sp = parameters.SpArray()[i]
	}

	//for i := range fpObjArray {
	//	tempTifArray[i].Fp = fpObjArray[i]
	//}
	//
	//for i := range toArray {
	//	tempTifArray[i].To = toArray[i]
	//}

	parameters.SetTaskInformation(tempTifArray)
}


func SpObjectArray(length int, key [][]string, value [][]string, parameters *Parameters.Parameter) {
// static parameter name, value array


	tempSpObjArray := make([]string, length)

	for i := range key {

		var sfs []reflect.StructField

		for j,v := range value[i] {
			t := reflect.TypeOf(v)
			sf := reflect.StructField{
				Name: key[i][j],
				Type: t,
				Tag:  reflect.StructTag(fmt.Sprintf(`json:"%s"`, strings.ToLower(key[i][j]))),
			}
			sfs = append(sfs, sf)
		}

		typ := reflect.StructOf(sfs)


		v := reflect.New(typ).Elem()
		for j := range value[i] {
			v.Field(j).SetString(value[i][j])
		}
		s := v.Addr().Interface()

		w := new(bytes.Buffer)
		if err := json.NewEncoder(w).Encode(s); err != nil {
			panic(err)
		}

		//fmt.Printf("value: %+v\n", s)
		//fmt.Printf("json:  %s", w.Bytes())
		//parameters.SetStaticTaskParameter(w.Bytes())

		tempSpObjArray[i] = string(w.Bytes())

	}

	parameters.SetSpArray(tempSpObjArray)

}

func Op(key []string, value []string, parameters *Parameters.Parameter) {

	var sfs []reflect.StructField

	for i := range value {
		//t := reflect.TypeOf(v)
		sf := reflect.StructField{
			Name: key[i],
			Type: reflect.TypeOf("s"),
			Tag:  reflect.StructTag(fmt.Sprintf(`json:"%s"`, strings.ToLower(key[i]))),
		}
		sfs = append(sfs, sf)
	}

	typ := reflect.StructOf(sfs)

	v := reflect.New(typ).Elem()
	for j := range value {
		v.Field(j).SetString(value[j])
	}
	s := v.Addr().Interface()

	w := new(bytes.Buffer)
	if err := json.NewEncoder(w).Encode(s); err != nil {
		panic(err)
	}

	//fmt.Printf("value: %+v\n", s)
	//fmt.Printf("json:  %s", w.Bytes())

	//parameters.SetStaticTaskParameter(w.Bytes())

	parameters.SetOutputParameter(string(w.Bytes()))


}


func FpObjectArray(key []string, value []string) {
// flexible parameter name, value array



}