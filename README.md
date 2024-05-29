**Mesher Algo Pseudo Code:**
1. Find bounds of point cloud
    * max and min y coords
    * max and min x coords
    * max and min z coords
2. Use those bound to calculate voxel grid
    * So that if x coords are -1 and +1 and y coords -1 and +1 and z coord are -1 and +1 with a voxel size of .5 we would get a voxel grid of 4 x 4 x 4
3. For every voxel calculate how many points fall into that voxel
    * if the amount is lower than some threshold(tbd) remove the voxel
    * else move on
4. Run a greedy mesher over the remaining voxels
    * greedy mesher should combine all planar faces to a singular face
    * greedy mesher should also smooth stair-stepping by putting a face at an angle to cover the stair-stepping
5. Export new mesh in some format(TBD)

**Normal Alog Pseudo Code**
1. Average all points to find a "center" of the point cloud
2. For each point find the two closest points and calculate the normal using cross product
    * Check the calculated normal and use dot product to make sure that it's point away from previously calculated center
    * if it's not point away flip the vector
3. Figure out how to bake to texture