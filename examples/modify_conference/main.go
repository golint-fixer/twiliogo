package main

import (
	"bufio"
	"fmt"
	"github.com/natebrennand/twiliogo"
	"github.com/natebrennand/twiliogo/common"
	"github.com/natebrennand/twiliogo/conference"
	"github.com/natebrennand/twiliogo/voice"
	"os"
	"time"
)

func makeCall(to string, act twiliogo.Account) string {
	fmt.Println("Add participant?")
	bufio.NewReader(os.Stdin).ReadString('\n')
	resp, err := act.Voice.Call(voice.Post{
		From: "+16162882901",
		To:   to,
		URL:  "http://twimlbin.com/de26e328",
	})

	if err != nil {
		fmt.Println("Error making call: ", err.Error())
	} else {
		fmt.Println("Participant added")
	}
	// fmt.Printf("%#v\n", resp)

	return resp.Sid
}

func muteParticipant(confSid string, callSid string, act twiliogo.Account) {
	var muted bool
	mutedPtr := &muted
	*mutedPtr = true
	_, err := act.Conferences.SetMute(confSid, callSid, conference.ParticipantAttr{
		Muted: mutedPtr,
	})
	if err != nil {
		fmt.Println("Error muting participant: ", err.Error())
	} else {
		fmt.Println("Participant muting")
	}
}

func kickParticipant(confSid string, callSid string, act twiliogo.Account) {
	err := act.Conferences.Kick(confSid, callSid)
	if err != nil {
		fmt.Println("Error kicking participant: ", err.Error())
	} else {
		fmt.Println("Participant kicked")
	}
}

func main() {
	act := twiliogo.NewAccountFromEnv()
	makeCall("+16164601267", act)
	sid2 := makeCall("+16628551523", act)
	fmt.Println("Mute participant?")
	bufio.NewReader(os.Stdin).ReadString('\n')
	resp, err := act.Conferences.List(conference.ListFilter{
		Status:      "in-progress",
		DateCreated: &(common.JSONTime{time.Date(2014, time.July, 11, 3, 45, 01, 0, &time.Location{})}),
	})
	if err != nil {
		fmt.Println("Error getting conferences: ", err.Error())
	}

	confSid := ""

	for _, c := range *resp.Conferences {
		confSid = c.Sid
	}

	muteParticipant(confSid, sid2, act)

	fmt.Println("Kick participant?")
	bufio.NewReader(os.Stdin).ReadString('\n')
	kickParticipant(confSid, sid2, act)

}
