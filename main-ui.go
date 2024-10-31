package main

import (
	"fmt"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
	"github.com/tofl/pngify/image"
)

func main() {
	a := app.New()
	w := a.NewWindow("PNGify UI")
	w.Resize(fyne.NewSize(600, 400))

	// Encode Tab
	encodeOption := widget.NewRadioGroup([]string{"Text", "File"}, func(string) {})
	encodeOption.SetSelected("Text")

	textEntry := widget.NewMultiLineEntry()
	textEntry.SetPlaceHolder("Enter text to encode")

	fileLabel := widget.NewLabel("No file selected")
	selectFileButton := widget.NewButton("Select File", func() {
		openFileDialog := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil || reader == nil {
				return
			}
			fileLabel.SetText(reader.URI().Path())
		}, w)
		openFileDialog.SetFilter(storage.NewExtensionFileFilter([]string{".txt", ".png"}))
		openFileDialog.Show()
	})

	encodeButton := widget.NewButton("Encode", func() {
		if encodeOption.Selected == "Text" {
			text := textEntry.Text
			if text == "" {
				dialog.ShowInformation("Error", "Please enter text to encode", w)
				return
			}
			img := image.NewImage([]byte(text))
			img.MakeImage()
			dialog.ShowInformation("Success", "Text encoded to output.png", w)
		} else if encodeOption.Selected == "File" {
			filePath := fileLabel.Text
			if filePath == "No file selected" {
				dialog.ShowInformation("Error", "Please select a file to encode", w)
				return
			}
			content, err := os.ReadFile(filePath)
			if err != nil {
				dialog.ShowInformation("Error", fmt.Sprintf("Failed to read file: %v", err), w)
				return
			}
			img := image.NewImage(content)
			img.MakeText([]byte("filename"), []byte(filePath))
			img.MakeImage()
			dialog.ShowInformation("Success", "File encoded to output.png", w)
		}
	})

	encodeContent := container.NewVBox(
		encodeOption,
		textEntry,
		container.NewHBox(selectFileButton, fileLabel),
		encodeButton,
	)

	// Decode Tab
	decodeFileLabel := widget.NewLabel("No file selected")
	selectDecodeFileButton := widget.NewButton("Select Image File", func() {
		openFileDialog := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil || reader == nil {
				return
			}
			decodeFileLabel.SetText(reader.URI().Path())
		}, w)
		openFileDialog.SetFilter(storage.NewExtensionFileFilter([]string{".png"}))
		openFileDialog.Show()
	})

	decodedTextEntry := widget.NewMultiLineEntry()
	decodedTextEntry.SetPlaceHolder("Decoded text will appear here")
	decodedTextEntry.Disable()

	decodeButton := widget.NewButton("Decode", func() {
		filePath := decodeFileLabel.Text
		if filePath == "No file selected" {
			dialog.ShowInformation("Error", "Please select an image file to decode", w)
			return
		}
		f, err := os.Open(filePath)
		if err != nil {
			dialog.ShowInformation("Error", fmt.Sprintf("Failed to open file: %v", err), w)
			return
		}
		defer f.Close()

		data, fileName := image.Decode(f)
		if fileName == "" {
			decodedTextEntry.SetText(data)
		} else {
			saveDialog := dialog.NewFileSave(func(writer fyne.URIWriteCloser, err error) {
				if err != nil || writer == nil {
					return
				}
				_, err = writer.Write([]byte(data))
				if err != nil {
					dialog.ShowInformation("Error", fmt.Sprintf("Failed to save file: %v", err), w)
					return
				}
				writer.Close()
				dialog.ShowInformation("Success", "File saved successfully", w)
			}, w)
			saveDialog.SetFileName(fileName)
			saveDialog.Show()
		}
	})

	decodeContent := container.NewVBox(
		selectDecodeFileButton,
		decodeFileLabel,
		decodeButton,
		decodedTextEntry,
	)

	tabs := container.NewAppTabs(
		container.NewTabItem("Encode", encodeContent),
		container.NewTabItem("Decode", decodeContent),
	)

	w.SetContent(tabs)
	w.ShowAndRun()
}
