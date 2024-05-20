package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type Config struct {
	NfUrls map[string]string `json:"nfUrls"`
}

var endpoint = "http://127.0.0.1:1000"
var client = &http.Client{
	Timeout: 10 * time.Second, // Timeout after 10 seconds
}
var config = Config{}

func removeEmpty(s []string) []string {
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}

func getEndpoint(url string) string {
	res := removeEmpty(strings.Split(url, "/"))
	if len(res) == 0 {
		fmt.Println("Path is empty, returning default endpoint")
		return endpoint
	}
	if res[0][0] != 'n' && res[0][0] != 'N' {
		fmt.Println("path does not begin with character n, returing default endpoint")
		return endpoint
	}
	part := res[0][1:]
	i := strings.LastIndex(part, "-")
	if i == -1 {
		fmt.Println("failed to parse url path", url)
		return endpoint
	}
	name := part[:i]
	finalUrl, found := config.NfUrls[strings.ToUpper(name)]
	if !found {
		fmt.Println("NF not found in the configuration, returing default endpoint")
		return endpoint
	}
	return finalUrl
}

func addCorsHeader(res http.ResponseWriter) {
	headers := res.Header()
	headers.Add("Access-Control-Allow-Origin", "*")
	headers.Add("Access-Control-Allow-Headers", "*")
	headers.Add("Access-Control-Allow-Methods", "*")
}

func handler(w http.ResponseWriter, r *http.Request) {
	addCorsHeader(w)
	fmt.Println("new request", r.Method, r.URL)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Create a new HTTP request with the same method, URL, and body as the original request
	targetURL := getEndpoint(r.URL.Path) + r.URL.String()
	fmt.Println("sending request to:", targetURL)
	proxyReq, err := http.NewRequest(r.Method, targetURL, r.Body)
	if err != nil {
		fmt.Println("failed to create request:", err)
		http.Error(w, "Error creating proxy request", http.StatusInternalServerError)
		return
	}

	// Copy the headers from the original request to the proxy request
	for name, values := range r.Header {
		for _, value := range values {
			proxyReq.Header.Add(name, value)
		}
	}

	// Send the proxy request using the custom transport
	resp, err := client.Do(proxyReq)
	if err != nil {
		fmt.Println("failed to send request:", err)
		http.Error(w, "Error sending proxy request", http.StatusInternalServerError)
		return
	}
	fmt.Println("sending reponse", resp.StatusCode)
	defer resp.Body.Close()

	// Copy the headers from the proxy response to the original response
	for name, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(name, value)
		}
	}

	// Set the status code of the original response to the status code of the proxy response
	w.WriteHeader(resp.StatusCode)

	// Copy the body of the proxy response to the original response
	io.Copy(w, resp.Body)
}

func retrieveApi(w http.ResponseWriter, r *http.Request) {
	addCorsHeader(w)
	fmt.Println("new request", r.Method, r.URL)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}
	Pathid := strings.TrimPrefix(r.URL.Path, "/api/v1/openapi/")
	fmt.Println("request for file", Pathid)
	dat, err := os.ReadFile("api/R15/" + Pathid)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(dat)
}

type OpenApiLink struct {
	Url  string `json:"url"`
	Name string `json:"name"`
}

func setEnv() {

	baseUri := "http://127.0.0.1:8080/api/v1/openapi/"
	entries, err := os.ReadDir("./api/R15")
	if err != nil {
		log.Fatal(err)
	}
	arr := make([]OpenApiLink, 0)
	for _, e := range entries {
		fullUri := baseUri + e.Name()
		basename := strings.Split(e.Name(), ".")[0]
		nameSlice := strings.Split(basename, "_")[1:]
		fullName := strings.Join(nameSlice, " ")
		obj := OpenApiLink{
			Url:  fullUri,
			Name: fullName,
		}
		arr = append(arr, obj)
	}
	final, err := json.MarshalIndent(arr, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	final2 := "URLS='" + string(final) + "'"

	fo, err := os.Create("../swagger-docker/.env")
	if err != nil {
		log.Fatal(err)
	}
	fo.WriteString(final2)
}

func loadConfig(dat []byte) {
	err := json.Unmarshal(dat, &config)
	if err != nil {
		fmt.Println("error parsing config: ", err)
	}
}

func main() {
	args := os.Args[1:]
	if len(args) != 0 {
		setEnv()
		return
	}

	dat, err := os.ReadFile("config.json")
	if err == nil {
		loadConfig(dat)
	} else {
		fmt.Println("error opening config file: ", err)
	}

	http.HandleFunc("/", handler)
	http.HandleFunc("/api/v1/openapi/", retrieveApi)
	fmt.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", nil)
}
