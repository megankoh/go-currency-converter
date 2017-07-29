/* mainDriver.go is the driver for the currency converter program
 *
 * Student name: Megan Koh
 * Date: 04/01/17
 */

package main

import (
	"currencies"
	"github.com/shopspring/decimal"
	"net/http"
	"html/template"
	"regexp"
	"strconv"
	"strings"
	"fmt"
)

/* Page struct builds the Page type
 * Elements: Title, a string: the title of the page; 
 *           Body, a byte array: displays the page/conversion results; 
 *           AllCurrencies, an array: displays conversion into the 6 main currencies;
 */ 
type Page struct {
	Title			string
	Body			[]byte
	AllCurrencies 	[6]string
}

/* loadPage() creates the main page
 * Input: amount, a string
 * Returns: a Page type; or an error
 */
func loadPage(amount string) (*Page, error) {
    return &Page{Title: amount, Body: []byte(amount)}, nil
}

/* viewHandler() loads the file "view.html" onto the server
 * Input: w, a ResponseWriter: loads the page
 *        r, a http.Request: reads any input from the page
 */
func viewHandler(w http.ResponseWriter, r *http.Request) {
	p, err := loadPage("Test")
	if err != nil {
		http.Redirect(w, r, "/view/", http.StatusFound)
		return
	}
	renderTemplate(w, "view", p)
}

/* editHandler() loads the file "edit.html" onto the server
 * Input: w, a ResponseWriter: loads the page
 *        r, a http.Request: reads any input from the page
 */
func editHandler(w http.ResponseWriter, r *http.Request) {
	p, err := loadPage("Test")
	if err != nil {
		p = &Page{Title: "Test"}
	}
	renderTemplate(w, "edit", p)
}

/* getAmount() takes the user's input and returns it as a Decimal type
 * Input: amount, a string
 * Output: a Decimal type
 */
func getAmount(amount string) decimal.Decimal {
	// function to check if there are any letters in the entered amount
	f := func(r rune) bool {
        return r < 'A' || r > 'z'
    }

    // if yes, then the default number will be 1
    if strings.IndexFunc(amount, f) == -1 {
        return decimal.NewFromFloat(1.0)
    } else {
    	monies, err := strconv.ParseFloat(amount, 64)
		if err != nil {
	        panic(err)
	    }
	    return decimal.NewFromFloat(monies)
    }	
}

/* getCurrency() takes the entered currency and returns it as a Currency type
 * Input: currency, a string
 * Output: a Currency type
 */
func getCurrency(currency string) currencies.Currency {
	return currencies.NewCurrency(currency)
}

/* getAll() checks if the checkbox is checked, returning true or false
 * Input: w, a ResponseWriter: loads the page
 *        r, a http.Request: reads any input from the page
 * Output: a boolean
 */
func getAll(w http.ResponseWriter, r *http.Request) bool {
	r.ParseForm()
	all := r.FormValue("allCurrencies")
	if all == "all" {
		return true
	}

	return false
}

/* convertAll() converts the entered amount and original currency into 6 main currencies
 * Input: c1, a Currency;
 		  amount, a Decimal;
 * Output: an array of 6 strings
 */
func convertAll(c1 currencies.Currency, amount decimal.Decimal) [6]string {
	var allCurrencies [6]string
	curList := currencies.MainCurrencies()

	i := 0
	for _, c2 := range curList {
		// converts the byte array into a string
		conv := string(Conversion(c1, c2, amount)[:])
		allCurrencies[i] = fmt.Sprint(conv)
		i++
	}

	return allCurrencies
}

/* Conversion() performs the currency conversion, returning a byte array of results
 * Input: c1, the original Currency;
 *        c2, the desired Currency;
 *        amount, a Decimal type;
 * Output: a byte array
 */
func Conversion(c1, c2 currencies.Currency, amount decimal.Decimal) []byte {
	converted := currencies.ConvertCurrency(c1, c2, amount)
	return converted
}

/* saveHandler() loads converted info onto the file "edit.html"
 * Input: w, a ResponseWriter: loads the page
 *        r, a http.Request: reads any input from the page
 */
func saveHandler(w http.ResponseWriter, r *http.Request) {
	// reads the values the user entered
	amount := r.FormValue("amount")
	c1 := r.FormValue("currency1")
	c2 := r.FormValue("currency2")

	// if the checkbox is checked
	if getAll(w, r) {
		p := &Page{AllCurrencies: convertAll(getCurrency(c1), getAmount(amount))}
		t, _ := template.ParseFiles("edit.html")
    	t.Execute(w, p)
	} else {
		converted := Conversion(getCurrency(c1), getCurrency(c2), getAmount(amount))
		p := &Page{Title: "Converted", Body: converted}
		t, _ := template.ParseFiles("edit.html")
   		t.Execute(w, p)
	}	
}

var templates = template.Must(template.ParseFiles("edit.html", "view.html"))

/* renderTemplate() renders the tempalate read from the provided HTML file
 * Input: w, a ResponseWriter: loads the page
 *        r, a http.Request: reads any input from the page
 */
func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

var validPath = regexp.MustCompile("^/(convert|view)")

/* makeHandler() is a wrapper function for the other handler functions
 * Input: fn, a function
 * Output: http.HandlerFunc type
 */
func makeHandler(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fn(w, r)
	}
}


/* main() executes all the code above when necessary
 */
func main() {
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/convert/", makeHandler(saveHandler))

	http.ListenAndServe(":8080", nil)
}