package main

import (
	"fmt"
	"modules/pointCloudDecoder"
	"strconv"
)

var jsonFilePath string = "./example_files/SeaverSignPCJSON.json"

func main() {
	var data pointCloudDecoder.PointCloud = pointCloudDecoder.DecodeFromPath(jsonFilePath)
	fmt.Println("Data Structure Size: " + strconv.Itoa(len(data.Structures)))
}
