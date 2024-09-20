package cmdproc

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/devldavydov/myfood/internal/common/messages"
)

func (r *CmdProcessor) calcCalCommand(cmdParts []string) []CmdResponse {
	if len(cmdParts) != 4 {
		return NewSingleCmdResponse(messages.MsgErrInvalidCommand)
	}

	gender := cmdParts[0]
	if !(gender == "m" || gender == "f") {
		return NewSingleCmdResponse(messages.MsgErrInvalidCommand)
	}

	weight, err := strconv.ParseFloat(cmdParts[1], 64)
	if err != nil || weight <= 0 {
		return NewSingleCmdResponse(messages.MsgErrInvalidCommand)
	}

	height, err := strconv.ParseFloat(cmdParts[2], 64)
	if err != nil || height <= 0 {
		return NewSingleCmdResponse(messages.MsgErrInvalidCommand)
	}

	age, err := strconv.ParseFloat(cmdParts[3], 64)
	if err != nil || age <= 0 {
		return NewSingleCmdResponse(messages.MsgErrInvalidCommand)
	}

	ubm := 10*weight + 6.25*height - 5*age
	if gender == "m" {
		ubm += 5
	} else {
		ubm -= 161
	}

	var sb strings.Builder
	sb.WriteString("<b>Уровень Базального Метаболизма (УБМ)</b>\n")
	sb.WriteString(fmt.Sprintf("%d ккал\n\n", int64(ubm)))

	sb.WriteString("<b>Усредненные значения по активностям</b>\n\n")
	for _, i := range []struct {
		name string
		k    float64
	}{
		{name: "Сидячая активность", k: 1.2},
		{name: "Легкая активность", k: 1.375},
		{name: "Средняя активность", k: 1.55},
		{name: "Полноценная активность", k: 1.725},
		{name: "Супер активность", k: 1.9},
	} {
		sb.WriteString(fmt.Sprintf("<b>%s</b>\n", i.name))
		norm := int64(ubm * i.k)
		sb.WriteString(fmt.Sprintf("ККал: %d\n", norm))
		sb.WriteString(fmt.Sprintf("Из них активных: %d\n", norm-int64(ubm)))
		sb.WriteString("\n")
	}

	return NewSingleCmdResponse(sb.String(), optsHTML)
}
