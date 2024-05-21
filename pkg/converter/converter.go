package converter

import (
	"bytes"
	"errors"
	"image"
	"image/color"
	"image/png"
	"io/ioutil"
	"os"

	"github.com/suyashkumar/dicom"
	"github.com/suyashkumar/dicom/pkg/tag"
)

// ConvertToPNG converts a DICOM file to a PNG image.
func ConvertToPNG(dicomFilePath string) (string, error) {
	data, err := ioutil.ReadFile(dicomFilePath)
	if err != nil {
		return "", err
	}

	parsedData, err := dicom.Parse(bytes.NewReader(data), int64(len(data)), nil)
	if err != nil {
		return "", err
	}

	img, err := extractImage(&parsedData)
	if err != nil {
		return "", err
	}

	pngPath := dicomFilePath + ".png"
	file, err := os.Create(pngPath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	if err := png.Encode(file, img); err != nil {
		return "", err
	}

	return pngPath, nil
}

func extractImage(parsedData *dicom.Dataset) (image.Image, error) {
	// Extract the PixelData element
	pixelDataElement, err := parsedData.FindElementByTag(tag.PixelData)
	if err != nil {
		return nil, errors.New("no pixel data found in DICOM data")
	}

	pixelData, ok := pixelDataElement.Value.GetValue().([]byte)
	if !ok || len(pixelData) == 0 {
		return nil, errors.New("no frames found in pixel data")
	}

	// Assuming Rows and Columns are present and valid in the dataset
	rowsElement, err := parsedData.FindElementByTag(tag.Rows)
	if err != nil {
		return nil, errors.New("no rows data found in DICOM data")
	}
	rows := rowsElement.Value.GetValue().(int)

	columnsElement, err := parsedData.FindElementByTag(tag.Columns)
	if err != nil {
		return nil, errors.New("no columns data found in DICOM data")
	}
	columns := columnsElement.Value.GetValue().(int)

	img := image.NewGray(image.Rect(0, 0, columns, rows))
	for y := 0; y < rows; y++ {
		for x := 0; x < columns; x++ {
			pixelValue := pixelData[y*columns+x]
			img.SetGray(x, y, color.Gray{Y: pixelValue})
		}
	}

	return img, nil
}
