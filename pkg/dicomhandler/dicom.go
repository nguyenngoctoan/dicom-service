package dicomhandler

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"dicom-service/pkg/converter"
	"dicom-service/pkg/storage"
	"github.com/gin-gonic/gin"
	"github.com/suyashkumar/dicom"
	"github.com/suyashkumar/dicom/pkg/tag"
)

// UploadDICOM handles the web request to upload a DICOM file.
func UploadDICOM(c *gin.Context) {
	if c.MustGet("role").(string) != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file provided"})
		return
	}

	filePath := filepath.Join(os.TempDir(), file.Filename)
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	fileID, err := storage.SaveFile(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store file"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "File uploaded successfully", "file_id": fileID})
}

// GetDICOMAttribute retrieves a specific DICOM attribute.
func GetDICOMAttribute(c *gin.Context) {
	fileID := c.Query("file_id")
	dicomTag := c.Query("tag")

	if fileID == "" || dicomTag == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing file_id or tag"})
		return
	}

	attribute, err := extractAttribute(fileID, dicomTag)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"attribute": attribute})
}

// ConvertDICOMToPNG converts a DICOM file to a PNG image.
func ConvertDICOMToPNG(c *gin.Context) {
	fileID := c.Query("file_id")
	if fileID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing file_id"})
		return
	}

	pngPath, err := convertToPNG(fileID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.File(pngPath)
}

// UploadDICOMCLI handles the CLI uploading of a DICOM file.
func UploadDICOMCLI(filePath string) (string, error) {
	fileID, err := storage.SaveFile(filePath)
	if err != nil {
		return "", err
	}

	return fileID, nil
}

// GetAttributeCLI is the CLI version to get a DICOM attribute.
func GetAttributeCLI(fileID, dicomTag string) (string, error) {
	return extractAttribute(fileID, dicomTag)
}

// ConvertToPNGCLI is the CLI version to convert a DICOM file to PNG.
func ConvertToPNGCLI(fileID string) (string, error) {
	return convertToPNG(fileID)
}

// HealthCheck provides health status of the service.
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "service is running"})
}

func extractAttribute(fileID, dicomTag string) (string, error) {
	filePath, err := storage.RetrieveFile(fileID)
	if err != nil {
		return "", err
	}

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	parsedData, err := dicom.Parse(bytes.NewReader(data), int64(len(data)), nil)
	if err != nil {
		return "", err
	}

	dicomTagParsed, err := parseDICOMTag(dicomTag)
	if err != nil {
		return "", err
	}

	element, err := parsedData.FindElementByTag(dicomTagParsed)
	if err != nil {
		return "", err
	}

	return fmt.Sprint(element.Value), nil
}

func convertToPNG(fileID string) (string, error) {
	filePath, err := storage.RetrieveFile(fileID)
	if err != nil {
		return "", err
	}

	pngPath, err := converter.ConvertToPNG(filePath)
	if err != nil {
		return "", err
	}

	return pngPath, nil
}

func parseDICOMTag(tagString string) (tag.Tag, error) {
	parts := strings.Split(tagString, ",")
	if len(parts) != 2 {
		return tag.Tag{}, fmt.Errorf("invalid DICOM tag format")
	}

	group, err := strconv.ParseUint(parts[0], 16, 16)
	if err != nil {
		return tag.Tag{}, fmt.Errorf("invalid group number in DICOM tag")
	}

	element, err := strconv.ParseUint(parts[1], 16, 16)
	if err != nil {
		return tag.Tag{}, fmt.Errorf("invalid element number in DICOM tag")
	}

	return tag.Tag{Group: uint16(group), Element: uint16(element)}, nil
}
