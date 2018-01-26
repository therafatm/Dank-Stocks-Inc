package test

import (
	"fmt"
	"github.com/gavv/httpexpect"
	"net/http"
	"os"
	"strconv"
	"test/workload_generator/commands"
	"testing"
)

var (
	url      string
	username string
)

const testPrice = 20000
const testSymbol = "TST"

var trans = 0

func getTrans() string {
	trans += 1
	return strconv.Itoa(trans)
}

func init() {
	host := os.Getenv("TRANS_HOST")
	port := os.Getenv("TRANS_PORT")
	url = fmt.Sprintf("http://%s:%s", host, port)
	username = "testuser"
}

func initTest(t *testing.T) (e *httpexpect.Expect) {
	e = httpexpect.New(t, url)
	e.GET("/api/clearUsers").
		Expect().
		Status(http.StatusOK)
	return e
}

func checkAvailableBalance(e *httpexpect.Expect, username string, expected int) (obj *httpexpect.Object) {
	cmd := commands.Command{Name: "BALANCE", Username: username, Tnum: getTrans()}
	endpoint := commands.FormatCommandEndpoint(cmd)
	obj = e.GET(endpoint).
		Expect().
		Status(http.StatusOK).
		JSON().Object()

	obj.Keys().ContainsOnly("balance")
	obj.ValueEqual("balance", expected)
	return
}

func checkAvailableShares(e *httpexpect.Expect, username string, symbol string, expected int) (obj *httpexpect.Object) {
	cmd := commands.Command{Name: "SHARES", Username: username, Symbol: symbol, Tnum: getTrans()}
	endpoint := commands.FormatCommandEndpoint(cmd)
	obj = e.GET(endpoint).
		Expect().
		Status(http.StatusOK).
		JSON().Object()

	obj.Keys().ContainsOnly("shares")
	obj.ValueEqual("shares", expected)
	return
}

func add(e *httpexpect.Expect, username string, amount int, status int) (obj *httpexpect.Object) {
	cmd := commands.Command{Name: "ADD", Username: username, Amount: amount, Tnum: getTrans()}
	endpoint := commands.FormatCommandEndpoint(cmd)
	obj = e.GET(endpoint).
		Expect().
		Status(status).
		JSON().Object()
	return
}

func buy(e *httpexpect.Expect, username string, symbol string, amount int, status int) (obj *httpexpect.Object) {
	cmd := commands.Command{Name: "BUY", Username: username, Symbol: testSymbol, Amount: amount, Tnum: getTrans()}
	endpoint := commands.FormatCommandEndpoint(cmd)
	obj = e.GET(endpoint).
		Expect().
		Status(status).
		JSON().Object()
	return
}

func sell(e *httpexpect.Expect, username string, symbol string, amount int, status int) (obj *httpexpect.Object) {
	cmd := commands.Command{Name: "SELL", Username: username, Symbol: testSymbol, Amount: amount, Tnum: getTrans()}
	endpoint := commands.FormatCommandEndpoint(cmd)
	obj = e.GET(endpoint).
		Expect().
		Status(status).
		JSON().Object()
	return
}

func commitBuy(e *httpexpect.Expect, username string, status int) (obj *httpexpect.Object) {
	cmd := commands.Command{Name: "COMMIT_BUY", Username: username, Tnum: getTrans()}
	endpoint := commands.FormatCommandEndpoint(cmd)
	obj = e.GET(endpoint).
		Expect().
		Status(status).
		JSON().Object()
	return
}

func cancelBuy(e *httpexpect.Expect, username string, status int) (obj *httpexpect.Object) {
	cmd := commands.Command{Name: "CANCEL_BUY", Username: username, Tnum: getTrans()}
	endpoint := commands.FormatCommandEndpoint(cmd)
	obj = e.GET(endpoint).
		Expect().
		Status(status).
		JSON().Object()
	return
}

func commitSell(e *httpexpect.Expect, username string, status int) (obj *httpexpect.Object) {
	cmd := commands.Command{Name: "COMMIT_SELL", Username: username, Tnum: getTrans()}
	endpoint := commands.FormatCommandEndpoint(cmd)
	obj = e.GET(endpoint).
		Expect().
		Status(status).
		JSON().Object()
	return
}

