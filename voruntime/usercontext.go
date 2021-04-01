package voruntime

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

var Users UserRC
type UserContext struct{
	Name string
	Group string
	Permission []string   //internal,exec,plugins
}
type PermissoinToken struct{
	Internal string `json:"internal""`
	Exec string `json:"exec""`
	Plugin string `json:"plugins"`
}
type User struct{
	PasswordEncrypted string `json:"password_encrypted"`
}
type UserRC struct{
	Groups map[string]map[string]User      `json:"user_groups"`
	Permissions map[string]PermissoinToken `json:"group_permission"`
}

func InitUser(){
	file,err:=os.Open("users.json")
	if err!=nil{
		log.Fatal(err)
	}
	jcont,err:=ioutil.ReadAll(file)
	if err!=nil{
		log.Fatal(err)
	}
	json.Unmarshal(jcont, &Users)
}
func Login(name string,group string,password string)(UserContext,error){
	h:=sha256.New()
	h.Write([]byte(password))
	p_has:=hex.EncodeToString(h.Sum(nil))
	p:=""
	if g,ok:=Users.Groups[group];ok{
		if u,ok:=g[name];ok{
			p=u.PasswordEncrypted
		}else{
			return UserContext{},errors.New("user not found")
		}
	}else{
		return UserContext{},errors.New("group not found")
	}
	if p==p_has{
		return UserContext{
			Name:       name,
			Group:      group,
			Permission: []string{Users.Permissions[group].Internal,Users.Permissions[group].Exec,Users.Permissions[group].Plugin},
		},nil
	}else{
		return UserContext{},errors.New("password incorrect")
	}
}

func PermissionFilter(commandname string,permission string) (bool,string){
	if permission == ""{
		return true,"All commands access are granted of this user group."
	}
	permissionToken:=strings.Split(permission,",")
	tmpPerm:=true
	tmpReason:="Permission restrictions of command '"+commandname+"' not found."
	for _,v:=range permissionToken{
		if v=="-"{
			tmpPerm=false
			tmpReason="Permission denied for command '"+commandname+"' of this user group."
		}else{
			if v=="-"+commandname {
				tmpPerm=false
				tmpReason="Permission denied for command '"+commandname+"' of this user group."
			}else if v==commandname{
				tmpPerm=true
				tmpReason="Command '"+commandname+"' is granted of this user group."
			}else{
				continue
			}
		}
	}
	return tmpPerm,tmpReason
}

func PermissionVisualize(uctx *UserContext) string{
	i:=uctx.Permission[0]
	e:=uctx.Permission[1]
	p:=uctx.Permission[2]
	com:=""
	if i==""&&e==""&&p==""{
		com=": all granted"
	}
	return "I("+i+"),E("+e+"),P("+p+")"+com
}