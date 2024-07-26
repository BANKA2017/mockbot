package command

import (
	"bytes"
	"image/color"
	"image/png"
	"strings"

	"github.com/BANKA2017/mockbot/dao/model"
	"github.com/BANKA2017/mockbot/share"
	"github.com/psykhi/wordclouds"
	"github.com/yanyiwu/gojieba"
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

	words := gojieba.NewJieba().Cut(strings.Join(textArray, "\n"), true)

	wordCounts := make(map[string]int)

	for _, word := range words {
		if _, ok := wordCounts[word]; !ok {
			wordCounts[word] = 0
		}
		wordCounts[word]++
	}

	// remove count <= 5
	for name, count := range wordCounts {
		if count <= 5 {
			delete(wordCounts, name)
		}
	}

	var DefaultColors = []color.RGBA{
		{0x1b, 0x1b, 0x1b, 0xff},
		{0x48, 0x48, 0x4B, 0xff},
		{0x59, 0x3a, 0xee, 0xff},
		{0x65, 0xCD, 0xFA, 0xff},
		{0x70, 0xD6, 0xBF, 0xff},
	}
	colors := make([]color.Color, 0)
	for _, c := range DefaultColors {
		colors = append(colors, c)
	}
	w := wordclouds.NewWordcloud(
		wordCounts,
		wordclouds.FontMaxSize(500),
		wordclouds.FontMinSize(50),
		wordclouds.FontFile("/root/mockbot/MiSans-Medium.ttf"),
		wordclouds.Height(1024),
		wordclouds.Width(1024),
		wordclouds.Colors(colors),
	)

	buf := new(bytes.Buffer)
	// Encode takes a writer interface and an image interface
	// We pass it the File and the RGBA
	png.Encode(buf, w.Draw())

	share.SaveTo("/root/mockbot/aaa.png", buf.Bytes())

	return nil
}
