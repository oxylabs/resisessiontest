package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/schollz/progressbar/v3"

	"github.com/oxylabs/resisessiontest/sessiontester"
)

const DefaultIPCheckTarget = "https://ipinfo.io/ip"

func main() {
	username := flag.String("u", "", "Oxylabs proxy username")
	password := flag.String("p", "", "Oxylabs proxy password")
	cc := flag.String("cc", "fr", "country parameter")
	city := flag.String("city", "nice", "city parameter")
	sessionCount := flag.Int("sessions", 100, "sessions to test concurrently")
	ipTarget := flag.String("iptarget", DefaultIPCheckTarget, "ip check target")
	flag.Parse()
	st := sessiontester.New(*sessionCount)

	if *username == "" || *password == "" {
		fmt.Println("missing auth parameters")
		os.Exit(1)
	}

	start := time.Now()
	go func() {
		bar := progressbar.Default(100)
		for i := 0; i < 100; i++ {
			bar.Add(1)
			time.Sleep(time.Minute / 10)
		}
	}()

	st.DoItAllWithParams(*username, *password, *cc, *city, *ipTarget)

	fmt.Println("Done at: ", time.Now().String(), " in: ", time.Now().Sub(start).Minutes(), " minutes")
}
