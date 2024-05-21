import json

jsonFilePath = "./example_files/footballPCJSON.json"
 
with open(jsonFilePath) as json_file:
    data = json.load(json_file)
 
    # for key in data.keys():
    #     print("%s: %s" % (key, type(data[key])))
    
    # for viewKey in data["views"][0].keys():
    #     print("%s: %s" % (viewKey, type(data["views"][0][viewKey])))
        
    # for intrinsicsKey in data["intrinsics"][0].keys():
    #     print("%s: %s" % (intrinsicsKey, type(data["intrinsics"][0][intrinsicsKey])))
    
    # for poseKey in data["poses"][0]["pose"]["transform"].keys():
    #     print("%s: %s" % (poseKey, type(data["poses"][0]["pose"]["transform"][poseKey])))
        
    # for structureKey in data["structure"][0]["observations"][0].keys():
    #     print("%s: %s" % (structureKey, type(data["structure"][0]["observations"][0][structureKey])))