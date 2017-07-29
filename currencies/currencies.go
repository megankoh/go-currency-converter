/* currency.go is the package where Currency types are created and conversion occurs
 *
 * Student name: Megan Koh
 * Date: 03/31/17
 ***********************************************************************************/

package currencies

import (
	"fmt"
	"github.com/shopspring/decimal"
	"net/http"
	"io/ioutil"
)

const availableCurrenciesURL string = "https://api.myjson.com/bins/1h6nt5"
// pulling data from http://www.hajanaone.com/free-currency-converter-api.php
const convertCurrencyURL string = "http://www.hajanaone.com/currency-api.php?amount=%s&from=%s&to=%s"


/* Currency struct builds the Currency type
 * Elements: ID, a string: the abbrev. currency name; Description, a string: the full currency name
 */ 
type Currency struct {
	ID          string `json:"id"`
	Description string `json:"description"`
}


/* newCurency() creates a new instance of the type Currency
 * Input: currency, a string
 * Output: type Currency
 */
func NewCurrency(currency string) Currency {
	return Currency{ID: currency}
}


/* MainCurrencies() returns an array with 6 of the most traded-with currencies in the world
 * Return: a Currency array
 */
func MainCurrencies() []Currency {
	usd := NewCurrency("USD")
	cad := NewCurrency("CAD")
	cny := NewCurrency("CNY")
	eur := NewCurrency("EUR")
	gbp := NewCurrency("GBP")
	jpy := NewCurrency("JPY")
	return []Currency{usd, cad, cny, eur, gbp, jpy}
}


/* ConvertCurrency() converts from one curreny to another
 * Input: from, a Currency: original currency; to, a Currency: desired currency; amount, a Decimal: a number
 * Return: a byte array
 */
func ConvertCurrency(from, to Currency, amount decimal.Decimal) ([]byte) {
	endpoint := fmt.Sprintf(convertCurrencyURL, amount, from.ID, to.ID)
	// web scraping the endpoint
	resp, err := http.Get(endpoint)
	// if there is an error, exit
	if err != nil {	
		panic(err)
	}

	// close the stream
	defer resp.Body.Close()

	// reads html as a slice of bytes
	html, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	return html
}
