package main

// PfsHandler is the core route for modifying the contents of the fileystem.
// Changes are not replicated until a call to CommitHandler.
func TestHandler(w http.ResponseWriter, r *http.Request) {
}

// TesterMux creates a multiplexer for a Tester.
func TesterMux() *http.ServeMux {
	mux := http.NewServeMux()

	commitHandler := func(w http.ResponseWriter, r *http.Request) {
		CommitHandler(w, r, fs)
	}

	pfsHandler := func(w http.ResponseWriter, r *http.Request) {
		PfsHandler(w, r, fs)
	}

	browseHandler := func(w http.ResponseWriter, r *http.Request) {
		BrowseHandler(w, r, fs)
	}

	mux.HandleFunc("/commit", commitHandler)
	mux.HandleFunc("/pfs/", pfsHandler)
	mux.HandleFunc("/browse", browseHandler)
	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) { fmt.Fprint(w, "pong\n") })

	return mux
}

// RunServer runs a master server listening on port 80
func RunServer(fs *btrfs.FS) {
	http.ListenAndServe(":80", TesterMux(fs))
}

func main() {
	log.SetFlags(log.Lshortfile)
	log.Print("Listening on port 80...")
	RunServer(fs)
}
