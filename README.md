# tpgl_inventory_srvr
Offline DB that supports CRUD operations over HTTP client

To use program, install 'main', 'model', 'view', 'ctrl', and 'util' folders under one parent directory and run 'main/main.go' and correct import paths if necessary. Once 'main.go' is running, go to 'http://localhost:8000/home' in browser to start using database. 

All db entries are saved offline for future use. Write operations require a lock and a write transaction to the offline database. Read operations do not require a lock, and data being read is served from memory, not disk. 
