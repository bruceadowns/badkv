package lib

import (
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

const (
	sep = "-"
)

var (
	input *Input
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("todo root handler"))
}

func statsHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("todo admin stats"))
}

func metricsHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("todo admin metrics"))
}

func membersHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("todo admin members"))
}

func keyGetHandler(w http.ResponseWriter, r *http.Request) {
	//
	// GET
	// /api/v1/keys/{tenant}/{name}
	//

	fields := strings.Split(r.URL.Path, "/")
	if len(fields) != 6 {
		log.Fatalf("Invalid GET path: %s", r.URL.Path)
	}

	tenant := fields[4]
	name := fields[5]
	k := Key(tenant + sep + name)
	log.Printf("GET %s", k)

	v, err := store.get(k)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
		return
	}

	w.Write(v)
}

func keyPostHandler(w http.ResponseWriter, r *http.Request) {
	//
	// POST
	// /api/v1/keys/{tenant}/{name}
	//

	fields := strings.Split(r.URL.Path, "/")
	if len(fields) != 6 {
		log.Fatalf("Invalid POST path: %s", r.URL.Path)
	}

	tenant := fields[4]
	name := fields[5]
	k := Key(tenant + sep + name)
	d, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatalf("Error reading POST body: %s", err)
	}
	tsv := &TimeStampValue{data: d, timestamp: time.Now()}

	// if leader exists, forward and return
	if input.leader != "" {
		if err := forwardLeaderPut(input.leader, r.URL.Path, tsv); err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		w.WriteHeader(http.StatusAccepted)
		return
	}

	// put locally; not atomic
	log.Printf("PUT %s '%s'", k, tsv)
	if err := store.put(k, tsv); err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	// if followers exist, forward
	for _, follower := range input.followers {
		if err := forwardFollowerPut(follower, r.URL.Path, tsv); err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
	}

	w.WriteHeader(http.StatusAccepted)
}

func keyPostTimestampHandler(w http.ResponseWriter, r *http.Request) {
	//
	// POST
	// /api/v1/keys/{tenant}/{name}/{timestamp}
	//

	fields := strings.Split(r.URL.Path, "/")
	if len(fields) != 7 {
		log.Fatalf("Invalid POST path: %s", r.URL.Path)
	}

	tenant := fields[4]
	name := fields[5]
	timestamp := fields[6]

	k := Key(tenant + sep + name)
	d, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatalf("Error reading POST body: %s", err)
	}
	nsecs, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		log.Fatal(err)
	}
	tsv := &TimeStampValue{data: d, timestamp: time.Unix(0, nsecs)}

	log.Printf("PUT %s '%s'", k, tsv.data)
	if err := store.put(k, tsv); err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func keyDeleteHandler(w http.ResponseWriter, r *http.Request) {
	//
	// DELETE
	// /api/v1/keys/{tenant}/{name}
	//

	fields := strings.Split(r.URL.Path, "/")
	if len(fields) != 6 {
		log.Fatalf("Invalid DELETE path: %s", r.URL.Path)
	}

	tenant := fields[4]
	name := fields[5]
	k := Key(tenant + sep + name)
	log.Printf("DELETE %s", k)

	// if leader exists, forward and return
	if input.leader != "" {
		forwardLeaderDelete(input.leader, r.URL.Path)
		return
	}

	// delete locally; not atomic
	if err := store.delete(k); err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
		return
	}

	// if followers exist, forward
	for _, follower := range input.followers {
		if err := forwardFollowerDelete(follower, r.URL.Path, &TimeStampValue{timestamp: time.Now()}); err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
	}

	w.WriteHeader(http.StatusAccepted)
}

func keyDeleteTimestampHandler(w http.ResponseWriter, r *http.Request) {
	//
	// DELETE
	// /api/v1/keys/{tenant}/{name}/{timestamp}
	//

	fields := strings.Split(r.URL.Path, "/")
	if len(fields) != 7 {
		log.Fatalf("Invalid DELETE path: %s", r.URL.Path)
	}

	tenant := fields[4]
	name := fields[5]
	k := Key(tenant + sep + name)
	log.Printf("DELETE %s", k)

	if err := store.delete(k); err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

// StartRootHandler starts the http handler
func StartRootHandler(wg *sync.WaitGroup, in *Input) {
	wg.Add(1)
	input = in

	go func() {
		defer wg.Done()

		//
		// /admin/{stats,metrics,members}
		// /api/v1/keys/{tenant}/{name} - GET
		// /api/v1/keys/{tenant}/{name}/{timestamp} - POST
		// /api/v1/keys/{tenant}/{name}/{timestamp} - DELETE
		// /api/v1/take/{timestamp} - GET
		//

		router := mux.NewRouter()

		// example admin endpoints
		router.HandleFunc("/admin/stats", statsHandler)
		router.HandleFunc("/admin/metrics", metricsHandler)
		router.HandleFunc("/admin/members", membersHandler)

		// get/put/remove/take endpoints
		router.HandleFunc("/api/v1/keys/{tenant}/{name}", keyGetHandler).Methods("GET")
		router.HandleFunc("/api/v1/keys/{tenant}/{name}", keyPostHandler).Methods("POST")
		router.HandleFunc("/api/v1/keys/{tenant}/{name}/{timestamp}", keyPostTimestampHandler).Methods("POST")
		router.HandleFunc("/api/v1/keys/{tenant}/{name}", keyDeleteHandler).Methods("DELETE")
		router.HandleFunc("/api/v1/keys/{tenant}/{name}/{timestamp}", keyDeleteTimestampHandler).Methods("DELETE")
		//router.HandleFunc("/api/v1/take/{timestamp}", takeHandler).Methods("GET")

		// catch-all handler
		router.HandleFunc("/", rootHandler)

		log.Print("listen and serve mux router")
		server := &http.Server{
			Handler:      router,
			Addr:         in.listen,
			WriteTimeout: 5 * time.Second,
			ReadTimeout:  5 * time.Second,
		}
		log.Fatal(server.ListenAndServe())
	}()
}
