# Welcome to dicom-service project

### Using the endpoints

upload:
```
curl -X POST http://localhost:8080/upload \                                                                           
     -H "Authorization: valid-token" \
     -F "file=@/Users/toan/dev/pocket-health/dicom-service/test/PA000001/ST000001/SE000001/IM000006"
```

convert:
```
curl -X GET "http://localhost:8080/convert?file_id=1a2cec607017bbd05e45dc5e1438372614320c414f3c50a462d35c7e3ff24ae4" \
     -H "Authorization: valid-token" \
     -o output.png
```

attribute:
```
curl -X GET "http://localhost:8080/attribute?file_id=61257f1650b81bc10ac3dce55d23680b89e9c2086f9a7cb187b64dd672937b99&tag=0010,0010" \                                 
     -H "Authorization: valid-token"
``` 

health:
```
curl -X GET http://localhost:8080/health
```        


### cmd/api/main.go:

This file will initialize and start the HTTP server.
It will setup the routes and handlers for the web API.
This main file will heavily rely on the dicomhandler package for processing requests.

### cmd/cli/main.go:

This file will parse command line arguments and execute the necessary operations based on the input commands.
It can be used for operations like uploading DICOM files from a CLI, fetching attributes, or converting files without a web interface.

### pkg/dicomhandler/dicom.go:

Contains logic for reading, processing, and extracting information from DICOM files.
Functions here will be used both by the API and CLI to ensure DRY (Don't Repeat Yourself) principles.

### pkg/storage/storage.go:

Abstraction for storage operations, which could be local file system interactions or cloud storage (like Azure Blob Storage).
This makes it easier to switch storage backends without affecting other parts of the application.

### pkg/converter/converter.go:

Contains logic to convert DICOM images to PNG format.
This separation allows the conversion logic to be maintained independently of the main application logic.

### internal/config/config.go:

Manages configuration and environment variables.
This is useful for setting up and reading configuration that affects the entire application, like API keys, database connection info, or external service URLs.