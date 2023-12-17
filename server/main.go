package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
)

func exists(path string) bool {
	if _, err := os.Stat(path); err == nil {
		// exist
		return true
	}
	// not exist
	return false
}

func createSymLinks(fileSrc string, folderDests []string, fileDest string) {

	for i := range folderDests {

		err := os.MkdirAll(folderDests[i], os.ModePerm)
		if err != nil {
			log.Println(err)
		}

		fileLink := folderDests[i] + fileDest

		if exists(fileLink) == false {
			fmt.Println("[INFO] symlink needs to be created: " + fileLink)
			err = os.Symlink(fileSrc, fileLink)
			if err != nil {
				log.Fatal(err)
			}
		}

	}

}

func getFileName(file string) string {
	fileExtensions := "jpg|jpeg|png|avi|mp4"
	rgx := regexp.MustCompile(`[^/]+\.(` + fileExtensions + `)`)
	fileName := rgx.FindString(file)
	fmt.Println("fileName: ", fileName)

	return fileName
}

func getDateTags(fileDict interface{}) map[string]string {

	dateTags := make(map[string]string)

	dateTags["year"] = fileDict.(map[string]interface{})["date"].(map[string]interface{})["year"].(string)
	dateTags["month"] = fileDict.(map[string]interface{})["date"].(map[string]interface{})["month"].(string)
	dateTags["day"] = fileDict.(map[string]interface{})["date"].(map[string]interface{})["day"].(string)

	fmt.Println("===>", "date/year:", dateTags["year"])
	fmt.Println("===>", "date/month:", dateTags["month"])
	fmt.Println("===>", "date/day:", dateTags["day"])

	return dateTags
}

func getTimeTags(fileDict interface{}) map[string]string {

	timeTags := make(map[string]string)

	timeTags["hour"] = fileDict.(map[string]interface{})["time"].(map[string]interface{})["hour"].(string)
	timeTags["minute"] = fileDict.(map[string]interface{})["time"].(map[string]interface{})["minute"].(string)
	timeTags["second"] = fileDict.(map[string]interface{})["time"].(map[string]interface{})["second"].(string)

	fmt.Println("===>", "time/hour:", timeTags["hour"])
	fmt.Println("===>", "time/minute:", timeTags["minute"])
	fmt.Println("===>", "time/second:", timeTags["second"])

	return timeTags
}

func getTagMap() map[string]interface{} {

	tagMapFile := "tag_map.json"

	jsonFile, err := os.Open(tagMapFile)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("[DEBUG] Successfully opened " + tagMapFile)
	defer jsonFile.Close()

	jsonStr, _ := ioutil.ReadAll(jsonFile)

	tagMap := map[string]interface{}{}
	err = json.Unmarshal([]byte(jsonStr), &tagMap)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(tagMap)

	return tagMap
}

