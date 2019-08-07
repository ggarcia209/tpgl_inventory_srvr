Package view contains components to pass user-input data to controller and data from controller to viewer.
package view

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"
)

// HTMLData formats data for use with HTML template
type HTMLData struct {
	Name  string
	Price string
}

// TmplMap maps html template paths to shortnames
var TmplMap = map[string]string{
	"home":       "view/html/home.html",
	"list":       "view/html/list.html",
	"price":      "view/html/price.html",
	"priceFail":  "view/html/price_fail.html",
	"update":     "view/html/update.html",
	"updateFail": "view/html/update_fail.html",
	"delete":     "view/html/delete.html",
	"deleteFail": "view/html/delete_fail.html",
	"dbFail":     "view/html/db_fail.html",
}

// GetFields returns 'name' and/or 'price' variables from corresponding HTML form fields
func GetFields(r *http.Request) (string, string) {
	r.ParseForm()
	name := strings.Join(r.Form["name"], "")
	price := strings.Join(r.Form["price"], "")
	return name, price
}

// TemplExe executes for the specified template for the given io.Writer and data interface
func TemplExe(tmpl string, w http.ResponseWriter, data interface{}) error {
	t := template.Must(template.ParseFiles(tmpl))
	if err := t.Execute(w, data); err != nil {
		fmt.Printf("template execution failed: %v", err)
		return fmt.Errorf("template execution failed: %v", err)
	}
	return nil
}
