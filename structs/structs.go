package structs

import "net/mail"

type Message struct {
	From    []*mail.Address
	To      []*mail.Address
	Subject string
}
