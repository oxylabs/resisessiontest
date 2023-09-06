package sessiontester

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

const ProxyHost = "pr.oxylabs.io:7777"

type SessTester struct {
	creds                  string
	sessCount              int
	sessMap                map[string]time.Duration
	doubleCreateMismatches int
	sessMu                 *sync.Mutex
	wg                     *sync.WaitGroup
}

type IPInfoResponse struct {
	IP      string `json:"ip"`
	City    string `json:"city"`
	Country string `json:"country"`
}

func New(sessCount int) *SessTester {
	st := &SessTester{
		sessCount: sessCount,
		sessMap:   map[string]time.Duration{},
		sessMu:    &sync.Mutex{},
		wg:        &sync.WaitGroup{},
	}
	return st
}

func (st *SessTester) TestSessions(username, password, cc, city, iptarget string) {
	for i := 0; i < st.sessCount; i++ {
		sessID := fmt.Sprintf("test%d_%d", time.Now().Unix(), i)
		go st.testSession(sessID, username, password, cc, city, iptarget)
	}
	time.Sleep(time.Second)
	st.wg.Wait()
	secs := 0.0
	sessions := 0
	over9Mins30Secs := 0
	badSessions := []string{}
	for k, v := range st.sessMap {
		fmt.Println("Session: ", k, " lasted: ", v.Minutes())
		secs += v.Seconds()
		sessions++
		if v.Seconds() > 575.0 {
			over9Mins30Secs++
		}
		if v.Seconds() < 450.0 {
			badSessions = append(badSessions, k)
		}
	}
	fmt.Println("Sessions tested: ", sessions)
	fmt.Println("Lasted over 9 mins 30 secs: ", over9Mins30Secs)
	fmt.Println("Avg sess duration: ", secs/float64(sessions), " seconds")
	fmt.Println("Bad sessions (<450s): ", strings.Join(badSessions, ","))
	reportFile, err := os.Create("session_test_report_" + strconv.Itoa(int(time.Now().UnixMilli())) + ".txt")
	defer reportFile.Close()
	if err != nil {
		fmt.Println("failed to create report file")
		return
	}
	reportFile.Write([]byte("Completed on:" + time.Now().String() + "\n"))
	reportFile.Write([]byte("Sessions tested:" + strconv.Itoa(sessions) + "\n"))
	reportFile.Write([]byte("Lasted over 9 mins 30 secs:" + strconv.Itoa(over9Mins30Secs) + "\n"))
	reportFile.WriteString(fmt.Sprintf("Avg sess duration (seconds):, %.2f\n\n", secs/float64(sessions)))
	for k, v := range st.sessMap {
		reportFile.WriteString("Session: " + k + " lasted (mins): " + fmt.Sprintf("%.2f\n", v.Minutes()))
	}
}

func (st *SessTester) testSession(sessID, username, password, cc, city, iptarget string) {
	st.wg.Add(1)
	defer st.wg.Done()

	proxyUrlStr := "http://customer-" + username + "-cc-" + cc + "-city-" + city + "-sessid-" + sessID + ":" + password + "@" + ProxyHost
	proxyUrl, _ := url.Parse(proxyUrlStr)
	t := &http.Transport{
		Proxy: http.ProxyURL(proxyUrl),
	}
	c := &http.Client{
		Transport: t,
		Timeout:   30 * time.Second,
	}
	res, err := c.Get(iptarget)
	if err != nil {
		fmt.Println("GET err: ", err)
		if strings.Contains(err.Error(), "Proxy Authentication Required") {
			fmt.Println("\nLikely invalid proxy auth provided, exiting")
			os.Exit(1)
		}
		return
	}
	if res.StatusCode != 200 {
		fmt.Println("Status: ", res.StatusCode)
		return
	}
	startTime := time.Now()
	b, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("body read err: ", err)
		return
	}
	initialIP := string(b)
	errs := 0
	for {
		timePassed := time.Now().Sub(startTime)
		if errs > 5 {
			fmt.Println("\nToo many errors, session ", sessID, " considered dropped after ", timePassed.Minutes(), " minutes.")
			st.sessMu.Lock()
			st.sessMap[sessID] = timePassed
			st.sessMu.Unlock()
			return
		}
		time.Sleep(5 * time.Second)
		res, err := c.Get(iptarget)
		if err != nil {
			errs++
			continue
		}
		if res.StatusCode != 200 {
			errs++
			continue
		}
		b, err := io.ReadAll(res.Body)
		if err != nil {
			errs++
			fmt.Println("body read err: ", err)
			continue
		}
		if string(b) != initialIP {
			fmt.Println("\nIP changed from ", initialIP, " to ", string(b), " after ", timePassed.Minutes(), " minutes.")
			st.sessMu.Lock()
			st.sessMap[sessID] = timePassed
			st.sessMu.Unlock()
			return
		}
	}
}
