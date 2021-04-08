// Copyright 2018 Traffic Label. All rights reserved.
// Source Code Written for public use.
// @author Noy Hillel

package utils

import (
	"bytes"
	"encoding/json"
	"github.com/dustin/go-humanize"
	"github.com/fatih/color"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

func ConvertToFloat(s string) float64 {
	num, err := strconv.ParseFloat(s, 64)
	if err != nil {
		log.Printf("Could not convert %v to a float!", s)
		return 0
	}
	return num
}

func PrintEmoji(emoji string, emojiMap map[string]string) string {
	for key, value := range emojiMap {
		if ":"+emoji+":" == key {
			return value
		}
	}
	return ""
}

func Commaf(v float64) string {
	return humanize.Commaf(v)
}

func Comma(v int64) string {
	return humanize.Comma(v)
}

func BubbleSortDesc(arr []string) []string {
	temp := ""
	for i := 0; i < len(arr); i++ {
		for j := 1; j < len(arr)-i; j++ {
			if arr[j-1] < arr[j] {
				temp = arr[j-1]
				arr[j-1] = arr[j]
				arr[j] = temp
			}
		}
	}
	return arr
}

func CheckDBErr(err error, db string) {
	if err != nil {
		color.Red("Error connecting to database: %v\nError:%v ", db, err.Error())
		return
	}
}

func String(n int) string {
	buf := [11]byte{}
	pos := len(buf)
	i := int64(n)
	signed := i < 0
	if signed {
		i = -i
	}
	for {
		pos--
		buf[pos], i = '0'+byte(i%10), i/10
		if i == 0 {
			if signed {
				pos--
				buf[pos] = '-'
			}
			return string(buf[pos:])
		}
	}
}

func ProperlyFormatDate(date string, headerMap map[string]interface{}) (string, error) {
	properDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		headerMap["Err"] = "There was an error loading the date!"
		return "", err
	}
	return properDate.Format("January 2 2006"), err
}

func FormatFloat(num float64) string {
	return strconv.FormatFloat(num, 'f', 2, 64)
}

func SendHTTPError(reason string, rw http.ResponseWriter) {
	http.Error(rw, reason, http.StatusBadRequest)
}

func RemoveDuplicates(elements []string) []string {
	encountered := map[string]bool{}
	var result []string
	for v := range elements {
		if encountered[elements[v]] == true {
		} else {
			encountered[elements[v]] = true
			result = append(result, elements[v])
		}
	}
	return result
}

func TrimCompletelyAfter(s, remover string) string {
	if idx := strings.Index(s, remover); idx != -1 {
		return s[:idx]
	}
	return s
}

func Interface(f interface{}) string {
	return f.(string)
}

func MonthInSlice(a interface{}, list []time.Month) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func DaysIn(m time.Month, year int) int {
	return time.Date(year, m+1, 0, 0, 0, 0, 0, time.UTC).Day()
}

func GetRealAddr(r *http.Request) string {
	remoteIP := ""
	if parts := strings.Split(r.RemoteAddr, ":"); len(parts) == 2 {
		remoteIP = parts[0]
	}
	if xff := strings.Trim(r.Header.Get("X-Forwarded-For"), ","); len(xff) > 0 {
		addrs := strings.Split(xff, ",")
		lastFwd := addrs[len(addrs)-1]
		if ip := net.ParseIP(lastFwd); ip != nil {
			remoteIP = ip.String()
		}
	} else if xri := r.Header.Get("X-Real-Ip"); len(xri) > 0 {
		if ip := net.ParseIP(xri); ip != nil {
			remoteIP = ip.String()
		}
	}
	return remoteIP
}

func DenyAccess(w http.ResponseWriter, data string) {
	color.HiRed("[%v] Access Denied From: %v", time.Now().Format(time.RFC850), color.HiBlueString(data))
	http.Error(w, "You can't access this from "+data, 401)
}

func RedirectToHome(w http.ResponseWriter, req *http.Request) {
	http.Redirect(w, req, "/", http.StatusFound)
}

func OpenFile(path string) (file io.ReadWriteCloser, err error) {
	file, err = os.Open(path)
	return
}

func Mode(mode []string) string {
	if len(mode) == 0 {
		return ""
	}
	var modeMap = map[string]int{}
	var maxEl = mode[0]
	var maxCount = 1
	for i := 0; i < len(mode); i++ {
		var el = mode[i]
		if modeMap[el] == 0 {
			modeMap[el] = 1
		} else {
			modeMap[el]++
		}
		if modeMap[el] > maxCount {
			maxEl = el
			maxCount = modeMap[el]
		}
	}
	return maxEl
}

