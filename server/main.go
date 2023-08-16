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

func createSymLinks(fileSrc string, folderDests []string, fileDest string) {

  for i := range folderDests {

    err := os.MkdirAll(folderDests[i], os.ModePerm)
    if err != nil {
      log.Println(err)
    }

    fileLink := folderDests[i]+fileDest

    if exists(fileLink) == false {
      fmt.Println("[INFO] symlink needs to be created: "+fileLink)
      err = os.Symlink(fileSrc, fileLink)
      if err != nil {
          log.Fatal(err)
      }

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
                    "place":  "frankfurt",
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
                    "place":  "frankfurt",
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
                    "place":  "frankfurt",
                    "person": "familie"
                  },
                  "/tmp/test/src/geburtstagsfeier-asd454asa.jpg": {
                    "date":  {
                      "year":   "2023",
                      "month":  "09",
                      "day":    "02"
                    },
                    "time": {
                      "hour":   "16",
                      "minute": "10",
                      "second": "01"
                    },
                    "event":  "geburtstagsfeier",
                    "place":  "frankfurt",
                    "person": "claudi maier"
                  },
                  "/tmp/test/src/geburtstagsfeier-sdf4wfasa.jpg": {
                    "date":  {
                      "year":   "2023",
                      "month":  "07",
                      "day":    "03"
                    },
                    "time": {
                      "hour":   "23",
                      "minute": "09",
                      "second": "01"
                    },
                    "event":  "geburtstagsfeier",
                    "place":  "berlin",
                    "person": "martin schuster"
                  },
                  "/tmp/test/src/geburtstagsfeier-d8zz9ad.jpg": {
                    "date":  {
                      "year":   "2023",
                      "month":  "04",
                      "day":    "01"
                    },
                    "time": {
                      "hour":   "22",
                      "minute": "07",
                      "second": "01"
                    },
                    "event":  "geburtstagsfeier",
                    "place":  "linz",
                    "person": "max mustermann"
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


      tag_event   := fileDict[file].(map[string]interface{})["event"].(string)
      tag_person  := fileDict[file].(map[string]interface{})["person"].(string)
      tag_place   := fileDict[file].(map[string]interface{})["place"].(string)

      fmt.Println("===>", "event:", tag_event)
      fmt.Println("===>", "person:", tag_person)
      fmt.Println("===>", "place:", tag_place)


      datePrefix  := dateTags["year"]+dateTags["month"]+dateTags["day"]+"-"
      timePrefix  := timeTags["hour"]+timeTags["minute"]+timeTags["second"]+"_-_"

      specialFolderPrefix := "0000-#-"

      fileLinkFolders := []string {
        ///////////////////////////////////// when
        fileDestPrefix+"/"+specialFolderPrefix+"when/"+dateTags["year"]+"/"+dateTags["month"]+"/"+dateTags["day"]+"/",
        fileDestPrefix+"/"+specialFolderPrefix+"when/"+dateTags["year"]+"/"+dateTags["month"]+"/"+specialFolderPrefix+"all/",
        fileDestPrefix+"/"+specialFolderPrefix+"when/"+dateTags["year"]+"/"+specialFolderPrefix+"all/",

        //////////  when/what
        fileDestPrefix+"/"+specialFolderPrefix+"when/"+dateTags["year"]+"/"+dateTags["month"]+"/"+dateTags["day"]+"/"+specialFolderPrefix+"what/"+tag_event+"/",
        fileDestPrefix+"/"+specialFolderPrefix+"when/"+dateTags["year"]+"/"+dateTags["month"]+"/"+specialFolderPrefix+"all/"+specialFolderPrefix+"what/"+tag_event+"/",
        fileDestPrefix+"/"+specialFolderPrefix+"when/"+dateTags["year"]+"/"+specialFolderPrefix+"all/"+specialFolderPrefix+"what/"+tag_event+"/",
        /////       when/what/where
        fileDestPrefix+"/"+specialFolderPrefix+"when/"+dateTags["year"]+"/"+dateTags["month"]+"/"+dateTags["day"]+"/"+specialFolderPrefix+"what/"+tag_event+"/"+specialFolderPrefix+"where/"+tag_place+"/",
        fileDestPrefix+"/"+specialFolderPrefix+"when/"+dateTags["year"]+"/"+dateTags["month"]+"/"+specialFolderPrefix+"all/"+specialFolderPrefix+"what/"+tag_event+"/"+specialFolderPrefix+"where/"+tag_place+"/",
        fileDestPrefix+"/"+specialFolderPrefix+"when/"+dateTags["year"]+"/"+specialFolderPrefix+"all/"+specialFolderPrefix+"what/"+tag_event+"/"+specialFolderPrefix+"where/"+tag_place+"/",
        /////       when/what/who
        fileDestPrefix+"/"+specialFolderPrefix+"when/"+dateTags["year"]+"/"+dateTags["month"]+"/"+dateTags["day"]+"/"+specialFolderPrefix+"what/"+tag_event+"/"+specialFolderPrefix+"who/"+tag_person+"/",
        fileDestPrefix+"/"+specialFolderPrefix+"when/"+dateTags["year"]+"/"+dateTags["month"]+"/"+specialFolderPrefix+"all/"+specialFolderPrefix+"what/"+tag_event+"/"+specialFolderPrefix+"who/"+tag_person+"/",
        fileDestPrefix+"/"+specialFolderPrefix+"when/"+dateTags["year"]+"/"+specialFolderPrefix+"all/"+specialFolderPrefix+"what/"+tag_event+"/"+specialFolderPrefix+"who/"+tag_person+"/",

        //////////  when/where
        fileDestPrefix+"/"+specialFolderPrefix+"when/"+dateTags["year"]+"/"+dateTags["month"]+"/"+dateTags["day"]+"/"+specialFolderPrefix+"where/"+tag_place+"/",
        fileDestPrefix+"/"+specialFolderPrefix+"when/"+dateTags["year"]+"/"+dateTags["month"]+"/"+specialFolderPrefix+"all/"+specialFolderPrefix+"where/"+tag_place+"/",
        fileDestPrefix+"/"+specialFolderPrefix+"when/"+dateTags["year"]+"/"+specialFolderPrefix+"all/"+specialFolderPrefix+"where/"+tag_place+"/",
        /////       when/where/what
        fileDestPrefix+"/"+specialFolderPrefix+"when/"+dateTags["year"]+"/"+dateTags["month"]+"/"+dateTags["day"]+"/"+specialFolderPrefix+"where/"+tag_place+"/"+specialFolderPrefix+"what/"+tag_event+"/",
        fileDestPrefix+"/"+specialFolderPrefix+"when/"+dateTags["year"]+"/"+dateTags["month"]+"/"+specialFolderPrefix+"all/"+specialFolderPrefix+"where/"+tag_place+"/"+specialFolderPrefix+"what/"+tag_event+"/",
        fileDestPrefix+"/"+specialFolderPrefix+"when/"+dateTags["year"]+"/"+specialFolderPrefix+"all/"+specialFolderPrefix+"where/"+tag_place+"/"+specialFolderPrefix+"what/"+tag_event+"/",
        /////       when/where/wwho
        fileDestPrefix+"/"+specialFolderPrefix+"when/"+dateTags["year"]+"/"+dateTags["month"]+"/"+dateTags["day"]+"/"+specialFolderPrefix+"where/"+tag_place+"/"+specialFolderPrefix+"who/"+tag_person+"/",
        fileDestPrefix+"/"+specialFolderPrefix+"when/"+dateTags["year"]+"/"+dateTags["month"]+"/"+specialFolderPrefix+"all/"+specialFolderPrefix+"where/"+tag_place+"/"+specialFolderPrefix+"who/"+tag_person+"/",
        fileDestPrefix+"/"+specialFolderPrefix+"when/"+dateTags["year"]+"/"+specialFolderPrefix+"all/"+specialFolderPrefix+"where/"+tag_place+"/"+specialFolderPrefix+"who/"+tag_person+"/",

        //////////  when/who
        fileDestPrefix+"/"+specialFolderPrefix+"when/"+dateTags["year"]+"/"+dateTags["month"]+"/"+dateTags["day"]+"/"+specialFolderPrefix+"who/"+tag_person+"/",
        fileDestPrefix+"/"+specialFolderPrefix+"when/"+dateTags["year"]+"/"+dateTags["month"]+"/"+specialFolderPrefix+"all/"+specialFolderPrefix+"who/"+tag_person+"/",
        fileDestPrefix+"/"+specialFolderPrefix+"when/"+dateTags["year"]+"/"+specialFolderPrefix+"all/"+specialFolderPrefix+"who/"+tag_person+"/",
        /////       when/who/what
        fileDestPrefix+"/"+specialFolderPrefix+"when/"+dateTags["year"]+"/"+dateTags["month"]+"/"+dateTags["day"]+"/"+specialFolderPrefix+"who/"+tag_person+"/"+specialFolderPrefix+"what/"+tag_event+"/",
        fileDestPrefix+"/"+specialFolderPrefix+"when/"+dateTags["year"]+"/"+dateTags["month"]+"/"+specialFolderPrefix+"all/"+specialFolderPrefix+"who/"+tag_person+"/"+specialFolderPrefix+"what/"+tag_event+"/",
        fileDestPrefix+"/"+specialFolderPrefix+"when/"+dateTags["year"]+"/"+specialFolderPrefix+"all/"+specialFolderPrefix+"who/"+tag_person+"/"+specialFolderPrefix+"what/"+tag_event+"/",
        /////       when/who/where
        fileDestPrefix+"/"+specialFolderPrefix+"when/"+dateTags["year"]+"/"+dateTags["month"]+"/"+dateTags["day"]+"/"+specialFolderPrefix+"who/"+tag_person+"/"+specialFolderPrefix+"where/"+tag_place+"/",
        fileDestPrefix+"/"+specialFolderPrefix+"when/"+dateTags["year"]+"/"+dateTags["month"]+"/"+specialFolderPrefix+"all/"+specialFolderPrefix+"who/"+tag_person+"/"+specialFolderPrefix+"where/"+tag_place+"/",
        fileDestPrefix+"/"+specialFolderPrefix+"when/"+dateTags["year"]+"/"+specialFolderPrefix+"all/"+specialFolderPrefix+"who/"+tag_person+"/"+specialFolderPrefix+"where/"+tag_place+"/",
        /////////////////////////////////////

        ///////////////////////////////////// what
        fileDestPrefix+"/"+specialFolderPrefix+"what/"+tag_event+"/"+specialFolderPrefix+"all/",

        //////////  what/when
        fileDestPrefix+"/"+specialFolderPrefix+"what/"+tag_event+"/"+specialFolderPrefix+"when/"+dateTags["year"]+"/"+specialFolderPrefix+"all/",
        fileDestPrefix+"/"+specialFolderPrefix+"what/"+tag_event+"/"+specialFolderPrefix+"when/"+dateTags["year"]+"/"+dateTags["month"]+"/"+specialFolderPrefix+"all/",
        fileDestPrefix+"/"+specialFolderPrefix+"what/"+tag_event+"/"+specialFolderPrefix+"when/"+dateTags["year"]+"/"+dateTags["month"]+"/"+dateTags["day"]+"/",

        //////////  what/who
        fileDestPrefix+"/"+specialFolderPrefix+"what/"+tag_event+"/"+specialFolderPrefix+"who/"+tag_person+"/"+specialFolderPrefix+"all/",
        /////       what/who/when
        fileDestPrefix+"/"+specialFolderPrefix+"what/"+tag_event+"/"+specialFolderPrefix+"who/"+tag_person+"/"+specialFolderPrefix+"when/"+dateTags["year"]+"/"+specialFolderPrefix+"all/",
        fileDestPrefix+"/"+specialFolderPrefix+"what/"+tag_event+"/"+specialFolderPrefix+"who/"+tag_person+"/"+specialFolderPrefix+"when/"+dateTags["year"]+"/"+dateTags["month"]+"/"+specialFolderPrefix+"all/",
        fileDestPrefix+"/"+specialFolderPrefix+"what/"+tag_event+"/"+specialFolderPrefix+"who/"+tag_person+"/"+specialFolderPrefix+"when/"+dateTags["year"]+"/"+dateTags["month"]+"/"+dateTags["day"]+"/",

        //////////  what/where
        fileDestPrefix+"/"+specialFolderPrefix+"what/"+tag_event+"/"+specialFolderPrefix+"where/"+tag_place+"/"+specialFolderPrefix+"all/",
        /////       what/where/when
        fileDestPrefix+"/"+specialFolderPrefix+"what/"+tag_event+"/"+specialFolderPrefix+"where/"+tag_place+"/"+specialFolderPrefix+"when/"+dateTags["year"]+"/"+specialFolderPrefix+"all/",
        fileDestPrefix+"/"+specialFolderPrefix+"what/"+tag_event+"/"+specialFolderPrefix+"where/"+tag_place+"/"+specialFolderPrefix+"when/"+dateTags["year"]+"/"+dateTags["month"]+"/"+specialFolderPrefix+"all/",
        fileDestPrefix+"/"+specialFolderPrefix+"what/"+tag_event+"/"+specialFolderPrefix+"where/"+tag_place+"/"+specialFolderPrefix+"when/"+dateTags["year"]+"/"+dateTags["month"]+"/"+dateTags["day"]+"/",
        /////////////////////////////////////

        ///////////////////////////////////// where
        fileDestPrefix+"/"+specialFolderPrefix+"where/"+tag_place+"/"+specialFolderPrefix+"all/",

        //////////  where/when
        fileDestPrefix+"/"+specialFolderPrefix+"where/"+tag_place+"/"+specialFolderPrefix+"when/"+dateTags["year"]+"/"+specialFolderPrefix+"all/",
        fileDestPrefix+"/"+specialFolderPrefix+"where/"+tag_place+"/"+specialFolderPrefix+"when/"+dateTags["year"]+"/"+dateTags["month"]+"/"+specialFolderPrefix+"all/",
        fileDestPrefix+"/"+specialFolderPrefix+"where/"+tag_place+"/"+specialFolderPrefix+"when/"+dateTags["year"]+"/"+dateTags["month"]+"/"+dateTags["day"]+"/",

        //////////  where/what
        fileDestPrefix+"/"+specialFolderPrefix+"where/"+tag_place+"/"+specialFolderPrefix+"what/"+tag_event+"/"+specialFolderPrefix+"all/",
        /////       where/what/when
        fileDestPrefix+"/"+specialFolderPrefix+"where/"+tag_place+"/"+specialFolderPrefix+"what/"+tag_event+"/"+specialFolderPrefix+"when/"+dateTags["year"]+"/"+specialFolderPrefix+"all/",
        fileDestPrefix+"/"+specialFolderPrefix+"where/"+tag_place+"/"+specialFolderPrefix+"what/"+tag_event+"/"+specialFolderPrefix+"when/"+dateTags["year"]+"/"+dateTags["month"]+"/"+specialFolderPrefix+"all/",
        fileDestPrefix+"/"+specialFolderPrefix+"where/"+tag_place+"/"+specialFolderPrefix+"what/"+tag_event+"/"+specialFolderPrefix+"when/"+dateTags["year"]+"/"+dateTags["month"]+"/"+dateTags["day"]+"/",

        //////////  where/who
        fileDestPrefix+"/"+specialFolderPrefix+"where/"+tag_place+"/"+specialFolderPrefix+"who/"+tag_person+"/"+specialFolderPrefix+"all/",
        /////       where/who/when
        fileDestPrefix+"/"+specialFolderPrefix+"where/"+tag_place+"/"+specialFolderPrefix+"who/"+tag_person+"/"+specialFolderPrefix+"when/"+dateTags["year"]+"/"+specialFolderPrefix+"all/",
        fileDestPrefix+"/"+specialFolderPrefix+"where/"+tag_place+"/"+specialFolderPrefix+"who/"+tag_person+"/"+specialFolderPrefix+"when/"+dateTags["year"]+"/"+dateTags["month"]+"/"+specialFolderPrefix+"all/",
        fileDestPrefix+"/"+specialFolderPrefix+"where/"+tag_place+"/"+specialFolderPrefix+"who/"+tag_person+"/"+specialFolderPrefix+"when/"+dateTags["year"]+"/"+dateTags["month"]+"/"+dateTags["day"]+"/",
        /////////////////////////////////////

        ///////////////////////////////////// who
        fileDestPrefix+"/"+specialFolderPrefix+"who/"+tag_person+"/"+specialFolderPrefix+"all/",

        //////////  who/when
        fileDestPrefix+"/"+specialFolderPrefix+"who/"+tag_person+"/"+specialFolderPrefix+"when/"+dateTags["year"]+"/"+specialFolderPrefix+"all/",
        fileDestPrefix+"/"+specialFolderPrefix+"who/"+tag_person+"/"+specialFolderPrefix+"when/"+dateTags["year"]+"/"+dateTags["month"]+"/"+specialFolderPrefix+"all/",
        fileDestPrefix+"/"+specialFolderPrefix+"who/"+tag_person+"/"+specialFolderPrefix+"when/"+dateTags["year"]+"/"+dateTags["month"]+"/"+dateTags["day"]+"/",

        //////////  who/what
        fileDestPrefix+"/"+specialFolderPrefix+"who/"+tag_person+"/"+specialFolderPrefix+"what/"+tag_event+"/"+specialFolderPrefix+"all/",
        /////       who/what/when
        fileDestPrefix+"/"+specialFolderPrefix+"who/"+tag_person+"/"+specialFolderPrefix+"what/"+tag_event+"/"+specialFolderPrefix+"when/"+dateTags["year"]+"/"+specialFolderPrefix+"all/",
        fileDestPrefix+"/"+specialFolderPrefix+"who/"+tag_person+"/"+specialFolderPrefix+"what/"+tag_event+"/"+specialFolderPrefix+"when/"+dateTags["year"]+"/"+dateTags["month"]+"/"+specialFolderPrefix+"all/",
        fileDestPrefix+"/"+specialFolderPrefix+"who/"+tag_person+"/"+specialFolderPrefix+"what/"+tag_event+"/"+specialFolderPrefix+"when/"+dateTags["year"]+"/"+dateTags["month"]+"/"+dateTags["day"]+"/",

        //////////  who/where
        fileDestPrefix+"/"+specialFolderPrefix+"who/"+tag_person+"/"+specialFolderPrefix+"where/"+tag_place+"/"+specialFolderPrefix+"all/",
        /////       who/where/when
        fileDestPrefix+"/"+specialFolderPrefix+"who/"+tag_person+"/"+specialFolderPrefix+"where/"+tag_place+"/"+specialFolderPrefix+"when/"+dateTags["year"]+"/"+specialFolderPrefix+"all/",
        fileDestPrefix+"/"+specialFolderPrefix+"who/"+tag_person+"/"+specialFolderPrefix+"where/"+tag_place+"/"+specialFolderPrefix+"when/"+dateTags["year"]+"/"+dateTags["month"]+"/"+specialFolderPrefix+"all/",
        fileDestPrefix+"/"+specialFolderPrefix+"who/"+tag_person+"/"+specialFolderPrefix+"where/"+tag_place+"/"+specialFolderPrefix+"when/"+dateTags["year"]+"/"+dateTags["month"]+"/"+dateTags["day"]+"/",
        /////////////////////////////////////

      }

      createSymLinks(file, fileLinkFolders, datePrefix+timePrefix+fileName)

    }
}
