package main

import (
	"fmt"
	"modules/pointCloudDecoder"
	"runtime"
	"strconv"
)

var jsonFilePath string = "./example_files/SeaverSignPCJSON.json"
var convertedJsonFileName string = "pointCloud.JSON"
var convertedJsonFilePath string = "./pointCloud.JSON"

func main() {
	//commented line for going from meshroom json data to a cleaned up version this program uses
	//var convertedJsonFilePath string = pointCloudDecoder.GenerateFloatJson(jsonFilePath, convertedJsonFileName)
	var pointData pointCloudDecoder.PointData = pointCloudDecoder.DecodeFromFloatJsonFromPath(convertedJsonFilePath)
	fmt.Println(strconv.Itoa(len(pointData.Points)))

	runtime.GC()
}
