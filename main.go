package main

import (
	//"fmt"
	"fmt"
	"furniture/utils"
	"html/template"
	"net/http"
)
type ProdImg struct {
	Price []string
	Bio []string
	Url []string
}
type User struct {
	Name string
	Password string
	//ProfUrl string
	Authorized bool
	//Urls []string
}
type Home struct {
	Pursl []string
	BioH []string
	PrecH []string
	PuSmall string
	BSmall string
	PrSmall string
	Autho bool
	Btnclicked bool

}
type Prof struct {
	ImgUrl []string
	Price []string
	Bio []string
	UName string 
}
var u User
var p ProdImg
var h Home

func HomePage(w http.ResponseWriter, r *http.Request){
	var err error
	db, _ := utils.Connect()
	defer db.Close()
	p.Url, p.Bio, p.Price,  err = utils.GetProdImag(db)
	if err != nil {
		fmt.Println(err)
	}
	h.Pursl = p.Url
	h.BioH = p.Bio
	h.PrecH = p.Price
	h.Autho = u.Authorized
	fmt.Println(p.Url)
	fmt.Println(h.Pursl)
	t, _ := template.ParseFiles("templates/index.html")
	t.Execute(w, h)
	h.Btnclicked = false
}
func ProfPage(w http.ResponseWriter, r *http.Request){
	var prof Prof
	//var err error
	db,_ := utils.Connect()
	prof.ImgUrl, prof.Bio, prof.Price, _ = utils.GetProdUse(db, u.Name)
	prof.UName = u.Name 
	fmt.Println(prof.ImgUrl)
	t, _ := template.ParseFiles("templates/profile.html")
	t.Execute(w, prof)
}
func InsertUs(w http.ResponseWriter, r *http.Request) {
	Name := r.FormValue("Name")
	Password := r.FormValue("Password")
	HashedPass, _ := utils.HashPassword(Password)
	
	db, err := utils.Connect()
	if err != nil {
		fmt.Println("err")
	}
	defer db.Close()
	
	err = utils.InserUser(db, Name, HashedPass)
	if err != nil {
		fmt.Println(err)
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
	

}
func Verify(w http.ResponseWriter, r *http.Request) {
	Name := r.FormValue("Name")
	Password := r.FormValue("Password")
	db, _ := utils.Connect()
	defer db.Close()
	exists, err := utils.CheckUser(db, Name)
	if err != nil {
		fmt.Println(err)
	}
	if utils.CheckPasswordHash(Password, exists){
		u.Name = Name
		u.Authorized = true
		fmt.Println("YES")
		http.Redirect(w, r, "/profile", http.StatusSeeOther)
		
	}
	
}
func BuyProd(w http.ResponseWriter, r *http.Request){
	urs := r.FormValue("Url")
	nam := u.Name
	db, _ := utils.Connect()
	defer db.Close()
	h.Btnclicked = false
	err, err2, id := utils.AddProdUse(db, nam, urs)
	h.Btnclicked = true
    h.PuSmall, h.BSmall, h.PrSmall, _ = utils.ShowSmallObj(db, id)
        // Здесь вы можете выполнить дополнительные задачи,
        // такие как запись логов, отправка уведомлений и т.д.

	if err != nil {
		fmt.Println(err)
		} else if err2 != nil{
			fmt.Println(err2)
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
	
func HandleReq () {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	//http.Handle("https://drive.usercontent.google.com/",http.StripPrefix("https://drive.usercontent.google.com/", http.FileServer(http.Dir("https://drive.usercontent.google.com/"))))
	http.HandleFunc("/", HomePage)
	http.HandleFunc("/profile", ProfPage)
	http.HandleFunc("/save", InsertUs)
	http.HandleFunc("/verify", Verify)
	http.HandleFunc("/buy", BuyProd)
	http.ListenAndServe(":8080", nil)
}
func main() {
	HandleReq()
}