func cancelSell(e *httpexpect.Expect, username string, status int) (obj *httpexpect.Object) {
	cmd := commands.Command{Name: "CANCEL_SELL", Username: username, Tnum: getTrans()}
	endpoint := commands.FormatCommandEndpoint(cmd)
	obj = e.GET(endpoint).
		Expect().
		Status(status).
		JSON().Object()
	return
}

func setBuyAmount(e *httpexpect.Expect, username string, symbol string, amount int, status int) (obj *httpexpect.Object) {
	cmd := commands.Command{Name: "SET_BUY_AMOUNT", Username: username, Symbol: testSymbol, Amount: amount, Tnum: getTrans()}
	endpoint := commands.FormatCommandEndpoint(cmd)
	obj = e.GET(endpoint).
		Expect().
		Status(status).
		JSON().Object()
	return
}

func setBuyTrigger(e *httpexpect.Expect, username string, symbol string, amount int, status int) (obj *httpexpect.Object) {
	cmd := commands.Command{Name: "SET_BUY_TRIGGER", Username: username, Symbol: testSymbol, Amount: amount, Tnum: getTrans()}
	endpoint := commands.FormatCommandEndpoint(cmd)
	obj = e.GET(endpoint).
		Expect().
		Status(status).
		JSON().Object()
	return
}

func cancelBuyTrigger(e *httpexpect.Expect, username string, symbol string, status int) (obj *httpexpect.Object) {
	cmd := commands.Command{Name: "CANCEL_SET_BUY", Username: username, Symbol: testSymbol, Tnum: getTrans()}
	endpoint := commands.FormatCommandEndpoint(cmd)
	obj = e.GET(endpoint).
		Expect().
		Status(status).
		JSON().Object()
	return
}

func setSellAmount(e *httpexpect.Expect, username string, symbol string, amount int, status int) (obj *httpexpect.Object) {
	cmd := commands.Command{Name: "SET_SELL_AMOUNT", Username: username, Symbol: testSymbol, Amount: amount, Tnum: getTrans()}
	endpoint := commands.FormatCommandEndpoint(cmd)
	obj = e.GET(endpoint).
		Expect().
		Status(status).
		JSON().Object()
	return
}

func cancelSellTrigger(e *httpexpect.Expect, username string, symbol string, status int) (obj *httpexpect.Object) {
	cmd := commands.Command{Name: "CANCEL_SET_SELL", Username: username, Symbol: testSymbol, Tnum: getTrans()}
	endpoint := commands.FormatCommandEndpoint(cmd)
	obj = e.GET(endpoint).
		Expect().
		Status(status).
		JSON().Object()
	return
}

func setSellTrigger(e *httpexpect.Expect, username string, symbol string, amount int, status int) (obj *httpexpect.Object) {
	cmd := commands.Command{Name: "SET_SELL_TRIGGER", Username: username, Symbol: testSymbol, Amount: amount, Tnum: getTrans()}
	endpoint := commands.FormatCommandEndpoint(cmd)
	obj = e.GET(endpoint).
		Expect().
		Status(status).
		JSON().Object()
	return
}

func executeTriggers(e *httpexpect.Expect, username string, status int) (obj *httpexpect.Array) {
	cmd := commands.Command{Name: "EXECUTE_TRIGGERS", Username: username, Tnum: getTrans()}
	endpoint := commands.FormatCommandEndpoint(cmd)
	obj = e.GET(endpoint).
		Expect().
		Status(status).
		JSON().Array()
	return
}

func TestAddUser(t *testing.T) {
	e := initTest(t)

	amount := 8786332

	//new user
	obj := add(e, username, amount, http.StatusOK)
	obj.Keys().ContainsOnly("id", "username", "money")
	obj.ValueEqual("username", username)
	obj.ValueEqual("money", amount)

	checkAvailableBalance(e, username, amount)

	//update
	addAmount := 20023
	obj = add(e, username, addAmount, http.StatusOK)
	obj.Keys().ContainsOnly("id", "username", "money")
	obj.ValueEqual("username", username)
	obj.ValueEqual("money", amount+addAmount)

	checkAvailableBalance(e, username, amount+addAmount)

	//bad amount
	fAmount := 200.23
	endpoint := fmt.Sprintf("/api/add/%s/%f/0", username, fAmount)
	obj = e.GET(endpoint).
		Expect().
		Status(http.StatusInternalServerError).
		JSON().Object()

	obj.Keys().ContainsOnly("error", "message")
}

