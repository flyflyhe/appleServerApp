package component

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

func searchOrderView(_ fyne.Window) fyne.CanvasObject {
	transactionId := widget.NewEntry()
	transactionId.SetPlaceHolder("请输入用户提供的订单号")

	largeText := widget.NewMultiLineEntry()
	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "用户订单号", Widget: transactionId, HintText: "用户提供的订单号"},
		},
		OnCancel: func() {
			fmt.Println("Cancelled")
		},
		OnSubmit: func() {
			fmt.Println("Form submitted")
			fyne.CurrentApp().SendNotification(&fyne.Notification{
				Title:   "Form for: " + transactionId.Text,
				Content: largeText.Text,
			})
		},
	}
	form.Append("Message", largeText)
	return form
}

func searchOrderView2(_ fyne.Window) fyne.CanvasObject {
	name := widget.NewEntry()
	name.SetPlaceHolder("John Smith")

	largeText := widget.NewMultiLineEntry()

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Name", Widget: name, HintText: "Your full name"},
		},
		OnCancel: func() {
			fmt.Println("Cancelled")
		},
		OnSubmit: func() {
			fmt.Println("Form submitted")
			fyne.CurrentApp().SendNotification(&fyne.Notification{
				Title:   "Form for: " + name.Text,
				Content: largeText.Text,
			})
		},
	}
	form.Append("Message", largeText)
	return form
}
