package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    "os"
    "strings"
)

type Certificate struct {
    CommonName    string `json:"common_name"`
    NameValue string `json:"name_value"`
}

func main() {
    // Check if input is coming from a pipe or file
    fi, _ := os.Stdin.Stat()
    if (fi.Mode() & os.ModeCharDevice) == 0 {
        // input is coming from a pipe or file
        bytes, _ := ioutil.ReadAll(os.Stdin)
        input := string(bytes)

        // Iterate over each line of input
        for _, domain := range strings.Split(input, "\n") {
            // Skip empty lines
            if domain == "" {
                continue
            }

            url := fmt.Sprintf("https://crt.sh/?q=%%25.%s&output=json", domain)

            response, err := http.Get(url)
            if err != nil {
                fmt.Printf("The HTTP request for %s failed with error %s\n", domain, err)
            } else {
                data, _ := ioutil.ReadAll(response.Body)
                var certificates []Certificate
                json.Unmarshal(data, &certificates)
                for _, certificate := range certificates {
                    fmt.Println(certificate.CommonName)
                    fmt.Println(certificate.NameValue)
                }
            }
        }
    } else {
        // input is coming from command line
        if len(os.Args) < 2 {
            fmt.Println("Please provide a domain name")
            return
        }

        domain := os.Args[1]
        url := fmt.Sprintf("https://crt.sh/?q=%%25.%s&output=json", domain)

        response, err := http.Get(url)
        if err != nil {
            fmt.Printf("The HTTP request for %s failed with error %s\n", domain, err)
        } else {
            data, _ := ioutil.ReadAll(response.Body)
            var certificates []Certificate
            json.Unmarshal(data, &certificates)
            for _, certificate := range certificates {
                fmt.Println(certificate.CommonName)
                fmt.Println(certificate.NameValue)
            }
        }
    }
}
