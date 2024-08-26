package utils

import (
	"database/sql"
	"fmt"

	//"log"
	_ "github.com/denisenkom/go-mssqldb"
	"golang.org/x/crypto/bcrypt"
)
var (
    server   = "DESKTOP-5SHM15R\\SQLEXPRESS"
    database = "Fullstack"
)
func Connect() (*sql.DB, error) {
    connString := fmt.Sprintf("server=%s;database=%s;integrated security=true;encrypt=disable",
        server, database)

    db, err := sql.Open("sqlserver", connString)
    if err != nil {
        return nil, fmt.Errorf("error creating connection pool: %v", err)
    }

    err = db.Ping()
    if err != nil {
        db.Close()
        return nil, fmt.Errorf("error pinging database: %v", err)
    }

    fmt.Println("Successfully connected to SQL Server")
    return db, nil
}

func InserUser(db *sql.DB, Name string, Password string) error {
	query := "INSERT INTO dbo.[User] (Name, Password) VALUES(@p1, @p2)"
	_, err := db.Exec(query, sql.Named("p1", Name), sql.Named("p2", Password))
	return err
}
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
func CheckUser(db *sql.DB, Name string) (string, error){
	var hashed string
	query := "SELECT Password FROM dbo.[User] WHERE Name = @name "
	err := db.QueryRow(query, sql.Named("name", Name)).Scan(&hashed)
    return hashed, err
}
func GetProdImag(db *sql.DB) ([]string, []string, []string, error) {
	query := "SELECT ImagUrl, Bio, Price FROM ProductImag"
	res, _ := db.Query(query)
	defer res.Close()
    var err error
	var urls = []string{}
    var bios = []string{}
    var prices = []string{}
    for res.Next(){
        var url string
        var bio string
        var price string
        if err = res.Scan(&url, &bio, &price); err != nil {
            return nil, nil, nil, err
        }
        urls = append(urls, url)
        bios = append(bios, bio)
        prices = append(prices, price)
    }
	return urls, bios, prices, err
}
func AddProdUse(db *sql.DB, Name, Url string) (error, error, int) {
    var id int
    query := "SELECT id FROM ProductImag WHERE ImagUrl = @url"
    err1 := db.QueryRow(query, sql.Named("url", Url)).Scan(&id)
    
    query = "INSERT INTO UsProd (Name, ProdUrl, Urlid) VALUES(@p1, @p2, @p3)"
    _, err := db.Exec(query, sql.Named("p1", Name), sql.Named("p2", Url), sql.Named("p3", id))
    return err, err1, id
}
func GetProdUse(db *sql.DB, Name string) ([]string, []string, []string, error) {
    query := "SELECT P.ImagUrl, P.Bio, P.Price FROM UsProd U JOIN ProductImag P ON U.Urlid = P.id WHERE U.Name = @Name"
    res, _:= db.Query(query, sql.Named("Name", Name))

    defer res.Close()
    var err error
    var urls = []string{}
    var bios = []string{}
    var prices = []string{}
    for res.Next(){
        var url string
        var bio string
        var price string
        err = res.Scan(&url, &bio, &price)
        urls = append(urls, url)
        bios = append(bios, bio)
        prices = append(prices, price)
    }    
	return urls, bios, prices, err
}
func ShowSmallObj(db *sql.DB, Id int) (string, string, string, error){
    // var urlid int
    // query := `
    //     SELECT produrl 
    //     FROM UsProd 
    //     WHERE id = (SELECT MAX(id) FROM UsProd WHERE name = @name)
    // `
    
    // err := db.QueryRow(query, sql.Named("name", Name)).Scan(&urlid)
    var (
        url string
        bio string
        price string
    )
    query := "SELECT ImagUrl, Bio, Price FROM ProductImag WHERE id = @urlid"
    err2 := db.QueryRow(query, sql.Named("urlid", Id)).Scan(&url, &bio, &price)
    
    return url, bio, price, err2
    
}