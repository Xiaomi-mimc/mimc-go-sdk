package strutil

import (
	"bytes"
	"container/list"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/Xiaomi-mimc/mimc-go-sdk/util/log"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
	"time"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

var logger = log.GetLogger()

func Substring(str *string, from int) string {
	source := []rune(*str)
	length := len(source)
	if from < 0 {
		from = length + from
	} else if from == 0 {
		return *str
	}
	return string(source[from:length])
}

func Substr(str *string, from, to int) string {
	source := []rune(*str)
	if from >= 0 && from <= to {
		length := len(source)
		if to >= length {
			to = length - 1
		}
		return string(source[from:to])
	}
	return ""
}

func Concat(str1, str2 *string) string {
	newStr := bytes.Buffer{}
	newStr.WriteString(*str1)
	newStr.WriteString(*str2)
	return newStr.String()
}

func ConcatStrs(strs ...*string) string {
	size := len(strs)
	newStr := bytes.Buffer{}
	for i := 0; i < size; i++ {
		newStr.WriteString(*strs[i])
	}
	return newStr.String()
}

func RandomStrWithLength(size int) string {
	randomStr := make([]rune, size)
	length := len(letters)
	rand.Seed(time.Now().UnixNano())
	for i := range randomStr {
		randomStr[i] = letters[rand.Intn(length)]
	}
	return string(randomStr)
}

func Bytes(str *string) []byte {
	bytes := &bytes.Buffer{}
	bytes.WriteString(*str)
	return bytes.Bytes()
}
func ConcatStrsByStr(strs *list.List, str *string) string {
	first := true
	newStr := ""
	for key := strs.Front(); key != nil; key = key.Next() {
		//val = key.Value
		val, _ := key.Value.(string)
		//val = &string{key.Value.(string)}
		if !first {
			newStr = ConcatStrs(&newStr, str, &val)
		} else {
			first = false
			newStr = ConcatStrs(&newStr, &val)
		}
	}
	return newStr
}

func Sha1(str *string) *string {
	sha := sha1.New()
	sha.Write([]byte(*str))
	shaStr := sha.Sum(nil)
	digest := base64.StdEncoding.EncodeToString(shaStr)
	return &digest
}

func CreateFile(pathfile *string) (bool, error) {
	pathFile, err := os.Create(*pathfile)
	if err != nil {
		return false, err
	} else {
		defer pathFile.Close()
		return true, nil
	}
}
func CrateDir(path *string) (bool, error) {
	err := os.Mkdir(*path, os.ModePerm)
	if err != nil {
		return false, err
	}
	return true, nil
}
func PathExists(path *string) (bool, error) {
	_, err := os.Stat(*path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func SynchronizeResource(root, dir, file, key, value *string) *string {
	rot := Substr(root, 0, strings.LastIndex(*root, "/"))
	resourcePath := rot + *dir
	resourceText := resourcePath + *file
	//fmt.Printf("rot:%v, pat:%v, text:%v\n", rot, resourcePath, resourceText)
	result, _ := PathExists(&resourcePath)
	if result {
		result, _ = PathExists(&resourceText)
		if !result {
			CreateFile(&resourceText)
		}
	} else {
		result, _ = CrateDir(&resourcePath)
		if result {
			CreateFile(&resourceText)
		} else {
			fmt.Printf("1.\n")
			panic("can not create dir.")
		}
	}
	return SynchronizeWithFile(key, value, &resourceText)
}

func Transfer(origin []interface{}) []string {
	var target []string
	for _, value := range origin {
		target = append(target, value.(string))
	}
	return target
}

func check(root, dir, file *string) (bool, *string) {
	rot := Substr(root, 0, strings.LastIndex(*root, "/"))
	resourcePath := rot + *dir
	resourceText := resourcePath + *file
	result, _ := PathExists(&resourcePath)
	if result {
		result, _ = PathExists(&resourceText)
		if !result {
			result, _ = CreateFile(&resourceText)
			return result, &resourceText
		}
		return true, &resourceText
	} else {
		result, _ = CrateDir(&resourcePath)
		if result {
			result, _ = CreateFile(&resourceText)
			return result, &resourceText
		} else {
			fmt.Printf("can not create dir %s.\n", *file)
			return false, nil
		}
	}
	return false, nil
}

func FlushUserInfo(root, dir, file *string, userInfo map[string]interface{}) bool {
	result, path := check(root, dir, file)
	if !result {
		logger.Info("check file%v failed", *file)
		return false
	}
	f, err := os.OpenFile(*path, os.O_RDWR, 0666)
	defer f.Close()
	if err != nil {
		logger.Info("open file:%s failed，%v", *file, err)
		return false
	}
	data, jsonErr := json.Marshal(userInfo)
	if jsonErr != nil {
		logger.Info("marshal user Info failed")
		return false
	}
	size, ferr := f.WriteAt(data, 0)
	if ferr != nil {
		logger.Warn("write user Info failed")
		return false
	} else if size != len(data) {
		logger.Info("write user Info length error")
		return false
	}
	return true
}

func FetchUserInfo(root, dir, file *string, userInfo map[string]interface{}) bool {
	result, path := check(root, dir, file)
	if !result {
		return false
	}
	f, err := os.OpenFile(*path, os.O_RDWR, 0666)
	defer f.Close()
	if err != nil {
		return false
	}
	data, _ := ioutil.ReadAll(f)
	logger.Info("path: " + *path + " data: " + string(data))
	var localInfo map[string]interface{} = make(map[string]interface{})
	if len(data) != 0 {
		decoder := json.NewDecoder(strings.NewReader(string(data)))
		decoder.UseNumber()
		if err = decoder.Decode(&localInfo); err != nil {
			panic(err)
		}
		updated := false
		for key, value := range localInfo {
			userInfo[key] = value
			updated = true
		}
		return updated
	}
	return false
}

func SynchronizeWithFile(key, value, file *string) *string {
	f, err := os.OpenFile(*file, os.O_RDWR, 0666)
	defer f.Close()
	if err != nil {
		return nil
	}
	data, _ := ioutil.ReadAll(f)
	var kvs map[string]interface{} = make(map[string]interface{})

	if len(data) == 0 {
		kvs[*key] = *value
	} else {
		if err = json.Unmarshal(data, &kvs); err != nil {
			panic(err)
		}
		val, ok := kvs[*key]
		//fmt.Printf("kvs: %v\n", kvs)
		if ok {
			old := val.(string)
			return &(old)
		} else {
			kvs[*key] = *value
		}
	}
	data, _ = json.Marshal(kvs)
	size, ferr := f.WriteAt(data, 0)
	if ferr != nil {
		panic(ferr)
	} else if size != len(data) {
		panic("write length err.")
	}
	return value
}
