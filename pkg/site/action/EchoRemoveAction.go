package action

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/vit1251/golden/pkg/site/widgets"
	"log"
	"net/http"
)

type EchoRemoveAction struct {
	Action
}

func NewEchoRemoveAction() *EchoRemoveAction {
	ra := new(EchoRemoveAction)
	return ra
}

func (self *EchoRemoveAction) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	mapperManager := self.restoreMapperManager()
	echoAreaMapper := mapperManager.GetEchoAreaMapper()

	//
	vars := mux.Vars(r)
	echoTag := vars["echoname"]
	log.Printf("echoTag = %v", echoTag)

	//
	area, err1 := echoAreaMapper.GetAreaByName(echoTag)
	if err1 != nil {
		panic(err1)
	}
	log.Printf("area = %+v", area)

	/* Render question */
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

	//<h1>Delete message?</h1>
	headerWidget := widgets.NewHeaderWidget().
		SetTitle("Delete area?")
	containerVBox.Add(headerWidget)

	//
	formWidget := widgets.NewFormWidget().
		SetMethod("POST").
		SetAction(fmt.Sprintf("/echo/%s/remove/complete", area.GetName()))
	formVBox := widgets.NewVBoxWidget()
	formWidget.SetWidget(formVBox)
	containerVBox.Add(formWidget)

	qustionWidget := widgets.NewDivWidget().
		SetContent(fmt.Sprintf("A you sure to remove '%s' area?", area.GetName()))
	formVBox.Add(qustionWidget)

	buttonWidget := widgets.NewFormButtonWidget().
		SetTitle("Remove")
	formVBox.Add(buttonWidget)

	if err := bw.Render(w); err != nil {
		status := fmt.Sprintf("%+v", err)
		http.Error(w, status, http.StatusInternalServerError)
		return
	}

}
