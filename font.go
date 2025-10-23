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
	Font *qt.QFont
	Size int
	W    int
	H    int
}

func NewFontData() MyFontData {
	switch runtime.GOOS {
	case "darwin":
		return MyFontData{qt.NewQFont6(DefFont, 12), 12, 0, 0}
	case "windows":
		return MyFontData{qt.NewQFont6(DefFont, 10), 10, 0, 0}
	case "linux", "freebsd":
		return MyFontData{qt.NewQFont6(DefFont, 10), 10, 0, 0}
	default:
		return MyFontData{qt.NewQFont6(DefFont, 10), 10, 0, 0}
	}
}

func (fd *MyFontData) UpdateMetrics(rowSep int) {
	metrics := qt.NewQFontMetrics(fd.Font)
	fd.W = metrics.Width("M")
	fd.H = metrics.Height() + rowSep*2
}

func (fd *MyFontData) UpdateSize(rowSep, amount int) bool {
	if fd.Size+amount >= 4 {
		fd.Size += amount
		fd.Font = qt.NewQFont6(DefFont, fd.Size)
		fd.UpdateMetrics(rowSep)
		return true
	}
	return false
}
