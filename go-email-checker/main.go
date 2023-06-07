package main

import (
	"bufio" // buffer package to parse whatever we parse in the terminal
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func main(){
	scanner:= bufio.NewScanner(os.Stdin)
	fmt.Printf("domain, hasMX, hasSPF, spfRecord, hasDMARC, dmarcRecord\n")
	/* 
		where,
	 	`domain`: be like example.com,
		`MX`: is a DNS(Domain Name System) 'mail exchange'(MX) record directs email to a mail server. 
				The `MX` record indicates how email messages should be routed in accordance with the Simple Mail Transfer Protocol (SMTP, the standard protocol for all email). 
				An MX record must always point to another domain.
		`SPF`: A 'Sender Policy Framework' (SPF) record is a type of "DNS TXT" record that lists all the servers authorized to send emails from a particular domain.
		`DMARC`: 'Domain-based Message Authentication, Reporting and Conformance (DMARC)' is an email authentication protocol. 
				It is designed to give email domain owners the ability to protect their domain from unauthorized use, commonly known as email spoofing.
				A DMARC record stores a domain's DMARC policy. DMARC records are stored in the Domain Name System (DNS) as DNS TXT records.
				A DNS TXT record can contain almost any text a domain administrator wants to associate with their domain. 
				One of the ways DNS TXT records are used is to store DMARC policies.
	*/
	for scanner.Scan(){
		checkDomain(scanner.Text()) // whatever we enter in the terminal will be send using this `checkDomain()` function's parameter
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error: could not read from input: %v\n", err)
	}
}

func checkDomain(domain string){
	var hasMX, hasSPF, hasDMARC bool
	var spfRecord, dmarcRecord string

	// MX records
	mxRecords, err := net.LookupMX(domain)

	if err != nil {
		log.Printf("Error: %v\n", err)
	}

	if len(mxRecords) > 0 {
		hasMX = true
	}

	
	// SPF records
	txtRecords, err := net.LookupTXT(domain)

	if err != nil {
		log.Printf("Error: %v\n", err)
	}
	// for loop
	for _, record := range txtRecords {
		if strings.HasPrefix(record, "v=spf1") {
			hasSPF = true
			spfRecord = record
			break
		}
	}

	
	// DMARC records
	dmarcRecords, err := net.LookupTXT("_dmarc." + domain)
	if err != nil {
		log.Printf("Error: %v\n", err)
	}
	// for loop
	for _, record := range dmarcRecords {
		if strings.HasPrefix(record, "v=DMARC1") {
			hasDMARC = true
			dmarcRecord = record 
			break
		}
	}

	fmt.Printf("%v, %v, %v, %v, %v, %v", domain, hasMX, hasSPF, spfRecord, hasDMARC, dmarcRecord)
}