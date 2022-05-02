# reference_manager

This is a browser-based reference manager, written using Go's standard `net/http` and `html/template` libraries, with the ability to generate formatted bibliographies. Data is stored in a postgres database, and bibliographies are output in one of two ways: 1) to a text file in the user's downloads folder, or 2) as a plain text response in the browser.

You can see a live version of the app here <ref-manager.michaelfatigati.dev>. Note, however, that accessing the app there means you can only output a bibliography in the browser (since the Go `os` library interacts with the server's OS, not the client's).

To get full functionality, download the Go package, and run it on your computer, as follows:

- create a `.env` file in the project root, with the following variables:
  - `USER` (your postgres database username)
  - `PASSWORD` (your postgres database password)
  - `DEFAULT_DB` (the default database, normally just "postgres")
    - the application first connects to the default database, before creating the database that it will use.
  - `PORT` (the port on which you want the application to run)