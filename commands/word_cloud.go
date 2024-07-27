package command

import (
	"bytes"
	"fmt"
	"image/color"
	"image/png"
	"regexp"
	"strings"

	"github.com/BANKA2017/mockbot/dao/model"
	"github.com/BANKA2017/mockbot/share"
	"github.com/psykhi/wordclouds"
)

func WordCloud(bot_info map[string]string, chat_id int64) error {
	// latest 24 hours
	dateOffset := share.Now.Unix() - 60*60*24

	messages := new([]model.GroupMessage)
	share.GormDB.R.Model(&model.GroupMessage{}).Where("chat_id = ? AND date >= ?", chat_id, dateOffset).Find(messages)

	//log.Println(messages)

	textArray := []string{}

	for _, v := range *messages {
		textArray = append(textArray, v.Content)
	}
	words := share.JiebaPtr.Tag(strings.Join(textArray, "\n"))

	wordCounts := make(map[string]int)

	// /n
	for _, word := range words {
		if strings.HasSuffix(word, "/n") {
			word = regexp.MustCompile(`(?m)/(n)$`).ReplaceAllString(word, "")
			if _, ok := wordCounts[word]; !ok {
				wordCounts[word] = 0
			}
			wordCounts[word]++
		}

	}
	fmt.Println(wordCounts)
	var DefaultColors = []color.RGBA{
		{0x2e, 0xc7, 0xc9, 0xff},
		{0xb6, 0xa2, 0xde, 0xff},
		{0x5a, 0xb1, 0xef, 0xff},
		{0xff, 0xb9, 0x80, 0xff},
		{0xd8, 0x7a, 0x80, 0xff},
		{0x8d, 0x98, 0xb3, 0xff},
		{0xe5, 0xcf, 0x0d, 0xff},
		{0x97, 0xb5, 0x52, 0xff},
		{0x95, 0x70, 0x6d, 0xff},
		{0xdc, 0x69, 0xaa, 0xff},
		{0x07, 0xa2, 0xa4, 0xff},
		{0x9a, 0x7f, 0xd1, 0xff},
		{0x58, 0x8d, 0xd5, 0xff},
		{0xf5, 0x99, 0x4e, 0xff},
		{0xc0, 0x50, 0x50, 0xff},
		{0x59, 0x67, 0x8c, 0xff},
		{0xc9, 0xab, 0x00, 0xff},
		{0x7e, 0xb0, 0x0a, 0xff},
		{0x6f, 0x55, 0x53, 0xff},
		{0xc1, 0x40, 0x89, 0xff},
	}
	colors := make([]color.Color, 0)
	for _, c := range DefaultColors {
		colors = append(colors, c)
	}
	w := wordclouds.NewWordcloud(
		wordCounts,
		wordclouds.FontMaxSize(200),
		wordclouds.FontMinSize(10),
		wordclouds.FontFile("/root/mockbot/MiSans-Medium.ttf"),
		wordclouds.Height(1024),
		wordclouds.Width(1024),
		wordclouds.Colors(colors),
	)

	buf := new(bytes.Buffer)
	// Encode takes a writer interface and an image interface
	// We pass it the File and the RGBA
	png.Encode(buf, w.Draw())

	share.SaveTo("/root/mockbot/commands/aaa.png", buf.Bytes())

	//	wordCloudContentTemplate := `☁️ 07-27 热门话题 #WordCloud
	//⏰ 截至今天 22:03
	//🗣️ 本群 20 位朋友共产生 200 条发言
	//🔍 看下有没有你感兴趣的关键词？
	//
	//活跃用户排行榜：
	//
	//    🥇111 贡献: 11
	//    🥈222 贡献: 22
	//    🥉333 贡献: 33
	//    🎖444 贡献: 44
	//    🎖555 贡献: 55
	//
	//🎉感谢这些朋友今天的分享!🎉`

	return nil
}
