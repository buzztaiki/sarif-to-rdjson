package testcase

import (
	"database/sql"
	"fmt"
	"net/http"
)

func ReflectedXSS() {
	http.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		username := r.Form.Get("username")
		fmt.Fprintf(w, "%q is an unknown user", username)
	})
	http.ListenAndServe(":80", nil)
}

func SQLInjection(db *sql.DB, req *http.Request) {
	q := fmt.Sprintf("SELECT ITEM,PRICE FROM PRODUCT WHERE ITEM_CATEGORY='%s' ORDER BY PRICE",
		req.URL.Query()["category"])
	db.Query(q)
}
