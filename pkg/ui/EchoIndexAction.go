package ui

import (
	"fmt"
	"github.com/vit1251/golden/pkg/area"
	"github.com/vit1251/golden/pkg/ui/widgets"
	"log"
	"net/http"
)

type EchoIndexAction struct {
	Action
}

func NewEchoIndexAction() *EchoIndexAction {
	aa := new(EchoIndexAction)
	return aa
}

func (self *EchoIndexAction) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	var areaManager *area.AreaManager
	self.Container.Invoke(func(am *area.AreaManager) {
		areaManager = am
	})

	/* Get message area */
	areas, err1 := areaManager.GetAreas()
	if err1 != nil {
		response := fmt.Sprintf("Fail on GetAreas")
		http.Error(w, response, http.StatusInternalServerError)
		return
	}

	//             <a class="table-row" href="/echo/{{ $area.Name }}" style="color: white">
	//                <div class="table-cell">{{ $area.Name }}</div>
	//                <div class="table-cell">{{ $area.Summary }}</div>
	//                <div class="table-cell" style="text-align: right">
	//{{ if $area.NewMessageCount }}
	//                        <span style="font-weight: bold">{{ $area.NewMessageCount }}</span>
	//                        <span> / </span>
	//{{ end }}
	//                        <span class="">{{ $area.MessageCount }}</span>
	//                </div>
	//            </a>

	bw := widgets.NewBaseWidget()

	vBox := widgets.NewVBoxWidget()
	bw.SetWidget(vBox)

	mmw := widgets.NewMainMenuWidget()
	vBox.Add(mmw)

	indexTable := widgets.NewTableWidget().
		SetClass("table")

	indexTable.
		SetClass("echo-index-items").
		AddRow(widgets.NewTableRowWidget().
			SetClass("echo-index-header").
			AddCell(widgets.NewTableCellWidget().SetWidget(widgets.NewTextWidgetWithText("Name"))).
			AddCell(widgets.NewTableCellWidget().SetWidget(widgets.NewTextWidgetWithText("Summary"))).
			AddCell(widgets.NewTableCellWidget().SetWidget(widgets.NewTextWidgetWithText("Count"))).
			AddCell(widgets.NewTableCellWidget().SetWidget(widgets.NewTextWidgetWithText("Action"))))

	for _, area := range areas {
		log.Printf("area = %+v", area)
		row := widgets.NewTableRowWidget()

		if area.NewMessageCount > 0 {
			row.SetClass("echo-index-item-new")
		} else {
			row.SetClass("echo-index-item")
		}

		row.AddCell(widgets.NewTableCellWidget().SetWidget(widgets.NewTextWidgetWithText(area.Name())))
		row.AddCell(widgets.NewTableCellWidget().SetWidget(widgets.NewTextWidgetWithText(area.Summary)))

		if area.NewMessageCount > 0 {
			cell := widgets.NewTableCellWidget().SetWidget(widgets.NewTextWidgetWithText(fmt.Sprintf("%d / %d", area.NewMessageCount, area.MessageCount)))
			cell.SetClass("echo-index-item-count-new")
			row.AddCell(cell)
		} else {
			cell := widgets.NewTableCellWidget().SetWidget(widgets.NewTextWidgetWithText(fmt.Sprintf("%d", area.MessageCount)))
			cell.SetClass("echo-index-item-count")
			row.AddCell(cell)
		}

		row.AddCell(widgets.NewTableCellWidget().SetWidget(widgets.NewLinkWidget().
				SetContent("View").
				SetLink(fmt.Sprintf("/echo/%s", area.Name()))))


		indexTable.AddRow(row)
	}

	vBox.Add(indexTable)

	if err := bw.Render(w); err != nil {
		status := fmt.Sprintf("%+v", err)
		http.Error(w, status, http.StatusInternalServerError)
		return
	}


}
