package component

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func searchOrderView(_ fyne.Window) fyne.CanvasObject {
	transactionId := widget.NewEntry()
	transactionId.SetPlaceHolder("请输入用户提供的订单号")

	resultLabel := widget.NewLabel("")

	form := &widget.Form{
		Items: []*widget.FormItem{},
		OnCancel: func() {
			transactionId.SetText("")
			fmt.Println("Cancelled")
		},
		OnSubmit: func() {
			fmt.Println("Form submitted")
			resultLabel.SetText(transactionId.Text)
			fyne.CurrentApp().SendNotification(&fyne.Notification{
				Title:   "Form for: " + transactionId.Text,
				Content: "",
			})
		},
		CancelText: "重置",
		SubmitText: "查找",
	}

	form.Append("用户订单号", transactionId)

	//return form
	return container.NewVBox(form, resultLabel)
}
