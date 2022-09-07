package ioutils

import "encoding/xml"

func GetPrettyPrintOfHtmlBody(body []byte) string {
	html := body
	type node struct {
		Attr     []xml.Attr
		XMLName  xml.Name
		Children []node `xml:",any"`
		Text     string `xml:",chardata"`
	}
	x := node{}
	_ = xml.Unmarshal([]byte(html), &x)
	buf, _ := xml.MarshalIndent(x, "", "   ")
	return string(buf)
}
