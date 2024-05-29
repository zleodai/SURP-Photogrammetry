package main

import (
	"fmt"
	"math"
	"modules/greedyMesher"
	"modules/objExporter"
	"modules/pointCloudDecoder"
	"modules/pointSorter"
	"modules/voxelMesher"
	"runtime"
	"strconv"
)

// var jsonFilePath string = "./example_files/SeaverSignPCJSON.json"
// var convertedJsonFileName string = "pointCloud.JSON"
var convertedJsonFilePath string = "./pointCloud.JSON"

var defaultVoxelSize float64 = 0.5

func main() {
	fmt.Printf("\n Running Program %s", strconv.Itoa(0))
	//commented line for going from meshroom json data to a cleaned up version this program uses
	//var convertedJsonFilePath string = pointCloudDecoder.GenerateFloatJson(jsonFilePath, convertedJsonFileName)
	var pointData pointCloudDecoder.PointData = pointCloudDecoder.DecodeFromFloatJsonFromPath(convertedJsonFilePath)

	runtime.GC()

	greedyMesher.Test()
	objExporter.Test()
	voxelMesher.Test()
	pointSorter.Test()

	xArray, yArray, zArray := pointSorter.SortPointData(pointData)

	// fmt.Printf("\nxMinValue: %s", strconv.FormatFloat(xArray[0].X, 'f', -1, 64))
	// fmt.Printf("\nxMaxValue: %s", strconv.FormatFloat(xArray[len(xArray)-1].X, 'f', -1, 64))
	// fmt.Printf("\nyMinValue: %s", strconv.FormatFloat(yArray[0].Y, 'f', -1, 64))
	// fmt.Printf("\nyMaxValue: %s", strconv.FormatFloat(yArray[len(yArray)-1].Y, 'f', -1, 64))
	// fmt.Printf("\nzMinValue: %s", strconv.FormatFloat(zArray[0].Z, 'f', -1, 64))
	// fmt.Printf("\nzMaxValue: %s", strconv.FormatFloat(zArray[len(zArray)-1].Z, 'f', -1, 64))

	var xSize int = int(math.Floor(math.Abs(xArray[0].X-xArray[len(xArray)-1].X) / defaultVoxelSize))
	var ySize int = int(math.Floor(math.Abs(yArray[0].Y-yArray[len(yArray)-1].Y) / defaultVoxelSize))
	var zSize int = int(math.Floor(math.Abs(zArray[0].Z-zArray[len(zArray)-1].Z) / defaultVoxelSize))
}
