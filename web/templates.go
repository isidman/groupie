// templates.go
// Υπεύθυνο για το φόρτωμα και το rendering των HTML templates.
// Όλα τα templates φορτώνονται κατά την εκκίνηση του server.
//
// Σε περίπτωση σφάλματος στο rendering,
// επιστρέφεται σελίδα σφάλματος αντί για crash.

package web

import (
	"html/template"
	"net/http"
)

type Renderer struct {
	tpls *template.Template
}

func NewRenderer(glob string) (*Renderer, error) {
	t, err := template.ParseGlob(glob)
	if err != nil {
		return nil, err
	}
	return &Renderer{tpls: t}, nil
}

func (r *Renderer) Render(w http.ResponseWriter, name string, data any) error {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	return r.tpls.ExecuteTemplate(w, name, data)
}