func TestBuy(t *testing.T) {
	e := initTest(t)

	// initial buy
	sharesToBuy := 10

	amount := testPrice * sharesToBuy
	obj := add(e, username, amount, http.StatusOK)

	// buy make sure remainder is discarded
	actualShares := sharesToBuy - 1
	obj = buy(e, username, testSymbol, amount-1, http.StatusOK)
	obj.Keys().ContainsOnly("id", "username", "symbol", "type", "shares", "amount", "time")
	obj.ValueEqual("username", username)
	obj.ValueEqual("symbol", testSymbol)
	obj.ValueEqual("type", "BUY")
	obj.ValueEqual("shares", actualShares)
	obj.ValueEqual("amount", actualShares*testPrice)

	checkAvailableBalance(e, username, amount)
	// commit buy order
	obj = commitBuy(e, username, http.StatusOK)
	obj.Keys().ContainsOnly("id", "username", "symbol", "shares")
	obj.ValueEqual("username", username)
	obj.ValueEqual("symbol", testSymbol)
	obj.ValueEqual("shares", actualShares)

	checkAvailableBalance(e, username, amount-(testPrice*actualShares))
	checkAvailableShares(e, username, testSymbol, actualShares)

	// not enough money buy
	obj = buy(e, username, testSymbol, testPrice*sharesToBuy, http.StatusInternalServerError)
	obj.Keys().ContainsOnly("error", "message")
}

func TestBuyCommitNotEnoughMoney(t *testing.T) {
	e := initTest(t)

	// initial buy
	sharesToBuy := 10
	amount := testPrice * sharesToBuy
	obj := add(e, username, amount, http.StatusOK)

	obj = buy(e, username, testSymbol, amount, http.StatusOK)
	obj.Keys().ContainsOnly("id", "username", "symbol", "type", "shares", "amount", "time")
	obj.ValueEqual("username", username)
	obj.ValueEqual("symbol", testSymbol)
	obj.ValueEqual("type", "BUY")
	obj.ValueEqual("shares", sharesToBuy)
	obj.ValueEqual("amount", amount)

	obj = buy(e, username, testSymbol, amount, http.StatusOK)
	obj.Keys().ContainsOnly("id", "username", "symbol", "type", "shares", "amount", "time")
	obj.ValueEqual("username", username)
	obj.ValueEqual("symbol", testSymbol)
	obj.ValueEqual("type", "BUY")
	obj.ValueEqual("shares", sharesToBuy)
	obj.ValueEqual("amount", amount)

	// commit bad buy order
	obj = commitBuy(e, username, http.StatusOK)
	obj.Keys().ContainsOnly("id", "username", "symbol", "shares")
	obj.ValueEqual("username", username)
	obj.ValueEqual("symbol", testSymbol)
	obj.ValueEqual("shares", sharesToBuy)

	checkAvailableBalance(e, username, 0)
	checkAvailableShares(e, username, testSymbol, sharesToBuy)

	obj = commitBuy(e, username, http.StatusInternalServerError)
	obj.Keys().ContainsOnly("error", "message")
}

