package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"

	"github.com/jahvon/flow/internal/io/styles"
	"github.com/jahvon/flow/internal/io/ui/types"
)

func ContextStr(ws, ns string) string {
	if ws == "" {
		ws = styles.RenderUnknown("unk")
	}
	if ns == "" {
		ns = "*"
	}
	ctxStr := styles.RenderContext(styles.RenderBold("ctx:"), false) +
		styles.RenderContext(fmt.Sprintf("%s/%s", ws, ns), true)
	return ctxStr
}

func NoticeStr(notice string, lvl types.NoticeLevel) string {
	if notice == "" {
		return ""
	}

	padded := lipgloss.Style{}.MarginLeft(2).Render
	switch lvl {
	case types.NoticeLevelInfo:
		return padded(styles.RenderSuccess(notice))
	case types.NoticeLevelWarning:
		return padded(styles.RenderWarning(notice))
	case types.NoticeLevelError:
		return padded(styles.RenderError(notice))
	default:
		return padded(styles.RenderUnknown(notice))
	}
}
