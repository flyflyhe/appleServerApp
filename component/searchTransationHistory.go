package component

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/flyflyhe/appleServerApp/layouts"
	"github.com/flyflyhe/appleServerApp/services/apple"
	"github.com/flyflyhe/appleServerApp/services/jsonHelper"
)

func searchTransactionHistoryView(_ fyne.Window) fyne.CanvasObject {
	transactionId := widget.NewEntry()
	transactionId.SetPlaceHolder("请输入用户原始订单号")

	resultText := widget.NewMultiLineEntry()

	form := &widget.Form{
		Items: []*widget.FormItem{},
		OnCancel: func() {
			transactionId.SetText("")
			resultText.SetText("")
			fmt.Println("Cancelled")
		},
		OnSubmit: func() {
			fmt.Println("Form submitted")
			orderArr, err := apple.GetTransactionHistory(transactionId.Text, apple.GetAppleJwtToken(), "")
			if err != nil {
				resultText.SetText(err.Error())
			} else {
				orderInfoJsonPretty := ""
				for _, v := range orderArr {
					t, err := jsonHelper.PrettyString(v)
					if err != nil {
						resultText.SetText(err.Error())
						return
					} else {
						orderInfoJsonPretty += t
					}
				}
				resultText.SetText(orderInfoJsonPretty)
			}
			// fyne.CurrentApp().SendNotification(&fyne.Notification{
			// 	Title:   "Form for: " + transactionId.Text,
			// 	Content: "",
			// })
		},
		CancelText: "重置",
		SubmitText: "查找",
	}

	form.Append("用户原始订单号", transactionId)

	//return form
	layout := layouts.NewVBoxLayout()
	layouts.SetObjConfigMap(resultText, &layouts.Size{Width: 0, Height: 300})
	//return form
	return container.New(layout, form, resultText)
}
