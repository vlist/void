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
	Permission []string   //I,E,P
}
type PermissionToken struct{
	Internal string `json:"internal""`
	Exec string `json:"exec""`
	Plugin string `json:"plugins"`
}
type User struct{
	PasswordEncrypted string `json:"password_encrypted"`
	Permission PermissionToken `json:"permission"`
}
type UserRC struct{
	Groups map[string]map[string]User      `json:"user_groups"`
	Permissions map[string]PermissionToken `json:"group_permission"`
}

func InitUserRC(){
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
		return CastUser(name,group),nil
	}else{
		return UserContext{},errors.New("password incorrect")
	}
}
func CastUser(name string,group string)UserContext{
	permPrimary:=Users.Permissions[group]
	permSecondary:=Users.Groups[group][name].Permission
	permInternal:=permPrimary.Internal+","+permSecondary.Internal
	permExec:=permPrimary.Exec+","+permSecondary.Exec
	permPlugin:=permPrimary.Plugin+","+permSecondary.Plugin
	return UserContext{
		Name:       name,
		Group:      group,
		Permission: []string{permInternal,permExec,permPlugin},
	}
}
func PermissionFilter(commandname string,permission string) (bool,string){
	if commandname[0]=='_'{
		return true, "internal"
	}
	if permission == ""{
		return true,"All commands access are granted of this user group."
	}
	permissionToken:=strings.Split(permission,",")
	tmpPerm:=true
	tmpReason:="Permission restrictions of command '"+commandname+"' not found."
	for _,v:=range permissionToken{
		if v=="" {continue}
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
	if i==","&&e==","&&p==","{
		com=": all granted"
	}
	return "I("+i+"),E("+e+"),P("+p+")"+com
}