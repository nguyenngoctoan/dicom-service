package main

import (
    "flag"
    "fmt"
    "log"
    "os"
    "dicom-service/pkg/dicomhandler"
)

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Expected 'upload', 'attribute', or 'convert' subcommands")
        os.Exit(1)
    }

    switch os.Args[1] {
    case "upload":
        handleUpload(os.Args[2:])
    case "attribute":
        handleAttribute(os.Args[2:])
    case "convert":
        handleConvert(os.Args[2:])
    default:
        fmt.Println("Expected 'upload', 'attribute', or 'convert' subcommands")
    }
}

func handleUpload(args []string) {
    uploadCmd := flag.NewFlagSet("upload", flag.ExitOnError)
    filePath := uploadCmd.String("file", "", "Path to the DICOM file")
    uploadCmd.Parse(args)

    if *filePath == "" {
        fmt.Println("You must specify a file path with -file")
        uploadCmd.PrintDefaults()
        os.Exit(1)
    }

    fmt.Println("Uploading file:", *filePath)
    result, err := dicomhandler.UploadDICOMCLI(*filePath)
    if err != nil {
        log.Fatalf("Error uploading file: %v", err)
    }
    fmt.Println("Upload successful, file ID:", result)
}

func handleAttribute(args []string) {
    attributeCmd := flag.NewFlagSet("attribute", flag.ExitOnError)
    fileID := attributeCmd.String("id", "", "File ID of the DICOM file")
    tag := attributeCmd.String("tag", "", "DICOM tag to retrieve (e.g., 0010,0010)")
    attributeCmd.Parse(args)

    if *fileID == "" || *tag == "" {
        fmt.Println("You must specify both -id and -tag")
        attributeCmd.PrintDefaults()
        os.Exit(1)
    }

    fmt.Printf("Retrieving attribute %s for file ID %s\n", *tag, *fileID)
    attribute, err := dicomhandler.GetAttributeCLI(*fileID, *tag)
    if err != nil {
        log.Fatalf("Error retrieving attribute: %v", err)
    }
    fmt.Println("Attribute value:", attribute)
}

func handleConvert(args []string) {
    convertCmd := flag.NewFlagSet("convert", flag.ExitOnError)
    fileID := convertCmd.String("id", "", "File ID of the DICOM file to convert to PNG")
    convertCmd.Parse(args)

    if *fileID == "" {
        fmt.Println("You must specify a file ID with -id")
        convertCmd.PrintDefaults()
        os.Exit(1)
    }

    fmt.Println("Converting file ID:", *fileID)
    path, err := dicomhandler.ConvertToPNGCLI(*fileID)
    if err != nil {
        log.Fatalf("Error converting file: %v", err)
    }
    fmt.Println("Converted file saved to:", path)
}
