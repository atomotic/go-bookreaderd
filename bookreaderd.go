package main

import (
	"bitbucket.org/tebeka/nrsc"
	"fmt"
	"github.com/codegangsta/martini"
	"io"
	"log"
	"mime"
	"net/http"
	"path/filepath"
)

func main() {
	router := martini.NewRouter()
	router.Get("/", func() string {
		return "ia-bookreaderd"
	})

	router.Get("/book/:name", func(params martini.Params) string {
		return index(params["name"])
	})

	mr := martini.New()
	mr.Use(nrscStatic("assets"))
	mr.Use(martini.Static("books"))
	mr.Action(router.Handle)

	fmt.Println("running on http://127.0.0.1:8088")
	http.ListenAndServe("127.0.0.1:8088", mr)
}

func nrscStatic(path string) interface{} {
	if err := nrsc.Initialize(); err != nil {
		panic("Unable to initialize nrsc: " + err.Error())
	}
	return func(res http.ResponseWriter, req *http.Request, log *log.Logger) {
		file := req.URL.Path

		// nrsc expects there not to be a leading slash
		if file[0] == '/' {
			file = file[1:]
		}

		f := nrsc.Get(file)
		if f == nil {
			return
		}

		rdr, err := f.Open()
		if err != nil {
			http.Error(res, "Internal Server Error", http.StatusInternalServerError)
		}
		defer rdr.Close()

		mtype := mime.TypeByExtension(filepath.Ext(req.URL.Path))
		if len(mtype) != 0 {
			res.Header().Set("Content-Type", mtype)
		}
		res.Header().Set("Content-Size", fmt.Sprintf("%d", f.Size()))
		res.Header().Set("Last-Modified", f.ModTime().UTC().Format(http.TimeFormat))

		io.Copy(res, rdr)
	}
}

func index(path string) string {
	html := `
<!DOCTYPE HTML PUBLIC "-//W3C//DTD HTML 4.01//EN" "http://www.w3.org/TR/html4/strict.dtd">
<html>
<head>
    <title>bookreader demo</title>
    <link rel="stylesheet" type="text/css" href="/BookReader/BookReader.css"/>
    <style>#BRtoolbar .embed, .print { display: none; }</style>
    <script type="text/javascript" src="/js/jquery-1.4.2.min.js"></script>
    <script type="text/javascript" src="/js/jquery-ui-1.8.5.custom.min.js"></script>
    <script type="text/javascript" src="/js/dragscrollable.js"></script>
    <script type="text/javascript" src="/js/jquery.colorbox-min.js"></script>
    <script type="text/javascript" src="/js/jquery.ui.ipad.js"></script>
    <script type="text/javascript" src="/js/jquery.bt.min.js"></script>
    <script type="text/javascript" src="/BookReader/BookReader.js"></script>
</head>
<body style="background-color: ##939598;">

<div id="BookReader">
    Internet Archive BookReader Demo    <br/>

    <noscript>
    <p>
        The BookReader requires JavaScript to be enabled. Please check that your browser supports JavaScript and that it is enabled in the browser settings.  You can also try one of the <a href="http://www.archive.org/details/goodytwoshoes00newyiala"> other formats of the book</a>.
    </p>
    </noscript>
</div>

<script>
    br = new BookReader();
    br.getPageWidth = function(index) {
        return 800;
    }
    br.getPageHeight = function(index) {
        return 1200;
    }
    br.getPageURI = function(index, reduce, rotate) {
        var leafStr = '000';
        var imgStr = (index+1).toString();
        var re = new RegExp("0{"+imgStr.length+"}$");
        var url = '/%s/page'+leafStr.replace(re, imgStr) + '.jpg';
        return url;
    }

    br.getPageSide = function(index) {
        if (0 == (index & 0x1)) {
            return 'R';
        } else {
            return 'L';
        }
    }

    br.getSpreadIndices = function(pindex) {
        var spreadIndices = [null, null];
        if ('rl' == this.pageProgression) {
            // Right to Left
            if (this.getPageSide(pindex) == 'R') {
                spreadIndices[1] = pindex;
                spreadIndices[0] = pindex + 1;
            } else {
                // Given index was LHS
                spreadIndices[0] = pindex;
                spreadIndices[1] = pindex - 1;
            }
        } else {
            // Left to right
            if (this.getPageSide(pindex) == 'L') {
                spreadIndices[0] = pindex;
                spreadIndices[1] = pindex + 1;
            } else {
                // Given index was RHS
                spreadIndices[1] = pindex;
                spreadIndices[0] = pindex - 1;
            }
        }

        return spreadIndices;
    }

    br.getPageNum = function(index) {
        return index+1;
    }

    br.numLeafs = 15;
    br.bookTitle= 'ia-bookreaderd';
    br.bookUrl  = 'http://openlibrary.org';
    br.imagesBaseURL = '/BookReader/images/';

    br.getEmbedCode = function(frameWidth, frameHeight, viewParams) {
        return "Embed code not supported in bookreader demo.";
    }

    br.init();

    $('#BRtoolbar').find('.read').hide();
    $('#textSrch').hide();
    $('#btnSrch').hide();
</script>

</body>
</html>

`
	return fmt.Sprintf(html, path)

}
