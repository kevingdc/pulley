package app

type Color int

const (
	ColorRed      Color = 0xf47b67
	ColorYellow   Color = 0xf8a532
	ColorGreen    Color = 0x48b784
	ColorCyan     Color = 0x45ddc0
	ColorDarkGrey Color = 0x36393f
	ColorGrey     Color = 0x99aab5
	ColorBlue     Color = 0x3e70dd
	ColorPurple   Color = 0x9c84ef
	ColorPink     Color = 0xf47fff
	ColorOrange   Color = 0xfc964b
	ColorWhite    Color = 0xffffff
)

type MessageAuthor struct {
	URL       string
	Name      string
	AvatarURL string
}

type MessageContent struct {
	URL       string
	Title     string
	Subtitle  string
	Body      string
	Color     Color
	Thumbnail string
	Author    *MessageAuthor
	Header    string
	Footer    string
}

type Message struct {
	User    *User
	Content *MessageContent
}

func NewSimpleMessage(user *User, content string) *Message {
	return &Message{
		User:    user,
		Content: NewSimpleMessageContent(content),
	}
}

func NewSimpleMessageContent(content string) *MessageContent {
	return &MessageContent{Body: content}
}
