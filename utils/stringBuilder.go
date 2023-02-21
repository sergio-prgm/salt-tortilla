package utils

import (
	"fmt"
)

func URLString(textInput func() string) string {
	return fmt.Sprintf("Input the URL:\n\n%s", textInput())
}

// Not yet fully implemented
func HttpVerbString(url string) string {
	return fmt.Sprintf("Input the HTTP Verb:\n\nURL: %s\n\n", url)
}

func HeadersString(url, httpVerb string, headers []string, textInput string) string {

	head := fmt.Sprintf(
		"Input the Headers:\n\nURL: %s\nHTTP Verb: %s\nHeaders:%s\n\n%s",
		url,
		httpVerb,
		PrintSlice(headers),
		textInput,
	)
	bottom := "\n\n(press tab to input new Header, press enter to input Body)"
	return head + bottom
}
