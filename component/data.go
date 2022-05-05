package component

import (
	"fyne.io/fyne/v2"
)

type AppView struct {
	Title, Intro string
	View         func(w fyne.Window) fyne.CanvasObject
}

var (
	AppViews = map[string]AppView{
		"welcome": {"Welcome", "", welcomeScreen},
		"canvas": {"Canvas",
			"See the canvas capabilities.",
			canvasScreen,
		},
		"searchOrder":              {"查询订单", "", searchOrderView},
		"searchTransactionHistory": {"查询历史", "", searchTransactionHistoryView},
	}

	//index tree
	AppViewsIndex = map[string][]string{
		//"":            {"welcome", "searchOrder", "canvas", "animations", "icons", "widgets", "collections", "containers", "dialogs", "windows", "binding", "advanced"},
		"": {"welcome", "searchOrder", "searchTransactionHistory"},
		//"collections": {"list", "table", "tree"},
		//"containers":  {"apptabs", "border", "box", "center", "doctabs", "grid", "scroll", "split"},
		//"widgets":     {"accordion", "button", "card", "entry", "form", "input", "progress", "text", "toolbar"},
	}
)