func TestMultiBuy(t *testing.T) {
	e := initTest(t)

	// buy 3 and commit 3 should work like stack 1 2 3 -> 3 2 1

	sharesForMoney := 100
	numShares1 := 1
	numShares2 := 2
	numShares3 := 3

	amount := testPrice * sharesForMoney
	add(e, username, amount, http.StatusOK)
	buy(e, username, testSymbol, numShares1*testPrice, http.StatusOK)
	buy(e, username, testSymbol, numShares2*testPrice, http.StatusOK)
	buy(e, username, testSymbol, numShares3*testPrice, http.StatusOK)

	checkAvailableShares(e, username, testSymbol, 0)
	checkAvailableBalance(e, username, amount)

	obj := commitBuy(e, username, http.StatusOK)
	obj.Keys().ContainsOnly("id", "username", "symbol", "shares")
	obj.ValueEqual("username", username)
	obj.ValueEqual("symbol", testSymbol)
	obj.ValueEqual("shares", numShares3)

	checkAvailableBalance(e, username, amount-((numShares3)*testPrice))
	checkAvailableShares(e, username, testSymbol, numShares3)

	obj = commitBuy(e, username, http.StatusOK)
	obj.Keys().ContainsOnly("id", "username", "symbol", "shares")
	obj.ValueEqual("username", username)
	obj.ValueEqual("symbol", testSymbol)
	obj.ValueEqual("shares", numShares3+numShares2)

	checkAvailableBalance(e, username, amount-((numShares3+numShares2)*testPrice))
	checkAvailableShares(e, username, testSymbol, numShares3+numShares2)

	obj = commitBuy(e, username, http.StatusOK)
	obj.Keys().ContainsOnly("id", "username", "symbol", "shares")
	obj.ValueEqual("username", username)
	obj.ValueEqual("symbol", testSymbol)
	obj.ValueEqual("shares", numShares3+numShares2+numShares1)

	checkAvailableBalance(e, username, amount-((numShares3+numShares2+numShares1)*testPrice))
	checkAvailableShares(e, username, testSymbol, numShares3+numShares2+numShares1)
}

func TestCancelBuy(t *testing.T) {
	e := initTest(t)

	// buy 3 and cancel 3 should work like stack 1 2 3 -> 3 2 1

	sharesForMoney := 100
	numShares1 := 1
	numShares2 := 2
	numShares3 := 3

	amount := testPrice * sharesForMoney
	add(e, username, amount, http.StatusOK)
	buy(e, username, testSymbol, numShares1*testPrice, http.StatusOK)
	buy(e, username, testSymbol, numShares2*testPrice, http.StatusOK)
	buy(e, username, testSymbol, numShares3*testPrice, http.StatusOK)

	obj := cancelBuy(e, username, http.StatusOK)
	obj.Keys().ContainsOnly("id", "username", "symbol", "type", "shares", "amount", "time")
	obj.ValueEqual("username", username)
	obj.ValueEqual("symbol", testSymbol)
	obj.ValueEqual("shares", numShares3)
	obj.ValueEqual("amount", numShares3*testPrice)
	obj.ValueEqual("type", "BUY")

	checkAvailableBalance(e, username, amount)
	checkAvailableShares(e, username, testSymbol, 0)

	obj = cancelBuy(e, username, http.StatusOK)
	obj.Keys().ContainsOnly("id", "username", "symbol", "type", "shares", "amount", "time")
	obj.ValueEqual("username", username)
	obj.ValueEqual("symbol", testSymbol)
	obj.ValueEqual("shares", numShares2)
	obj.ValueEqual("amount", numShares2*testPrice)
	obj.ValueEqual("type", "BUY")

	checkAvailableBalance(e, username, amount)
	checkAvailableShares(e, username, testSymbol, 0)

	obj = cancelBuy(e, username, http.StatusOK)
	obj.Keys().ContainsOnly("id", "username", "symbol", "type", "shares", "amount", "time")
	obj.ValueEqual("username", username)
	obj.ValueEqual("symbol", testSymbol)
	obj.ValueEqual("shares", numShares1)
	obj.ValueEqual("amount", numShares1*testPrice)
	obj.ValueEqual("type", "BUY")

	checkAvailableBalance(e, username, amount)
	checkAvailableShares(e, username, testSymbol, 0)
}

