package whatsapp

import (
	"context"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
	"google.golang.org/protobuf/proto"
)

func Main(phoneNumber string, debug bool) {
	ctx := context.Background()

	dbLog := waLog.Stdout("Database", "INFO", true)

	container, err := sqlstore.New(ctx, "sqlite3", fmt.Sprintf("file:%s.db?_foreign_keys=on", phoneNumber), dbLog)
	if err != nil {
		panic(err)
	}

	deviceStore, err := container.GetFirstDevice(ctx)
	if err != nil {
		panic(err)
	}

	clientLog := waLog.Stdout("Client", "INFO", true)

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
		switch v := evt.(type) {
		case *events.Message:
			var messageBody = v.Message.GetConversation()
			if messageBody != "" {
				switch messageBody {
				case "ping":
					client.SendMessage(context.Background(), v.Info.Chat, &waE2E.Message{
						Conversation: proto.String("pong"),
					})

					client.MarkRead([]types.MessageID{v.Info.ID}, v.Info.Timestamp, v.Info.Chat, v.Info.Sender)

					log.Println("Received a message:", messageBody)
				default:
					fmt.Println("Received a message:", messageBody)
				}
			}
		case *events.PairSuccess:
			log.Println("Device paired")
		case *events.Connected:
			log.Println("Connected")
		}
	})

	select{}
}
