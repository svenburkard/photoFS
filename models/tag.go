package models

import (
    "fmt"

    "fyne.io/fyne/v2/data/binding"
)

type Tag struct {
    Name      string
    Selected  bool
}

func NewTag(name string) Tag {
    return Tag{name, false}
}

func NewTagFromDataItem(di binding.DataItem) Tag {
    v, _ := di.(binding.Untyped).Get()
    return v.(Tag)
}

func (t Tag) String() string {
    return fmt.Sprintf("%s  - %t", t.Name, t.Selected)
}


