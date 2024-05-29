package voxelMesher

import (
	"fmt"
	"math"
	"modules/pointCloudDecoder"
	"os"
	"runtime"
	"strconv"
)

func Test() byte {
	return 0
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func Mesh(xArray, yArray, zArray []pointCloudDecoder.Point, voxelSize float64) {
	var xSize int = int(math.Floor(math.Abs(xArray[0].X-xArray[len(xArray)-1].X) / voxelSize))
	var ySize int = int(math.Floor(math.Abs(yArray[0].Y-yArray[len(yArray)-1].Y) / voxelSize))
	var zSize int = int(math.Floor(math.Abs(zArray[0].Z-zArray[len(zArray)-1].Z) / voxelSize))
	fmt.Println("Voxel Grid X, Y, Z:", xSize, ySize, zSize)
	fmt.Println("Total Voxel Amount:", xSize*ySize*zSize)

	voxels := make([][][]int, xSize+1)
	for yIndex, _ := range voxels {
		voxels[yIndex] = make([][]int, ySize+1)
		for zIndex, _ := range voxels[yIndex] {
			voxels[yIndex][zIndex] = make([]int, zSize+1)
		}
	}

	for _, point := range xArray {
		xIndex := int(math.Floor(math.Abs(xArray[0].X-point.X) / voxelSize))
		yIndex := int(math.Floor(math.Abs(yArray[0].Y-point.Y) / voxelSize))
		zIndex := int(math.Floor(math.Abs(zArray[0].Z-point.Z) / voxelSize))
		voxels[xIndex][yIndex][zIndex] += 1
	}

	// GenerateVoxelJson(voxels, voxelSize)

	runtime.GC()
}

func GenerateVoxelJson(voxels [][][]int, voxelSize float64) {
	file, errs := os.Create("VoxelMatrix.JSON")
	if errs != nil {
		panic("Failed to write to file:" + errs.Error())
	}
	defer file.Close()

	_, errs = file.WriteString("{\"points\": \n	[")
	check(errs)

	for xIndex, xArray := range voxels {
		for yIndex, yArray := range xArray {
			for zIndex, Value := range yArray {
				if Value != 0 || xIndex == len(voxels)-1 && yIndex == len(xArray)-1 && zIndex == len(yArray)-1 {
					_, errs = file.WriteString("\n		{")
					check(errs)

					var xValue float64 = float64(xIndex) * voxelSize
					var yValue float64 = float64(yIndex) * voxelSize
					var zValue float64 = float64(zIndex) * voxelSize

					_, errs = file.WriteString("\n			\"x\": " + strconv.FormatFloat(xValue, 'f', -1, 64) + ",")
					check(errs)
					_, errs = file.WriteString("\n			\"y\": " + strconv.FormatFloat(yValue, 'f', -1, 64) + ",")
					check(errs)
					_, errs = file.WriteString("\n			\"z\": " + strconv.FormatFloat(zValue, 'f', -1, 64) + ",")
					check(errs)
					_, errs = file.WriteString("\n			\"Value\": " + strconv.Itoa(Value))
					check(errs)

					if xIndex == len(voxels)-1 && yIndex == len(xArray)-1 && zIndex == len(yArray)-1 {
						_, errs = file.WriteString("\n		}")
						check(errs)
					} else {
						_, errs = file.WriteString("\n		},")
						check(errs)
					}
				}
			}
		}
	}

	_, errs = file.WriteString("\n	]\n}")
	check(errs)
}
