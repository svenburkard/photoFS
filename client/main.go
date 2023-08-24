package main

import (
    "fmt"
    "os"

    "tag/models"

    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/app"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/data/binding"
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

    w.Resize(fyne.NewSize(1280, 1280))
    w.CenterOnScreen()

    data := []models.Tag{
        models.NewTag("party/birthday/smith/steven"),
        models.NewTag("party/birthday/smith/julia"),
        models.NewTag("party/christmas/company"),
    }

    tags := binding.NewUntypedList()
    for _, t := range data {
      tags.Append(t)
    }


    newtagNameTxt := widget.NewEntry()
    newtagNameTxt.PlaceHolder = "New Tag..."

    addBtn := widget.NewButton("Add", func() {
      tags.Append(models.NewTag(newtagNameTxt.Text))
      newtagNameTxt.Text = ""
    })
    addBtn.Disable()

    newtagNameTxt.OnChanged = func(s string) {
        addBtn.Disable()

        if len(s) >= 3 {
            addBtn.Enable()
        }
    }


    selectedTags := widget.NewListWithData(
        // the binding.List type
        tags,
        // func that returns the component structure of the List Item
        // exactly the same as the Simple List
        func() fyne.CanvasObject {
            return container.NewBorder(
              nil, // TOP
              nil, // BOTTOM
              nil, // Left
              // LEFT
              widget.NewCheck("", func(b bool) {}),
              // REST
              widget.NewLabel(""),
            )
        },
        // func that is called for each item in the list and allows
        // but this time we get the actual DataItem we need to cast
        func(di binding.DataItem, o fyne.CanvasObject) {
            ctr, _ := o.(*fyne.Container)
            // ideally we should check `ok` for each one of those casting
            // but we know that they are those types for sure
            l := ctr.Objects[0].(*widget.Label)
            c := ctr.Objects[1].(*widget.Check)
            /*
              diu, _ := di.(binding.Untyped).Get()
              tag := diu.(models.Tag)
            */
            tag := models.NewTagFromDataItem(di)
            l.SetText(tag.Name)
            c.SetChecked(tag.Selected)

//             fmt.Println(tag.Selected)

            if(tag.Selected == true){
                fmt.Println(tag.Name)
            }
        },
    )

    tagBtn := widget.NewButton("tag    (with selected tags and files)", func() {
//       tags.Append(models.NewTag(newtagNameTxt.Text))
//       fmt.Printf("Type of Arg = %s\n", newtagNameTxt.Text)
//       fmt.Println(tags.Text)
//       newtagsNameTxt.Text = ""

//         fmt.Println(selectedTags)

//         // Does not work
//         for _, tag := range tags {
//             if(tag.Selected == true){
//                 fmt.Println(tag.Name)
//             }
//         }

//         // Works, but only gives the initial data, not any changes
//         for _, tag := range data {
//             fmt.Println(tag)
//         }

//         fmt.Println(tags.Get())

//////////////////////
        list, _ := tags.Get()
//         fmt.Println(list)

        for _, tag := range list {
//             fmt.Println(tag)
            fmt.Println(tag)
        }
//////////////////////

    })

    delBtn := widget.NewButton("delete selected tags", func() {
//       tags.Append(models.NewTag(newtagNameTxt.Text))
      newtagNameTxt.Text = ""
//////////////////////
      // to remove all tags
      list, _ := tags.Get()
      list = list[:0]
      fmt.Println(list)
      tags.Set(list)
//////////////////////
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
                    delBtn, // BOTTOM
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
                        selectedTags,
                        widget.NewLabel("placeholder"),
                        widget.NewLabel("placeholder"),
                        widget.NewLabel("placeholder"),
                    ),
                    container.NewBorder(
                        // TOP
                        tagBtn,
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
