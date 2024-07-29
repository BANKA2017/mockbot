package command

import (
	"bytes"
	"fmt"
	"image/color"
	"image/png"
	"sort"
	"strconv"
	"strings"

	"github.com/BANKA2017/mockbot/dao/model"
	"github.com/BANKA2017/mockbot/share"
	"github.com/psykhi/wordclouds"
)

func WordCloud(bot_info share.BotSettingsType, bot_request *share.BotRequest, content string) error {
	chat_id := bot_request.Message.Chat.ID

	// latest 24 hours
	now := share.Now
	dateOffset := now.Unix() - 60*60*24

	messages := new([]model.GroupMessage)
	share.GormDB.R.Model(&model.GroupMessage{}).Where("chat_id = ? AND date >= ?", chat_id, dateOffset).Find(messages)

	//log.Println(messages)

	messageTotal := 0
	userTotal := 0
	userData := make(map[string]int64)
	userNameKV := make(map[string]string)

	textArray := []string{}

	for _, v := range *messages {
		if _, ok := userData[v.UserID]; !ok {
			userData[v.UserID] = 0
			userTotal++
		}
		messageTotal++
		userData[v.UserID]++
		userNameKV[v.UserID] = v.FullName
		textArray = append(textArray, v.Text)
	}
	words := share.JiebaPtr.Tag(strings.Join(textArray, "\n"))

	wordCounts := make(map[string]int)

	// /n
	for _, _word := range words {
		i := strings.LastIndex(_word, "/")
		if i == -1 {
			continue
		}
		word := _word[:i]
		pos := _word[i+1:]

		if strings.HasPrefix(pos, "n") || pos == "名詞" {
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
		wordclouds.FontMaxSize(300),
		wordclouds.FontMinSize(20),
		wordclouds.FontFile("/root/mockbot/MiSans-Medium.ttf"),
		wordclouds.Height(1024),
		wordclouds.Width(1024),
		wordclouds.Colors(colors),
	)

	buf := new(bytes.Buffer)
	// Encode takes a writer interface and an image interface
	// We pass it the File and the RGBA
	err := png.Encode(buf, w.Draw())

	if err != nil {
		return err
	}

	type UserKV struct {
		Name  string
		Count int64
	}
	var userRank []UserKV
	for k, v := range userData {
		userRank = append(userRank, UserKV{
			Name:  userNameKV[k],
			Count: v,
		})
	}
	sort.Slice(userRank, func(i, j int) bool {
		return userRank[i].Count > userRank[j].Count
	})

	rankList := ""
	for i, user := range userRank {
		if i >= 5 {
			break
		}
		switch i {
		case 0:
			rankList += "🥇"
		case 1:
			rankList += "🥈"
		case 2:
			rankList += "🥉"
		default:
			rankList += "🎖"
		}
		rankList += fmt.Sprintf("`%s` 贡献: %d\n", share.FixMarkdownV2(user.Name), user.Count)
	}

	wordCloudContentTemplate := fmt.Sprintf("☁️ %s 热门话题 \\#WordCloud\n⏰ 截至今天 %s\n🗣️ 本群 %d 位朋友共产生 %d 条发言\n🔍 看下有没有你感兴趣的关键词？\n\n活跃用户排行榜：\n\n%s\n🎉感谢这些朋友今天的分享\\!🎉", strings.ReplaceAll(now.Format("01-02"), "-", "\\-"), now.Format("15:04"), userTotal, messageTotal, rankList)

	_, err = share.SendPhoto(bot_info, strconv.Itoa(int(chat_id)), buf.Bytes(), map[string]any{
		"caption":              wordCloudContentTemplate,
		"parse_mode":           "MarkdownV2",
		"disable_notification": "true",
	})
	return err
}
