# go-data-collector
A structured data collector that pulls data from various sources and exposes them in html format. Each provider will be listed together with its files. 

New data sources can be added easily.

Currently the program is able to scan the local file system for .csv and .prn files.

### Quick setup
    - go get github.com/g-leon/go-data-collector
    - cd $GOPATH/src/github.com/g-leon/go-data-collector
    - go test ./...
    - go run main.go
    
### Access data list
    - localhost:3000

### Future improvements 
    - a table caching system (making a file load on each call can result in major slow down) 
    - invalidate caching only when files change;
    - if tables are big a pagination system is needed;
    - if table schemas will modify a dynamic schema loaded from file header will be needed (with fallback if there is no header);
    - if this service becomes an API to be consumed by other services then each field should be represented more accurate as opposed to a string;
    - it would be useful to detect on case by case basis if the table has a header to be skipped on loading it;
    - read data using correct encoding;
    - find a smarter way to deal with .prn files since they are not tab formated;