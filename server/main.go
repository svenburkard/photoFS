package main

import (
    "os"
    "log"
    "fmt"
    "encoding/json"
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

func createSymLink(fileSrc string, folderDest string, fileDest string) {

  err := os.MkdirAll(folderDest, os.ModePerm)
  if err != nil {
    log.Println(err)
  }

  fileLink := folderDest+fileDest

  if exists(fileLink) == false {
    fmt.Println("[INFO] symlink needs to be created: "+fileLink)
    err = os.Symlink(fileSrc, fileLink)
    if err != nil {
        log.Fatal(err)
    }

  }
}

func getFileName(file string) string {
  fileExtensions  := "jpg|jpeg|png|avi|mp4"
  rgx             := regexp.MustCompile(`[^/]+\.(`+fileExtensions+`)`)
  fileName        := rgx.FindString(file)
  fmt.Println("fileName: ", fileName)

  return fileName
}

func getDateTags(fileDict interface{}) map[string]string {

  dateTags := make(map[string]string)

  dateTags["year"]  = fileDict.(map[string]interface{})["date"].(map[string]interface{})["year"].(string)
  dateTags["month"] = fileDict.(map[string]interface{})["date"].(map[string]interface{})["month"].(string)
  dateTags["day"]   = fileDict.(map[string]interface{})["date"].(map[string]interface{})["day"].(string)

  fmt.Println("===>", "date/year:",   dateTags["year"])
  fmt.Println("===>", "date/month:",  dateTags["month"])
  fmt.Println("===>", "date/day:",    dateTags["day"])

  return dateTags
}

func getTimeTags(fileDict interface{}) map[string]string {

  timeTags := make(map[string]string)

  timeTags["hour"]    = fileDict.(map[string]interface{})["time"].(map[string]interface{})["hour"].(string)
  timeTags["minute"]  = fileDict.(map[string]interface{})["time"].(map[string]interface{})["minute"].(string)
  timeTags["second"]  = fileDict.(map[string]interface{})["time"].(map[string]interface{})["second"].(string)

  fmt.Println("===>", "time/hour:",   timeTags["hour"])
  fmt.Println("===>", "time/minute:", timeTags["minute"])
  fmt.Println("===>", "time/second:", timeTags["second"])

  return timeTags
}

func main() {

    fileDestPrefix := "/tmp/test/dst"
    fmt.Printf("fileDestPrefix: %s\n", fileDestPrefix)

    jsonStr := `{
                  "/tmp/test/src/baustelle-asd.jpg": {
                    "date":  {
                      "year":   "2023",
                      "month":  "10",
                      "day":    "02"
                    },
                    "time": {
                      "hour":   "15",
                      "minute": "10",
                      "second": "01"
                    },
                    "event":  "baustelle",
                    "place":  "estenfeld",
                    "person": "familie"
                  },
                  "/tmp/test/src/baustelle-dgfh.jpg": {
                    "date":  {
                      "year":   "2023",
                      "month":  "10",
                      "day":    "02"
                    },
                    "time": {
                      "hour":   "15",
                      "minute": "11",
                      "second": "01"
                    },
                    "event":  "baustelle",
                    "place":  "estenfeld",
                    "person": "familie"
                  },
                  "/tmp/test/src/baustelle-tzuzti.jpg": {
                    "date":  {
                      "year":   "2023",
                      "month":  "09",
                      "day":    "02"
                    },
                    "time": {
                      "hour":   "15",
                      "minute": "12",
                      "second": "01"
                    },
                    "event":  "baustelle",
                    "place":  "estenfeld",
                    "person": "familie"
                  }
                }`

    fileDict := map[string]interface{}{}
    err := json.Unmarshal([]byte(jsonStr), &fileDict)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(fileDict)


    for file := range fileDict {

      fmt.Println("file:", file, "=>", "tag-map:", fileDict[file])

      fileName := getFileName(file)

      dateTags := getDateTags(fileDict[file])
      timeTags := getTimeTags(fileDict[file])


      tag_event       := fileDict[file].(map[string]interface{})["event"].(string)
      tag_person      := fileDict[file].(map[string]interface{})["person"].(string)
      tag_place         := fileDict[file].(map[string]interface{})["place"].(string)

      fmt.Println("===>", "event:", tag_event)
      fmt.Println("===>", "person:", tag_person)
      fmt.Println("===>", "place:", tag_place)


      datePrefix  := dateTags["year"]+dateTags["month"]+dateTags["day"]+"-"
      timePrefix  := timeTags["hour"]+timeTags["minute"]+timeTags["second"]+"_-_"

      //// date Top Level
      fileLinkFolder := fileDestPrefix+"/date/"+dateTags["year"]+"/"+dateTags["month"]+"/"+dateTags["day"]+"/"
      createSymLink(file, fileLinkFolder, datePrefix+timePrefix+fileName)

      //// event Top Level
      fileLinkFolder = fileDestPrefix+"/"+"event/"+tag_event+"/"
      createSymLink(file, fileLinkFolder, datePrefix+timePrefix+fileName)

      //// place Top Level
      fileLinkFolder = fileDestPrefix+"/"+"place/"+tag_place+"/"
      createSymLink(file, fileLinkFolder, datePrefix+timePrefix+fileName)

      //// person Top Level
      fileLinkFolder = fileDestPrefix+"/"+"person/"+tag_person+"/"
      createSymLink(file, fileLinkFolder, datePrefix+timePrefix+fileName)
      /////////////////////

    }
}
