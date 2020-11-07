package action

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/vit1251/golden/pkg/site/widgets"
	"log"
	"net/http"
)

type EchoComposeAction struct {
	Action
}

func NewEchoComposeAction() *EchoComposeAction {
	ca := new(EchoComposeAction)
	return ca
}

func (self *EchoComposeAction) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	mapperManager := self.restoreMapperManager()
	echoAreaMapper := mapperManager.GetEchoAreaMapper()

	/* Parse URL parameters */
	vars := mux.Vars(r)
	echoTag := vars["echoname"]
	log.Printf("echoTag = %v", echoTag)

	/* Search echo area */

	area, err1 := echoAreaMapper.GetAreaByName(echoTag)
	if err1 != nil {
		response := fmt.Sprintf("Fail on GetAreaByName")
		http.Error(w, response, http.StatusInternalServerError)
		return
	}
	log.Printf("area = %+v", area)

	bw := widgets.NewBaseWidget()

	vBox := widgets.NewVBoxWidget()
	bw.SetWidget(vBox)

	mmw := self.makeMenu()
	vBox.Add(mmw)


	container := widgets.NewDivWidget()
	container.SetClass("container")

	containerVBox := widgets.NewVBoxWidget()

	container.SetWidget(containerVBox)

	vBox.Add(container)


	formVBox := widgets.NewVBoxWidget()

	formWidget := widgets.NewFormWidget()
	formWidget.
		SetMethod("POST").
		SetAction(fmt.Sprintf("/echo/%s/message/compose/complete", area.GetName())).
		SetWidget(formVBox)

	formVBox.Add(widgets.NewFormInputWidget().SetTitle("TO").SetName("to"))
	formVBox.Add(widgets.NewFormInputWidget().SetTitle("SUBJ").SetName("subject"))
	formVBox.Add(widgets.NewFormTextWidget().SetClass("echomail-text").SetName("body"))
	formVBox.Add(widgets.NewFormButtonWidget().SetTitle("Compose").SetType("submit"))

	containerVBox.Add(formWidget)

	if err := bw.Render(w); err != nil {
		status := fmt.Sprintf("%+v", err)
		http.Error(w, status, http.StatusInternalServerError)
		return
	}

}
