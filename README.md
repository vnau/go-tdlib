# go-tdlib
Golang Telegram [TDLib](https://core.telegram.org/tdlib) JSON bindings

**ATTENTION!** This is a fork of rezam90's [fork](https://github.com/rezam90/go-tdlib) of [original package](https://github.com/Arman92/go-tdlib) by Arman92. Thanks to these guys for their work.

[![GoDoc](https://godoc.org/github.com/dvz1/go-tdlib?status.svg)](https://godoc.org/github.com/dvz1/go-tdlib)

## Introduction
Telegram TDlib is a complete library for creating telegram clients, it laso has a simple tdjson ready-to-use library to ease
the integration with different programming languages and platforms.

**go-tdlib** is a complete tdlib-tdjson binding package to help you create your own Telegram clients.

**NOTE:** basic tdjson-golang binding is inspired from this package: [go-tdjson](https://github.com/L11R/go-tdjson)

All the classes and functions declared in [TDLib TypeLanguage schema](https://github.com/tdlib/td/blob/master/td/generate/scheme/td_api.tl)
file have been exported using the autogenerate tool [tl-parser](https://github.com/Arman92/go-tl-parser).
So you can use every single type and method in Tdlib.

## Key features:
* Autogenerated golang structs and methods of tdlib .tl schema
* Custom event receivers defined by user (e.g. get only text messages from a specific user)
* Supports all tdjson functions: Send(), Execute(), Receive(), Destroy(), SetFilePath(), SetLogVerbosityLevel()
* Supports all tdlib functions and types

## Installation

First of all you need to clone the TDlib repo and build it:
```bash
git clone git@github.com:tdlib/td.git
cd td
mkdir build
cd build
cmake -DCMAKE_BUILD_TYPE=Release ..
cmake --build . -- -j5
make install

# -j5 refers to number of your cpu cores + 1 for multi-threaded build.
```

If hit any build errors, refer to [TDLib build instructions](https://github.com/tdlib/td#building)
I'm using static linking against tdlib so it won't require to build the whole tdlib source files.

## Docker
You can use prebuilt tdlib with following Docker image:

***Windows:***
``` shell
docker pull mihaildemidoff/tdlib-go
```

## Example
Here is a simple example for authorization and fetching updates:
```golang
package main

import (
	"fmt"

	"github.com/dvz1/go-tdlib"
)

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

	// Main loop
	for update := range client.GetRawUpdatesChannel(1) {
		// Show all updates
		fmt.Println(update.Data)
		fmt.Print("\n\n")
	}

}

```

More examples can be found on [examples folder](https://github.com/dvz1/go-tdlib/tree/master/examples)
