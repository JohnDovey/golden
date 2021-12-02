package action

import (
	"fmt"
	"github.com/vit1251/golden/pkg/mapper"
	"github.com/vit1251/golden/pkg/site/utils"
	"github.com/vit1251/golden/pkg/site/widgets"
	"log"
	"net/http"
	"strings"
)

type NetmailIndexAction struct {
	Action
}

func NewNetmailIndexAction() *NetmailIndexAction {
	nm := new(NetmailIndexAction)
	return nm
}

func (self NetmailIndexAction) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	mapperManager := self.restoreMapperManager()
	netmailMapper := mapperManager.GetNetmailMapper()

	/* Message headers */
	msgHeaders, err1 := netmailMapper.GetMessageHeaders()
	if err1 != nil {
		response := fmt.Sprintf("Fail on GetMessageHeaders on netmailMapper: err = %+v", err1)
		http.Error(w, response, http.StatusInternalServerError)
		return
	}
	log.Printf("msgHeader = %+v", msgHeaders)

	bw := widgets.NewBaseWidget()

	vBox := widgets.NewVBoxWidget()
	bw.SetWidget(vBox)

	mmw := self.makeMenu()
	vBox.Add(mmw)

	container := widgets.NewDivWidget()
	container.SetClass("container")

	containerVBox := widgets.NewVBoxWidget()

	amw := widgets.NewActionMenuWidget().
		Add(widgets.NewMenuAction().
			SetLink("/netmail/compose").
			SetIcon("icofont-edit").
			SetLabel("Compose"))
	containerVBox.Add(amw)

	container.AddWidget(containerVBox)
	vBox.Add(container)

	indexTable := widgets.NewDivWidget().
		SetClass("echo-msg-index-table").
		SetStyle("width: 100%")

	for _, msg := range msgHeaders {

		/* Render message row */
		msgRow := self.renderRow(msg)
		indexTable.AddWidget(msgRow)

	}
	containerVBox.Add(indexTable)

	/* Render */
	if err := bw.Render(w); err != nil {
		status := fmt.Sprintf("%+v", err)
		http.Error(w, status, http.StatusInternalServerError)
	}

}

func (self *NetmailIndexAction) renderRow(m *mapper.NetmailMsg) widgets.IWidget {

	/* Make message row container */
	rowWidget := widgets.NewDivWidget().
		SetStyle("display: flex").
		SetStyle("direction: column").
		SetStyle("align-items: center")

	var classNames []string
	classNames = append(classNames, "echo-msg-index-item")
	if m.ViewCount == 0 {
		classNames = append(classNames, "echo-msg-index-item-new")
	}
	//if strings.EqualFold(m.From, myName) || strings.EqualFold(m.To, myName) {
	classNames = append(classNames, "echo-msg-index-item-own")
	//}
	rowWidget.SetClass(strings.Join(classNames, " "))

	/* Render user pic */
	nameTitle := utils.TextHelper_makeNameTitle(m.From)
	nameColor := utils.TextHelper_makeColorByName(m.From)
	userpicWidget := widgets.NewDivWidget().
		SetWidth("30px").
		SetHeight("30px").
		SetStyle("flex-shrink: 0").
		//SetStyle("border: 1px solid yellow").
		SetStyle(fmt.Sprintf("background-color: %s", nameColor)).
		SetStyle("border-radius: 50%").
		SetContent(nameTitle)
	rowWidget.AddWidget(userpicWidget)

	/* Render sender name */
	sourceWidget := widgets.NewDivWidget().
		SetWidth("190px").
		SetHeight("38px").
		SetStyle("flex-shrink: 0").
		SetStyle("white-space: nowrap").
		SetStyle("overflow: hidden").
		SetStyle("text-overflow: ellipsis").
		//SetStyle("border: 1px solid green").
		SetContent(m.From)
	rowWidget.AddWidget(sourceWidget)
	// TODO - add `m.To` under m.From ....

	subjectWidget := widgets.NewDivWidget().
		SetStyle("min-width: 350px").
		SetHeight("38px").
		SetStyle("flex-grow: 1").
		SetStyle("white-space: nowrap").
		SetStyle("overflow: hidden").
		SetStyle("text-overflow: ellipsis").
		//SetStyle("border: 1px solid red").
		SetContent(m.Subject)

	rowWidget.AddWidget(subjectWidget)

	msgDate := utils.DateHelper_renderDate(m.DateWritten)
	dateWidget := widgets.NewDivWidget().
		SetHeight("38px").
		SetWidth("160px").
		SetStyle("flex-shrink: 0").
		//SetStyle("border: 1px solid blue").
		SetContent(msgDate)
	rowWidget.AddWidget(dateWidget)

	/* Link container */
	navigateAddress := fmt.Sprintf("/netmail/%s/view", m.Hash)
	navigateItem := widgets.NewLinkWidget().
		SetLink(navigateAddress).
		AddWidget(rowWidget)

	return navigateItem

}
