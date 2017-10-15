package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "strconv"

    "github.com/siddontang/go-mysql/client"
    "github.com/gorilla/mux"
)

func Index(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "Welcome!")
}

func PageIndex(w http.ResponseWriter, r *http.Request) {
    pages := Pages{
        Page{Name: "Markdown numero uno"},
        Page{Name: "Markdown numer dos"},
    }

    if err := json.NewEncoder(w).Encode(pages); err != nil {
        panic(err)
    }
}

func PageShow(w http.ResponseWriter, r *http.Request) {
     w.Header().Set("Access-Control-Allow-Origin", "*")

     vars := mux.Vars(r)
     userID := vars["userID"]
     pageID := vars["pageID"]
     if !Sanitized(userID, pageID) {
          return;
     }

     conn, _ := client.Connect("104.236.141.69:3306", "gu-port", "gu-port", "gu-port")

     conn.Ping()

     m, _ := conn.Execute(`SELECT * FROM pages WHERE pageID=` + pageID + ` AND userID=` + userID)

     name, _ := m.GetString(0,2)
     body, _ := m.GetString(0,4)
     author, _ := m.GetString(0,3)
     page :=  Page{
                    Name: name,
                    Body: body,
                    Author: author,
                    PageID: pageID,
                    UserID: userID,
               }

     if err := json.NewEncoder(w).Encode(page); err != nil {
         panic(err)
     }
}

func PagePagination(w http.ResponseWriter, r *http.Request){
     w.Header().Set("Access-Control-Allow-Origin", "*")

     vars := mux.Vars(r)

     if !Sanitized(vars["count"], vars["offset"]) {
          return;
     }

     count, _ := strconv.Atoi(vars["count"])
     offset, _ := strconv.Atoi(vars["offset"])
     offset *= count


     conn, _ := client.Connect("104.236.141.69:3306", "gu-port", "gu-port", "gu-port")

     conn.Ping()
     m, _ := conn.Execute(`SELECT * FROM pages LIMIT ` + strconv.Itoa(offset) + `,` +  strconv.Itoa(count))


     type Book struct{
          Pages     []Page     `json:"pages"`
          Count     int       `json:"count"`
          Offset    int       `json:"offset"`
          Total     int       `json:"total"`
     }

     book := Book{
               Pages: []Page{},
               Count: count,
               Offset: offset,
               Total: 0,
          }

     for i := 0; i < count; i++ {
          name, _ := m.GetString(i,2)
          author, _ := m.GetString(i,3)
          pageID, _ := m.GetString(i, 0)
          userID, _ := m.GetString(i, 1)
          page := Page{
                    Name: name,
                    Author: author,
                    PageID: pageID,
                    UserID: userID,
               }
          book.Pages = append(book.Pages, page);
	}

     l, _ := conn.Execute(`SELECT COUNT(pageID) FROM pages`)
     total, _ := l.GetString(0,0)
     book.Total, _ = strconv.Atoi(total);

     if err := json.NewEncoder(w).Encode(book); err != nil {
         panic(err)
     }
}

func Sanitized(inputs... string) bool{
     for _, i := range inputs {
         _, e := strconv.Atoi(i)
         if e != nil{
              return false;
         }
     }
     return true;
}
