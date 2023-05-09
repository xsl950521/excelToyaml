package import_data

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Game struct {
	TenantID int       `json:"tenant_id"`
	AreaID   int       `json:"area_id"`
	GameID   int       `json:"game_id"`
	Name     string    `json:"name"`
	Sessions []Session `json:"sessions"`
}

type Session struct {
	SessionID    int    `json:"session_id"`
	SessionLevel int    `json:"session_level"`
	SessionName  string `json:"session_name"`
	SessionSort  int    `json:"session_sort"`
	SessionFlag  int    `json:"session_flag"`
	MinScore     int    `json:"min_score"`
	MaxScore     int    `json:"max_score"`
	Cost         int    `json:"cost"`
	CostMode     int    `json:"cost_mode"`
	BaseScore    int    `json:"base_score"`
	BaseOnline   int    `json:"base_online"`
	ChairCnt     int    `json:"chair_cnt"`
	GameRule     string `json:"game_rule"`
}

var (
	sourcePath string = "./conf/game_config_out/source.json"
	outPath    string = "./conf/game_config_out"
	// sourceData []byte
)

func ImportData(filename string) {
	// 读取原有的 JSON 数据
	jsonData, err := ioutil.ReadFile(sourcePath)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 解析 JSON 数据为切片
	var games []Game
	err = json.Unmarshal(jsonData, &games)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 读取要添加的新 JSON 对象
	newGameData, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 解析新 JSON 对象并添加到现有数组末尾
	var newGame Game
	err = json.Unmarshal(newGameData, &newGame)
	if err != nil {
		fmt.Println(err)
		return
	}
	games = append(games, newGame)

	// 将切片转换为 JSON 数据并写入文件
	jsonData, err = json.MarshalIndent(games, "", "  ")
	if err != nil {
		fmt.Println(err)
		return
	}
	err = ioutil.WriteFile(sourcePath, jsonData, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("New game added to game_data.json")
	outpath := fmt.Sprintf("%s/game_config.5.7160.yaml", outPath)
	err = ioutil.WriteFile(outpath, jsonData, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
}
