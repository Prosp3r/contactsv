package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"net"
	"os"
	"regexp"
	"strings"
)

type Contact struct {
	ID     int64  `json:"id"`
	Email  string `json:"email"`
	Domain string `json:"domain"`
	Valid  string `json:"Valid"`
}

//files array of diff file types
//new_prospects002.txt => comma separated,
//.csv new line separated

var Contacts []Contact
var Domains map[string]string

//
var MContacts = make(map[string]Contact)

func main() {
	fileSource := "contacts/Oaken_blast_database2.csv"
	readFile(fileSource)
}

//Read contact file
func readFile(fileSource string) {
	dataFile, _ := os.Open(fileSource)
	reader := csv.NewReader(bufio.NewReader(dataFile))
	n := int64(0)
	for {
		newcontact := Contact{}
		line, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Printf("Could not read file source: %v system returned error %v \n", fileSource, err.Error())
		}

		if len(line[0]) > 0 {

			email := trimString(line[0])
			emailCount := strings.Split(email, ",")
			if len(emailCount) > 1 {
				emailAdd := ""

				// fmt.Println(email)
				// fmt.Println(emailCount[0])
				// fmt.Println(emailCount[1])

				//correct comma ,com in place of dot .com error
				if len(emailCount[1]) > 0 && len(emailCount[1]) <= 4 {
					emailAdd = emailCount[0] + "." + emailCount[1]
					fmt.Println(trimString(emailAdd))
					newcontact = Contact{
						n,
						emailAdd,
						getDomain(emailAdd),
						"DOUBLE",
					}
					Contacts = append(Contacts, newcontact)
					MContacts["email"] = newcontact
					n++
				}

				//solve for malformed email prefix
				prefixSplit := strings.Split(emailCount[0], "@")
				if len(prefixSplit) < 2 {
					fmt.Println("===>", prefixSplit)
					emailAdd = strings.Trim(prefixSplit[0], ",") + "." + strings.Trim(emailCount[1], ",")
					fmt.Println("===><===>", emailAdd)
					newcontact = Contact{
						n,
						emailAdd,
						getDomain(emailAdd),
						"DOUBLE",
					}
					Contacts = append(Contacts, newcontact)
					MContacts["email"] = newcontact
					n++
				}

				for i := 0; i < len(emailCount); i++ {
					newcontact = Contact{
						n,
						emailCount[i],
						getDomain(emailCount[i]),
						"DOUBLE",
					}
					Contacts = append(Contacts, newcontact)
					MContacts["email"] = newcontact
					n++
				}
			}

			//
			newcontact = Contact{
				n,
				email,
				getDomain(email),
				"NULL",
			}
			Contacts = append(Contacts, newcontact)
			MContacts["email"] = newcontact
			n++
		}
	}
	// for i, c := range Contacts {
	// 	fmt.Printf("Count: %v - ID: %v : %v : %v : %v \n\n", i+1, c.ID, c.Email, c.Domain, c.Valid)
	// }
	fmt.Println("Total ", n)
}

func trimString(em string) string {
	em = strings.TrimSpace(em)
	em = strings.Trim(em, "\t \n")
	em = strings.TrimLeft(em, "\t \n")
	em = strings.TrimRight(em, "\t \n")
	em = strings.Trim(em, "\n")

	return em
}

func getDomain(email string) string {
	em_frag := strings.Split(email, "@")
	if len(em_frag) > 1 && len(em_frag[1]) > 4 {
		return em_frag[1]
	}
	return ""
}

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func isEmailValid(e string) bool {
	if len(e) < 3 && len(e) > 254 {
		return false
	}
	if !emailRegex.MatchString(e) {
		return false
	}
	parts := strings.Split(e, "@")
	// mx, err := net.LookupMX(parts[1])
	mx, err := net.LookupMX(parts[1])
	if err != nil || len(mx) == 0 {
		return false
	}
	return true
}

//sanitze contacts

//Check for duplicates

//Check domain name availability

//Save good contacts

//Publish records
