package lib

import (
	"encoding/json"
	"fmt"
	"slices"
	"time"

	bolt "go.etcd.io/bbolt"
)

type TagNameList struct {
	TagNames []string
}

type TagsOfFile struct {
	What  []string `json:"what,omitempty"`
	Where []string `json:"where,omitempty"`
	Who   []string `json:"who,omitempty"`
	When  []string `json:"when,omitempty"`
}

var mainBucket = "photoFS"
var subBuckets = []string{"tagNames", "tags"}

func InitDB() (*bolt.DB, error) {

	db, err := bolt.Open("data.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, fmt.Errorf("could not open db, %w", err)
	}
	err = db.Update(func(tx *bolt.Tx) error {
		main, err := tx.CreateBucketIfNotExists([]byte(mainBucket))
		if err != nil {
			return fmt.Errorf("could not create main bucket(%v): %w", mainBucket, err)
		}
		for subBucket := range subBuckets {
			_, err = main.CreateBucketIfNotExists([]byte(string(subBucket)))
			if err != nil {
				return fmt.Errorf("could not create %v/%v bucket: %w", mainBucket, subBucket, err)
			}
		}
		return nil
	})
	fmt.Println("[DEBUG] DB Init is Done")

	return db, err
}

func verifyBucketName(bucketName string) error {

	if slices.Contains(subBuckets, bucketName) == false {
		fmt.Printf("subBuckets (%v) does not contain bucketName (%v)\n", subBuckets, bucketName)
		return fmt.Errorf("unknown bucketName: %v", bucketName)
	}

	return nil
}

func jsonToMap(jsonStr string) (map[string][]string, error) {

	result := make(map[string][]string)
	err := json.Unmarshal([]byte(jsonStr), &result)

	return result, err
}

func getKVfromDB(db *bolt.DB, bucketName string) (map[string]map[string][]string, error) {

	err := verifyBucketName(bucketName)
	if err != nil {
		return nil, fmt.Errorf("bucketName verification failed: %w", err)
	}

	kvMap := make(map[string]map[string][]string)

	err = db.View(func(tx *bolt.Tx) error {

		fmt.Println("mainBucket: ", mainBucket)
		fmt.Println("bucketName: ", bucketName)
		bucket := tx.Bucket([]byte(mainBucket)).Bucket([]byte(bucketName))
		fmt.Println("bucket: ", bucket)
		if bucket == nil {
			return fmt.Errorf("failed to get '%v/%v' bucket", mainBucket, bucketName)
		}

		c := bucket.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			jsonMap, err := jsonToMap(string(v))
			if err != nil {
				return fmt.Errorf("jsonToMap failed for string %v: %w", string(v), err)
			}
			kvMap[string(k)] = jsonMap
		}
		return nil
	})

	return kvMap, err
}

func addKVtoDB(db *bolt.DB, bucketName string, key string, value []byte) error {

	err := verifyBucketName(bucketName)
	if err != nil {
		return fmt.Errorf("bucketName verification failed: %w", err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		err := tx.Bucket([]byte(mainBucket)).Bucket([]byte(bucketName)).Put([]byte(key), []byte(value))
		if err != nil {
			return fmt.Errorf("could not add key (%v) with value (%v): %w", key, value, err)
		}
		return nil
	})

	return err
}

func GetTagNames(db *bolt.DB) (map[string][]string, error) {

	bucketName := "tagNames"
	kvMap, err := getKVfromDB(db, bucketName)
	if err != nil {
		return nil, fmt.Errorf("could not get '%v' kvMap from db: %w", bucketName, err)
	}

	tagNames := make(map[string][]string)

	for tagType, tagNameMap := range kvMap {
		tagNames[tagType] = tagNameMap["TagNames"]
	}

	return tagNames, err
}

func GetTags(db *bolt.DB) (map[string]map[string][]string, error) {

	bucketName := "tags"
	kvMap, err := getKVfromDB(db, bucketName)

	return kvMap, err
}

func AddTagNames(db *bolt.DB, tagNames map[string]TagNameList) error {

	for tagType, tagNameList := range tagNames {

		tagNamesBytes, err := json.Marshal(tagNameList)
		if err != nil {
			return fmt.Errorf("could not marshal tagNames of '%v' json: %w", tagType, err)
		} else if string(tagNamesBytes) == "{}" {
			return fmt.Errorf("tagNamesBytes marshal result of '%v' is an empty json: %v", tagType, string(tagNamesBytes))
		}

		err = addKVtoDB(db, "tagNames", tagType, tagNamesBytes)
		if err != nil {
			return fmt.Errorf("failed to add key-value to db: %w", err)
		}
	}

	return nil
}

func AddTestTagNames(db *bolt.DB) error {
	tagNamesToAdd := map[string]TagNameList{
		"What":  TagNameList{[]string{"party/birthday", "cake", "party/christmas", "party/halloween", "party/wedding", "sunrise", "sunset"}},
		"Where": TagNameList{[]string{"Germany/Frankfurt", "Germany/Berlin", "bar", "club", "restaurant"}},
	}
	err := AddTagNames(db, tagNamesToAdd)
	if err != nil {
		return fmt.Errorf("failed to add test tag names to db: %w", err)
	}

	return nil
}

func AddTagsOfFile(db *bolt.DB, srcFile string, tagsOfFile *TagsOfFile) error {

	tagsBytes, err := json.Marshal(tagsOfFile)

	if err != nil {
		return fmt.Errorf("could not marshal tags json: %w", err)
	} else if string(tagsBytes) == "{}" {
		return fmt.Errorf("tagsBytes marshal result is an empty json: %v", string(tagsBytes))
	}

	err = addKVtoDB(db, "tags", srcFile, tagsBytes)

	return err
}
