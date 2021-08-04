package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	_ "github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/shopspring/decimal"
)

var mySigningKey = []byte("johenews")

type Funds struct {
	Id               int             `json:"id"`
	Name             string          `json:"name"`
	Ticker           string          `json:"ticker"`
	Amount           int64           `json:"amount"`
	PricePerItem     decimal.Decimal `json:"priceperitem"`
	PurchasePrice    decimal.Decimal `json:"purchaseprice"`
	PriceCurrent     decimal.Decimal `json:"pricecurrent"`
	PercentChanges   decimal.Decimal `json:"percentchanges"`
	YearlyInvestment decimal.Decimal `json:"yearlyinvestment"`
	ClearMoney       decimal.Decimal `json:"clearmoney"`
	DatePurchase     time.Time       `json:"datepurchase"`
	DateLastUpdate   time.Time       `json:"datelastupdate"`
	Type             string          `json:"type"`
}

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var user = User{
	Username: "1",
	Password: "1",
}

func isAuthorized(endpoint func(http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Connection", "close")
		defer r.Body.Close()

		if r.Header["Token"] != nil {
			token, err := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("There was an error")
				}
				return mySigningKey, nil
			})

			if err != nil {
				w.WriteHeader(http.StatusForbidden)
				w.Header().Add("Content-Type", "application/json")
				return
			}

			if token.Valid {
				endpoint(w, r)
			}

		} else {
			fmt.Fprintf(w, "Not Authorized")
		}
	})
}

func main() {
	fmt.Println("My Simple Server")

	r := mux.NewRouter()

	//////////////////////////////////////////////////
	//////////////////// FUNDS RUB ///////////////////
	//////////////////////////////////////////////////

	//GET
	r.Handle("/funds/rub/shares", isAuthorized(getRUBFundsShares)).Methods("GET")
	r.Handle("/funds/rub/bonds", isAuthorized(getRUBFundsBonds)).Methods("GET")

	r.Handle("/funds/usd/shares", isAuthorized(getUSDFundsShares)).Methods("GET")
	r.Handle("/funds/usd/bonds", isAuthorized(getUSDFundsBonds)).Methods("GET")

	//POST
	r.Handle("/funds/rub", isAuthorized(createFund)).Methods("POST")

	//////////////////////////////////////////////////
	////////////////////// Login /////////////////////
	//////////////////////////////////////////////////

	//POST
	r.HandleFunc("/login", login).Methods("POST")

	log.Fatal(http.ListenAndServe(":8000", r))
}

func getRUBFundsShares(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	var ArrShares = myRUBCurrentFunds("share")
	json.NewEncoder(w).Encode(ArrShares)
}

func getRUBFundsBonds(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	var ArrShares = myRUBCurrentFunds("bond")
	json.NewEncoder(w).Encode(ArrShares)
}

func getUSDFundsShares(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	var ArrShares = myUSDCurrentFunds("share")
	json.NewEncoder(w).Encode(ArrShares)
}

func getUSDFundsBonds(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	var ArrShares = myUSDCurrentFunds("bond")
	json.NewEncoder(w).Encode(ArrShares)
}

func myRUBCurrentFunds(fundType string) []Funds {
	var amountShares []Funds

	db, err := sql.Open("postgres", "postgres://postgres:1234@localhost/fin?sslmode=disable")

	if err != nil {
		log.Fatal(err)
	}

	rows, err := db.Query("SELECT * FROM funds WHERE type = $1 ORDER BY ticker ASC", fundType)

	for rows.Next() {
		bk := Funds{}
		err = rows.Scan(&bk.Id, &bk.Name, &bk.Ticker, &bk.Amount, &bk.PricePerItem, &bk.PurchasePrice, &bk.PriceCurrent, &bk.PercentChanges, &bk.YearlyInvestment, &bk.ClearMoney, &bk.DatePurchase, &bk.DateLastUpdate, &bk.Type)
		if err != nil {
			fmt.Println(err)
		}
		amountShares = append(amountShares, bk)
	}

	defer rows.Close()

	return amountShares
}

func myUSDCurrentFunds(fundType string) []Funds {
	var amountShares []Funds
	db, err := sql.Open("postgres", "postgres://postgres:1234@localhost/fin?sslmode=disable")

	if err != nil {
		log.Fatal(err)
	}

	rows, err := db.Query("SELECT * FROM fundsusd WHERE type = $1 ORDER BY ticker ASC", fundType)

	for rows.Next() {
		bk := Funds{}
		err = rows.Scan(&bk.Id, &bk.Name, &bk.Ticker, &bk.Amount, &bk.PricePerItem, &bk.PurchasePrice, &bk.PriceCurrent, &bk.PercentChanges, &bk.YearlyInvestment, &bk.ClearMoney, &bk.DatePurchase, &bk.DateLastUpdate, &bk.Type)
		if err != nil {
			fmt.Println(err)
		}
		amountShares = append(amountShares, bk)
	}

	defer rows.Close()

	return amountShares
}

func createFund(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	var fund Funds
	json.NewDecoder(r.Body).Decode(&fund)
	addNewFunds(fund)
}

func addNewFunds(data Funds) {

	db, err := sql.Open("postgres", "postgres://postgres:1234@localhost/fin?sslmode=disable")

	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("INSERT INTO funds (name, ticker, amount, priceperitem, purchaseprice, pricecurrent, percentchanges, yearlyinvestment, clearmoney, datePurchase, dateLastUpdate, type) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)",
		data.Name,
		data.Ticker,
		data.Amount,
		data.PricePerItem,
		data.PurchasePrice,
		data.PriceCurrent,
		data.PercentChanges,
		data.YearlyInvestment,
		data.ClearMoney,
		data.DatePurchase,
		time.Now(),
		data.Type)

	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("addNewF crashed")
		return
	}
}

func login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	var u User
	json.NewDecoder(r.Body).Decode(&u)
	checkLogin(u)
}

func checkLogin(u User) string {

	if user.Username != u.Username || user.Password != u.Password {
		fmt.Println("NOT CORRECT")
		err := "error"
		return err
	}

	validToken, err := GenerateJWT()
	//fmt.Println(validToken)

	if err != nil {
		fmt.Println(err)
	}

	return validToken
}

func GenerateJWT() (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["user"] = "Elliot Forbes"
	claims["exp"] = time.Now().Add(time.Hour * 2160).Unix()

	tokenString, err := token.SignedString(mySigningKey)

	if err != nil {
		fmt.Errorf("Something went wrong: %s", err.Error())
	}

	return tokenString, nil
}
