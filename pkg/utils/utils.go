package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"
)

func CheckingFileExistence(logPath, filename string) (*os.File, error) {
	if err := os.MkdirAll(logPath, 0777); err != nil {
		return nil, fmt.Errorf("ошибка создания Log-директории: %v", err)
	}

	logFileName, err := getLastDailyFile(logPath, filename)
	if err != nil {
		panic(err)
	}

	var file *os.File

	if isMoreThan100MB(logFileName) {
		logFileName = fmt.Sprintf("%s_%d.log", logPath+filename+"_"+time.Now().Format("02_01_2006"), getNextFileNum(logPath, filename))
		file, err = os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0777)
		if err == nil {
			writeLogHeader(file, filename)
		}
	} else {
		file, err = os.OpenFile(logFileName, os.O_WRONLY|os.O_APPEND, 0777)
		if err != nil {
			if os.IsNotExist(err) {
				file, err = os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0777)
				if err == nil {
					writeLogHeader(file, filename)
				}
			}
		}
	}

	return file, nil
}

func writeLogHeader(file *os.File, appName string) {
	appName = "AppName = " + appName
	logDate := "LogDate = " + time.Now().Format("02.01.2006")
	format := "Format = [\"Date\",\"Time\",\"M\",\"Level\",\"Message\"]"
	header := fmt.Sprintf("{\n\t%s\n\t%s\n\t%s\n}", appName, logDate, format)

	if _, err := file.Write([]byte(header + "\n")); err != nil {
		fmt.Printf("Ошибка записи заголовка лога: %v\n", err)
	}
}

func findAllDailyFiles(path, filename string) ([]string, error) {
	currentDate := time.Now().Format("02_01_2006")
	prefix := filepath.Join(path, filename+"_"+currentDate+"*", "/")
	files, err := filepath.Glob(prefix)
	if err != nil {
		return nil, err
	}

	return files, nil
}

func getNextFileNum(path, filename string) int {
	files, err := findAllDailyFiles(path, filename)
	if err != nil {
		return 0
	}

	return len(files)
}

func getLastDailyFile(path, filename string) (string, error) {
	files, err := findAllDailyFiles(path, filename)
	if err != nil {
		return "", fmt.Errorf("не найден последний log-файл за сутки")
	}

	if len(files) == 0 {
		return fmt.Sprintf(filepath.Join(path, filename)+"_%s.log", time.Now().Format("02_01_2006")), nil
	}

	sort.Strings(files)

	return files[len(files)-1], nil
}

func isMoreThan100MB(filename string) bool {
	info, err := os.Stat(filename)
	if err != nil {
		return false
	}

	return info.Size() > 100*1024*1024
}
