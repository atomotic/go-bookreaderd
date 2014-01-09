# go-bookreaderd

**internet archive bookreader** in a standalone server, built in **Go**  
static assets from [BookReader](https://github.com/openlibrary/bookreader/tree/master/BookReader) are self-contained into the binary with [nrsc](http://bitbucket.org/tebeka/nrsc) 

## setup

### install required packages

	go get bitbucket.org/tebeka/nrsc
	go get bitbucket.org/tebeka/nrsc/nrsc
	go get github.com/codegangsta/martini
	
### compile

	go build bookreaderd.go
	nrsc bookreaderd assets
	
### run

	./bookreaderd
	open http://localhost:8088/book/example


## publish books

put in ```books/{DIRECTORY}``` a bunch of jpg files in the form ```page00n.jpg```
and open ```http://localhost:8088/book/{DIRECTORY}```


	
## todo
- use html/template to render index and not fmt.Sprintf
- automatic index at / with list of available books
- automatic generation of ```br.numLeafs = ;```
- serve jpg in .zip containers, as internetarchive (maybe with automatic [jp2->jpg](https://github.com/openlibrary/bookreader/blob/master/BookReaderIA/datanode/BookReaderImages.inc.php))


## credits
```func nrscStatic()``` is inspired by [syncthing](https://github.com/calmh/syncthing/blob/master/gui.go#L103)

## License

* CC0
