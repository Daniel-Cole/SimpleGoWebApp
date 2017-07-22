package utils

import "github.com/boltdb/bolt"

type Context struct {
	DBConn *bolt.DB
	DBTimeout int
	DBBucketEnv []byte
	DBBucketApp []byte
}