package util

import (
	"bufio"
	"io"
	"os"
	"strings"
	"sync"
)

type Config map[string]string

var once sync.Once
var instance *singleton

type singleton struct {
	config Config
}

func GetConfig() *singleton {
	once.Do(func() {
		instance = &singleton{}
	})
	return instance
}

func (s *singleton) Get(key string) string {
	return s.config[key]
}

func (s *singleton) Read(filePath string) error {
	if len(s.config) == 0 {
		s.config = Config{}
	}

	file, err := os.Open(filePath)
	if err != nil {
		LogTool(err)
		return err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	for {
		str, err := reader.ReadString('\n')

		if strings.Contains(str, ":") {
			kvList := ConvertConfigParam(str)
			s.config[kvList[0]] = kvList[1]
		}

		if err == io.EOF {
			break
		}
	}

	return nil
}
