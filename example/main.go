package main

// go build -ldflags="-s -w"

import (
	"fmt"
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

	// Boxes
	buttonsBox1 := qt.NewQHBoxLayout2()
	buttonsBox1.SetContentsMargins(0, 0, 0, 0)
	boxWidget1 := qt.NewQWidget2()
	boxWidget1.SetLayout(buttonsBox1.QLayout)
	vbox.AddWidget3(boxWidget1, 0, qt.AlignBottom)

	buttonsBox2 := qt.NewQHBoxLayout2()
	buttonsBox2.SetContentsMargins(0, 0, 0, 0)
	boxWidget2 := qt.NewQWidget2()
	boxWidget2.SetLayout(buttonsBox2.QLayout)
	vbox.AddWidget3(boxWidget2, 0, qt.AlignBottom)

	// Buttons 1
	buttonOne := qt.NewQPushButton3("Add Row")
	buttonOne.SetMinimumWidth(120)
	buttonsBox1.AddWidget3(buttonOne.QWidget, 0, qt.AlignCenter)

	buttonTwo := qt.NewQPushButton3("Remove Row")
	buttonTwo.SetMinimumWidth(120)
	buttonsBox1.AddWidget3(buttonTwo.QWidget, 0, qt.AlignCenter)

	buttonThree := qt.NewQPushButton3("Edit Cell")
	buttonThree.SetMinimumWidth(120)
	buttonsBox1.AddWidget3(buttonThree.QWidget, 0, qt.AlignCenter)

	buttonOne.OnClicked(func() {
		area.AddRow([]string{"Uno6", "Dos6", "Tres6", "Cuatro6", "Cinco6"})
	})

	buttonTwo.OnClicked(func() {
		// area.RemoveRow(0)
		area.RemoveActiveRow()
	})

	buttonThree.OnClicked(func() {
		area.EditCell(0, 0, "Cambio!")
	})

	// Buttons 2
	buttonFour := qt.NewQPushButton3("Row Texts")
	buttonFour.SetMinimumWidth(120)
	buttonsBox2.AddWidget3(buttonFour.QWidget, 0, qt.AlignCenter)

	buttonFour.OnClicked(func() {
		texts := area.RowTexts(0)
		fmt.Println(texts)
		text := area.CellText(1, 2)
		fmt.Println(text)
	})

	window.Show()
}
