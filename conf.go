package tools

import (
	"encoding/json"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"runtime"
)

type GetName struct {
	Name string
}

type PlugRequest struct {
	Body     string
	Header   http.Header
	Form     url.Values
	PostForm url.Values
	Url      string
	Method   string
	HeadVals map[string]string
	Status   int
}

type ReturnMsg struct {
	Method string
	Err    string
	Plugin string
	Email  string
}

type Message struct {
	Method    string
	Name      string
	Email     string
	Activated string
	Sam       string
	Password  string
}

func ReadMergeConf(out interface{}, filename string) ([]byte, error) {
	d, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Println("No Configuration file found in ~/.Config/nanocloud, now looking in /etc/nanocloud")
		return nil, err
	}
	return d, nil
}

func WriteConf(in interface{}, filename string) error {
	d, err := yaml.Marshal(in)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, d, 0644)
}

func InitConf(in interface{}) []byte {
	name, err := json.Marshal(in)
	if err != nil {
		log.Println(err)
	}
	var gname GetName
	err = json.Unmarshal(name, &gname)

	if err != nil {
		log.Println(err)
	}
	usr, err := user.Current()
	if err != nil {
		log.Println(err)
	}
	home := usr.HomeDir
	f := gname.Name + ".yaml"
	if runtime.GOOS == "linux" {
		d := home + "/.config/nanocloud/" + gname.Name + "/"
		err := os.MkdirAll(d, 0755)
		if err == nil {
			f = d + f
		} else {
			log.Println(err)
		}
	}

	Newconf, err := ReadMergeConf(in, f)
	if err != nil {
		alt := "/etc/nanocloud/" + gname.Name + "/" + gname.Name + ".yaml"
		Newconf, err = ReadMergeConf(in, alt)
		if err != nil {
			log.Println("No Configuration file found in /etc/nanocloud, using default Configuration")
		}
	}
	if err := WriteConf(in, f); err != nil {
		log.Println(err)
	}
	return Newconf
}
