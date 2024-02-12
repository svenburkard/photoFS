package main

import (
	"fmt"
	"log"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

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

	// db: add some test tag names
	err = photofs.AddTestTagNames(db)
	if err != nil {
		log.Fatal(err)
	}

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

	newtagNameTxt := widget.NewEntry()
	newtagNameTxt.PlaceHolder = "New Tag..."

	addBtn := widget.NewButton("Add", func() {
		//     tags.Append(newtagNameTxt.Text) // TODO: needs to be changed to the new tagWidgets map
		newtagNameTxt.Text = ""
	})
	addBtn.Disable()

	newtagNameTxt.OnChanged = func(s string) {
		addBtn.Disable()

		if len(s) >= 3 {
			addBtn.Enable()
		}
	}

	updateBtn := widget.NewButton("Update Tags", func() {
		if len(selectedTags) > 0 {
			fmt.Println("[client] TODO: make DB update of selectedTags: ", selectedTags)
		} else {
			fmt.Println("[client] DB will not be updated, because selectedTags map is empty")
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
							addBtn,
							// take the rest of the space
							newtagNameTxt,
						),
						container.NewBorder(
							nil, // TOP
							nil, // BOTTOM
							nil, // Left
							// RIGHT ↓
							addBtn,
							// take the rest of the space
							newtagNameTxt,
						),
						container.NewBorder(
							nil, // TOP
							nil, // BOTTOM
							nil, // Left
							// RIGHT ↓
							addBtn,
							// take the rest of the space
							newtagNameTxt,
						),
						container.NewBorder(
							nil, // TOP
							nil, // BOTTOM
							nil, // Left
							// RIGHT ↓
							addBtn,
							// take the rest of the space
							newtagNameTxt,
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
						// 						tagWidgets["Who"],
						container.NewScroll(widget.NewLabel("placeholder")), // who
						container.NewScroll(tagWidgets["Where"]),
						container.NewScroll(widget.NewLabel("placeholder")), // misc
					),
					container.NewBorder(
						// TOP
						// 						widget.NewLabel("Update Tags"), // TOP
						updateBtn,
						// 						nil,
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
