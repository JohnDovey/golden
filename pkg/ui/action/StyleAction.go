package action

import (
	"fmt"
	"github.com/vit1251/golden/pkg/ui/action/style"
	"net/http"
)

type StyleAction struct {
	Action
}

func NewStyleAction() *StyleAction {
	return new(StyleAction)
}

func (self *StyleAction) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	css1 := style.NewCSSStyleSheet()
	rule1 := style.NewCSSRule()

	// Message preview box
	rule1.SetSelectorText(".message-preview")
	rule1.Set("border", "1px solid red")
	rule1.Set("flex-shrink", "0")
	rule1.Set("flex-grow", "0")
	rule1.Set("white-space", "pre-wrap")
	rule1.Set("font-family", "\"Courier New\", monospace")

	css1.InsertRule(rule1)

	content := css1.String()

	fmt.Printf("style = %+v\n", content)

	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(content)))
	w.Header().Set("Content-Type", " text/css; charset=utf-8")
	w.WriteHeader(200)

	w.Write([]byte(content))

}
