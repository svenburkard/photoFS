package main

import (
	"fmt"
	"log"
	"os"
	"slices"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	bolt "go.etcd.io/bbolt"

	photofs "photofs/lib"
)

func getTagWidgets(tagNames map[string][]string) (map[string]fyne.CanvasObject, map[string][]string, error) {
	tagWidgets := make(map[string]fyne.CanvasObject)
	selectedTags := make(map[string][]string)

	for tagType, tagNameList := range tagNames {
		localTagType := tagType // Create a local copy of tagType
		tagWidgets[localTagType] = widget.NewCheckGroup(tagNameList, func(selected []string) {
			selectedTags[localTagType] = selected
			fmt.Printf("selectedTags map: %v\n", selectedTags)
		})
	}

	if len(tagWidgets) == 0 {
		return nil, nil, fmt.Errorf("failed to build map of tagWidgets in getTagWidgets func")
	}

	return tagWidgets, selectedTags, nil
}

func addTagNameToCheckGroup(tagWidgets map[string]fyne.CanvasObject, tagType string, tagName string) {
	checkGroup := tagWidgets[tagType].(*widget.CheckGroup)

	oldOptions := checkGroup.Options
	newOptions := append(checkGroup.Options, tagName)

	if (newOptions == nil) || slices.Equal(oldOptions, newOptions) {
		log.Fatal("no new tagName could be added to checkGroup.Options in addTagNameToCheckGroup")
	}

	checkGroup.Options = newOptions
	checkGroup.Refresh()
}

func getTagNameEntriesAndButtons(db *bolt.DB, tagNames map[string][]string, tagWidgets map[string]fyne.CanvasObject) (map[string]*widget.Entry, map[string]*widget.Button, error) {
	tagNameEntries := make(map[string]*widget.Entry)
	tagNameButtons := make(map[string]*widget.Button)

	for tagType, _ := range tagNames {
		localTagType := tagType // Create a local copy of tagType

		tagNameEntries[localTagType] = widget.NewEntry()
		tagNameEntries[localTagType].PlaceHolder = "New Tag..."

		tagNameButtons[localTagType] = widget.NewButton("Add", func() {

			addTagNameToCheckGroup(tagWidgets, localTagType, tagNameEntries[localTagType].Text)
			tagNames[localTagType] = append(tagNames[localTagType], tagNameEntries[localTagType].Text)

			tagNamesToAdd := map[string]photofs.TagNameList{
				localTagType: photofs.TagNameList{tagNames[localTagType]},
			}
			err := photofs.AddTagNames(db, tagNamesToAdd)
			if err != nil {
				log.Fatal(err)
			}

			tagNameEntries[localTagType].Text = ""

		})

		tagNameButtons[localTagType].Disable()

		tagNameEntries[localTagType].OnChanged = func(s string) {
			tagNameButtons[localTagType].Disable()

			if len(s) >= 3 {
				tagNameButtons[localTagType].Enable()
			}
		}
	}

	if len(tagNameEntries) == 0 || len(tagNameButtons) == 0 {
		return nil, nil, fmt.Errorf("failed to build the map of tagNameEntries and tagNameButtons in getTagNameEntriesAndButtons func")
	}

	return tagNameEntries, tagNameButtons, nil
}

