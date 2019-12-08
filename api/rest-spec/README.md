# Kiosk REST APIs

## API Specification

All successful API calls result in HTTP/1.1 200 OK status code.

#### Echo request. Appropriate for health checking of kiosk service.

|Method        |Path                                           |Headers                        |
|---           |---                                            |---                            |
|POST          |/v1/echo                                       |Content-Type: application/json |

###### Request Body
```json
{
    "content": "Hi"
}
```

###### Response Body
```json
{
    "content": "Hi"
}
```

#### Create ticket request. Creates a new ticket with provided values.

|Method        |Path                                           |Headers                        |
|---           |---                                            |---                            |
|POST          |/v1/tickets                                    |Content-Type: application/json |

###### Request Body
```json
{
    "issuer": "Ticketing System",
    "owner": "user@example.com",
    "subject": "Technical Issue!",
    "content": "Hello, i have some technical issue with your API documentation. Please help!",
    "metadata": "{\"ip\": \"185.186.187.188\"}",
    "ticket_importance_level": "LOW",
    "ticket_status": "NEW"
}
```

###### Notes
- issuer: The name of your company or name of microservice who is creating the ticket.
- owner: Who is this ticket for?
- ticket_importance_level: Can be LOW, MEDIUM, HIGH or CRITICAL.
- ticket_status: Must be NEW.

###### Response Body
```json
{}
```

#### Read ticket request. Returns back a ticket and all associated comments by using its id.

|Method        |Path                                           |
|---           |---                                            |
|GET           |/v1/tickets/{id}                               |

###### Response Body
```json
{
    "id": "1",
    "issuer": "Ticketing System",
    "owner": "user@example.com",
    "subject": "Technical Issue!",
    "content": "Hello, i have some technical issue with your API documentation. Please help!",
    "metadata": "{\"ip\": \"185.186.187.188\"}",
    "ticket_importance_level": "LOW",
    "ticket_status": "NEW",
    "comments": [
        {
            "id": "1",
            "ticket_id": "1",
            "owner": "user@example.com",
            "content": "Hello, i have some technical issue with your API documentation. Please help!",
            "metadata": "{\"ip\": \"185.186.187.188\"}",
            "created_at": "2019-12-08T10:16:41.635862Z",
            "updated_at": "2019-12-08T10:16:41.635862Z"
        }
    ],
    "issued_at": "2019-12-08T10:15:28.726382Z",
    "updated_at": "2019-12-08T10:16:41.635862Z"
}
```

#### Update ticket request. Updates an already exists ticket with provided values.

|Method        |Path                                           |Headers                        |
|---           |---                                            |---                            |
|PUT           |/v1/tickets                                    |Content-Type: application/json |

###### Request Body
```json
{
    "id": "1",
    "ticket_status": "CLOSED"
}
```

###### Notes
- ticket_status: Can be REPLIED, RESOLVED, CLOSED or BLOCKED.

###### Response Body
```json
{}
```

#### Delete ticket request. Deletes ticket and all associated comments by using its id.

|Method        |Path                                           |
|---           |---                                            |
|DELETE        |/v1/tickets/{id}                               |

###### Response Body
```json
{}
```

#### Filter tickets request.

|Method        |Path                                           |
|---           |---                                            |
|GET           |/v1/tickets?page_number=1&page_size=10         |

###### Notes
- Other posibble query strings: issuer, owner, ticket_importance_level, ticket_status, from_date, to_data.

###### Response Body
```json
{
    "tickets": [
        {
            "id": "1",
            "issuer": "Ticketing System",
            "owner": "user@example.com",
            "subject": "Technical Issue!",
            "content": "Hello, i have some technical issue with your API documentation. Please help!",
            "metadata": "{\"ip\": \"185.186.187.188\"}",
            "ticket_importance_level": "LOW",
            "ticket_status": "CLOSED",
            "comments": [
                {
                    "id": "1",
                    "ticket_id": "1",
                    "owner": "user@example.com",
                    "content": "Hello, i have some technical issue with your API documentation. Please help!",
                    "metadata": "{\"ip\": \"185.186.187.188\"}",
                    "created_at": "2019-12-08T10:16:41.635862Z",
                    "updated_at": "2019-12-08T10:16:41.635862Z"
                }
            ],
            "issued_at": "2019-12-08T10:15:28.726382Z",
            "updated_at": "2019-12-08T10:18:16.591607Z"
        },
        {
            "id": "3",
            "issuer": "Ticketing System",
            "owner": "user@example.com",
            "subject": "Technical Issue!",
            "content": "Hello, i have some technical issue with your API documentation. Please help!",
            "metadata": "{\"ip\": \"185.186.187.188\"}",
            "ticket_importance_level": "HIGH",
            "ticket_status": "NEW",
            "comments": [],
            "issued_at": "2019-12-08T10:15:36.509418Z",
            "updated_at": "2019-12-08T10:15:36.509418Z"
        },
        {
            "id": "2",
            "issuer": "Ticketing System",
            "owner": "user@example.com",
            "subject": "Technical Issue!",
            "content": "Hello, i have some technical issue with your API documentation. Please help!",
            "metadata": "{\"ip\": \"185.186.187.188\"}",
            "ticket_importance_level": "MEDIUM",
            "ticket_status": "NEW",
            "comments": [],
            "issued_at": "2019-12-08T10:15:33.527792Z",
            "updated_at": "2019-12-08T10:15:33.527792Z"
        }
    ],
    "page": {
        "number": 1,
        "size": 10,
        "has_next": false
    }
}
```

## Error Handling

As you may know 4xx and 5xx statuses indicate an error and the response body is as follow:

```json
    {
        "errors":[
            {
                "code":"create_ticket.empty_issuer"
            }
        ]
    }
```

The code field and its possible values are described on the following list:

- create_ticket.empty_issuer

- create_ticket.empty_owner

- create_ticket.empty_subject

- create_ticket.empty_content

- create_ticket.invalid_ticket_importance_level

- create_ticket.invalid_ticket_status

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
