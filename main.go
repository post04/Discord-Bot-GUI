package main

import (
	"fmt"
	"strings"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	"github.com/bwmarrin/discordgo"
	"github.com/fatih/color"
)

var (
	username       string
	cmdlog, msglog = true, true
	cmds, msgs     = []string{""}, []string{}
)
var mainapp = app.New()

var usernameOutput = widget.NewLabel("N/A")

var msgLoggingEntry = widget.NewLabel(strings.Join(msgs, "\n"))
var msgLoggingOutput = widget.NewScrollContainer(widget.NewHBox(msgLoggingEntry))
var cmdLoggingEntry = widget.NewLabel(strings.Join(cmds, "\n"))
var cmdLoggingOutput = widget.NewScrollContainer(widget.NewHBox(cmdLoggingEntry))

var commands = []string{"ping"}

func ready(s *discordgo.Session, event *discordgo.Ready) {
	username := event.User.Username + "#" + event.User.Discriminator
	usernameOutput.SetText(username)

}

func timeformat(t time.Time) string {
	return fmt.Sprintf("[%d-%02d-%02d %02d:%02d]", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute())
}

func contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}

func contains1(arr []string, str string) bool {
	for _, a := range arr {
		if strings.Contains(a, str) {
			return true
		}
	}
	return false
}

func makeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	cmd := strings.ToLower(m.Content)
	cmd = strings.Split(cmd, " ")[0]
	cmd = strings.Replace(cmd, ".", "", 1)
	if contains(commands, cmd) {
		cmds = append(cmds, "[CMD]"+"["+m.ID+"] "+m.Author.Username+"#"+m.Author.Discriminator+": "+strings.ReplaceAll(m.Content, "\n", "\\n"))
		cmdLoggingEntry.SetText(strings.Join(cmds, "\n"))
		if cmd == "ping" {
			before := makeTimestamp()
			msg, err := s.ChannelMessageSend(m.ChannelID, "Pong....")
			if err != nil {
				fmt.Println(err)
			}
			after := makeTimestamp()
			s.ChannelMessageEdit(m.ChannelID, msg.ID, "Ping: "+fmt.Sprint(after-before))

		}
	} else {
		msgs = append(msgs, "[MSG]"+"["+m.ID+"] "+m.Author.Username+"#"+m.Author.Discriminator+": "+strings.ReplaceAll(m.Content, "\n", "\\n"))
		msgLoggingEntry.SetText(strings.Join(msgs, "\n"))

	}
}

func editLog(id string, newcontent string) {

	if contains1(msgs, id) {
		for i := 0; i < len(msgs); i++ {
			if strings.Contains(msgs[i], id) {
				lol := strings.Split(msgs[i], ": ")[0]
				lol1 := strings.Replace(lol, "MSG", "MSG-EDIT", 1)
				var newstr = lol1 + ": " + strings.ReplaceAll(newcontent, "\n", "\\n") + " (edited)"
				msgs[i] = msgs[i] + "\n" + newstr
			}
		}
		msgLoggingEntry.SetText(strings.Join(msgs, "\n"))
	} else if contains1(cmds, id) {
		for i := 0; i < len(cmds); i++ {
			if strings.Contains(cmds[i], id) {
				lol := strings.Split(cmds[i], ": ")[0]
				lol1 := strings.Replace(lol, "CMD", "CMD-EDIT", 1)
				var newstr = lol1 + ": " + strings.ReplaceAll(newcontent, "\n", "\\n") + " (edited)"
				cmds[i] = cmds[i] + "\n" + newstr

			}
		}
		cmdLoggingEntry.SetText(strings.Join(cmds, "\n"))
	}
}

func messageEdit(s *discordgo.Session, m *discordgo.MessageUpdate) {
	editLog(m.ID, m.Content)
}

func main() {
	var token string
	color.Red("Token: ")
	fmt.Scanln(&token)
	bot, err := discordgo.New("Bot " + token)
	if err != nil {
		color.HiRed("Error Logging into provided token!")
		return
	}
	bot.AddHandler(ready)
	bot.AddHandler(messageCreate)
	bot.AddHandler(messageEdit)
	err = bot.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}
	mainwindow := mainapp.NewWindow("Discord Bot Managment!")
	//sendmsglabel := widget.NewLabel("Send Message!")
	// channelIDentry := widget.NewEntry()
	// channelIDentry.SetText("Channel ID")
	// msgContententry := widget.NewEntry()
	// msgContententry.SetText("Message Content!")
	// sendMsgButton := widget.NewButton("", func() {
	// 	_, err := bot.ChannelMessageSend(channelIDentry.Text, msgContententry.Text)
	// 	channelIDentry.SetText("Channel ID")
	// 	msgContententry.SetText("Message Content!")
	// 	if err != nil {

	// 		fmt.Println(err)
	// 	}

	// })
	mainwindow.SetContent(
		fyne.NewContainerWithLayout(
			layout.NewBorderLayout(
				usernameOutput,
				nil,
				nil,
				nil,
			),
			usernameOutput,
			// fyne.NewContainerWithLayout(
			// 	layout.NewGridLayoutWithRows(1),
			// 	channelIDentry,
			// 	msgContententry,
			// 	sendMsgButton,
			// ),
			fyne.NewContainerWithLayout(
				layout.NewAdaptiveGridLayout(2),
				msgLoggingOutput,
				cmdLoggingOutput,
			),
		),
	)
	mainwindow.Resize(fyne.NewSize(980, 620))
	mainwindow.ShowAndRun()

}
