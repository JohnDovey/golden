package action

import (
	"fmt"
	"github.com/vit1251/golden/pkg/site/widgets"
	"net/http"
)

type EchoAreaCreateAction struct {
	Action
}

func NewEchoAreaCreateAction() *EchoAreaCreateAction {
	return new(EchoAreaCreateAction)
}

func (self *EchoAreaCreateAction) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	/* Render */
	bw := widgets.NewBaseWidget()

	vBox := widgets.NewVBoxWidget()
	bw.SetWidget(vBox)

	mmw := self.makeMenu()
	vBox.Add(mmw)

	container := widgets.NewDivWidget()
	container.SetClass("container")

	containerVBox := widgets.NewVBoxWidget()

	container.AddWidget(containerVBox)

	vBox.Add(container)

	setupForm := widgets.NewFormWidget().
		SetMethod("POST").
		SetAction("/echo/create/complete")

	/* Add custom param field */
	setupFormBox := widgets.NewVBoxWidget()

	self.createInputField(setupFormBox, "echoname", "Area name", "?")

	setupFormBox.Add(widgets.NewFormButtonWidget().SetTitle("Save").SetType("submit"))
	setupForm.SetWidget(setupFormBox)

	containerVBox.Add(setupForm)

	if err := bw.Render(w); err != nil {
		status := fmt.Sprintf("%+v", err)
		http.Error(w, status, http.StatusInternalServerError)
		return
	}

}

func (self *EchoAreaCreateAction) createInputField(box *widgets.VBoxWidget, name string, summary string, value string) {

	mainDiv := widgets.NewDivWidget().
		SetClass("form-group row")

	mainDivBox := widgets.NewVBoxWidget()
	mainDiv.AddWidget(mainDivBox)

	mainTitle := widgets.NewDivWidget().
		SetClass("col-sm-2 col-form-label").
		SetContent(name)

	mainDivBox.Add(mainTitle)

	mainInput := widgets.NewFormInputWidget().
		SetTitle(summary).
		SetName(name).
		SetValue(value)

	mainDivBox.Add(mainInput)

	box.Add(mainDiv)

}
