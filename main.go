package tuifp

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"io/ioutil"
	"path/filepath"
	"strings"
)

type TuiFilePicker struct {
	infoView     *tview.TextView
	footerView   *tview.TextView
	headerView   *tview.Flex
	listView     *tview.List
	currentPath  string
	selectedPath string
}

func NewTuiFilePicker() *TuiFilePicker {
	return &TuiFilePicker{}
}

func (f *TuiFilePicker) changeDir(dir string) error {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}
	f.listView.Clear()
	for _, file := range files {
		str := file.Name()
		if file.IsDir() {
			str += "/"
		}
		f.listView = f.listView.AddItem(str, "", 0, nil)
	}
	f.currentPath = dir
	f.infoView.SetText("current dir: " + dir)
	return nil
}
func (f *TuiFilePicker) Pick(startDir string) (string, error) {
	var err error
	app := tview.NewApplication()

	// UI components
	f.listView = tview.NewList().
		ShowSecondaryText(false).
		SetSelectedFocusOnly(true)
	f.footerView = tview.NewTextView().SetText("Press 'q' to quit")
	f.headerView = tview.NewFlex()
	f.infoView = tview.NewTextView()

	button := tview.NewButton(".. <parent dir>")
	button.Box = tview.NewBox()
	button.SetSelectedFunc(func() {
		err = f.changeDir(filepath.Dir(f.currentPath))
		if err != nil {
			app.Stop()
		}
	})
	f.headerView.AddItem(button, 16, 0, true)

	// build UI
	pages := tview.NewPages()
	body := tview.NewFlex().SetDirection(tview.FlexRow).AddItem(f.listView, 0, 1, true)
	splitterView := tview.NewTextView().SetText("----------------------------")
	mainFlex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(f.infoView, 1, 0, false).
		AddItem(splitterView, 1, 0, false).
		AddItem(f.headerView, 1, 0, true).
		AddItem(body, 0, 1, true).
		AddItem(f.footerView, 1, 0, false)
	pages.AddPage("main", mainFlex, true, true)

	// if down pressed on header
	f.headerView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyDown {
			app.SetFocus(f.listView)
			return nil
		}
		return event
	})

	// if up pressed on body
	body.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyUp {
			if f.listView.GetCurrentItem() == 0 {
				app.SetFocus(button)
				return nil
			}
		}
		return event
	})

	// if q pressed
	pages.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyRune {
			switch event.Rune() {
			case 'q':
				app.Stop()
			}
		}
		return event
	})

	// if item selected
	f.listView.SetSelectedFunc(func(index int, selectedStr string, secondary string, code rune) {
		f.infoView.SetText("Selected: " + selectedStr)
		if strings.HasSuffix(selectedStr, "/") {
			err = f.changeDir(filepath.Join(f.currentPath, selectedStr))
		} else {
			f.selectedPath = filepath.Join(f.currentPath, selectedStr)
			app.Stop()
		}
	})

	// -------------------
	// main
	// -------------------
	err = f.changeDir(startDir)

	err = app.SetRoot(pages, true).Run()
	if err != nil {
		return "", err
	}
	return f.selectedPath, nil
}
