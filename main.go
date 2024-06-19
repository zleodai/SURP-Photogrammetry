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

var jsonFilePath string = "./example_files/footballPCJSON.json"
var convertedJsonFileName string = "pointCloud.JSON"
var convertedJsonFilePath string = "./pointCloud.JSON"

var defaultVoxelSize float64 = 0.01

func main() {
	// commented line for going from meshroom json data to a cleaned up version this program uses
	//var convertedJsonFilePath string = pointCloudDecoder.GenerateFloatJson(jsonFilePath, convertedJsonFileName)
	runtime.GC()
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
	pointData.Points = voxelMesher.PointcloudPreprocessFilter(xMinMax, yMinMax, zMinMax, pointData.Points, 100, .01)
	voxels := voxelMesher.MinMaxMesh(xMinMax, yMinMax, zMinMax, pointData.Points, defaultVoxelSize, true)
	runtime.GC()
	// voxelMesher.GenerateVoxelJson(voxels, defaultVoxelSize)

	// xSize := 10
	// ySize := 10
	// zSize := 10

	// voxels := make([][][]uint8, xSize)
	// for i := 0; i < len(voxels); i++ {
	// 	voxels[i] = make([][]uint8, ySize)
	// 	for j := 0; j < len(voxels[i]); j++ {
	// 		voxels[i][j] = make([]uint8, zSize)
	// 	}
	// }

	// voxels[5][5][5] = 10
	// voxels[6][5][5] = 10
	// voxels[5][6][5] = 10
	// voxels[6][6][5] = 10
	// voxels[5][6][6] = 10
	// voxels[6][6][6] = 10
	// voxels[5][6][7] = 10
	// voxels[6][6][7] = 10
	// voxels[5][6][8] = 10
	// voxels[6][6][8] = 10

	faces := greedyMesher.GreedyMesh(voxels, 2)
	vertices := objExporter.GetVerticesFromFaces(faces)
	fmt.Print(vertices)
	runtime.GC()

	// voxelMesher.GenerateVoxelJson(voxels, defaultVoxelSize)
}
