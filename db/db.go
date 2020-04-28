package db

import (
	"fmt"
	"io/ioutil"
	"log"
	//"os"

	"github.com/xujiajun/nutsdb"
	"github.com/sirupsen/logrus"
	
	config "jirm.cz/gwc-server/config"
)

var (
	db     		*nutsdb.DB
	bucket 		string
	getvalue 	string
)

// InitDB exported
func InitDB(log *logrus.Logger, config config.Configs) {
	opt := nutsdb.DefaultOptions
	fileDir := "./" + config.Database.Name
	var noDb bool = true
	files, _ := ioutil.ReadDir(fileDir)
	for _, f := range files {
		name := f.Name()
		if name != "" {
			log.Info("Using existing DB: " + fileDir + "/" + name)
			noDb = false
			//err := os.RemoveAll(fileDir + "/" + name)
			// if err != nil {
			// 	panic(err)
			// }
		} else {
			log.Info("Creating new DB: " + fileDir + "/" + name)
		}
	}
	opt.Dir = fileDir
	opt.SegmentSize = 1024 * 1024 // 1MB
	db, _ = nutsdb.Open(opt)
	if noDb {
		PutValue("options", "ipcounter", config.Wireguard.WgIPBegin)
		PutValue("options", "ipspace", config.Wireguard.WgIPSpace)
		PutValue("options", "server", config.Wireguard.WgServer)
	}
	//bucket = "bucketForString"
}

// DbTest test
func DbTest() {
	// //insert
	// put()
	// //read
	// read()

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
	//read
	//read()

	//PutValue("pubkeys", "jiri.matejicek","ybdJGJZGjkgFjkFKJghjb")
	PutValue("pubkeys", "john","ahoj")
	PutValue("pubkeys", "mary","svete")
	value := GetValue("pubkeys","jiri.matejicek")
	fmt.Println("JM Pubkey: " + value)
	value = GetValue("pubkeys","jhn")
	fmt.Println("JM Pubkey: " + value)

}

// PutValue putvalue
func PutValue(bucket, keyname, value string) {

	if err := db.Update(
		func(tx *nutsdb.Tx) error {
			key := []byte(keyname)
			val := []byte(value)
			return tx.Put(bucket, key, val, 0)
		}); err != nil {
		log.Fatal(err)
	}
}



// GetValue getvaluae
func GetValue(bucket, keyname string) (string) {
	var getvalue string
	if err := db.View(
		func(tx *nutsdb.Tx) error {
			key := []byte(keyname)
			e, err := tx.Get(bucket, key)
			if err != nil {
				return err
			}
			getvalue = string(e.Value)
			return nil

		}); err != nil {
		log.Println(err)
		return ""
	}
	return getvalue
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