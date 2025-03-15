package main

import (
	//"github.com/gdamore/tcell/v2"
	//"github.com/rivo/tview"
	//"github.com/go-git/go-git/v5"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// Work on menu
// Make an app that gets the user requested golang github package form github
// and downloads and imports it automatically

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Create the connection handler that upgrades to websocket
// for constant real-time updating
func handleConnections(w http.ResponseWriter, r *http.Request) {
	//Upgrade initial GET request to websocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	for {
		messageType, p, err := ws.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Printf("Message was recieved: %s\n", p)

		if err := ws.WriteMessage(messageType, p); err != nil {
			log.Println(p)
			return
		}

	}

}

func setupRoutes() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Server is up!!")
	})

	http.HandleFunc("/ws", handleConnections)
}

func main() {

	setupRoutes()
	log.Fatal(http.ListenAndServe(":8080", nil))

	/*
		menu := tview.NewBox().SetBorder(true).SetBorderColor(tcell.ColorWhite).SetTitle("Menu")
		cryptoMarket := tview.NewBox().SetBorder(true).SetBorderColor(tcell.ColorGreen.TrueColor()).SetTitle("Cyrpto Market")
		crypto := tview.NewBox().SetBorder(true).SetBorderColor(tcell.ColorDarkRed.TrueColor()).SetTitle("Crypto")
		box := tview.NewTextView().SetText("Terminal").SetBorder(true)

		// Menu stuff
		currencyName := tview.NewInputField().SetLabel("Enter Name or Ticker: ")

		form := tview.NewForm().
			AddFormItem(currencyName).
			AddButton("Find Crypto", nil)

		horizontalFlex := tview.NewFlex().
			AddItem(menu, 0, 1, false).
			AddItem(form, 0, 1, false).
			AddItem(cryptoMarket, 0, 1, false).
			AddItem(crypto, 0, 1, false)

		verticalFlex := tview.NewFlex().
			SetDirection(tview.FlexRow).
			AddItem(box, 3, 0, false).
			AddItem(horizontalFlex, 0, 1, true)

		if err := tview.NewApplication().SetRoot(verticalFlex, true).Run(); err != nil {
			panic(err)
		}
	*/

}
