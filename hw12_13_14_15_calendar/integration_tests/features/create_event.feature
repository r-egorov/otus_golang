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

    Scenario: successfully create event for the same week
                When I send a "POST" request to "http://localhost:8080/events" with JSON-body:
                """
                {
                    "event": {
                        "title": "Last Tuesday before vacation",
                        "datetime": "2022-12-20T16:00:00Z",
                        "duration": 7200,
                        "description": "Enjoying the week before my vacation",
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
                    "title": "Last Tuesday before vacation",
                    "datetime": "2022-12-20T16:00:00Z",
                    "duration": 7200,
                    "description": "Enjoying the week before my vacation",
                    "ownerId": "12e354aa-4cfa-449b-9ceb-e5b34cb6d84a"
                }
                """

    Scenario: successfully create event for the same month
                When I send a "POST" request to "http://localhost:8080/events" with JSON-body:
                """
                {
                    "event": {
                        "title": "December Monday",
                        "datetime": "2022-12-05T16:00:00Z",
                        "duration": 7200,
                        "description": "Some Monday in December",
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
                    "title": "December Monday",
                    "datetime": "2022-12-05T16:00:00Z",
                    "duration": 7200,
                    "description": "Some Monday in December",
                    "ownerId": "12e354aa-4cfa-449b-9ceb-e5b34cb6d84a"
                }
                """

    Scenario: successfully create event for the same year
                When I send a "POST" request to "http://localhost:8080/events" with JSON-body:
                """
                {
                    "event": {
                        "title": "Some day in May",
                        "datetime": "2022-05-05T16:00:00Z",
                        "duration": 7200,
                        "description": "Some day in May",
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
                    "title": "Some day in May",
                    "datetime": "2022-05-05T16:00:00Z",
                    "duration": 7200,
                    "description": "Some day in May",
                    "ownerId": "12e354aa-4cfa-449b-9ceb-e5b34cb6d84a"
                }
                """

    Scenario: successfully create event for the different year
                When I send a "POST" request to "http://localhost:8080/events" with JSON-body:
                """
                {
                    "event": {
                        "title": "2023 FUTURE",
                        "datetime": "2023-01-07T16:00:00Z",
                        "duration": 7200,
                        "description": "Some day in future year",
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
                    "title": "2023 FUTURE",
                    "datetime": "2023-01-07T16:00:00Z",
                    "duration": 7200,
                    "description": "Some day in future year",
                    "ownerId": "12e354aa-4cfa-449b-9ceb-e5b34cb6d84a"
                }
                """