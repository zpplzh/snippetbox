package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/golangcollege/sessions"
	"github.com/zappel/snippetbox/pkg/models/mysql"
)

type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	session       *sessions.Session             // nama library session ada func Session
	snippets      *mysql.SnippetModel           //ini di file mysql trus ambil struct Snippetmodel diatas sama juga begitu
	templateCache map[string]*template.Template //buat ambil dari library template yang dalaman nya ada Template

}

func main() {
	adr := flag.String("addr", ":1310", "HTTP network address") // address nya bisa di ganti pas go run di command tambahin -addr=":port"
	dsn := flag.String("dsn", "web:pass@/snippetbox?parseTime=true", "MYSQL data source name")
	// ga bole pake web@localhost pas bikin user nya karna kita pake docker jadi aplikasi yang akses ke dalam docker di anggap dari luar (punya ip bukan localhost)
	//kecuali kalau deploy golang nya di dalam docker yang sama dengan docker sql nya baru bisa pake localhost
	secret := flag.String("secret", "s6Ndh+pPbnzHbS*+9Pk8qGWhTzbpa@ge", "Secret key")
	flag.Parse()

	//log
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err1 := openDB(*dsn)
	if err1 != nil {
		errorLog.Fatal(err1)
	}
	defer db.Close()

	templateCache, err := newTemplateCache("./ui/html/")
	if err != nil {
		errorLog.Fatal(err)
	}

	session := sessions.New([]byte(*secret))
	session.Lifetime = 12 * time.Hour
	session.Secure = true

	app := &application{
		errorLog:      errorLog, // ini yang kanan di tampung oleh yang kiri
		infoLog:       infoLog,
		session:       session,
		snippets:      &mysql.SnippetModel{DB: db},
		templateCache: templateCache,
	}

	tlsConfig := &tls.Config{
		PreferServerCipherSuites: true,
		CurvePreferences:         []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	//log http
	srv := &http.Server{ // buat input ListenAndServer jadi ga usah di masukin ke () nya si Listen and serve
		Addr:         *adr,
		ErrorLog:     errorLog,
		Handler:      app.routes(),
		TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	infoLog.Printf("Starting server on port %s", *adr) //ini di print di terminal bukan di web nya

	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem") // ini buat define port nya dan mux ini unutk kasi tau handler nya si mux yang di atas yang udah di define path" nya
	//ada err di atas tampung http.ListenandServer karna untuk tampung err kalau server nya ga bisa start
	errorLog.Fatal(err) // untuk munculin error nya di terminal

}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
