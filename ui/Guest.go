package ui

import (
	"image"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
)

func Guest() {
	w = App.NewWindow("Screem :: Guest")
	w.Resize(fyne.NewSize(800, 600))
	img = canvas.NewImageFromImage(image.NewRGBA(image.Rect(0, 0, 200, 200)))
	w.SetContent(img)
	w.ShowAndRun()
}

func UpdateScreen(newImg *image.Image) {
    img.Image = *newImg
    img.Refresh()
}
