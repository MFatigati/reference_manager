package dbController

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"regexp"

	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

var apiRoute = regexp.MustCompile("api")

var db *sql.DB
var newDBName string = "ref_app"

type List struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

type NewList struct {
	Title string `json:"title"`
}

type Reference struct {
	ID               int    `json:"id"`
	Title            string `json:"title"`
	Author_last      string `json:"author_last"`
	Author_first     string `json:"author_first"`
	Publication_date int    `json:"publication_date"`
	List             int    `json:"list"`
}

type NewReference struct {
	Title            string `json:"title"`
	Author_last      string `json:"author_last"`
	Author_first     string `json:"author_first"`
	Publication_date int    `json:"publication_date"`
	List             int    `json:"list"`
}

func runCommand(cmd *exec.Cmd) {
	stderr, err := cmd.StderrPipe()
	if err != nil {
		panic(err)
	}

	if err := cmd.Start(); err != nil {
		panic(err)
	}

	errout, _ := ioutil.ReadAll(stderr)
	if err := cmd.Wait(); err != nil {
		fmt.Println(errout)
		panic(err)
	}
}

func ConnectDefault() {
	user := os.Getenv("USER")
	password := os.Getenv("PASSWORD")
	defaultDB := os.Getenv("DEFAULT_DB")

	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	connStr := fmt.Sprintf("user=%v dbname=%v password=%v sslmode=disable", user, defaultDB, password)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	//defer db.Close()

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	log.Printf("Connected to default database: %v", defaultDB)
}

func ConnectNew() {
	user := os.Getenv("USER")
	password := os.Getenv("PASSWORD")

	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	connStr := fmt.Sprintf("user=%v dbname=%v password=%v sslmode=disable", user, newDBName, password)

	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	//defer db.Close()

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	log.Printf("Switched to database: %v", newDBName)
}

func ResetDB() {
	if db != nil {
		db.Close()
	}
	ConnectDefault()
	defaultDB := os.Getenv("DEFAULT_DB")
	newDBName := "ref_app"

	dropDatabase := exec.Command("psql", "-d", defaultDB, "-c", "drop database if exists ref_app;")
	runCommand(dropDatabase)
	createDatabase := exec.Command("psql", "-d", defaultDB, "-c", "create database ref_app;")
	runCommand(createDatabase)
	populateDatabase := exec.Command("psql", "-d", newDBName, "-f", "./lib/pgSchema.sql")
	runCommand(populateDatabase)
	log.Printf("%v database reset", newDBName)
	ConnectNew()
}

func Insert(query string) {

}

func SelectAllLists() ([]List, error) {
	var lists []List
	rows, err := db.Query("select * from lists")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var list List
		err := rows.Scan(&list.ID, &list.Title)
		if err != nil {
			log.Fatal(err)
		}
		lists = append(lists, list)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	return lists, nil
}

func SelectAllReferences(orderBy string, ascending bool) ([]Reference, error) {
	var sortDirection string
	if ascending {
		sortDirection = "asc"
	} else {
		sortDirection = "desc"
	}
	var references []Reference
	query := fmt.Sprintf("select * from texts order by %v %v", orderBy, sortDirection)
	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var reference Reference
		err := rows.Scan(&reference.ID,
			&reference.Title,
			&reference.Author_last,
			&reference.Author_first,
			&reference.Publication_date,
			&reference.List)
		if err != nil {
			log.Fatal(err)
		}
		references = append(references, reference)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	return references, nil
}

func AddReference(ref NewReference) (int, error) {
	var id int
	err := db.QueryRow(`INSERT INTO texts
		(title, author_last, author_first, publication_date, list)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id`, ref.Title, ref.Author_last, ref.Author_first, ref.Publication_date, ref.List).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("addReference: %v", err)
	}
	return id, nil
}

func AddList(list NewList) (int, error) {
	var id int
	err := db.QueryRow(`INSERT INTO lists
		(title)
		VALUES ($1)
		RETURNING id`, list.Title).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("addReference: %v", err)
	}
	return id, nil
}

func DeleteList(list List) (string, error) {
	_, err := db.Exec(`DELETE FROM lists WHERE id = $1`, list.ID)
	if err != nil {
		return "", fmt.Errorf("deleteList: %v", err)
	}
	return list.Title, nil
}

func DeleteRef(ref Reference) (string, error) {
	_, err := db.Exec(`DELETE FROM texts WHERE id = $1`, ref.ID)
	if err != nil {
		return "", fmt.Errorf("deleteRef: %v", err)
	}
	return ref.Title, nil
}

func EditEntry(table string, column string, value string, id int) (int, error) {
	fmt.Println("attempting update")
	sqlStatement := fmt.Sprintf(`UPDATE %v SET %v = '%v' WHERE id = %v`, table, column, value, id)
	fmt.Println(sqlStatement)
	res, err := db.Exec(sqlStatement)
	if err != nil {
		panic(err)
	}
	count, err := res.RowsAffected()
	if err != nil {
		panic(err)
	}
	fmt.Println(count)
	return 1, nil
}

func resetDBapi(w http.ResponseWriter, r *http.Request) {
	ResetDB()
	message := "Successfully reset database."
	w.Write([]byte(message))
}

func resetNoApi(w http.ResponseWriter, r *http.Request) {
	ResetDB()
	http.Redirect(w, r, "/lists", http.StatusSeeOther)
}

func DBHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		if apiRoute.FindStringSubmatch(r.URL.Path) != nil {
			resetDBapi(w, r)
		} else {
			resetNoApi(w, r)
		}
	default:
		fmt.Fprintf(w, "Sorry, only GET method is supported for this route.")
	}
}
