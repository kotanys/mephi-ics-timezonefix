package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type State struct {
	Client    *http.Client
	FmtString string
}

func startServer(server *http.Server) {
	log.Println("starting server on port", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		log.Println("server closed with status:", err)
	}
}

func getURLParameters(u *url.URL) map[string]string {
	urlParams, _ := url.ParseQuery(u.RawQuery)
	mapParams := make(map[string]string)
	for k, v := range urlParams {
		mapParams[k] = v[0]
	}
	return mapParams
}

func (s *State) getICS(id int) (string, error) {
	url := fmt.Sprintf(s.FmtString, id)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0")

	resp, err := s.Client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("call to %v returned status code: %v", url, resp.StatusCode)
	}
	ics, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(ics), nil
}

func getTimezone(ics string) (string, error) {
	re := regexp.MustCompilePOSIX("^TZID:(.*)$")
	matches := re.FindAllString(ics, -1)
	foundCount := len(matches)
	if foundCount == 0 {
		return "", errors.New("no timezone found, using none")
	}
	if foundCount > 1 {
		return "", errors.New("multiple timezones found, using none")
	}
	timezone := re.ReplaceAllString(matches[0], "$1")
	timezone = strings.TrimSuffix(timezone, "\r") // in case file is in DOS format
	return timezone, nil
}

func fixUpICS(ics string) string {
	timezone, err := getTimezone(ics)
	if err != nil {
		log.Println(err)
		return ics
	}

	re := regexp.MustCompilePOSIX("DT(START|END):")
	icsFixed := re.ReplaceAllString(ics, fmt.Sprintf("DT$1;TZID=%v:", timezone))
	return icsFixed
}

func (s *State) fixHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	params := getURLParameters(req.URL)
	idstr, exist := params["id"]
	if !exist {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "id not found in URL params")
		return
	}
	id, err := strconv.Atoi(idstr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "id is not an int: %v\n", idstr)
		return
	}

	ics, err := s.getICS(id)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, err.Error())
		return
	}
	ics = fixUpICS(ics)

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, ics)
}

func main() {
	if len(os.Args) <= 1 {
		fmt.Printf("usage: %v port\n", os.Args[0])
		fmt.Println("provide port")
		os.Exit(1)
	}

	port, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Printf("usage: %v port\n", os.Args[0])
		fmt.Println("provide port")
		os.Exit(1)
	}

	s := &State{
		Client:    &http.Client{},
		FmtString: "https://home.mephi.ru/study_groups/%v/schedule.ics",
	}

	mux := &http.ServeMux{}
	mux.HandleFunc("/fix", s.fixHandler)
	server := &http.Server{Addr: fmt.Sprintf(":%v", port), Handler: mux}
	go startServer(server)

	stopSignal := make(chan os.Signal, 1)
	signal.Notify(stopSignal, os.Interrupt)
	<-stopSignal

	log.Println("SIGINT recieved. Shutting down the server.")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalln(err)
	}
}
