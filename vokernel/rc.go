package vokernel

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

var RC map[string]string
func InitRC(){
	file,err:=os.Open("vsrc.json")
	if err!=nil{
		log.Fatal(err)
	}
	jcont,err:=ioutil.ReadAll(file)
	file.Close()
	if err!=nil{
		log.Fatal(err)
	}
	json.Unmarshal(jcont, &RC)
}
