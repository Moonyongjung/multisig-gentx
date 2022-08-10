package util

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/logrusorgru/aurora"
)

func ConvertConfigParam(str string) []string {
	var strList []string
	if strings.Contains(str, "http") {
		str = strings.Replace(str, " ", "", -1)
		str = strings.Replace(str, "\r", "", -1)
		str = strings.Replace(str, "\n", "", -1)
		str = strings.Replace(str, ",", "", -1)
		str = strings.Replace(str, "\"", "", -1)
		strList = strings.Split(str, ":")
		strList[1] = strings.Join(strList[1:], ":")
	} else {
		str = strings.Replace(str, " ", "", -1)
		str = strings.Replace(str, "\r", "", -1)
		str = strings.Replace(str, "\n", "", -1)
		str = strings.Replace(str, ",", "", -1)
		str = strings.Replace(str, "\"", "", -1)
		strList = strings.Split(str, ":")
	}

	return strList
}

func ToString(value interface{}, defaultValue string) string {
	str := strings.TrimSpace(fmt.Sprintf("%v", value))
	if str == "" {
		return defaultValue
	} else {
		return str
	}
}

func FromStringToUint64(value string) uint64 {
	number, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		LogTool(err)
	}

	return number
}

func SaveJsonPretty(jsonByte []byte, saveTxPath string) error{
	var prettyJson bytes.Buffer	
	err := json.Indent(&prettyJson, jsonByte, "", "    ")
	if err != nil {
		LogErr(err)
	}

	err = ioutil.WriteFile(saveTxPath, prettyJson.Bytes(), 0660)
	if err != nil {
		LogErr(err)
	}

	return nil
}

func LogTool(log ...interface{}) {
	str := ToString(log, "")
	fmt.Println(aurora.White("GenTxTool ").String() + str)
}

func LogErr(log ...interface{}) error {
	str := ToString(log, "")
	fmt.Println(aurora.Red("Err       ").String() + str)

	return errors.New(str)
}
