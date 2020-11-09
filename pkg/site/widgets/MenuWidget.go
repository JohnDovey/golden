package widgets

import (
	"fmt"
	"io"
)

type MenuWidget struct {
	actions []*MenuAction
}

func (self*MenuWidget) Add(manuAction *MenuAction) {
	self.actions = append(self.actions, manuAction)
}

func (self*MenuWidget) Actions() []*MenuAction {
	return self.actions
}

func NewMenuWidget() *MenuWidget {
	mw := new(MenuWidget)
	return mw
}

func (self *MenuWidget) Render(w io.Writer) error {

	fmt.Fprintf(w, "<div class=\"panel\">\n")

	fmt.Fprintf(w, "\t<div class=\"tab-group\">\n")
	for _, item := range self.actions {
		fmt.Fprintf(w, "\t\t<a class=\"nav-link\" href=\"%s\">\n", item.Link)
		fmt.Fprintf(w, "\t\t\t<div class=\"tab\">\n")
		fmt.Fprintf(w, "\t\t\t\t<span class=\"tab-label\">%s</span>\n", item.Label)
		if item.Metric > 0 {
			fmt.Fprintf(w, "\t\t\t\t<span class=\"badge\" id=\"%s\">%d</span>\n", item.ID, item.Metric)
		} else {
			fmt.Fprintf(w, "\t\t\t\t<span class=\"badge hidden\" id=\"%s\"></span>\n", item.ID)
		}
		fmt.Fprintf(w, "\t\t\t</div>\n")
		fmt.Fprintf(w, "\t\t</a>\n")
	}

	/* Clock */
	fmt.Fprintf(w, "<div class=\"nav-link\" style=\"margin-left: auto\">\n")
	fmt.Fprintf(w, "\t<div class=\"tab\" id=\"clock\"></div>\n")
	fmt.Fprintf(w, "</div>\n")

	fmt.Fprintf(w, "\t</div>\n")
	fmt.Fprintf(w, "</div>\n")

	return nil
}
