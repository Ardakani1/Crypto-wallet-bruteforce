package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
	"github.com/tyler-smith/go-bip39"
)

const (
	dbFile      = "bip39_wallets.db"
	totalCombos = 15000000 // تولید ۱۵ میلیون عبارت بازیابی
	wordsCount  = 12
)

func main() {
	// اتصال به دیتابیس و ساخت جدول
	db := setupDatabase()
	defer db.Close()

	// تولید و ذخیره عبارات بازیابی معتبر
	storeMnemonics(db)
}

func setupDatabase() *sql.DB {
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		log.Fatal("خطا در باز کردن دیتابیس:", err)
	}

	query := `CREATE TABLE IF NOT EXISTS mnemonics (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		phrase TEXT UNIQUE
	);`
	_, err = db.Exec(query)
	if err != nil {
		log.Fatal("خطا در ایجاد جدول:", err)
	}

	_, err = db.Exec("CREATE INDEX IF NOT EXISTS idx_phrase ON mnemonics(phrase);")
	if err != nil {
		log.Fatal("خطا در ایجاد ایندکس:", err)
	}

	return db
}

func storeMnemonics(db *sql.DB) {
	stmt, err := db.Prepare("INSERT OR IGNORE INTO mnemonics(phrase) VALUES(?)")
	if err != nil {
		log.Fatal("خطا در آماده‌سازی کوئری:", err)
	}
	defer stmt.Close()

	for i := 0; i < totalCombos; i++ {
		mnemonic := generateMnemonic()
		_, err = stmt.Exec(mnemonic)
		if err != nil {
			log.Println("خطا در ذخیره‌سازی ترکیب:", err)
		}
		if i%100000 == 0 {
			fmt.Printf("%d ترکیب ذخیره شد...\n", i)
		}
	}
}

func generateMnemonic() string {
	entropy, err := bip39.NewEntropy(128) // تولید 128 بیت Entropy استاندارد
	if err != nil {
		log.Fatal("خطا در تولید Entropy:", err)
	}

	mnemonic, err := bip39.NewMnemonic(entropy) // تولید عبارت بازیابی
	if err != nil {
		log.Fatal("خطا در تولید Mnemonic:", err)
	}

	return mnemonic
}
