package main

import (
	"github.com/gorilla/sessions"
	"html/template"
	"net/http"
	"os/exec"
)

var (
	store = sessions.NewCookieStore([]byte("your-very-secret-key"))
)

func main() {
	http.HandleFunc("/", formHandler)
	http.HandleFunc("/execute", executeCommand)
	http.ListenAndServe(":8080", nil)
}

func formHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session-name")

	// Retrieve the last output from the session
	commandOutput, _ := session.Values["output"].(string)

	tmpl := template.Must(template.New("form").Parse(`
        <html>
            <head><title>Command Executor</title></head>
            <body>
            <h1>Command Executor</h1>
                <form action="/execute" method="post">
                    <input type="text" name="command" placeholder="">
                    <input type="submit" value="Run">
                </form>
                <h2>Output</h2>
                <pre>{{.}}</pre>
            </body>
        </html>
    `))

	tmpl.Execute(w, commandOutput)
}

func executeCommand(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	r.ParseForm()
	command := r.FormValue("command")

	// Security Warning: Executing arbitrary commands from user input is dangerous
	cmd := exec.Command("/bin/sh", "-c", command)
	output, err := cmd.CombinedOutput()
	if err != nil {
		// http.Error(w, err.Error(), http.StatusInternalServerError)
		session, _ := store.Get(r, "session-name")
		session.Values["output"] = string(output)
		session.Save(r, w)

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	session, _ := store.Get(r, "session-name")
	session.Values["output"] = string(output)
	session.Save(r, w)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
