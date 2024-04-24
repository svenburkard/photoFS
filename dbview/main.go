package main

import (
    "fmt"
    "log"
    "time"

    bolt "go.etcd.io/bbolt"
)

var mainBucket = "photoFS"
var subBuckets = []string{"tagNames", "tags"}

func main() {
    db, err := bolt.Open("data.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
    if err != nil {
      log.Fatal("could not open db, %w", err)
    }

    defer db.Close()

    err = db.View(func(tx *bolt.Tx) error {
        mainBucket := tx.Bucket([]byte(mainBucket))
        if mainBucket == nil {
            return fmt.Errorf("main bucket %v not found", mainBucket)
        }

        for _, subBucketName := range subBuckets {
            fmt.Printf("Contents of sub-bucket: %s\n", subBucketName)
            subBucket := mainBucket.Bucket([]byte(subBucketName))
            if subBucket == nil {
                fmt.Printf("Sub-bucket '%s' not found\n", subBucketName)
                continue
            }

            err := subBucket.ForEach(func(k, v []byte) error {
                fmt.Printf("Key: %s, Value: %s\n", k, v)
                return nil
            })
            if err != nil {
                return fmt.Errorf("error reading sub-bucket '%s': %v", subBucketName, err)
            }
            fmt.Println()
        }
        return nil
    })

    if err != nil {
        log.Fatal(err)
    }
}