func TestSell(t *testing.T) {
	e := initTest(t)

	sharesToBuy := 10
	sharesToSell := sharesToBuy - 2

	amount := testPrice * sharesToBuy
	add(e, username, amount, http.StatusOK)
	buy(e, username, testSymbol, amount, http.StatusOK)
	commitBuy(e, username, http.StatusOK)
	checkAvailableBalance(e, username, amount-(testPrice*sharesToBuy))
	checkAvailableShares(e, username, testSymbol, sharesToBuy)

	obj := sell(e, username, testSymbol, sharesToSell*testPrice, http.StatusOK)
	obj.Keys().ContainsOnly("id", "username", "symbol", "type", "shares", "amount", "time")
	obj.ValueEqual("username", username)
	obj.ValueEqual("symbol", testSymbol)
	obj.ValueEqual("shares", sharesToSell)
	obj.ValueEqual("amount", sharesToSell*testPrice)
	obj.ValueEqual("type", "SELL")

	checkAvailableShares(e, username, testSymbol, sharesToBuy)
	checkAvailableBalance(e, username, amount-(testPrice*sharesToBuy))

	// commit sell order
	obj = commitSell(e, username, http.StatusOK)
	obj.Keys().ContainsOnly("id", "username", "symbol", "shares")
	obj.ValueEqual("username", username)
	obj.ValueEqual("symbol", testSymbol)
	obj.ValueEqual("shares", sharesToBuy-sharesToSell)

	checkAvailableShares(e, username, testSymbol, sharesToBuy-sharesToSell)
	checkAvailableBalance(e, username, amount-(testPrice*sharesToBuy)+(sharesToSell*testPrice))

	// not enough shares
	obj = sell(e, username, testSymbol, testPrice*sharesToBuy, http.StatusInternalServerError)
	obj.Keys().ContainsOnly("error", "message")

	//sell remaining shares
	obj = sell(e, username, testSymbol, (sharesToBuy-sharesToSell)*testPrice, http.StatusOK)
	obj.Keys().ContainsOnly("id", "username", "symbol", "type", "shares", "amount", "time")
	obj.ValueEqual("username", username)
	obj.ValueEqual("symbol", testSymbol)
	obj.ValueEqual("shares", (sharesToBuy - sharesToSell))
	obj.ValueEqual("amount", (sharesToBuy-sharesToSell)*testPrice)
	obj.ValueEqual("type", "SELL")

	// commit sell order
	obj = commitSell(e, username, http.StatusOK)
	obj.Keys().ContainsOnly("id", "username", "symbol", "shares")
	obj.ValueEqual("username", username)
	obj.ValueEqual("symbol", testSymbol)
	obj.ValueEqual("shares", 0)

	checkAvailableShares(e, username, testSymbol, 0)
	checkAvailableBalance(e, username, amount)
}

func TestMultiSell(t *testing.T) {
	e := initTest(t)
	// sell 3 and commit 3 should work like stack 1 2 3 -> 3 2 1

	sharesForMoney := 100
	numShares1 := 1
	numShares2 := 2
	numShares3 := 3

	amount := testPrice * sharesForMoney
	add(e, username, amount, http.StatusOK)
	buy(e, username, testSymbol, numShares1*testPrice, http.StatusOK)
	buy(e, username, testSymbol, numShares2*testPrice, http.StatusOK)
	buy(e, username, testSymbol, numShares3*testPrice, http.StatusOK)
	commitBuy(e, username, http.StatusOK)
	commitBuy(e, username, http.StatusOK)
	commitBuy(e, username, http.StatusOK)
	checkAvailableShares(e, username, testSymbol, numShares3+numShares2+numShares1)
	checkAvailableBalance(e, username, amount-((numShares3+numShares2+numShares1)*testPrice))
	sell(e, username, testSymbol, numShares1*testPrice, http.StatusOK)
	sell(e, username, testSymbol, numShares2*testPrice, http.StatusOK)
	sell(e, username, testSymbol, numShares3*testPrice, http.StatusOK)

	checkAvailableShares(e, username, testSymbol, (numShares3 + numShares2 + numShares1))
	checkAvailableBalance(e, username, amount-((numShares3+numShares2+numShares1)*testPrice))

	obj := commitSell(e, username, http.StatusOK)
	obj.Keys().ContainsOnly("id", "username", "symbol", "shares")
	obj.ValueEqual("username", username)
	obj.ValueEqual("symbol", testSymbol)
	obj.ValueEqual("shares", numShares2+numShares1)

	obj = commitSell(e, username, http.StatusOK)
	obj.Keys().ContainsOnly("id", "username", "symbol", "shares")
	obj.ValueEqual("username", username)
	obj.ValueEqual("symbol", testSymbol)
	obj.ValueEqual("shares", numShares1)

	obj = commitSell(e, username, http.StatusOK)
	obj.Keys().ContainsOnly("id", "username", "symbol", "shares")
	obj.ValueEqual("username", username)
	obj.ValueEqual("symbol", testSymbol)
	obj.ValueEqual("shares", 0)

	checkAvailableShares(e, username, testSymbol, 0)
	checkAvailableBalance(e, username, amount)
}

