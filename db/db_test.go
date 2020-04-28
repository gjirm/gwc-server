package db

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/xujiajun/nutsdb"
)

var (
	db     *nutsdb.DB
	bucket string
)

func init() {
	opt := nutsdb.DefaultOptions
	fileDir := "./nutsdb_test"

	files, _ := ioutil.ReadDir(fileDir)
	for _, f := range files {
		name := f.Name()
		if name != "" {
			fmt.Println(fileDir + "/" + name)
			err := os.RemoveAll(fileDir + "/" + name)
			if err != nil {
				panic(err)
			}
		}
	}
	opt.Dir = fileDir
	opt.SegmentSize = 1024 * 1024 // 1MB
	db, _ = nutsdb.Open(opt)
	bucket = "bucketForString"
}

// Test test
func Test() {
	//insert
	// put()
	//read
	read()

	// //delete
	// delete()
	// //read
	// read()

	// //insert
	// put()
	// //read
	// read()

	// //update
	// put2()
	// //read
	// read()
}

func delete() {
	if err := db.Update(
		func(tx *nutsdb.Tx) error {
			key := []byte("name1")
			return tx.Delete(bucket, key)
		}); err != nil {
		log.Fatal(err)
	}
}

func put() {
	if err := db.Update(
		func(tx *nutsdb.Tx) error {
			key := []byte("name1")
			val := []byte("val1")
			return tx.Put(bucket, key, val, 0)
		}); err != nil {
		log.Fatal(err)
	}
}
func put2() {
	if err := db.Update(
		func(tx *nutsdb.Tx) error {
			key := []byte("name1")
			val := []byte("val2")
			return tx.Put(bucket, key, val, 0)
		}); err != nil {
		log.Fatal(err)
	}
}
func read() {
	if err := db.View(
		func(tx *nutsdb.Tx) error {
			key := []byte("name1")
			e, err := tx.Get(bucket, key)
			if err != nil {
				return err
			}
			fmt.Println("val:", string(e.Value))

			return nil
		}); err != nil {
		log.Println(err)
	}
}