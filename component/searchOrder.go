package component

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

func searchOrderView(_ fyne.Window) fyne.CanvasObject {
	transactionId := widget.NewEntry()
	transactionId.SetPlaceHolder("请输入用户提供的订单号")

	disabled := widget.NewRadioGroup([]string{"Option 1", "Option 2"}, func(string) {})
	disabled.Horizontal = true
	disabled.Disable()
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
				Content: "",
			})
		},
	}
	form.Append("", disabled)
	return form
}
