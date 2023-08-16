package main

import (
    "os"
    "log"
    "fmt"
    "encoding/json"
    "regexp"
)

func createSymLink(fileSrc string, folderDest string, fileDest string) {

  err := os.MkdirAll(folderDest, os.ModePerm)
  if err != nil {
    log.Println(err)
  }

  fileLink := folderDest+fileDest

  err = os.Symlink(fileSrc, fileLink)
  if err != nil {
      log.Fatal(err)
  }
}


func main() {

    fileDestPrefix := "/tmp/test/dst"
    fmt.Printf("fileDestPrefix: %s\n", fileDestPrefix)

    jsonStr := `{
                  "/tmp/test/src/baustelle-asd.jpg": {
                    "datum":  {
                      "jahr": "2023",
                      "monat": "10",
                      "tag": "02"
                    },
                    "zeit": {
                      "stunde":   "15",
                      "minute":   "10",
                      "sekunde":  "01"
                    },
                    "event":  "baustelle",
                    "ort":    "estenfeld",
                    "person": "familie"
                  },
                  "/tmp/test/src/baustelle-dgfh.jpg": {
                    "datum":  {
                      "jahr": "2023",
                      "monat": "10",
                      "tag": "02"
                    },
                    "zeit": {
                      "stunde":   "15",
                      "minute":   "11",
                      "sekunde":  "01"
                    },
                    "event":  "baustelle",
                    "ort":    "estenfeld",
                    "person": "familie"
                  },
                  "/tmp/test/src/baustelle-tzuzti.jpg": {
                    "datum":  {
                      "jahr": "2023",
                      "monat": "09",
                      "tag": "02"
                    },
                    "zeit": {
                      "stunde":   "15",
                      "minute":   "12",
                      "sekunde":  "01"
                    },
                    "event":  "baustelle",
                    "ort":    "estenfeld",
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

      fmt.Println("file1:", file, "=>", "tag-map:", fileDict[file])

      fileExtensions  := "jpg|jpeg|png|avi|mp4"
      rgx             := regexp.MustCompile(`[^/]+\.(`+fileExtensions+`)`)
      fileName        := rgx.FindString(file)
      fmt.Println("fileName: ", fileName)

      tag_event       := fileDict[file].(map[string]interface{})["event"].(string)
      tag_person      := fileDict[file].(map[string]interface{})["person"].(string)
      tag_ort         := fileDict[file].(map[string]interface{})["ort"].(string)
      tag_datum_jahr  := fileDict[file].(map[string]interface{})["datum"].(map[string]interface{})["jahr"].(string)
      tag_datum_monat := fileDict[file].(map[string]interface{})["datum"].(map[string]interface{})["monat"].(string)
      tag_datum_tag   := fileDict[file].(map[string]interface{})["datum"].(map[string]interface{})["tag"].(string)

      tag_zeit_stunde  := fileDict[file].(map[string]interface{})["zeit"].(map[string]interface{})["stunde"].(string)
      tag_zeit_minute  := fileDict[file].(map[string]interface{})["zeit"].(map[string]interface{})["minute"].(string)
      tag_zeit_sekunde := fileDict[file].(map[string]interface{})["zeit"].(map[string]interface{})["sekunde"].(string)

      fmt.Println("===>", "event:", tag_event)
      fmt.Println("===>", "person:", tag_person)
      fmt.Println("===>", "ort:", tag_ort)

      fmt.Println("===>", "datum/jahr:", tag_datum_jahr)
      fmt.Println("===>", "datum/monat:", tag_datum_monat)
      fmt.Println("===>", "datum/tag:", tag_datum_tag)

      fmt.Println("===>", "zeit/stunde:", tag_zeit_stunde)
      fmt.Println("===>", "zeit/minute:", tag_zeit_minute)
      fmt.Println("===>", "zeit/sekunde:", tag_zeit_sekunde)

//       datumPrefix := tag_datum_jahr+tag_datum_monat+tag_datum_tag+"_-_"
      datumPrefix := tag_datum_jahr+tag_datum_monat+tag_datum_tag+"-"
      zeitPrefix  := tag_zeit_stunde+tag_zeit_minute+tag_zeit_sekunde+"_-_"

      //// Datum Top Level
      fileLinkFolder := fileDestPrefix+"/datum/"+tag_datum_jahr+"/"+tag_datum_monat+"/"+tag_datum_tag+"/"
      createSymLink(file, fileLinkFolder, datumPrefix+zeitPrefix+fileName)

      //// Event Top Level
      fileLinkFolder = fileDestPrefix+"/"+"event/"+tag_event+"/"
      createSymLink(file, fileLinkFolder, datumPrefix+zeitPrefix+fileName)

      //// Ort Top Level
      fileLinkFolder = fileDestPrefix+"/"+"ort/"+tag_ort+"/"
      createSymLink(file, fileLinkFolder, datumPrefix+zeitPrefix+fileName)

      //// Person Top Level
      fileLinkFolder = fileDestPrefix+"/"+"person/"+tag_person+"/"
      createSymLink(file, fileLinkFolder, datumPrefix+zeitPrefix+fileName)
      /////////////////////

    }
}