func TestCancelSell(t *testing.T) {
	e := initTest(t)

	// buy 3 and cancel 3 should work like stack 1 2 3 -> 3 2 1
	sharesForMoney := 100
	numShares1 := 1
	numShares2 := 2
	numShares3 := 3

	amount := testPrice * sharesForMoney
	add(e, username, amount, http.StatusOK)
	buy(e, username, testSymbol, numShares1*testPrice, http.StatusOK)
	buy(e, username, testSymbol, numShares2*testPrice, http.StatusOK)
	buy(e, username, testSymbol, numShares3*testPrice, http.StatusOK)

	obj := cancelBuy(e, username, http.StatusOK)
	obj.Keys().ContainsOnly("id", "username", "symbol", "type", "shares", "amount", "time")
	obj.ValueEqual("username", username)
	obj.ValueEqual("symbol", testSymbol)
	obj.ValueEqual("shares", numShares3)
	obj.ValueEqual("amount", numShares3*testPrice)
	obj.ValueEqual("type", "BUY")

	checkAvailableBalance(e, username, amount)
	checkAvailableShares(e, username, testSymbol, 0)

	obj = cancelBuy(e, username, http.StatusOK)
	obj.Keys().ContainsOnly("id", "username", "symbol", "type", "shares", "amount", "time")
	obj.ValueEqual("username", username)
	obj.ValueEqual("symbol", testSymbol)
	obj.ValueEqual("shares", numShares2)
	obj.ValueEqual("amount", numShares2*testPrice)
	obj.ValueEqual("type", "BUY")

	checkAvailableBalance(e, username, amount)
	checkAvailableShares(e, username, testSymbol, 0)

	obj = cancelBuy(e, username, http.StatusOK)
	obj.Keys().ContainsOnly("id", "username", "symbol", "type", "shares", "amount", "time")
	obj.ValueEqual("username", username)
	obj.ValueEqual("symbol", testSymbol)
	obj.ValueEqual("shares", numShares1)
	obj.ValueEqual("amount", numShares1*testPrice)
	obj.ValueEqual("type", "BUY")

	checkAvailableBalance(e, username, amount)
	checkAvailableShares(e, username, testSymbol, 0)
}

