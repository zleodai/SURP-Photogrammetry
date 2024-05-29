package main

import (
	"fmt"
	"modules/pointCloudDecoder"
	"runtime"
	"strconv"
)

var jsonFilePath string = "./example_files/SeaverSignPCJSON.json"
var convertedJsonFileName string = "pointCloud.JSON"

func main() {
	var convertedJsonFilePath string = pointCloudDecoder.GenerateFloatJson(jsonFilePath, convertedJsonFileName)
	var pointData pointCloudDecoder.PointData = pointCloudDecoder.DecodeFromFloatJsonFromPath(convertedJsonFilePath)
	fmt.Println(strconv.Itoa(len(pointData.Points)))

	runtime.GC()
}
