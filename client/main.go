package main

import (
    "fmt"
    "os"

    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/app"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/widget"
    "fyne.io/fyne/v2/layout"
)

func main() {
    args := os.Args

    var files []string

    if len(args) >= 0 {
        for _, arg := range args[1:] {              // args[1:] to skip the element 0, which would be the script name
            fmt.Printf("Type of Arg = %T\n", arg)
            fmt.Println(arg)
            files = append(files, arg)
        }
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
    w.CenterOnScreen()

    tagsList := []string{
      "party/birthday/smith/steven",
      "party/birthday/smith/julia",
      "party/christmas/company",
    }

    tags := widget.NewCheckGroup(tagsList, func(selectedTags []string) {
      fmt.Println("tags check_group:")
      fmt.Println("  selected tags: ", selectedTags)
    })


    newtagNameTxt := widget.NewEntry()
    newtagNameTxt.PlaceHolder = "New Tag..."

    addBtn := widget.NewButton("Add", func() {
      tags.Append(newtagNameTxt.Text)
      newtagNameTxt.Text = ""
    })
    addBtn.Disable()

    newtagNameTxt.OnChanged = func(s string) {
        addBtn.Disable()

        if len(s) >= 3 {
            addBtn.Enable()
        }
    }



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
                        tags,
                        widget.NewLabel("placeholder"),
                        widget.NewLabel("placeholder"),
                        widget.NewLabel("placeholder"),
                    ),
                    container.NewBorder(
                        // TOP
                        nil,
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
