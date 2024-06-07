package main

import (
	"fmt"
	"modules/greedyMesher"
	"modules/objExporter"
	"modules/pointCloudDecoder"
	"modules/pointSorter"
	"modules/voxelMesher"
	"runtime"
	"strconv"
)

// var jsonFilePath string = "./example_files/footballPCJSON.json"
// var convertedJsonFileName string = "pointCloud.JSON"
var convertedJsonFilePath string = "./pointCloud.JSON"

var defaultVoxelSize float64 = 0.005

func main() {
	// commented line for going from meshroom json data to a cleaned up version this program uses
	// var convertedJsonFilePath string = pointCloudDecoder.GenerateFloatJson(jsonFilePath, convertedJsonFileName)
	var pointData pointCloudDecoder.PointData = pointCloudDecoder.DecodeFromFloatJsonFromPath(convertedJsonFilePath)
	fmt.Printf("Running Program\n%s Points Loaded\n", strconv.Itoa(len(pointData.Points)))

	runtime.GC()

	objExporter.Test()
	pointSorter.Test()

	// fmt.Printf("\nxMinValue: %s", strconv.FormatFloat(xArray[0].X, 'f', -1, 64))
	// fmt.Printf("\nxMaxValue: %s", strconv.FormatFloat(xArray[len(xArray)-1].X, 'f', -1, 64))
	// fmt.Printf("\nyMinValue: %s", strconv.FormatFloat(yArray[0].Y, 'f', -1, 64))
	// fmt.Printf("\nyMaxValue: %s", strconv.FormatFloat(yArray[len(yArray)-1].Y, 'f', -1, 64))
	// fmt.Printf("\nzMinValue: %s", strconv.FormatFloat(zArray[0].Z, 'f', -1, 64))
	// fmt.Printf("\nzMaxValue: %s", strconv.FormatFloat(zArray[len(zArray)-1].Z, 'f', -1, 64))
	// xArray, yArray, zArray := pointSorter.SortPointData(pointData)
	// voxelMesher.Mesh(xArray, yArray, zArray, defaultVoxelSize)

	xMinMax, yMinMax, zMinMax := pointSorter.MinMaxPoints(pointData)
	voxels := voxelMesher.MinMaxMesh(xMinMax, yMinMax, zMinMax, pointData.Points, defaultVoxelSize, true)

	greedyMesher.GreedyMesh(voxels, 2)

	// voxelMesher.GenerateVoxelJson(voxels, defaultVoxelSize)
}
