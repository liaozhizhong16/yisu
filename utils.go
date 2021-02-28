package yisu

import (
	"encoding/json"
	"errors"
	"github.com/dselans/dmidecode"
	"os"
)

/**
* @Time    : 21/2/2 下午4:38
* @Author  : liaozz
*
* 自我介绍一下
 */

// 判断所给路径文件/文件夹是否存在
func Exists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

// 判断所给路径是否为文件夹
func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// 判断所给路径是否为文件
func IsFile(path string) bool {
	return !IsDir(path)
}

func get_uuid() (string, error) {
	var uuid string
	dmi := dmidecode.New()
	if err := dmi.Run(); err != nil {
		return "", err
	}
	for _, v := range dmi.Data {
		for _, d := range v {
			for k, c := range d {
				if k == "UUID" {
					uuid = c
					return uuid, nil
				}

			}
		}
	}

	return "", errors.New("not found uuid")
}

func save_json(path string, data *interface{}) {
	if Exists(path) {
		//文件存在
	} else {
		//文件不存在
		filePtr, err := os.Create(path)
		if err == nil {
			defer filePtr.Close()
			encoder := json.NewEncoder(filePtr)
			err := encoder.Encode(data)
			if err != nil {
				panic(err)
			}
		}
	}

}
func read_json(path string) (resd *ResponseData, err error) {
	if Exists(path) {
		//文件存在
		filePtr, err := os.Open(path)
		if err != nil {
			panic(err)
		} else {
			defer filePtr.Close()
			// 创建json解码器
			decoder := json.NewDecoder(filePtr)
			resd := &ResponseData{}
			err = decoder.Decode(resd)
			if err != nil {
				return nil, err
			} else {
				return resd, nil
			}
		}

	} else {
		return nil, errors.New("file is not found")
	}
}
