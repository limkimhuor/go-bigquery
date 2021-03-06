package utilities

import (
	"encoding/json"
	"fmt"
	"github.com/euclid1990/go-bigquery/configs"
	"github.com/gocarina/gocsv"
	"github.com/joho/godotenv"
	"os"
	"path/filepath"
	"reflect"
	"time"
)

func LoadEnv(file string) {
	if file == "" {
		file = ".env"
	}
	err := godotenv.Load(file)
	if err != nil {
		Logf(configs.LOG_CRITICAL, "Error loading %v file", file)
	}
}

func Random(min, max int) int {
	rndSrc := CreateRndSrc()
	return rndSrc.Intn(max-min) + min
}

func RandomDateTimeBetween(afterObj time.Time, beforeObj time.Time) time.Time {
	afterUnix := afterObj.Unix()
	beforeUnix := beforeObj.Unix()

	diff := beforeUnix - afterUnix

	rndSrc := CreateRndSrc()
	sec := rndSrc.Int63n(diff) + afterUnix

	return time.Unix(sec, 0)
}

func GetUserJsonFilePath() string {
	return fmt.Sprintf(configs.DATA_FORMAT_FILE_NAME, configs.DATA_INPUT_PATH+configs.DATA_INPUT_USER, configs.DATA_TYPE_JSON)
}

func GetAccessJsonFilePath() string {
	return fmt.Sprintf(configs.DATA_FORMAT_FILE_NAME, configs.DATA_INPUT_PATH+configs.DATA_INPUT_ACCESS, configs.DATA_TYPE_JSON)
}

func WriteCsv(file string, data interface{}) (bool, error) {
	csvFile, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer csvFile.Close()
	if err := gocsv.MarshalFile(data, csvFile); err != nil {
		return false, err
	}
	return true, nil
}

// BigQuery expects newline-delimited JSON files to contain a single record per line.
func WriteJson(file string, data interface{}) (bool, error) {
	jsonFile, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer jsonFile.Close()

	breakLine := []byte("\n")

	switch reflect.TypeOf(data).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(data)
		for i := 0; i < s.Len(); i++ {
			buff, errJson := json.Marshal(s.Index(i).Interface())
			if errJson != nil {
				panic(errJson)
			}
			buff = append(buff, breakLine...)
			if _, err := jsonFile.Write(buff); err != nil {
				panic(err)
			}
		}
	}

	return true, nil
}

// Get directory total size (MB)
func DirSize(path string) (float64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			size += info.Size()
		}
		return err
	})
	// Ignore errors
	if err != nil {
		return 0, err
	}
	sizeMB := float64(size) / 1024.0 / 1024.0
	return sizeMB, err
}

func TimeTrack(start time.Time) float64 {
	elapsed := time.Since(start)
	return float64(elapsed) / float64(time.Second)
}