func main() {

	fileDestPrefix := "/tmp/test/dst"
	fmt.Printf("fileDestPrefix: %s\n", fileDestPrefix)

	tagMap := getTagMap()

	for file := range tagMap {

		fmt.Println("file:", file, "=>", "tag-map:", tagMap[file])

		fileName := getFileName(file)

		dateTags := getDateTags(tagMap[file])
		timeTags := getTimeTags(tagMap[file])

		tag_event := tagMap[file].(map[string]interface{})["event"].(string)
		tag_person := tagMap[file].(map[string]interface{})["person"].(string)
		tag_place := tagMap[file].(map[string]interface{})["place"].(string)

		fmt.Println("===>", "event:", tag_event)
		fmt.Println("===>", "person:", tag_person)
		fmt.Println("===>", "place:", tag_place)

		datePrefix := dateTags["year"] + dateTags["month"] + dateTags["day"] + "-"
		timePrefix := timeTags["hour"] + timeTags["minute"] + timeTags["second"] + "_-_"

		specialFolderPrefix := "0000-#-"

		fileLinkFolders := []string{
			///////////////////////////////////// when
			fileDestPrefix + "/" + specialFolderPrefix + "when/" + dateTags["year"] + "/" + dateTags["month"] + "/" + dateTags["day"] + "/",
			fileDestPrefix + "/" + specialFolderPrefix + "when/" + dateTags["year"] + "/" + dateTags["month"] + "/" + specialFolderPrefix + "all/",
			fileDestPrefix + "/" + specialFolderPrefix + "when/" + dateTags["year"] + "/" + specialFolderPrefix + "all/",

			//////////  when/what
			fileDestPrefix + "/" + specialFolderPrefix + "when/" + dateTags["year"] + "/" + dateTags["month"] + "/" + dateTags["day"] + "/" + specialFolderPrefix + "what/" + tag_event + "/" + specialFolderPrefix + "all/",
			fileDestPrefix + "/" + specialFolderPrefix + "when/" + dateTags["year"] + "/" + dateTags["month"] + "/" + specialFolderPrefix + "what/" + tag_event + "/" + specialFolderPrefix + "all/",
			fileDestPrefix + "/" + specialFolderPrefix + "when/" + dateTags["year"] + "/" + specialFolderPrefix + "what/" + tag_event + "/" + specialFolderPrefix + "all/",
			/////       when/what/where
			fileDestPrefix + "/" + specialFolderPrefix + "when/" + dateTags["year"] + "/" + dateTags["month"] + "/" + dateTags["day"] + "/" + specialFolderPrefix + "what/" + tag_event + "/" + specialFolderPrefix + "where/" + tag_place + "/",
			fileDestPrefix + "/" + specialFolderPrefix + "when/" + dateTags["year"] + "/" + dateTags["month"] + "/" + specialFolderPrefix + "what/" + tag_event + "/" + specialFolderPrefix + "where/" + tag_place + "/",
			fileDestPrefix + "/" + specialFolderPrefix + "when/" + dateTags["year"] + "/" + specialFolderPrefix + "what/" + tag_event + "/" + specialFolderPrefix + "where/" + tag_place + "/",
			/////       when/what/who
			fileDestPrefix + "/" + specialFolderPrefix + "when/" + dateTags["year"] + "/" + dateTags["month"] + "/" + dateTags["day"] + "/" + specialFolderPrefix + "what/" + tag_event + "/" + specialFolderPrefix + "who/" + tag_person + "/",
			fileDestPrefix + "/" + specialFolderPrefix + "when/" + dateTags["year"] + "/" + dateTags["month"] + "/" + specialFolderPrefix + "what/" + tag_event + "/" + specialFolderPrefix + "who/" + tag_person + "/",
			fileDestPrefix + "/" + specialFolderPrefix + "when/" + dateTags["year"] + "/" + specialFolderPrefix + "what/" + tag_event + "/" + specialFolderPrefix + "who/" + tag_person + "/",

			//////////  when/where
			fileDestPrefix + "/" + specialFolderPrefix + "when/" + dateTags["year"] + "/" + dateTags["month"] + "/" + dateTags["day"] + "/" + specialFolderPrefix + "where/" + tag_place + "/" + specialFolderPrefix + "all/",
			fileDestPrefix + "/" + specialFolderPrefix + "when/" + dateTags["year"] + "/" + dateTags["month"] + "/" + specialFolderPrefix + "where/" + tag_place + "/" + specialFolderPrefix + "all/",
			fileDestPrefix + "/" + specialFolderPrefix + "when/" + dateTags["year"] + "/" + specialFolderPrefix + "where/" + tag_place + "/" + specialFolderPrefix + "all/",
			/////       when/where/what
			fileDestPrefix + "/" + specialFolderPrefix + "when/" + dateTags["year"] + "/" + dateTags["month"] + "/" + dateTags["day"] + "/" + specialFolderPrefix + "where/" + tag_place + "/" + specialFolderPrefix + "what/" + tag_event + "/",
			fileDestPrefix + "/" + specialFolderPrefix + "when/" + dateTags["year"] + "/" + dateTags["month"] + "/" + specialFolderPrefix + "where/" + tag_place + "/" + specialFolderPrefix + "what/" + tag_event + "/",
			fileDestPrefix + "/" + specialFolderPrefix + "when/" + dateTags["year"] + "/" + specialFolderPrefix + "where/" + tag_place + "/" + specialFolderPrefix + "what/" + tag_event + "/",
			/////       when/where/who
			fileDestPrefix + "/" + specialFolderPrefix + "when/" + dateTags["year"] + "/" + dateTags["month"] + "/" + dateTags["day"] + "/" + specialFolderPrefix + "where/" + tag_place + "/" + specialFolderPrefix + "who/" + tag_person + "/",
			fileDestPrefix + "/" + specialFolderPrefix + "when/" + dateTags["year"] + "/" + dateTags["month"] + "/" + specialFolderPrefix + "where/" + tag_place + "/" + specialFolderPrefix + "who/" + tag_person + "/",
			fileDestPrefix + "/" + specialFolderPrefix + "when/" + dateTags["year"] + "/" + specialFolderPrefix + "where/" + tag_place + "/" + specialFolderPrefix + "who/" + tag_person + "/",

			//////////  when/who
			fileDestPrefix + "/" + specialFolderPrefix + "when/" + dateTags["year"] + "/" + dateTags["month"] + "/" + dateTags["day"] + "/" + specialFolderPrefix + "who/" + tag_person + "/" + specialFolderPrefix + "all/",
			fileDestPrefix + "/" + specialFolderPrefix + "when/" + dateTags["year"] + "/" + dateTags["month"] + "/" + specialFolderPrefix + "who/" + tag_person + "/" + specialFolderPrefix + "all/",
			fileDestPrefix + "/" + specialFolderPrefix + "when/" + dateTags["year"] + "/" + specialFolderPrefix + "who/" + tag_person + "/" + specialFolderPrefix + "all/",
			/////       when/who/what
			fileDestPrefix + "/" + specialFolderPrefix + "when/" + dateTags["year"] + "/" + dateTags["month"] + "/" + dateTags["day"] + "/" + specialFolderPrefix + "who/" + tag_person + "/" + specialFolderPrefix + "what/" + tag_event + "/",
			fileDestPrefix + "/" + specialFolderPrefix + "when/" + dateTags["year"] + "/" + dateTags["month"] + "/" + specialFolderPrefix + "who/" + tag_person + "/" + specialFolderPrefix + "what/" + tag_event + "/",
			fileDestPrefix + "/" + specialFolderPrefix + "when/" + dateTags["year"] + "/" + specialFolderPrefix + "who/" + tag_person + "/" + specialFolderPrefix + "what/" + tag_event + "/",
			/////       when/who/where
			fileDestPrefix + "/" + specialFolderPrefix + "when/" + dateTags["year"] + "/" + dateTags["month"] + "/" + dateTags["day"] + "/" + specialFolderPrefix + "who/" + tag_person + "/" + specialFolderPrefix + "where/" + tag_place + "/",
			fileDestPrefix + "/" + specialFolderPrefix + "when/" + dateTags["year"] + "/" + dateTags["month"] + "/" + specialFolderPrefix + "who/" + tag_person + "/" + specialFolderPrefix + "where/" + tag_place + "/",
			fileDestPrefix + "/" + specialFolderPrefix + "when/" + dateTags["year"] + "/" + specialFolderPrefix + "who/" + tag_person + "/" + specialFolderPrefix + "where/" + tag_place + "/",
			/////////////////////////////////////

			///////////////////////////////////// what
			fileDestPrefix + "/" + specialFolderPrefix + "what/" + tag_event + "/" + specialFolderPrefix + "all/",

			//////////  what/when
			fileDestPrefix + "/" + specialFolderPrefix + "what/" + tag_event + "/" + specialFolderPrefix + "when/" + dateTags["year"] + "/" + specialFolderPrefix + "all/",
			fileDestPrefix + "/" + specialFolderPrefix + "what/" + tag_event + "/" + specialFolderPrefix + "when/" + dateTags["year"] + "/" + dateTags["month"] + "/" + specialFolderPrefix + "all/",
			fileDestPrefix + "/" + specialFolderPrefix + "what/" + tag_event + "/" + specialFolderPrefix + "when/" + dateTags["year"] + "/" + dateTags["month"] + "/" + dateTags["day"] + "/",

			//////////  what/who
			fileDestPrefix + "/" + specialFolderPrefix + "what/" + tag_event + "/" + specialFolderPrefix + "who/" + tag_person + "/" + specialFolderPrefix + "all/",
			/////       what/who/when
			fileDestPrefix + "/" + specialFolderPrefix + "what/" + tag_event + "/" + specialFolderPrefix + "who/" + tag_person + "/" + specialFolderPrefix + "when/" + dateTags["year"] + "/" + specialFolderPrefix + "all/",
			fileDestPrefix + "/" + specialFolderPrefix + "what/" + tag_event + "/" + specialFolderPrefix + "who/" + tag_person + "/" + specialFolderPrefix + "when/" + dateTags["year"] + "/" + dateTags["month"] + "/" + specialFolderPrefix + "all/",
			fileDestPrefix + "/" + specialFolderPrefix + "what/" + tag_event + "/" + specialFolderPrefix + "who/" + tag_person + "/" + specialFolderPrefix + "when/" + dateTags["year"] + "/" + dateTags["month"] + "/" + dateTags["day"] + "/",

			//////////  what/where
			fileDestPrefix + "/" + specialFolderPrefix + "what/" + tag_event + "/" + specialFolderPrefix + "where/" + tag_place + "/" + specialFolderPrefix + "all/",
			/////       what/where/when
			fileDestPrefix + "/" + specialFolderPrefix + "what/" + tag_event + "/" + specialFolderPrefix + "where/" + tag_place + "/" + specialFolderPrefix + "when/" + dateTags["year"] + "/" + specialFolderPrefix + "all/",
			fileDestPrefix + "/" + specialFolderPrefix + "what/" + tag_event + "/" + specialFolderPrefix + "where/" + tag_place + "/" + specialFolderPrefix + "when/" + dateTags["year"] + "/" + dateTags["month"] + "/" + specialFolderPrefix + "all/",
			fileDestPrefix + "/" + specialFolderPrefix + "what/" + tag_event + "/" + specialFolderPrefix + "where/" + tag_place + "/" + specialFolderPrefix + "when/" + dateTags["year"] + "/" + dateTags["month"] + "/" + dateTags["day"] + "/",
			/////////////////////////////////////

			///////////////////////////////////// where
			fileDestPrefix + "/" + specialFolderPrefix + "where/" + tag_place + "/" + specialFolderPrefix + "all/",

			//////////  where/when
			fileDestPrefix + "/" + specialFolderPrefix + "where/" + tag_place + "/" + specialFolderPrefix + "when/" + dateTags["year"] + "/" + specialFolderPrefix + "all/",
			fileDestPrefix + "/" + specialFolderPrefix + "where/" + tag_place + "/" + specialFolderPrefix + "when/" + dateTags["year"] + "/" + dateTags["month"] + "/" + specialFolderPrefix + "all/",
			fileDestPrefix + "/" + specialFolderPrefix + "where/" + tag_place + "/" + specialFolderPrefix + "when/" + dateTags["year"] + "/" + dateTags["month"] + "/" + dateTags["day"] + "/",

			//////////  where/what
			fileDestPrefix + "/" + specialFolderPrefix + "where/" + tag_place + "/" + specialFolderPrefix + "what/" + tag_event + "/" + specialFolderPrefix + "all/",
			/////       where/what/when
			fileDestPrefix + "/" + specialFolderPrefix + "where/" + tag_place + "/" + specialFolderPrefix + "what/" + tag_event + "/" + specialFolderPrefix + "when/" + dateTags["year"] + "/" + specialFolderPrefix + "all/",
			fileDestPrefix + "/" + specialFolderPrefix + "where/" + tag_place + "/" + specialFolderPrefix + "what/" + tag_event + "/" + specialFolderPrefix + "when/" + dateTags["year"] + "/" + dateTags["month"] + "/" + specialFolderPrefix + "all/",
			fileDestPrefix + "/" + specialFolderPrefix + "where/" + tag_place + "/" + specialFolderPrefix + "what/" + tag_event + "/" + specialFolderPrefix + "when/" + dateTags["year"] + "/" + dateTags["month"] + "/" + dateTags["day"] + "/",

			//////////  where/who
			fileDestPrefix + "/" + specialFolderPrefix + "where/" + tag_place + "/" + specialFolderPrefix + "who/" + tag_person + "/" + specialFolderPrefix + "all/",
			/////       where/who/when
			fileDestPrefix + "/" + specialFolderPrefix + "where/" + tag_place + "/" + specialFolderPrefix + "who/" + tag_person + "/" + specialFolderPrefix + "when/" + dateTags["year"] + "/" + specialFolderPrefix + "all/",
			fileDestPrefix + "/" + specialFolderPrefix + "where/" + tag_place + "/" + specialFolderPrefix + "who/" + tag_person + "/" + specialFolderPrefix + "when/" + dateTags["year"] + "/" + dateTags["month"] + "/" + specialFolderPrefix + "all/",
			fileDestPrefix + "/" + specialFolderPrefix + "where/" + tag_place + "/" + specialFolderPrefix + "who/" + tag_person + "/" + specialFolderPrefix + "when/" + dateTags["year"] + "/" + dateTags["month"] + "/" + dateTags["day"] + "/",
			/////////////////////////////////////

			///////////////////////////////////// who
			fileDestPrefix + "/" + specialFolderPrefix + "who/" + tag_person + "/" + specialFolderPrefix + "all/",

			//////////  who/when
			fileDestPrefix + "/" + specialFolderPrefix + "who/" + tag_person + "/" + specialFolderPrefix + "when/" + dateTags["year"] + "/" + specialFolderPrefix + "all/",
			fileDestPrefix + "/" + specialFolderPrefix + "who/" + tag_person + "/" + specialFolderPrefix + "when/" + dateTags["year"] + "/" + dateTags["month"] + "/" + specialFolderPrefix + "all/",
			fileDestPrefix + "/" + specialFolderPrefix + "who/" + tag_person + "/" + specialFolderPrefix + "when/" + dateTags["year"] + "/" + dateTags["month"] + "/" + dateTags["day"] + "/",

			//////////  who/what
			fileDestPrefix + "/" + specialFolderPrefix + "who/" + tag_person + "/" + specialFolderPrefix + "what/" + tag_event + "/" + specialFolderPrefix + "all/",
			/////       who/what/when
			fileDestPrefix + "/" + specialFolderPrefix + "who/" + tag_person + "/" + specialFolderPrefix + "what/" + tag_event + "/" + specialFolderPrefix + "when/" + dateTags["year"] + "/" + specialFolderPrefix + "all/",
			fileDestPrefix + "/" + specialFolderPrefix + "who/" + tag_person + "/" + specialFolderPrefix + "what/" + tag_event + "/" + specialFolderPrefix + "when/" + dateTags["year"] + "/" + dateTags["month"] + "/" + specialFolderPrefix + "all/",
			fileDestPrefix + "/" + specialFolderPrefix + "who/" + tag_person + "/" + specialFolderPrefix + "what/" + tag_event + "/" + specialFolderPrefix + "when/" + dateTags["year"] + "/" + dateTags["month"] + "/" + dateTags["day"] + "/",

			//////////  who/where
			fileDestPrefix + "/" + specialFolderPrefix + "who/" + tag_person + "/" + specialFolderPrefix + "where/" + tag_place + "/" + specialFolderPrefix + "all/",
			/////       who/where/when
			fileDestPrefix + "/" + specialFolderPrefix + "who/" + tag_person + "/" + specialFolderPrefix + "where/" + tag_place + "/" + specialFolderPrefix + "when/" + dateTags["year"] + "/" + specialFolderPrefix + "all/",
			fileDestPrefix + "/" + specialFolderPrefix + "who/" + tag_person + "/" + specialFolderPrefix + "where/" + tag_place + "/" + specialFolderPrefix + "when/" + dateTags["year"] + "/" + dateTags["month"] + "/" + specialFolderPrefix + "all/",
			fileDestPrefix + "/" + specialFolderPrefix + "who/" + tag_person + "/" + specialFolderPrefix + "where/" + tag_place + "/" + specialFolderPrefix + "when/" + dateTags["year"] + "/" + dateTags["month"] + "/" + dateTags["day"] + "/",
			/////////////////////////////////////

		}

		createSymLinks(file, fileLinkFolders, datePrefix+timePrefix+fileName)

	}
}