func main() {
	args := os.Args

	var files []string

	if len(args) >= 0 {
		for _, arg := range args[1:] { // args[1:] to skip the element 0, which would be the script name
			fmt.Printf("Type of Arg = %T\n", arg)
			fmt.Println(arg)
			files = append(files, arg)
		}
	}

	// db: init
	db, err := photofs.InitDB()
	if err != nil {
		fmt.Errorf("[ERROR] failed to initialize DB: %v", err)
	}
	defer db.Close()

	// 	// db: add some test tag names
	// 	err = photofs.AddTestTagNames(db)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}

	fileList := widget.NewList(
		func() int {
			return len(files)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(files[i])
		},
	)

	a := app.New()
	w := a.NewWindow("photoFS Tag App")

	//     w.Resize(fyne.NewSize(1280, 1280))
	w.Resize(fyne.NewSize(1280, 600))
	// 	w.Resize(fyne.NewSize(1280, 800))
	w.CenterOnScreen()

	// db: get tagNames
	fmt.Println("[client] Get tagNames:")
	tagNames, err := photofs.GetTagNames(db)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("[client] tagNames: ", tagNames)

	tagWidgets, selectedTags, err := getTagWidgets(tagNames)
	if err != nil {
		log.Fatal(err)
	}

	tagNameEntries, tagNameButtons, err := getTagNameEntriesAndButtons(db, tagNames, tagWidgets)
	if err != nil {
		log.Fatal(err)
	}

	updateBtn := widget.NewButton("Update Tags", func() {
		if len(files) > 0 {
			if len(selectedTags) > 0 {

				for _, srcFile := range files {
					tagsOfFile, err := photofs.ConvertSelectedTagsToTagsOfFile(selectedTags)
					if err != nil {
						log.Fatal(err)
					}
					err = photofs.AddTagsOfFile(db, srcFile, tagsOfFile)
					if err == nil {
						fmt.Printf("[client] succesfully updated tags (%v) of srcFile (%v)", selectedTags, srcFile)
					} else {
						log.Fatal("[client] failed to add tags (%v) of file (%v) to db: %w", selectedTags, srcFile, err)
					}
				}

			} else {
				fmt.Println("[client] DB will not be updated, because selectedTags map is empty")
			}
		} else {
			fmt.Println("[client] DB will not be updated, because no files have been selected")
		}
	})

	w.SetContent(
		container.NewBorder(
			// TOP of the container
			container.New(
				layout.NewGridLayout(4),
				widget.NewLabel("What:"),
				widget.NewLabel("Who:"),
				widget.NewLabel("Where:"),
				widget.NewLabel("Misc:"),
			),
			nil, // BOTTOM
			nil, // Right
			nil, // Left
			// REST
			container.NewBorder(
				// TOP of the container

				container.NewBorder(
					// TOP of the container
					container.New(
						layout.NewGridLayout(4),
						container.NewBorder(
							nil, // TOP
							nil, // BOTTOM
							nil, // Left
							// RIGHT ↓
							tagNameButtons["What"],
							// take the rest of the space
							tagNameEntries["What"],
						),
						container.NewBorder(
							nil, // TOP
							nil, // BOTTOM
							nil, // Left
							// RIGHT ↓
							tagNameButtons["Who"],
							// take the rest of the space
							tagNameEntries["Who"],
						),
						container.NewBorder(
							nil, // TOP
							nil, // BOTTOM
							nil, // Left
							// RIGHT ↓
							tagNameButtons["Where"],
							// take the rest of the space
							tagNameEntries["Where"],
						),
						container.NewBorder(
							nil, // TOP
							nil, // BOTTOM
							nil, // Left
							// RIGHT ↓
							tagNameButtons["Misc"],
							// take the rest of the space
							tagNameEntries["Misc"],
						),
					),
					nil, // BOTTOM
					nil, // Right
					nil, // Left
					nil, // REST
				),

				nil, // BOTTOM
				nil, // Right
				nil, // Left
				// REST
				container.New(
					layout.NewGridLayout(1),
					container.New(
						layout.NewGridLayout(4),
						container.NewScroll(tagWidgets["What"]),
						container.NewScroll(tagWidgets["Who"]),
						container.NewScroll(tagWidgets["Where"]),
						container.NewScroll(tagWidgets["Misc"]),
					),
					container.NewBorder(
						// TOP
						updateBtn,
						nil, // BOTTOM
						nil, // RIGHT
						nil, // LEFT
						// REST
						container.NewBorder(
							// TOP
							widget.NewLabel("List of Files:"),
							nil, // BOTTOM
							nil, // RIGHT
							nil, // LEFT
							// REST
							container.New(
								layout.NewGridLayout(1),
								fileList,
							),
						),
					),
					container.New(
						layout.NewGridLayout(1),
						widget.NewLabel("Status Messages:"),
					),
				),
			),
		),
	)

	w.ShowAndRun()

}
