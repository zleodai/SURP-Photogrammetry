module github.com/zleodai/SURP-Photogrammetry
require modules/pointCloudDecoder v1.0.0
replace modules/pointCloudDecoder => ./modules/pointCloudDecoder
require modules/greedyMesher v1.0.0
replace modules/greedyMesher => ./modules/greedyMesher
require modules/objExporter v1.0.0
replace modules/objExporter => ./modules/objExporter
require modules/voxelMesher v1.0.0
replace modules/voxelMesher => ./modules/voxelMesher
require modules/pointSorter v1.0.0
replace modules/pointSorter => ./modules/pointSorter

go 1.22.3