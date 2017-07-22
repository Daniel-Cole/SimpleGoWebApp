package database

import
(
	"github.com/boltdb/bolt"
	"github.com/daniel-cole/SimpleGoWebApp/app/log"
	"time"
	"fmt"
)


func InitDB(dbName string) (dbConn *bolt.DB) {
	log.LogInfo.Printf("Initialising BoltDB using DB: %s", dbName)

	db, err := bolt.Open(dbName, 0600, &bolt.Options{Timeout: 10 * time.Second})
	if err != nil {
		log.LogFatal("", err)
	}
	return db
}

func InsertDBValue(dbConn *bolt.DB, bucket []byte, key []byte, value []byte) error {
	err := dbConn.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(bucket)
		if err != nil {
			return err
		}

		err = bucket.Put(key, value)
		if err != nil {
			return err
		}
		return nil
	})

	//if error is anything but nil then something went wrong
	return err
}

func ReadDBValue(dbConn *bolt.DB, bucket []byte, key []byte) ([]byte, error) {

	var val []byte

	err := dbConn.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(bucket)
		if bucket == nil {
			return fmt.Errorf("Bucket %q not found!", bucket)
		}

		val = bucket.Get(key)

		return nil
	})

	//if error is anything but nil then something went wrong
	return val, err
}

func ReadAllDBValues(dbConn *bolt.DB, bucket []byte) (map[string]string, error) {

	values := make(map[string]string)

	err := dbConn.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte(bucket))

		b.ForEach(func(k, v []byte) error {
			values[string(k)] = string(v)
			return nil
		})
		return nil
	})

	//if error is anything but nil then something went wrong
	return values, err
}

func DeleteDBValue(dbConn *bolt.DB, bucket []byte, key []byte) error {

	err := dbConn.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(bucket)
		if err != nil {
			return err
		}

		err = bucket.Delete(key)

		if err != nil {
			return err
		}
		return nil
	})

	//if error is anything but nil then something went wrong
	return err
}

//TODO: implement backup function
func BackupDB(dbConn *bolt.DB){

}
