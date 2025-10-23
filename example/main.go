package main

// go build -ldflags="-s -w"

import (
	// "fmt"
	"os"

	"github.com/errz99/qtablam"
	"github.com/mappu/miqt/qt"
)

var (
	app    *qt.QApplication
	window *qt.QMainWindow
)

func main() {
	initGui()
}

func initGui() {
	app = qt.NewQApplication(os.Args)
	mainWindow()

	qt.QApplication_Exec()
}

func mainWindow() {
	window = qt.NewQMainWindow2()

	vbox := qt.NewQVBoxLayout2()
	winWidget := qt.NewQWidget2()
	winWidget.SetLayout(vbox.QLayout)
	window.SetCentralWidget(winWidget)

	titles := []string{"Uno", "Dos", "Tres", "Cuatro", "Cinco"}
	data := [][]string{
		{"Uno1", "Dos1", "Tres1", "Cuatro1", "Cinco1"},
		{"Uno2", "Dos2", "Tres2", "Cuatro2", "Cinco2"},
		{"Uno3", "Dos3", "Tres3", "Cuatro3", "Cinco3"},
		{"Uno4", "Dos4", "Tres4", "Cuatro4", "Cinco4"},
		{"Uno5", "Dos5", "Tres5", "Cuatro5", "Cinco5"},
	}

	area := qtablam.NewQTablam(titles, data)
	vbox.AddWidget(area.QWidget)

	buttonsBox := qt.NewQHBoxLayout2()
	boxWidget := qt.NewQWidget2()
	boxWidget.SetLayout(buttonsBox.QLayout)
	vbox.AddWidget3(boxWidget, 0, qt.AlignBottom)

	buttonOne := qt.NewQPushButton3("Add Row")
	buttonOne.SetMinimumWidth(120)
	buttonsBox.AddWidget3(buttonOne.QWidget, 0, qt.AlignCenter)
	buttonTwo := qt.NewQPushButton3("Remove Row")
	buttonTwo.SetMinimumWidth(120)
	buttonsBox.AddWidget3(buttonTwo.QWidget, 0, qt.AlignCenter)

	buttonOne.OnClicked(func() {
		area.AddRow([]string{"Uno6", "Dos6", "Tres6", "Cuatro6", "Cinco6"})
	})

	buttonTwo.OnClicked(func() {
		// area.RemoveRow(0)
		area.RemoveActiveRow()
	})

	window.Show()
}
