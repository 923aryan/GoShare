package models

import (
	"sync"

	"fyne.io/fyne/v2/widget"
	"github.com/grandcat/zeroconf"
)

var ServiceEntries []*zeroconf.ServiceEntry
var Mu sync.Mutex
var ConnectionEstablished bool = false
var ConnectionFormed = make(chan bool, 0)
var ChatLog *widget.Entry = nil
var SendPressed = make(chan bool, 0)
var GotMessage = make(chan bool, 0)
var Message string = ""
var MessageReceived string = ""
