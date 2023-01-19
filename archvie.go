package main

import (
    "fmt"
    "net/http"
    "io/ioutil"
    "strings"
    "os"
)

func main() {
    var domain string
    if len(os.Args) > 1 {
        domain = os.Args[1]
    } else {
        // Get domain from standard input
        input, err := ioutil.ReadAll(os.Stdin)
        if err != nil {
            fmt.Println(err)
            return
        }
        domain = strings.TrimSpace(string(input))
    }

    url := fmt.Sprintf("http://web.archive.org/cdx/search/cdx?url=*.%s/*&output=text&fl=original&collapse=urlkey", domain)

    response, err := http.Get(url)
    if err != nil {
        fmt.Println(err)
        return
    }
    defer response.Body.Close()

    body, err := ioutil.ReadAll(response.Body)
    if err != nil {
        fmt.Println(err)
        return
    }

    lines := strings.Split(string(body), "\n")
    for _, line := range lines {
        if line != "" {
            parts := strings.Split(line, " ")
            hostname := strings.TrimSuffix(strings.TrimPrefix(parts[0], "http://"), ".")
            hostname = strings.TrimSuffix(strings.TrimPrefix(hostname, "https://"), ".")
            hostname = strings.TrimSuffix(strings.TrimPrefix(hostname, "ftp://"), ".")
            hostname = strings.Split(hostname, "/")[0]
            hostname = strings.Split(hostname, "?")[0]
            hostname = strings.Split(hostname, ":")[0]
            hostname = strings.Replace(hostname,"web.archive.org", "", -1)
            fmt.Println(hostname)
        }
    }
}
