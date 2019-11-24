# Kiosk REST APIs

## Error Handling

As you may know 4xx and 5xx statuses indicate an error and the response body is as follow:

```json
    {
        "errors":[
            {
                "code":"create_ticket.empty_issuer",
                "message": "Issuer field could not be empty."
            }
        ]
    }
```

The code field and its possible values are described on the following list:

- create_ticket.empty_issuer

- create_ticket.empty_owner

- create_ticket.empty_subject

- create_ticket.empty_content

- create_ticket.invalid_status

- create_ticket.failed

--

- read_ticket.invalid_id

- read_ticket.not_found

- read_ticket.failed

--

- update_ticket.invalid_id

- update_ticket.invalid_ticket_status

- update_ticket.not_found

- update_ticket.failed

--

- delete_ticket.invalid_id

- delete_ticket.failed

--

- filter_tickets.invalid_page_number

- filter_tickets.invalid_page_size

- filter_tickets.failed

--

- create_comment.empty_owner

- create_comment.empty_content

- create_comment.ticket_not_exists

- create_comment.failed

--

- update_comment.failed

- update_comment.ticket_not_found

--

- delete_comment.invalid_id

- delete_comment.failed
