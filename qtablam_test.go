package qtablam

import (
	"fmt"
	"os"
	"testing"

	// "github.com/errz99/qtablam"
	"github.com/mappu/miqt/qt"
)

func TestTablam(t *testing.T) {
	fmt.Println("testing")
	qt.NewQApplication(os.Args)
	// mainWindow()
	window := qt.NewQMainWindow2()
	window.SetMinimumSize2(450, 300)

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

	vbox.AddWidget3(NewQTablam(titles, data), 0, qt.AlignTop)

	window.Show()

	qt.QApplication_Exec()
}
