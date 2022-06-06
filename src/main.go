// RedNaga / Tim Strazzere (c) 2018-*

package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	raven "github.com/getsentry/raven-go"
	"github.com/gorilla/mux"
)

var (
	port         = 3000
	readTimeout  = time.Duration(10)
	writeTimeout = time.Duration(10)

	ravenConfig = "CHANGEME-sentry-url"
)

var stats Statistics

var (
	adminAPIKey = "CHANGEME"

	allowedAPIKeys = []Key{
		{ // Unused test key
			Key:   "CHANGEME2",
			Limit: 1000,
		},
	}

	disallowedAPIKeys = []string{}
)

func main() {
	RedisInit()
	err := Ping()
	if err != nil {
		log.Fatalf("Redis not found : Ping err ? %v\n", err)
		return
	}

	raven.SetDSN(ravenConfig)
	log.Printf("Starting distil proof of ignorance...\n")
	stats.StartTime = getTime().Format("2006-01-02 15:04:05")
	router := mux.NewRouter().StrictSlash(true)
	// /api/v1 is now handled by nginx forwarding
	router.Handle("/health", http.HandlerFunc(raven.RecoveryHandler(handleSimpleHealth)))
	router.Handle("/dump", http.HandlerFunc(raven.RecoveryHandler(handleHealth)))
	router.Handle("/session", http.HandlerFunc(raven.RecoveryHandler(handleSession))).Methods("POST")
	router.Handle("/stats", http.HandlerFunc(raven.RecoveryHandler(handleStats)))
	router.NotFoundHandler = http.HandlerFunc(notFound)

	log.Printf("Stats : %+v", stats)
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      router,
		ReadTimeout:  readTimeout * time.Second,
		WriteTimeout: writeTimeout * time.Second,
	}
	raven.CapturePanic(func() {
		err := server.ListenAndServe()
		raven.CaptureError(err, nil)
		log.Fatal(err)
	}, nil)
}

// curl https://127.0.0.1:3000/api/v1/stats
func handleStats(w http.ResponseWriter, r *http.Request) {
	data, _ := GetKeys("page:*")
	stats := make(map[string]int64)
	for _, key := range data {
		b, _ := Get(key)
		stats[key], _ = strconv.ParseInt(string(b), 10, 64)
	}

	data, _ = GetKeys("asset:*")
	assets := make(map[string]int64)
	for _, key := range data {
		b, _ := Get(key)
		assets[key], _ = strconv.ParseInt(string(b), 10, 64)
	}

	data, _ = GetKeys("auth:*")
	auth := make(map[string]int64)
	for _, key := range data {
		b, _ := Get(key)
		auth[key], _ = strconv.ParseInt(string(b), 10, 64)
	}

	data, _ = GetKeys("error:*")
	errors := make(map[string]int64)
	for _, key := range data {
		b, _ := Get(key)
		errors[key], _ = strconv.ParseInt(string(b), 10, 64)
	}

	data, _ = GetKeys("stat:*")
	sessions := make(map[string]int64)
	for _, key := range data {
		b, _ := Get(key)
		sessions[key], _ = strconv.ParseInt(string(b), 10, 64)
	}

	data, _ = GetKeys("usage:*")
	usageByKeys := make(map[string]int64)
	for _, key := range data {
		b, _ := Get(key)
		usageByKeys[key], _ = strconv.ParseInt(string(b), 10, 64)
	}

	data, _ = GetKeys("overage:*")
	overageByKeys := make(map[string]int64)
	for _, key := range data {
		b, _ := Get(key)
		overageByKeys[key], _ = strconv.ParseInt(string(b), 10, 64)
	}

	d := make(map[string]interface{})
	d["errors"] = errors
	d["auth"] = auth
	d["assets"] = assets
	d["stats"] = stats
	d["sessions"] = sessions
	d["usage"] = usageByKeys
	d["overage"] = overageByKeys

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(d)
}