func TestSetBuyTrigger(t *testing.T) {
	e := initTest(t)
	sharesForMoney := 10

	amount := testPrice * sharesForMoney
	amount1 := amount / 5
	amount2 := amount - amount1

	add(e, username, amount, http.StatusOK)
	obj := setBuyAmount(e, username, testSymbol, amount1, http.StatusOK)
	obj.Keys().ContainsOnly("id", "username", "symbol", "type", "amount", "triggerprice", "executable", "time")
	obj.ValueEqual("username", username)
	obj.ValueEqual("symbol", testSymbol)
	obj.ValueEqual("amount", amount1)
	obj.ValueEqual("type", "BUY")
	obj.ValueEqual("executable", false)

	checkAvailableBalance(e, username, amount-amount1)

	//duplicate buy amount
	setBuyAmount(e, username, testSymbol, amount2, http.StatusInternalServerError)

	obj = cancelBuyTrigger(e, username, testSymbol, http.StatusOK)
	obj.Keys().ContainsOnly("id", "username", "symbol", "type", "amount", "triggerprice", "executable", "time")
	obj.ValueEqual("username", username)
	obj.ValueEqual("symbol", testSymbol)
	obj.ValueEqual("amount", amount1)
	obj.ValueEqual("type", "BUY")
	obj.ValueEqual("executable", false)

	checkAvailableBalance(e, username, amount)

	obj = cancelBuyTrigger(e, username, testSymbol, http.StatusInternalServerError)
	obj.Keys().ContainsOnly("error", "message")

	obj = setBuyTrigger(e, username, testSymbol, testPrice, http.StatusInternalServerError)
	obj.Keys().ContainsOnly("error", "message")

	obj = setBuyAmount(e, username, testSymbol, amount1, http.StatusOK)
	obj.Keys().ContainsOnly("id", "username", "symbol", "type", "amount", "triggerprice", "executable", "time")
	obj.ValueEqual("username", username)
	obj.ValueEqual("symbol", testSymbol)
	obj.ValueEqual("amount", amount1)
	obj.ValueEqual("type", "BUY")
	obj.ValueEqual("executable", false)

	checkAvailableBalance(e, username, amount-amount1)

	obj = setBuyTrigger(e, username, testSymbol, testPrice, http.StatusOK)
	obj.Keys().ContainsOnly("id", "username", "symbol", "type", "amount", "triggerprice", "executable", "time")
	obj.ValueEqual("username", username)
	obj.ValueEqual("symbol", testSymbol)
	obj.ValueEqual("amount", amount1)
	obj.ValueEqual("type", "BUY")
	obj.ValueEqual("triggerprice", testPrice)
	obj.ValueEqual("executable", true)

	checkAvailableBalance(e, username, amount-amount1)

	obj = cancelBuyTrigger(e, username, testSymbol, http.StatusOK)
	obj.Keys().ContainsOnly("id", "username", "symbol", "type", "amount", "triggerprice", "executable", "time")
	obj.ValueEqual("username", username)
	obj.ValueEqual("symbol", testSymbol)
	obj.ValueEqual("amount", amount1)
	obj.ValueEqual("type", "BUY")
	obj.ValueEqual("triggerprice", testPrice)
	obj.ValueEqual("executable", true)

	obj = cancelBuyTrigger(e, username, testSymbol, http.StatusInternalServerError)
	obj.Keys().ContainsOnly("error", "message")

	checkAvailableBalance(e, username, amount)
}

func TestSetSellTrigger(t *testing.T) {
	e := initTest(t)
	sharesForMoney := 10
	numShares1 := sharesForMoney / 5
	numShares2 := sharesForMoney - numShares1

	amount := testPrice * sharesForMoney
	amount1 := numShares1 * testPrice
	amount2 := numShares2 * testPrice

	add(e, username, amount, http.StatusOK)
	buy(e, username, testSymbol, amount, http.StatusOK)
	commitBuy(e, username, http.StatusOK)
	checkAvailableBalance(e, username, amount-(testPrice*sharesForMoney))
	checkAvailableShares(e, username, testSymbol, sharesForMoney)

	obj := setSellAmount(e, username, testSymbol, amount1, http.StatusOK)
	obj.Keys().ContainsOnly("id", "username", "symbol", "type", "amount", "triggerprice", "executable", "time")
	obj.ValueEqual("username", username)
	obj.ValueEqual("symbol", testSymbol)
	obj.ValueEqual("amount", numShares1)
	obj.ValueEqual("type", "SELL")
	obj.ValueEqual("executable", false)

	checkAvailableShares(e, username, testSymbol, sharesForMoney-numShares1)
	checkAvailableBalance(e, username, amount-(testPrice*sharesForMoney))

	//duplicate
	setSellAmount(e, username, testSymbol, amount2, http.StatusInternalServerError)

	obj = cancelSellTrigger(e, username, testSymbol, http.StatusOK)
	obj.Keys().ContainsOnly("id", "username", "symbol", "type", "amount", "triggerprice", "executable", "time")
	obj.ValueEqual("username", username)
	obj.ValueEqual("symbol", testSymbol)
	obj.ValueEqual("amount", numShares1)
	obj.ValueEqual("type", "SELL")
	obj.ValueEqual("executable", false)

	checkAvailableShares(e, username, testSymbol, sharesForMoney)
	checkAvailableBalance(e, username, amount-(testPrice*sharesForMoney))

	obj = cancelSellTrigger(e, username, testSymbol, http.StatusInternalServerError)
	obj.Keys().ContainsOnly("error", "message")

	obj = setSellTrigger(e, username, testSymbol, testPrice, http.StatusInternalServerError)
	obj.Keys().ContainsOnly("error", "message")

	obj = setSellAmount(e, username, testSymbol, amount1, http.StatusOK)
	obj.Keys().ContainsOnly("id", "username", "symbol", "type", "amount", "triggerprice", "executable", "time")
	obj.ValueEqual("username", username)
	obj.ValueEqual("symbol", testSymbol)
	obj.ValueEqual("amount", numShares1)
	obj.ValueEqual("type", "SELL")
	obj.ValueEqual("executable", false)

	checkAvailableShares(e, username, testSymbol, sharesForMoney-numShares1)

	obj = setSellTrigger(e, username, testSymbol, testPrice, http.StatusOK)
	obj.Keys().ContainsOnly("id", "username", "symbol", "type", "amount", "triggerprice", "executable", "time")
	obj.ValueEqual("username", username)
	obj.ValueEqual("symbol", testSymbol)
	obj.ValueEqual("amount", numShares1)
	obj.ValueEqual("type", "SELL")
	obj.ValueEqual("triggerprice", testPrice)
	obj.ValueEqual("executable", true)

	checkAvailableShares(e, username, testSymbol, sharesForMoney-numShares1)

	obj = cancelSellTrigger(e, username, testSymbol, http.StatusOK)
	obj.Keys().ContainsOnly("id", "username", "symbol", "type", "amount", "triggerprice", "executable", "time")
	obj.ValueEqual("username", username)
	obj.ValueEqual("symbol", testSymbol)
	obj.ValueEqual("amount", numShares1)
	obj.ValueEqual("type", "SELL")
	obj.ValueEqual("triggerprice", testPrice)
	obj.ValueEqual("executable", true)

	obj = cancelSellTrigger(e, username, testSymbol, http.StatusInternalServerError)
	obj.Keys().ContainsOnly("error", "message")

	checkAvailableShares(e, username, testSymbol, sharesForMoney)
	checkAvailableBalance(e, username, amount-(testPrice*sharesForMoney))
}

