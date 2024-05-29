package main

import (
	"modules/pointCloudDecoder"
)

var jsonFilePath string = "./example_files/SeaverSignPCJSON.json"
var convertedJsonFileName string = "pointCloud.JSON"

func main() {
	pointCloudDecoder.GenerateFloatJson(jsonFilePath, convertedJsonFileName)

}