// curl -X POST -H "Auth-Key: bf9f83bfda32ce7afabbe54d3f2f846638b919c4" https://127.0.0.1:3000/api/v1/session -d @test-data/jsonpayload.json
func handleSession(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	connectingIP := r.Header.Get("Cf-Connecting-Ip")

	authed := false
	authKey := r.Header.Get("Auth-Key")
	for _, key := range allowedAPIKeys {
		if strings.EqualFold(key.Key, authKey) {

			todaysUsageBytes, _ := Get(fmt.Sprintf("usage:%s:%s", authKey, getTime().Format("2006-01-02")))
			if len(string(todaysUsageBytes)) > 0 {
				todaysUsage, err := strconv.ParseInt(string(todaysUsageBytes), 10, 64)
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(Exception{
						Code:    10001,
						Message: "Error occured getting proper key usage data",
					})
					raven.CaptureError(err, nil)
					return
				}

				if todaysUsage >= key.Limit {
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(Exception{
						Code:    10002,
						Message: "API Key has exceeded maximum daily usage!",
					})
					Incr(fmt.Sprintf("overage:%s:%s", authKey, getTime().Format("2006-01-02")))
					return
				}
			} else {
				todaysUsage := int64(0)
				if todaysUsage >= key.Limit {
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(Exception{
						Code:    10002,
						Message: "API Key has exceeded maximum daily usage!",
					})
					Incr(fmt.Sprintf("overage:%s:%s", authKey, getTime().Format("2006-01-02")))
					return
				}
			}

			authed = true
			break
		}
	}

	if !authed {
		bannedAuth := false
		for _, key := range disallowedAPIKeys {
			if strings.EqualFold(key, authKey) {
				bannedAuth = true
				break
			}
		}
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(Exception{
			Code:    10003,
			Message: fmt.Sprintf("Not authorized"),
		})
		if bannedAuth {
			Incr(fmt.Sprintf("auth:banned:%s:%s:%s", "session", getTime().Format("2006-01-02"), connectingIP))
		} else {
			Incr(fmt.Sprintf("auth:bad:%s:%s:%s", "session", getTime().Format("2006-01-02"), connectingIP))
		}
		return
	}

	Incr(fmt.Sprintf("usage:%s:%s", authKey, getTime().Format("2006-01-02")))
	Incr(fmt.Sprintf("stat:%s:%s", "session", getTime().Format("2006-01-02")))

	var sessionData SessionData
	byteBody := new(bytes.Buffer)
	_, err := io.Copy(byteBody, r.Body)
	body := byteBody.Bytes()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Exception{
			Code:    10004,
			Message: fmt.Sprintf("Error occured reading post body : %v", err),
		})

		Incr(fmt.Sprintf("error:400:%s:%s:%s", "session", getTime().Format("2006-01-02"), connectingIP))
		return
	}

	err = json.Unmarshal(body, &sessionData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Exception{
			Code:    10005,
			Message: fmt.Sprintf("Error occured unmarshaling json body : %v", err),
		})

		Incr(fmt.Sprintf("error:400:%s:%s:%s", "session", getTime().Format("2006-01-02"), connectingIP))
		return
	}

	decoded, err := base64.StdEncoding.DecodeString(sessionData.JsData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Exception{
			Code:    10006,
			Message: fmt.Sprintf("Error occured base64 decoding jsData : %v", err),
		})

		Incr(fmt.Sprintf("error:400:%s:%s:%s", "session", getTime().Format("2006-01-02"), connectingIP))
		return
	}

	hash := sha1.Sum(decoded)
	if !strings.EqualFold(sessionData.JsSha1, fmt.Sprintf("%x", hash)) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Exception{
			Code:    10007,
			Message: "Data did not match provided sha1",
		})

		Incr(fmt.Sprintf("error:400:%s:%s:%s", "session", getTime().Format("2006-01-02"), connectingIP))
		return
	}

	stats.SessionData = appendNoDupes(sessionData)
	Incr(fmt.Sprintf("asset:%s:%s", sessionData.JsSha1, getTime().Format("2006-01-02")))

	distilConfig, err := parseFingerprints(decoded)
	if err != nil {
		raven.CaptureError(err, nil)

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Exception{
			Code:    10008,
			Message: "Error getting correct data out of js blob",
		})

		Incr(fmt.Sprintf("error:400:%s:%s:%s", "session", getTime().Format("2006-01-02"), connectingIP))
		return
	}

	proof := getProof()
	proof, err = workOnProof(proof, 8)
	if err != nil {
		raven.CaptureError(err, nil)

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Exception{
			Code:    10009,
			Message: "Error generating proof of work",
		})

		Incr(fmt.Sprintf("error:400:%s:%s:%s", "session", getTime().Format("2006-01-02"), connectingIP))
		return
	}

	userAgent := r.Header.Get("User-Agent")
	proofData, err := getEscapedProofQuery(proof, userAgent)
	if err != nil {
		raven.CaptureError(err, nil)

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Exception{
			Code:    10010,
			Message: "Error generating escaped proof query",
		})

		Incr(fmt.Sprintf("error:400:%s:%s:%s", "session", getTime().Format("2006-01-02"), connectingIP))
		return
	}

	postTask := Task{
		URI:      distilConfig.Path,
		Method:   "POST",
		Headers:  []string{fmt.Sprintf("%s:%s", "X-Distil-Ajax", distilConfig.XDistilAjax)},
		Data:     proofData,
		Interval: 0,
	}

	headTask := Task{
		URI:      distilConfig.Path[0 : strings.Index(distilConfig.Path, ".js")+3],
		Method:   "HEAD",
		Headers:  []string{},
		Interval: distilConfig.HeartbeatInterval,
	}

	var tasks []Task
	tasks = append(tasks, postTask)
	tasks = append(tasks, headTask)

	var sessionDataResponse SessionDataResponse
	sessionDataResponse.Tasks = tasks
	sessionDataResponse.Headers = []string{fmt.Sprintf("%s:%s", "X-Distil-Ajax", distilConfig.XDistilAjax)}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(sessionDataResponse)

	Incr(fmt.Sprintf("page:%s:%s", "session", getTime().Format("2006-01-02")))
}

func handleSimpleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("OK")
	Incr(fmt.Sprintf("page:%s:%s", "health", getTime().Format("2006-01-02")))
}

// curl https://127.0.0.1:3000/api/v1/dump -H 'Auth-Key:3c2ddb9ededb236f42ae0f9ade351ad92b8f2f3f'
func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	authed := false
	authKey := r.Header.Get("Auth-Key")
	if strings.EqualFold(adminAPIKey, authKey) {
		authed = true
	}

	if !authed {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(Exception{
			Code:    10003,
			Message: fmt.Sprintf("Not authorized"),
		})

		connectingIP := r.Header.Get("Cf-Connecting-Ip")
		Incr(fmt.Sprintf("auth:bad:%s:%s:%s", "session", getTime().Format("2006-01-02"), connectingIP))
		return
	}

	jsonData, err := json.Marshal(stats)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Exception{
			Code:    10011,
			Message: fmt.Sprintf("Error occured unmarshaling json stats : %v", err),
		})

		raven.CaptureError(err, nil)
	} else {
		w.WriteHeader(http.StatusOK)
		w.Write(jsonData)
		Incr(fmt.Sprintf("page:%s:%s", "dump", getTime().Format("2006-01-02")))
	}
}

func notFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(Exception{
		Code:    404,
		Message: "Bad request",
	})
	connectingIP := r.Header.Get("Cf-Connecting-Ip")
	Incr(fmt.Sprintf("error:404:%s:%s", getTime().Format("2006-01-02"), connectingIP))
}
