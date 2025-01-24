package main

import (
	"fmt"
	"log"
)

const (
	BindUsername = ""
	BindPassword = ""
	FQDN         = ""
	BaseDN       = ""
	Filter       = ""
)

func main() {
	// TLS Connection
	//l, err := ConnectTLS()
	//if err != nil {
	//	log.Fatal(err)
	//}
	//defer l.Close()

	// Non-TLS Connection
	l, err := Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	// Anonymous Bind and Search
	//result, err := AnonymousBindAndSearch(l)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//result.Entries[0].Print()

	// Normal Bind and Search
	result, err := BindAndSearch(l)
	if err != nil {
		log.Fatal(err)
	}
	result.Entries[0].Print()
}

// Ldap Connection with TLS
//func ConnectTLS() (*ldap.Conn, error) {
//	// You can also use IP instead of FQDN
//	l, err := ldap.DialURL(fmt.Sprintf("ldaps://%s:636", FQDN)) // port 3268
//	if err != nil {
//		return nil, err
//	}
//
//	return l, nil
//}

// Ldap Connection without TLS
func Connect() (*ldap.Conn, error) {
	// You can also use IP instead of FQDN
	l, err := ldap.DialURL(fmt.Sprintf("ldap://%s:32", FQDN)) // port 389
	if err != nil {
		return nil, err
	}

	return l, nil
}

// Anonymous Bind and Search
//func AnonymousBindAndSearch(l *ldap.Conn) (*ldap.SearchResult, error) {
//	l.UnauthenticatedBind("")
//
//	anonReq := ldap.NewSearchRequest(
//		"",
//		ldap.ScopeBaseObject, // you can also use ldap.ScopeWholeSubtree
//		ldap.NeverDerefAliases,
//		0,
//		0,
//		false,
//		Filter,
//		[]string{},
//		nil,
//	)
//	result, err := l.Search(anonReq)
//	if err != nil {
//		return nil, fmt.Errorf("Anonymous Bind Search Error: %s", err)
//	}
//
//	if len(result.Entries) > 0 {
//		result.Entries[0].Print()
//		return result, nil
//	} else {
//		return nil, fmt.Errorf("Couldn't fetch anonymous bind search entries")
//	}
//}

// Normal Bind and Search
func BindAndSearch(l *ldap.Conn) (*ldap.SearchResult, error) {
	_ = l.Bind(BindUsername, BindPassword)

	searchReq := ldap.NewSearchRequest(
		BaseDN,
		ldap.ScopeBaseObject, // you can also use ldap.ScopeWholeSubtree | ScopeBaseObject
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		Filter,
		[]string{},
		nil,
	)
	result, err := l.Search(searchReq)
	if err != nil {
		return nil, fmt.Errorf("Search Error: %s", err)
	} else {
		log.Println("auth ok")
	}

	log.Printf("search %v", searchReq)

	if len(result.Entries) > 0 {
		return result, nil
	} else {
		return nil, fmt.Errorf("Couldn't fetch search entries")
	}

}
