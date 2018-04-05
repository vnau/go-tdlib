package main

import (
	"fmt"
	"math"
	"os"
	"os/signal"
	"syscall"
	"tg-tdlib/tdlib"
	"time"
)

var allChats []tdlib.Chat

func main() {
	tdlib.SetLogVerbosityLevel(1)
	tdlib.SetFilePath("./errors.txt")

	// Create new instance of client
	client := tdlib.NewClient(tdlib.Config{
		APIID:               "187786",
		APIHash:             "e782045df67ba48e441ccb105da8fc85",
		SystemLanguageCode:  "en",
		DeviceModel:         "Server",
		SystemVersion:       "1.0.0",
		ApplicationVersion:  "1.0.0",
		UseMessageDatabase:  true,
		UseFileDatabase:     true,
		UseChatInfoDatabase: true,
		UseTestDataCenter:   false,
		DatabaseDirectory:   "./tdlib-db",
		FileDirectory:       "./tdlib-files",
		IgnoreFileNames:     false,
	})

	// Handle Ctrl+C
	ch := make(chan os.Signal, 2)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-ch
		client.DestroyInstance()
		os.Exit(1)
	}()

	for {
		currentState, _ := client.Authorize()
		if currentState.GetAuthorizationStateEnum() == tdlib.AuthorizationStateWaitPhoneNumberType {
			fmt.Print("Enter phone: ")
			var number string
			fmt.Scanln(&number)
			_, err := client.SendPhoneNumber(number)
			if err != nil {
				fmt.Printf("Error sending phone number: %v", err)
			}
		} else if currentState.GetAuthorizationStateEnum() == tdlib.AuthorizationStateWaitCodeType {
			fmt.Print("Enter code: ")
			var code string
			fmt.Scanln(&code)
			_, err := client.SendAuthCode(code)
			if err != nil {
				fmt.Printf("Error sending auth code : %v", err)
			}
		} else if currentState.GetAuthorizationStateEnum() == tdlib.AuthorizationStateWaitPasswordType {
			fmt.Print("Enter Password: ")
			var password string
			fmt.Scanln(&password)
			_, err := client.SendAuthPassword(password)
			if err != nil {
				fmt.Printf("Error sending auth password: %v", err)
			}
		} else if currentState.GetAuthorizationStateEnum() == tdlib.AuthorizationStateReadyType {
			fmt.Println("Authorization Ready! Let's rock")
			break
		}
	}

	go func() {
		eventFilter := func(msg *tdlib.TdMessage) bool {
			updateMsg := (*msg).(*tdlib.UpdateNewMessage)
			if updateMsg.Message.SenderUserID == 41507975 {
				return true
			}
			return false
		}

		receiver := client.AddEventReceiver(&tdlib.UpdateNewMessage{}, eventFilter, 5)
		for newMsg := range receiver.Chan {
			fmt.Println(newMsg)
			updateMsg := (newMsg).(*tdlib.UpdateNewMessage)
			msgText := updateMsg.Message.Content.(*tdlib.MessageText)
			fmt.Println("MsgText:  ", msgText.Text)
			fmt.Print("\n\n")
		}

	}()

	// Main loop
	go func() {
		for update := range client.RawUpdates {
			// Show all updates
			// fmt.Println(update.Data)
			// fmt.Print("\n\n")
			_ = update
		}
	}()

	// see https://stackoverflow.com/questions/37782348/how-to-use-getchats-in-tdlib
	chats, err := client.GetChats(math.MaxInt64, 0, 100)
	allChats = make([]tdlib.Chat, 0, 1)
	if err != nil {
		fmt.Printf("Error getting chats, err: %v\n", err)
	} else {
		for _, chatID := range chats.ChatIDs {
			chat, err := client.GetChat(chatID)
			if err == nil {
				fmt.Println("Got chat info: ", *chat)
				allChats = append(allChats, *chat)
			}
		}
	}

	for {
		time.Sleep(1 * time.Second)
	}
}
