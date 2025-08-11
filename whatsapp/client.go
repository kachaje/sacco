package whatsapp

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sacco/forms"

	_ "github.com/mattn/go-sqlite3"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
	"google.golang.org/protobuf/proto"
)

var bot *forms.MembershipChatbot
var chatActivated bool
var chatTarget string

func handleChats(input string, v *events.Message, sendMessage func(message string, v *events.Message)) {
	if !chatActivated || chatTarget == "" {
		if input == "new kaso member" {
			chatActivated = true
			chatTarget = v.Info.Sender.User

			bot = forms.NewMembershipChatBot()

			log.Println("Initialised session for", chatTarget)
		} else {
			log.Println("Discarding", input)
			return
		}
	} else if v.Info.Sender.User != chatTarget {
		return
	} else if input == "abort" {
		chatActivated = false
		chatTarget = ""

		sendMessage("Session aborted", v)
		
		return
	}

	log.Println("New message from", v.Info.Sender.User)

	question := bot.ProcessInput(input)

	if question == "" {
		chatActivated = false
		chatTarget = ""

		payload, _ := json.MarshalIndent(bot.Data, "", "  ")

		sendMessage(string(payload), v)
	} else {
		sendMessage(question, v)
	}
}

func Main(phoneNumber string, debug bool) {
	ctx := context.Background()

	dbLog := waLog.Stdout("Database", "INFO", true)

	if !debug {
		dbLog = waLog.Noop
	}

	container, err := sqlstore.New(ctx, "sqlite3", fmt.Sprintf("file:%s.db?_foreign_keys=on", phoneNumber), dbLog)
	if err != nil {
		panic(err)
	}

	deviceStore, err := container.GetFirstDevice(ctx)
	if err != nil {
		panic(err)
	}

	clientLog := waLog.Stdout("Client", "INFO", true)

	if !debug {
		clientLog = waLog.Noop
	}

	client := whatsmeow.NewClient(deviceStore, clientLog)

	if client.Store.ID == nil {
		err = client.Connect()
		if err != nil {
			panic(err)
		}

		authCode, err := client.PairPhone(ctx, phoneNumber, true, whatsmeow.PairClientChrome, "Chrome (Linux)")
		if err != nil {
			panic(err)
		}

		fmt.Println(authCode)
	} else {
		err = client.Connect()
		if err != nil {
			panic(err)
		}
	}

	client.AddEventHandler(func(evt any) {
		sendMessage := func(message string, v *events.Message) {
			client.SendMessage(context.Background(), v.Info.Chat, &waE2E.Message{
				Conversation: proto.String(message),
			})
		}

		switch v := evt.(type) {
		case *events.Message:
			var messageBody = v.Message.GetConversation()
			if messageBody != "" {
				switch messageBody {
				case "ping":
					sendMessage("pong", v)

					client.MarkRead([]types.MessageID{v.Info.ID}, v.Info.Timestamp, v.Info.Chat, v.Info.Sender)
				default:
					handleChats(messageBody, v, sendMessage)
				}
			}
		case *events.PairSuccess:
			log.Println("Device paired")
		case *events.Connected:
			log.Println("Connected")
		}
	})

	select {}
}
