package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

func main() {
	// コマンドライン引数でファイルパスを取得
	if len(os.Args) < 2 {
		fmt.Println("エラー: ファイルパスを指定してください。")
		return
	}
	filePath := os.Args[1]

	// ファイルを開く
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("エラー: ファイルを開けませんでした (%v)\n", err)
		return
	}
	defer file.Close()

	// ファイル全体を1つの文字列に連結し、改行や無駄なスペースを削除
	var contentBuilder strings.Builder
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text()) // 各行の両端のスペースを削除
		contentBuilder.WriteString(line)          // 各行を結合
	}
	if err := scanner.Err(); err != nil {
		fmt.Printf("エラー: ファイルの読み取り中に問題が発生しました (%v)\n", err)
		return
	}

	// 全体の内容を取得
	content := contentBuilder.String()

	// 正規表現で会議時間とURLを抽出
	startTimeRegex := regexp.MustCompile(`DTSTART;TZID=[^:]+:(\d{4})(\d{2})(\d{2})T(\d{2})(\d{2})`)
	endTimeRegex := regexp.MustCompile(`DTEND;TZID=[^:]+:(\d{4})(\d{2})(\d{2})T(\d{2})(\d{2})`)
	urlRegex := regexp.MustCompile(`https://teams\.microsoft\.com/[^\s<>]+`)

	// 開始時間を抽出
	meetingStartTime := ""
	if matches := startTimeRegex.FindStringSubmatch(content); matches != nil {
		meetingStartTime = fmt.Sprintf("%s年%s月%s日 %s:%s",
			matches[1], matches[2], matches[3], matches[4], matches[5])
	}

	// 終了時間を抽出
	meetingEndTime := ""
	if matches := endTimeRegex.FindStringSubmatch(content); matches != nil {
		meetingEndTime = fmt.Sprintf("%s年%s月%s日 %s:%s",
			matches[1], matches[2], matches[3], matches[4], matches[5])
	}

	// 会議URLを抽出
	meetingURL := ""
	if matches := urlRegex.FindString(content); matches != "" {
		meetingURL = matches
	}

	// 結果を表示
	if meetingStartTime != "" {
		fmt.Printf("会議の開始時間: %s\n", meetingStartTime)
	} else {
		fmt.Println("会議の開始時間が見つかりませんでした。")
	}

	if meetingEndTime != "" {
		fmt.Printf("会議の終了時間: %s\n", meetingEndTime)
	} else {
		fmt.Println("会議の終了時間が見つかりませんでした。")
	}

	if meetingURL != "" {
		fmt.Printf("会議URL: %s\n", meetingURL)
	} else {
		fmt.Println("会議URLが見つかりませんでした。")
	}
}