func TestExecuteBuyTrigger(t *testing.T) {
	e := initTest(t)
	sharesForMoney := 10
	amount := testPrice * sharesForMoney

	add(e, username, amount, http.StatusOK)
	setBuyAmount(e, username, testSymbol, amount, http.StatusOK)
	checkAvailableBalance(e, username, 0)
	setBuyTrigger(e, username, testSymbol, testPrice, http.StatusOK)
	checkAvailableBalance(e, username, 0)

	arr := executeTriggers(e, username, http.StatusOK)
	obj := arr.First().Object()
	obj.Keys().ContainsOnly("id", "username", "symbol", "type", "amount", "triggerprice", "executable", "time")
	obj.ValueEqual("username", username)
	obj.ValueEqual("symbol", testSymbol)
	obj.ValueEqual("amount", amount)
	obj.ValueEqual("type", "BUY")
	obj.ValueEqual("triggerprice", testPrice)
	obj.ValueEqual("executable", true)

	checkAvailableBalance(e, username, 0)
	checkAvailableShares(e, username, testSymbol, sharesForMoney)
}

// func TestExecuteSellTrigger(t *testing.T) {
// 	e := initTest(t)
// 	sharesForMoney := 10
// 	amount := testPrice * sharesForMoney

// 	add(e, username, amount, http.StatusOK)
// 	setSellAmount(e, username, testSymbol, amount, http.StatusOK)
// 	checkAvailableBalance(e, username, 0)
// 	setBuyTrigger(e, username, testSymbol, testPrice, http.StatusOK)
// 	checkAvailableBalance(e, username, 0)

// 	arr := executeTriggers(e, username, http.StatusOK)
// 	obj := arr.First().Object()
// 	obj.Keys().ContainsOnly("id", "username", "symbol", "type", "amount", "triggerprice", "executable", "time")
// 	obj.ValueEqual("username", username)
// 	obj.ValueEqual("symbol", testSymbol)
// 	obj.ValueEqual("amount", amount)
// 	obj.ValueEqual("type", "BUY")
// 	obj.ValueEqual("triggerprice", testPrice)
// 	obj.ValueEqual("executable", true)

// 	checkAvailableBalance(e, username, 0)
// 	checkAvailableShares(e, username, testSymbol, sharesForMoney)
// }
