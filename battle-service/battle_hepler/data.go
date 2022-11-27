package battle_hepler

import (
	"battle-service/tools"
	"encoding/json"
	"os"
)

var (
	rolesMap map[string]*Role
	logger   = tools.Logger
	filePath = "./roles.json"
)

func setFilePathForTest() {
	filePath = "../roles.json"
}

func ReadData() {
	jsonFile, err := os.ReadFile(filePath)
	if err != nil {
		logger.Panicf("read roles.json error: %s", err)
	}
	var rolesArr []*Role
	err = json.Unmarshal(jsonFile, &rolesArr)
	if err != nil {
		logger.Panicf("unmarshal error: %s", err)
	}
	rolesMap = make(map[string]*Role, len(rolesMap))
	for i, v := range rolesArr {
		rolesMap[v.ID] = rolesArr[i]
	}
}

func GetRole(rid string) Role {
	return *rolesMap[rid]
}
