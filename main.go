package main

import (
	"encoding/json"
	"extoyaml/import_data"
	"fmt"
	"io/ioutil"
	"strconv"

	"github.com/xuri/excelize/v2"
	"gopkg.in/yaml.v3"
)

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

type Game struct {
	TenantID int       `json:"tenant_id"`
	AreaID   int       `json:"area_id"`
	GameID   int       `json:"game_id"`
	Name     string    `json:"name"`
	Sessions []Session `json:"sessions"`
}

var (
	templatePath string = "./conf/Template/roomConfig.xlsx"
)

func Transfer(sheetName string) {
	f, err := excelize.OpenFile(templatePath)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 获取sheet的内容
	rows, err := f.GetRows(sheetName)
	if err != nil {
		fmt.Println(err)
		return
	}

	var game Game
	var sessionList []Session
	for i, row := range rows {
		// 跳过表头
		if i == 0 {
			continue
		}

		// 将字符串转换为整数
		tenantID, _ := strconv.Atoi(row[0])
		areaID, _ := strconv.Atoi(row[1])
		gameID, _ := strconv.Atoi(row[2])
		sessionID, _ := strconv.Atoi(row[4])
		sessionLevel, _ := strconv.Atoi(row[5])
		sessionSort, _ := strconv.Atoi(row[7])
		sessionFlag, _ := strconv.Atoi(row[8])
		minScore, _ := strconv.Atoi(row[9])
		maxScore, _ := strconv.Atoi(row[10])
		cost, _ := strconv.Atoi(row[11])
		costMode, _ := strconv.Atoi(row[12])
		baseScore, _ := strconv.Atoi(row[13])
		baseOnline, _ := strconv.Atoi(row[14])
		chairCnt, _ := strconv.Atoi(row[15])
		gameRule := row[16]

		// 构造Session对象
		session := Session{
			SessionID:    sessionID,
			SessionLevel: sessionLevel,
			SessionName:  row[6],
			SessionSort:  sessionSort,
			SessionFlag:  sessionFlag,
			MinScore:     minScore,
			MaxScore:     maxScore,
			Cost:         cost,
			CostMode:     costMode,
			BaseScore:    baseScore,
			BaseOnline:   baseOnline,
			ChairCnt:     chairCnt,
			GameRule:     gameRule,
		}

		// 将Session添加到列表中
		sessionList = append(sessionList, session)

		// 如果已经处理到最后一行或者下一行的游戏ID与当前游戏ID不同，则构造Game对象并输出JSON格式数据
		if i == len(rows)-1 || gameID != atoi(rows[i+1][2]) {
			game = Game{
				TenantID: tenantID,
				AreaID:   areaID,
				GameID:   gameID,
				Name:     row[3],
				Sessions: sessionList,
			}

			printJson(game, gameID)

			sessionList = nil // 清空列表
		}
	}
}

func main() {
	f, err := excelize.OpenFile(templatePath)
	if err != nil {
		fmt.Println(err)
		return
	}
	sheetList := f.GetSheetList()
	for _, sheetName := range sheetList {
		Transfer(sheetName)
		fmt.Println(sheetName)
	}
}

func atoi(s string) int {
	n, _ := strconv.Atoi(s)
	return n
}

// 输出为yaml文件
func printYaml(game Game, gameID int) {
	// 将Game对象转换为YAML格式数据
	yamlData, err := yaml.Marshal(game)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 输出到文件中
	err = ioutil.WriteFile(fmt.Sprintf("game_%d.yaml", gameID), yamlData, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
}

// 输出为json文件
func printJson(game Game, gameID int) {
	jsonStr, _ := json.MarshalIndent(game, "", "  ")

	// 将 JSON 数据写入文件
	outPath := fmt.Sprintf("./conf/outJson/game_%d.json", gameID)
	err := ioutil.WriteFile(outPath, jsonStr, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
	import_data.ImportData(outPath)
	fmt.Printf("JSON data written to file: game_%d.json\n", gameID)
}
