package ui

import (
	"net/http"
	"github.com/gorilla/mux"
	"log"
	"github.com/vit1251/golden/pkg/packet"
	"github.com/vit1251/golden/pkg/msg"
	"github.com/vit1251/golden/pkg/mailer"
	"github.com/satori/go.uuid"
	"hash/crc32"
	"time"
	"fmt"
)

type ReplyCompleteAction struct {
	Action
}

type ReplyMessage struct {
	Subject string
	To string
	From string
	Body string
	AreaName string
}

func ProcessReplyMessage(um ReplyMessage) (error) {

	/* Message UUID */
	u1 := uuid.NewV4()
//	u1, err4 := uuid.NewV4()
//	if err4 != nil {
//		return err4
//	}

	/* Create packet name */
	pktName := fmt.Sprintf("%s.pkt", u1)

	/* Open outbound packet */
	pw, err1 := packet.NewPacketWriter(pktName)
	if err1 != nil {
		return err1
	}
	defer pw.Close()

	/* Write packet header */
	pktHeader := packet.NewPacketHeader()
	pktHeader.OrigAddr.SetAddr("2:5023/24.3752")
	pktHeader.DestAddr.SetAddr("2:5023/24")
	if err2 := pw.WritePacketHeader(pktHeader); err2 != nil {
		return err2
	}

	/* Prepare packet message */
	msgHeader := packet.NewPacketMessageHeader()
	msgHeader.OrigAddr.SetAddr("2:5023/24.3752")
	msgHeader.DestAddr.SetAddr("2:5023/24.0")
	msgHeader.SetAttribute("Direct")
	msgHeader.SetToUserName(um.To)
	msgHeader.SetFromUserName("Vitold Sedyshev")
	msgHeader.SetSubject(um.Subject)
	var now time.Time = time.Now()
	msgHeader.SetTime(&now)
	if err3 := pw.WriteMessageHeader(msgHeader); err3 != nil {
		return err3
	}

	/* Construct message content */
	msgContent := msg.NewMessageContent()
	msgContent.AddLine(um.Body)
	msgContent.AddLine("")
	msgContent.AddLine("--- Golden/LNX 1.2.0 2020-01-05 18:29:20 MSK (master)")
	msgContent.AddLine(" * Origin: Yo Adrian, I Did It! (c) Rocky II (2:5023/24.3752)")
	rawMsg := msgContent.Pack()

	/* Calculate checksumm */
	h := crc32.NewIEEE()
	h.Write(rawMsg)
	hs := h.Sum32()
	log.Printf("crc32 = %+v", hs)

	/* Write message body */
	msgBody := packet.NewMessageBody()
	//
	msgBody.SetArea(um.AreaName)
	//
	msgBody.AddKludge("TZUTC", "0300")
	msgBody.AddKludge("CHRS", "UTF-8 4")
	msgBody.AddKludge("MSGID", fmt.Sprintf("%s %08x", "2:5023/24.3752", hs))
	msgBody.AddKludge("UUID", fmt.Sprintf("%s", u1))
	msgBody.AddKludge("TID", "golden/lnx 1.2.1 2020-01-05 20:41 (master)")
	//
	msgBody.SetRaw(rawMsg)
	//
	if err5 := pw.WriteMessage(msgBody); err5 != nil {
		return err5
	}

	/* Close packet */
	pw.Close()

	/* Setup upload */
	m := mailer.NewMailerCompat()
	m.TransmitFile(pktName)

	return nil
}

func (self *ReplyCompleteAction) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//
	err := r.ParseForm()
	if err != nil {
		panic(err)
	}
	//
	vars := mux.Vars(r)
	echoTag := vars["echoname"]
	log.Printf("echoTag = %v", echoTag)

	//
	areaManager := self.Site.app.GetAreaManager()
	area, err1 := areaManager.GetAreaByName(echoTag)
	if (err1 != nil) {
		panic(err1)
	}
	log.Printf("area = %v", area)
	//
	to := r.Form.Get("to")
	subj := r.Form.Get("subject")
	body := r.Form.Get("body")
	//
	var um ReplyMessage
	um.Subject = subj
	um.To = to
	um.Body = body
	um.AreaName = area.Name
	//
	ProcessReplyMessage(um)
	//
	log.Printf("to = %s subj = %s body = %s", to, subj, body)
}