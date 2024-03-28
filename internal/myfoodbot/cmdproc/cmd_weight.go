package cmdproc

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/devldavydov/myfood/internal/storage"
	"go.uber.org/zap"
	tele "gopkg.in/telebot.v3"
)

func (r *CmdProcessor) processWeight(c tele.Context, cmdParts []string, userID int64) error {
	if len(cmdParts) == 0 {
		r.logger.Error(
			"invalid weight command",
			zap.String("reason", "len parts"),
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
		)
		return c.Send(msgErrInvalidCommand)
	}

	switch cmdParts[0] {
	case "set":
		return r.weightSetCommand(c, cmdParts[1:], userID)
	case "del":
		return r.weightDelCommand(c, cmdParts[1:], userID)
	case "list":
		return r.weightListCommand(c, cmdParts[1:], userID)
	}

	r.logger.Error(
		"invalid weight command",
		zap.String("reason", "unknown command"),
		zap.Strings("command", cmdParts),
		zap.Int64("userid", userID),
	)
	return c.Send(msgErrInvalidCommand)
}

func (r *CmdProcessor) weightSetCommand(c tele.Context, cmdParts []string, userID int64) error {
	if len(cmdParts) != 2 {
		r.logger.Error(
			"invalid weight set command",
			zap.String("reason", "len parts"),
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
		)
		return c.Send(msgErrInvalidCommand)
	}

	// Parse timestamp
	ts, err := parseTimestampAsUnix(cmdParts[0])
	if err != nil {
		r.logger.Error(
			"invalid weight set command",
			zap.String("reason", "ts format"),
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
			zap.Error(err),
		)
		return c.Send(msgErrInvalidCommand)
	}

	// Parse val
	val, err := strconv.ParseFloat(cmdParts[1], 64)
	if err != nil {
		r.logger.Error(
			"invalid weight set command",
			zap.String("reason", "val format"),
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
			zap.Error(err),
		)
		return c.Send(msgErrInvalidCommand)
	}

	// Save in DB
	ctx, cancel := context.WithTimeout(context.Background(), _stgOperationTimeout)
	defer cancel()

	if err := r.stg.SetWeight(ctx, userID, &storage.Weight{Timestamp: ts, Value: val}); err != nil {
		if errors.Is(err, storage.ErrWeightInvalid) {
			return c.Send(msgErrInvalidCommand)
		}

		r.logger.Error(
			"weight set command DB error",
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
			zap.Error(err),
		)

		return c.Send(msgErrInternal)
	}

	return c.Send(msgOK)
}

func (r *CmdProcessor) weightDelCommand(c tele.Context, cmdParts []string, userID int64) error {
	if len(cmdParts) != 1 {
		r.logger.Error(
			"invalid weight del command",
			zap.String("reason", "len parts"),
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
		)
		return c.Send(msgErrInvalidCommand)
	}

	// Parse timestamp
	ts, err := parseTimestampAsUnix(cmdParts[0])
	if err != nil {
		r.logger.Error(
			"invalid weight del command",
			zap.String("reason", "ts format"),
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
		)
		return c.Send(msgErrInvalidCommand)
	}

	// Delete from DB
	ctx, cancel := context.WithTimeout(context.Background(), _stgOperationTimeout)
	defer cancel()

	if err := r.stg.DeleteWeight(ctx, userID, ts); err != nil {
		r.logger.Error(
			"weight del command DB error",
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
			zap.Error(err),
		)

		return c.Send(msgErrInternal)
	}

	return c.Send(msgOK)
}

