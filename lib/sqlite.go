package ledger

import (
	// "gorm.io/driver/sqlite"
	// "gorm.io/gorm"
	// "gorm.io/gorm/logger"
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

func NewDB() {
	// dialector := sqlite.Open("local.sqlite")
	// db, _ := gorm.Open(dialector, &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})
	// db.AutoMigrate(&Transaction{})
	// os.Remove("foo.db")

	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqlStmt := `
	create table foo (id integer not null primary key, name text);
	delete from foo;
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}
}
