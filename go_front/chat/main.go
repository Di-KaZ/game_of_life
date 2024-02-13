package chat

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type ChatMessage struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

type ChatWidget struct {
	messages          []ChatMessage
	list              *widget.List
	send_message_form *widget.Form
	input             *widget.Entry
	Widget            *fyne.Container
}

type ChatWidgetProps struct {
	OnSubmit func(text string)
}

func Init(props ChatWidgetProps) *ChatWidget {
	chatWidget := &ChatWidget{}

	chatWidget.list = widget.NewList(
		func() int {
			return len(chatWidget.messages)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("Test")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			msg := chatWidget.messages[i]
			o.(*widget.Label).SetText(fmt.Sprintf("%s: %s", msg.Name, msg.Content))
		},
	)

	chatWidget.input = widget.NewEntry()

	chatWidget.send_message_form = &widget.Form{
		Items: []*widget.FormItem{
			{
				Text:   "New Message",
				Widget: chatWidget.input,
			},
		},
		OnSubmit: func() {
			props.OnSubmit(chatWidget.input.Text)
			chatWidget.Refresh()
		},
	}

	chatWidget.Widget = container.NewBorder(
		widget.NewLabel("Chat"),
		chatWidget.send_message_form,
		nil,
		nil,
		chatWidget.list,
	)

	return chatWidget
}

func (w *ChatWidget) Refresh() {
	w.input.Text = ""
	w.input.Refresh()
	w.list.Refresh()
}

func (w *ChatWidget) AddMessage(message ChatMessage) {
	w.messages = append(w.messages, message)
	w.Refresh()
}