func (r *CmdProcessor) weightListCommand(c tele.Context, cmdParts []string, userID int64) error {
	if len(cmdParts) != 2 {
		r.logger.Error(
			"invalid weight list command",
			zap.String("reason", "len parts"),
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
		)
		return c.Send(msgErrInvalidCommand)
	}

	// Parse timestamp
	tsFrom, err := parseTimestampAsUnix(cmdParts[0])
	if err != nil {
		r.logger.Error(
			"invalid weight list command",
			zap.String("reason", "ts from format"),
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
		)
		return c.Send(msgErrInvalidCommand)
	}

	tsTo, err := parseTimestampAsUnix(cmdParts[1])
	if err != nil {
		r.logger.Error(
			"invalid weight list command",
			zap.String("reason", "ts to format"),
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
		)
		return c.Send(msgErrInvalidCommand)
	}

	// List from DB
	ctx, cancel := context.WithTimeout(context.Background(), _stgOperationTimeout)
	defer cancel()

	lst, err := r.stg.GetWeightList(ctx, userID, tsFrom, tsTo)
	if err != nil {
		if errors.Is(err, storage.ErrWeightEmptyList) {
			return c.Send(msgErrEmptyList)
		}

		r.logger.Error(
			"weight list command DB error",
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
			zap.Error(err),
		)

		return c.Send(msgErrInternal)
	}

	// Report table
	tsFromStr, tsToStr := formatTimestamp(tsFrom), formatTimestamp(tsTo)

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(`
	<html>
	<link href="%s" rel="stylesheet">
	<body>
	<div class="container">
		<div class="accordion" id="accordionWeight">
			<div class="accordion-item">
				<h2 class="accordion-header">
					<button class="accordion-button" type="button" data-bs-toggle="collapse" data-bs-target="#collapseTbl"
							aria-expanded="true" aria-controls="collapseTbl">
					<b>Таблица веса за %s - %s</b>
					</button>
				</h2>
				<div id="collapseTbl" class="accordion-collapse collapse show" data-bs-parent="#accordionWeight">
					<div class="accordion-body">
						<table class="table table-bordered table-hover">
							<thead class="table-light">
								<tr>
									<th>Дата</th>
									<th>Вес</th>
								</tr>
							</thead>
	`,
		_cssBotstrapURL,
		tsFromStr,
		tsToStr))

	sb.WriteString("<tbody>")
	xlabels := make([]string, 0, len(lst))
	data := make([]float64, 0, len(lst))
	for _, w := range lst {
		sb.WriteString(fmt.Sprintf("<tr><td>%s</td><td>%.1f</td>", formatTimestamp(w.Timestamp), w.Value))
		xlabels = append(xlabels, formatTimestamp(w.Timestamp))
		data = append(data, w.Value)
	}
	sb.WriteString(`
							</tbody>
						</table>
					</div>
				</div>
			</div>
	`)

	// Chart
	chartSnip, err := GetChartSnippet(&ChardData{
		XLabels: xlabels,
		Data:    data,
		Label:   "Вес",
		Type:    "line",
	})
	if err != nil {
		r.logger.Error(
			"weight list command chart error",
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
			zap.Error(err),
		)

		return c.Send(msgErrInternal)
	}

	sb.WriteString(fmt.Sprintf(`
			<div class="accordion-item">
				<h2 class="accordion-header">
					<button class="accordion-button" type="button" data-bs-toggle="collapse" data-bs-target="#collapseChart"
							aria-expanded="false" aria-controls="collapseChart">
					<b>График веса за %s - %s</b>
					</button>
				</h2>
				<div id="collapseChart" class="accordion-collapse collapse" data-bs-parent="#accordionWeight">
					<div class="accordion-body">
						<canvas id="chart"></canvas>
					</div>
				</div>
			</div>
		</div>
	</div>
	</body>
	<script src="%s"></script>
	%s
	</html>
	`,
		tsFromStr,
		tsToStr,
		_jsBootstrapURL,
		chartSnip))

	return c.Send(&tele.Document{
		File:     tele.FromReader(bytes.NewBufferString(sb.String())),
		MIME:     "text/html",
		FileName: fmt.Sprintf("weight_%s_%s.html", tsFromStr, tsToStr),
	})
}
