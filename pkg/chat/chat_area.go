package chat

import (
	"go-chatroom/pkg/enum"
	"strings"
)

func ShowInOneArea(area enum.Area, args ...string) string {

	var showMsg = []string{area.AreaCheck()}

	showMsg = append(showMsg, args...)

	return strings.Join(showMsg, "")
}
