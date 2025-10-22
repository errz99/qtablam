package qtablam

import (
	"runtime"

	"github.com/mappu/miqt/qt"
)

func DefaultFont() string {
	switch runtime.GOOS {
	case "darwin":
		return "Menlo"
	case "windows":
		return "Cascadia Mono"
	case "linux", "freebsd":
		return "DejaVu Sans Mono"
	default:
		return "mono"
	}
}

type MyFontData struct {
	font *qt.QFont
	size int
	w    int
	h    int
}

func newFontData() MyFontData {
	switch runtime.GOOS {
	case "darwin":
		return MyFontData{qt.NewQFont6(defFont, 12), 12, 0, 0}
	case "windows":
		return MyFontData{qt.NewQFont6(defFont, 10), 10, 0, 0}
	case "linux", "freebsd":
		return MyFontData{qt.NewQFont6(defFont, 8), 8, 0, 0}
	default:
		return MyFontData{qt.NewQFont6(defFont, 10), 10, 0, 0}
	}
}

func (fd *MyFontData) UpdateMetrics(rowSep int) {
	metrics := qt.NewQFontMetrics(fd.font)
	fd.w = metrics.Width("M")
	fd.h = metrics.Height() + rowSep*2
}

func (fd *MyFontData) UpdateSize(amount int) bool {
	if fd.size+amount >= 4 {
		fd.size += amount
		fd.font = qt.NewQFont6(defFont, fd.size)
		fd.UpdateMetrics(centerArea.rowSep)
		return true
	}
	return false
}
