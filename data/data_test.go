package data

import (
	"fmt"
	"log"
	"testing"
)

const folderPath = "F:/Curation/Music/ytrip"

func createElementFiles() {
	FoldersAsTags([]string{}, folderPath)
}

var d Data

func TestGenerateSets(t *testing.T) {
	// This is a temporary fix, automated testing will be implemented later on.
	if folderPath == "" {
		log.Fatalln("Folder path isn't set. To run tests, describe folder to use.")
	}
	d.Init(10)
	createElementFiles()
	d.Generate()
	res := d.Stringify()
	fmt.Println(res)
}

func TestWriteSets(t *testing.T) {

}

func TestLoadSets(t *testing.T) {

}

func TestQuery(t *testing.T) {
	//	t1 := ""
	//	t2 := ""
	//	sets := GetArray(&d, []string{t1, t2})
	//	res := Junction(sets)
	//	if res == nil {
	//		panic("junction is empty")
	//	}
	//	strs := res.Stringify()
	//	for _, v := range strs {
	//		fmt.Println(v)
	//	}
}
