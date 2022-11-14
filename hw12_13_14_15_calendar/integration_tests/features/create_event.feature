Feature: create event
    In order to store events
    As an API user
    I need to create events

    Scenario: successfully create event
        When I send a "POST" request to "http://localhost:8080/events" with JSON-body:
        """
        {
            "event": {
                "title": "Christmas party",
                "datetime": "2022-12-25T16:00:00Z",
                "duration": 7200,
                "description": "Seeing Santa and all that",
                "ownerId": "78e624aa-4cfa-449b-9ceb-e2b12cb8d48a"
            }
        }
        """
        Then the response code is 201
        And the response contains an event
        And the event has an ID
        And other fields equal to following:
        """
        {
            "title": "Christmas party",
            "datetime": "2022-12-25T16:00:00Z",
            "duration": 7200,
            "description": "Seeing Santa and all that",
            "ownerId": "78e624aa-4cfa-449b-9ceb-e2b12cb8d48a"
        }
        """

    Scenario: time is busy for that user
            When I send a "POST" request to "http://localhost:8080/events" with JSON-body:
            """
            {
                "event": {
                    "title": "Another christmas party",
                    "datetime": "2022-12-25T16:00:00Z",
                    "duration": 7200,
                    "description": "Can't make it there",
                    "ownerId": "78e624aa-4cfa-449b-9ceb-e2b12cb8d48a"
                }
            }
            """
            Then the response code is 400

    Scenario: successfully create event with same date for different user
            When I send a "POST" request to "http://localhost:8080/events" with JSON-body:
            """
            {
                "event": {
                    "title": "Christmas party",
                    "datetime": "2022-12-25T16:00:00Z",
                    "duration": 7200,
                    "description": "Seeing Santa and all that",
                    "ownerId": "12e354aa-4cfa-449b-9ceb-e5b34cb6d84a"
                }
            }
            """
            Then the response code is 201
            And the response contains an event
            And the event has an ID
            And other fields equal to following:
            """
            {
                "title": "Christmas party",
                "datetime": "2022-12-25T16:00:00Z",
                "duration": 7200,
                "description": "Seeing Santa and all that",
                "ownerId": "12e354aa-4cfa-449b-9ceb-e5b34cb6d84a"
            }
            """

    Scenario: successfully create event with same date for different user
                When I send a "POST" request to "http://localhost:8080/events" with JSON-body:
                """
                {
                    "event": {
                        "title": "Christmas party",
                        "datetime": "2022-12-25T16:00:00Z",
                        "duration": 7200,
                        "description": "Seeing Santa and all that",
                        "ownerId": "12e354aa-4cfa-449b-9ceb-e5b34cb6d84a"
                    }
                }
                """
                Then the response code is 201
                And the response contains an event
                And the event has an ID
                And other fields equal to following:
                """
                {
                    "title": "Christmas party",
                    "datetime": "2022-12-25T16:00:00Z",
                    "duration": 7200,
                    "description": "Seeing Santa and all that",
                    "ownerId": "12e354aa-4cfa-449b-9ceb-e5b34cb6d84a"
                }
                """