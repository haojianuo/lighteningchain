package wallet

import (
	"bytes"
	"encoding/gob"
	"errors"
	"io/ioutil"
	"lighteningchain/constcoe"
	"lighteningchain/utils"
	"os"
	"path/filepath"
	"strings"
)

type RefList map[string]string

func (r *RefList) Save() {
	filename := constcoe.WalletsRefList + "ref_list.data"
	var content bytes.Buffer
	encoder := gob.NewEncoder(&content)
	err := encoder.Encode(r)
	utils.Handle(err)
	err = os.WriteFile(filename, content.Bytes(), 0644)
	utils.Handle(err)
}

func (r *RefList) Update() {
	err := filepath.Walk(constcoe.Wallets, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if f.IsDir() {
			return nil
		}
		fileName := f.Name()
		if strings.Compare(fileName[len(fileName)-4:], ".wlt") == 0 {
			_, ok := (*r)[fileName[:len(fileName)-4]]
			if !ok {
				(*r)[fileName[:len(fileName)-4]] = ""
			}
		}
		return nil
	})
	utils.Handle(err)
}

func LoadRefList() *RefList {
	filename := constcoe.WalletsRefList + "ref_list.data"
	var refList RefList
	if utils.FileExists(filename) {
		fileContent, err := ioutil.ReadFile(filename)
		utils.Handle(err)
		decoder := gob.NewDecoder(bytes.NewBuffer(fileContent))
		err = decoder.Decode(&refList)
		utils.Handle(err)
	} else {
		refList = make(RefList)
		refList.Update()
	}
	return &refList
}

func (r *RefList) BindRef(address, refName string) {
	(*r)[address] = refName
}

func (r *RefList) FindRef(refname string) (string, error) {
	temp := ""
	for key, val := range *r {
		if val == refname {
			temp = key
			break
		}
	}
	if temp == "" {
		err := errors.New("the refname is not found")
		return temp, err
	}
	return temp, nil
}
