package action

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/vit1251/golden/pkg/site/widgets"
	"log"
	"net/http"
)

type EchoMsgIndexAction struct {
	Action
}

func NewEchoMsgIndexAction() *EchoMsgIndexAction {
	ea := new(EchoMsgIndexAction)
	return ea
}

func (self *EchoMsgIndexAction) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	mapperManager := self.restoreMapperManager()
	echoAreaMapper := mapperManager.GetEchoAreaMapper()
	echoMapper := mapperManager.GetEchoMapper()

	/* Parse URL parameters */
	vars := mux.Vars(r)
	echoTag := vars["echoname"]
	log.Printf("echoTag = %v", echoTag)

	newArea, err1 := echoAreaMapper.GetAreaByName(echoTag)
	if err1 != nil {
		response := fmt.Sprintf("Fail on GetAreaByName where echoTag is %s: err = %+v", echoTag, err1)
		http.Error(w, response, http.StatusInternalServerError)
		return
	}
	log.Printf("area = %+v", newArea)

	/* Get message headers */
	msgHeaders, err2 := echoMapper.GetMessageHeaders(echoTag)
	if err2 != nil {
		response := fmt.Sprintf("Fail on GetMessageHeaders where echoTag is %s: err = %+v", echoTag, err2)
		http.Error(w, response, http.StatusInternalServerError)
		return
	}
	log.Printf("msgHeaders = %+v", msgHeaders)
	for _, msg := range msgHeaders {
		log.Printf("msg = %+v", msg)
	}

	// Views

	bw := widgets.NewBaseWidget()

	vBox := widgets.NewVBoxWidget()
	bw.SetWidget(vBox)

	mmw := self.makeMenu()
	vBox.Add(mmw)

	container := widgets.NewDivWidget()
	container.SetClass("container")
	vBox.Add(container)

	containerVBox := widgets.NewVBoxWidget()
	container.SetWidget(containerVBox)

	/* Context actions */
	amw := widgets.NewActionMenuWidget().
		Add(widgets.NewMenuAction().
			SetLink(fmt.Sprintf("/echo/%s/message/compose", newArea.GetName())).
			SetIcon("icofont-edit").
			SetLabel("Compose")).
		Add(widgets.NewMenuAction().
			SetLink(fmt.Sprintf("/echo/%s/update", newArea.GetName())).
			SetIcon("icofont-update").
			SetLabel("Settings"))

	containerVBox.Add(amw)

	indexTable := widgets.NewTableWidget().
		SetClass("table")

	indexTable.AddRow(widgets.NewTableRowWidget().
		AddCell(widgets.NewTableCellWidget().SetWidget(widgets.NewTextWidgetWithText("From"))).
		AddCell(widgets.NewTableCellWidget().SetWidget(widgets.NewTextWidgetWithText("To"))).
		AddCell(widgets.NewTableCellWidget().SetWidget(widgets.NewTextWidgetWithText("Subject"))).
		AddCell(widgets.NewTableCellWidget().SetWidget(widgets.NewTextWidgetWithText("Date"))).
		AddCell(widgets.NewTableCellWidget().SetWidget(widgets.NewTextWidgetWithText("Action"))))

	for _, msg := range msgHeaders {
		log.Printf("msg = %+v", msg)

		actions := widgets.NewVBoxWidget()
		actions.Add(
			widgets.NewLinkWidget().
				SetContent("View").
				SetClass("btn").
				SetLink(fmt.Sprintf("/echo/%s/message/%s/view", msg.Area, msg.Hash)))

		row := widgets.NewTableRowWidget().
			AddCell(widgets.NewTableCellWidget().SetWidget(widgets.NewTextWidgetWithText(msg.From))).
			AddCell(widgets.NewTableCellWidget().SetWidget(widgets.NewTextWidgetWithText(msg.To))).
			AddCell(widgets.NewTableCellWidget().SetWidget(widgets.NewTextWidgetWithText(msg.Subject))).
			AddCell(widgets.NewTableCellWidget().SetClass("echo-msg-index-date").SetWidget(widgets.NewTextWidgetWithText(msg.GetAge()))).
			AddCell(widgets.NewTableCellWidget().SetWidget(actions))
		//
		row.SetClass("")
		if msg.ViewCount == 0 {
			row.SetClass("message-item-new")
		}
		//
		indexTable.AddRow(row)
	}

	containerVBox.Add(indexTable)

	if err := bw.Render(w); err != nil {
		status := fmt.Sprintf("%+v", err)
		http.Error(w, status, http.StatusInternalServerError)
		return
	}

}
