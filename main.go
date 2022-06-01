package main

import (
	"crypto/md5"
	"encoding/hex"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
	"time"
)

var wg sync.WaitGroup

func main() {
	//get pathName from terminal
	//example: go run main.go -p pathName
	var pathName string
	flag.StringVar(&pathName, "p", "", "pathName")
	flag.Parse()

	//call the getFileName function
	fileName, err := getFileName(pathName)
	if err != nil {
		fmt.Println(err)
		return
	}
	length := len(fileName)

	//parallel run
	begin := time.Now()
	wg.Add(length)
	for i := 0; i < length; i++ {
		//create a new routine
		go getMD5(fileName[i], true)
	}
	wg.Wait()
	parallel := time.Since(begin).Seconds()

	//serial run
	begin = time.Now()
	for i := 0; i < length; i++ {
		getMD5(fileName[i], false)
	}
	serial := time.Since(begin).Seconds()

	fmt.Println("parallel time: ", parallel)
	fmt.Println("serial time: ", serial)
}

func getMD5(fileName string, isParallel bool) {
	if isParallel {
		defer wg.Done()
	}
	file, _ := os.ReadFile(fileName)
	sum := md5.Sum(file)
	fmt.Println(fileName, hex.EncodeToString(sum[:]))
}

func getFileName(pathName string) ([]string, error) {
	result := []string{}

	//try to read the path, if fail, return empty result
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
				//add the result of getFileName to result, it may be empty
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
