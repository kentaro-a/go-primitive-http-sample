package main

import (
	"time"

	"github.com/jessevdk/go-assets"
)

var _Templatesa3f99aa5f064b2bcd929fff273f9c0ed8264db02 = "<html>\n<head>\n\t<title>Internal Server Error</title>\n</head>\n<body>\n\t<div>\n\t\t<h1>Internal Server Error</h1>\n\t\t<div>\n\t\t\t<a href=\"/\">Top</a>\n\t\t</div>\n\t</div>\n</body>\n</html>\n"
var _Templates80c7dcf53562ef69508c11cc5c7148a765b21a2b = "<html>\n<head>\n\t<title>Home</title>\n</head>\n<body>\n\t<div>\n\t\t<h1>My Home</h1>\n\t\t<div>\n\t\t\taaaaaaaaaaaaa<br>\n\t\t\tbbbbbbbbbbbbb<br>\n\t\t</div>\n\t</div>\n</body>\n</html>\n"
var _Templates7152281708687c52b03bb629e6be2dffb809500f = "<html>\n<head>\n\t<title>Page Not Found</title>\n</head>\n<body>\n\t<div>\n\t\t<h1>Page Not Found</h1>\n\t\t<div>\n\t\t\t<a href=\"/\">Top</a>\n\t\t</div>\n\t</div>\n</body>\n</html>\n"

// Templates returns go-assets FileSystem
var Templates = assets.NewFileSystem(map[string][]string{}, map[string]*assets.File{
	"404.html": &assets.File{
		Path:     "404.html",
		FileMode: 0x1a4,
		Mtime:    time.Unix(1568171765, 1568171765000000000),
		Data:     []byte(_Templates7152281708687c52b03bb629e6be2dffb809500f),
	}, "500.html": &assets.File{
		Path:     "500.html",
		FileMode: 0x1a4,
		Mtime:    time.Unix(1568171755, 1568171755000000000),
		Data:     []byte(_Templatesa3f99aa5f064b2bcd929fff273f9c0ed8264db02),
	}, "home.html": &assets.File{
		Path:     "home.html",
		FileMode: 0x1a4,
		Mtime:    time.Unix(1567654826, 1567654826000000000),
		Data:     []byte(_Templates80c7dcf53562ef69508c11cc5c7148a765b21a2b),
	}}, "")
