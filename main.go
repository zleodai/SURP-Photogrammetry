package main

import (
	"fmt"
	"modules/greedyMesher"
	"modules/objExporter"
	"modules/pointCloudDecoder"
	"modules/voxelMesher"
	"runtime"
	"strconv"
)

var jsonFilePath string = "./example_files/SeaverSignPCJSON.json"
var convertedJsonFileName string = "pointCloud.JSON"

func main() {
	var convertedJsonFilePath string = pointCloudDecoder.GenerateFloatJson(jsonFilePath, convertedJsonFileName)
	var pointData pointCloudDecoder.PointData = pointCloudDecoder.DecodeFromFloatJsonFromPath(convertedJsonFilePath)
	fmt.Println(strconv.Itoa(len(pointData.Points)))

	greedyMesher.Test()
	objExporter.Test()
	voxelMesher.Test()

	runtime.GC()
}