func GetExchangeRates(currency string, fallback float64, apiKey string) float64 {
	req, err := http.Get("https://v6.exchangerate-api.com/v6/" + apiKey + "/latest/" + currency)
	if err != nil {
		log.Printf("Something went wrong getting the API: %v", err.Error())
		return fallback
	}
	type Rates struct {
		ConversionRates struct {
			GBP float64 `json:"GBP"`
		} `json:"conversion_rates"`
	}
	defer req.Body.Close()
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Printf("Error with reading body.. falling back on %v. Error: %v", fallback, err.Error())
		return fallback
	}
	var cR Rates
	if err = json.Unmarshal(body, &cR); err != nil {
		if err != nil {
			log.Printf("Error with unmarshal, falling back on %v. Error: %v", fallback, err.Error())
			return fallback
		}
	}
	return cR.ConversionRates.GBP
}

func JsonPrettyPrint(in string) string {
	var out bytes.Buffer
	err := json.Indent(&out, []byte(in), "", "\t")
	if err != nil {
		return in
	}
	return out.String()
}

func GetDaysInMonth(month string, year int) int {
	switch month {
	case time.January.String():
		return 31
		break
	case time.February.String():
		wholeYear := time.Date(year, time.December, 31, 0, 0, 0, 0, time.Local)
		days := wholeYear.YearDay()
		if days > 365 {
			return 29
			break
		} else {
			return 28
			break
		}
		break
	case time.March.String():
		return 31
		break
	case time.April.String():
		return 30
		break
	case time.May.String():
		return 31
		break
	case time.June.String():
		return 30
		break
	case time.July.String():
		return 31
		break
	case time.August.String():
		return 31
		break
	case time.September.String():
		return 30
		break
	case time.October.String():
		return 31
		break
	case time.November.String():
		return 30
		break
	case time.December.String():
		return 31
		break
	default:
		break
	}
	return 31
}

func GetMonthFromName(month string) time.Month {
	switch month {
	case time.January.String():
		return time.January
		break
	case time.February.String():
		return time.February
		break
	case time.March.String():
		return time.March
		break
	case time.April.String():
		return time.April
		break
	case time.May.String():
		return time.May
		break
	case time.June.String():
		return time.June
		break
	case time.July.String():
		return time.July
		break
	case time.August.String():
		return time.August
		break
	case time.September.String():
		return time.September
		break
	case time.October.String():
		return time.October
		break
	case time.November.String():
		return time.November
		break
	case time.December.String():
		return time.December
		break
	default:
		break
	}
	return time.January
}

func IntArrayToString(A []int, delim string) string {
	var buffer bytes.Buffer
	for i := 0; i < len(A); i++ {
		buffer.WriteString(strconv.Itoa(A[i]))
		if i != len(A)-1 {
			buffer.WriteString(delim)
		}
	}
	return buffer.String()
}

func ReverseList(reversed []interface{}) []interface{} {
	sort.SliceStable(reversed, func(i, j int) bool {
		return true
	})
	return reversed
}

func CurrencySymbol(country, currency string) string {
	if country == "United Kingdom" {
		return "£" + currency
	}
	if country == "Sweden" || country == "Norway" {
		return currency + "kr"
	}
	if country == "Canada" || country == "New Zealand" {
		return "$" + currency
	} else {
		return "€" + currency
	}
}

func ParseCSVFile(fileName string) *os.File {
	file, err := os.Open(fileName)
	if err != nil {
		log.Printf("err opening file, %v", err.Error())
	}
	return file
}

func FormatDateWithSuffix(t time.Time) string {
	suffix := "th"
	switch t.Day() {
	case 1, 21, 31:
		suffix = "st"
	case 2, 22:
		suffix = "nd"
	case 3, 23:
		suffix = "rd"
	}
	return t.Format("January 2" + suffix + ", 2006")
}

func DownloadAndSaveFile(filepath string, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, resp.Body)
	return err
}

func Bod(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, t.Location())
}

func Eod(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 23, 59, 59, 0, t.Location())
}

func GetRoad(lat, long, apiKey string) string {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.opencagedata.com/geocode/v1/json?q="+lat+"+"+long+"&key="+apiKey, nil)
	if err != nil {
		return err.Error()
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err.Error())
	}
	defer resp.Body.Close()
	type Results struct {
		Formatted string `json:"formatted"`
	}
	type Geo struct {
		Results []Results `json:"results"`
	}
	var r Geo
	respBody, err := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(respBody, &r); err != nil {
		log.Println(err.Error())
	}
	for _, c := range r.Results {
		return c.Formatted
	}
	return "Not found"
}

func GetExchangeRateFor(currency, toCurrency string) float64 {
	type CurrencyResult struct {
		Rates map[string]float64 `json:"rates"`
	}
	re, err := http.Get("https://api.exchangeratesapi.io/latest?base=" + currency)
	if err != nil {
		log.Printf("Something went wrong getting the API: %v", err.Error())
		return 0
	}
	defer re.Body.Close()
	body, err := ioutil.ReadAll(re.Body)
	if err != nil {
		log.Printf("Error with reading body.. falling back on %v. Error: %v", 0, err.Error())
		return 0
	}
	var cR CurrencyResult
	if err = json.Unmarshal(body, &cR); err != nil {
		log.Println(err.Error())
	}
	return cR.Rates[toCurrency]
}
