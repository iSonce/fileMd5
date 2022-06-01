package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
	"time"
)

var wg sync.WaitGroup

func main() {
	//get pathName from terminal
	var pathName string
	fmt.Scanln(&pathName)

	//call the getFileName function
	fileName, err := getFileName(pathName)
	if err != nil {
		fmt.Println(err)
		return
	}

	//parallel
	start := time.Now()
	wg.Add(len(fileName))
	for i := 0; i < len(fileName); i++ {
		go func(s string) {
			defer wg.Done()
			file, _ := os.ReadFile(s)
			sum := md5.Sum(file)
			fmt.Println(s, hex.EncodeToString(sum[:]))
		}(fileName[i])
	}
	wg.Wait()
	parallel := time.Since(start).Seconds()

	//serial
	start = time.Now()
	for i := 0; i < len(fileName); i++ {
		file, _ := os.ReadFile(fileName[i])
		sum := md5.Sum(file)
		fmt.Println(fileName[i], hex.EncodeToString(sum[:]))
	}
	serial := time.Since(start).Seconds()

	fmt.Println("parallel time: ", parallel)
	fmt.Println("serial time: ", serial)
}

func getFileName(pathName string) ([]string, error) {
	result := []string{}

	//try to read the path
	fis, err := ioutil.ReadDir(pathName)
	if err != nil {
		return result, err
	}

	for _, fi := range fis {
		fullname := pathName + "/" + fi.Name()
		if fi.IsDir() {
			//if fullname is a direction, get in this direction to read file name
			temp, err := getFileName(fullname)
			if err != nil {
				fmt.Println(err)
			} else {
				result = append(result, temp...)
			}
		} else {
			//if fullname is a file, try to read this file
			_, err := os.ReadFile(fullname)
			if err != nil {
				fmt.Println(err)
			} else {
				//succeed to read file, add it to result
				result = append(result, fullname)
			}
		}
	}
	return result, nil
}
