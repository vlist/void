package voruntime

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

var RC map[string]string
func Initrc(){
	file,err:=os.Open("vsrc.json")
	if err!=nil{
		log.Fatal(err)
	}
	jcont,err:=ioutil.ReadAll(file)
	if err!=nil{
		log.Fatal(err)
	}
	json.Unmarshal(jcont, &RC)
}
