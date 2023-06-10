package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/mdanialr/pwman_backend/pkg/migration"
	"github.com/mdanialr/pwman_backend/pkg/otp"
	"github.com/mdanialr/pwman_backend/pkg/twofa"
	"github.com/mdanialr/pwman_backend/server"
)

var (
	isGenerateSecret          bool
	isMigrate, isDrop, isSeed bool
	generateQR                string
	verify                    string
)

func init() {
	flag.BoolVar(&isMigrate, "migrate", false, "Run migrations")
	flag.BoolVar(&isSeed, "seed", false, "Run available seeders. This can only be used with -migrate")
	flag.BoolVar(&isDrop, "drop", false, "Drop all tables! WARNING! this will delete all data inside the database. This can only be used with -migrate")
	flag.BoolVar(&isGenerateSecret, "gen", false, "Generate secret that can be placed in app config")
	flag.StringVar(&generateQR, "qr", "", "Generate QR code to given readable directory or full path")
	flag.StringVar(&verify, "verify", "", "Verify the given code")
	flag.Parse()
}

func main() {
	if isGenerateSecret {
		sec, err := otp.NewSecret()
		if err != nil {
			log.Fatalln("failed to generate secret:", err)
		}
		fmt.Println("Your secret:", sec)
		return
	}
	if verify != "" {
		if !twofa.Verify(verify) {
			fmt.Println("ERR: INVALID")
			return
		}
		fmt.Println("VERIFIED")
		return
	}
	if generateQR != "" {
		qr := twofa.GenerateQR()
		os.WriteFile(strings.TrimSuffix(generateQR, "/")+"/qr.png", qr, 0660)
		return
	}
	if isMigrate {
		migration.Run(isSeed, isDrop)
		return
	}

	server.HTTP()
}
