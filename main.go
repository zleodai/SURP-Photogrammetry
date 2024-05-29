package main

import (
	"fmt"
	"modules/pointCloudDecoder"
	"os"
)

var jsonFilePath string = "./example_files/SeaverSignPCJSON.json"

func main() {
	file, errs := os.Create("pointCloud.JSON")
	if errs != nil {
		fmt.Println("Failed to create file:", errs)
		return
	}
	defer file.Close()

	var data pointCloudDecoder.PointCloud = pointCloudDecoder.DecodeFromPath(jsonFilePath)

	_, errs = file.WriteString("{\"points\": \n	[")
	if errs != nil {
		fmt.Println("Failed to write to file:", errs)
		return
	}

	for index, structure := range data.Structures {
		_, errs = file.WriteString("\n		{")
		if errs != nil {
			fmt.Println("Failed to write to file:", errs)
			return
		}
		for index, x := range structure.X {
			if index == 0 {
				_, errs = file.WriteString("\n			\"x\": " + x + ",")
				if errs != nil {
					fmt.Println("Failed to write to file:", errs)
					return
				}
			} else if index == 1 {
				_, errs = file.WriteString("\n			\"y\": " + x + ",")
				if errs != nil {
					fmt.Println("Failed to write to file:", errs)
					return
				}
			} else {
				_, errs = file.WriteString("\n			\"z\": " + x)
				if errs != nil {
					fmt.Println("Failed to write to file:", errs)
					return
				}
			}
		}
		if index != len(data.Structures)-1 {
			_, errs = file.WriteString("\n		}, ")
			if errs != nil {
				fmt.Println("Failed to write to file:", errs)
				return
			}
		} else {
			_, errs = file.WriteString("\n		}")
			if errs != nil {
				fmt.Println("Failed to write to file:", errs)
				return
			}
		}
	}

	_, errs = file.WriteString("\n	]\n}")
	if errs != nil {
		fmt.Println("Failed to write to file:", errs)
		return
	}
}
