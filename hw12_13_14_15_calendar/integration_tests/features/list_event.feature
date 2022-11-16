Feature: list event
    In order to list events
    As an API user
    I need to get list of events from server

    Scenario: get list of events for the day
        When I send a "GET" request to "http://calendar:8080/events?datetime=2022-12-25T00:00:00Z&period=day"
        Then the response code is 200
        And the response contains list of events
        And the events have IDs
        And events fields are as following:
        """
        [
            {
                "title": "Christmas party",
                "datetime": "2022-12-25T16:00:00Z",
                "duration": 7200,
                "description": "Seeing Santa and all that",
                "ownerId": "12e354aa-4cfa-449b-9ceb-e5b34cb6d84a"
            },
            {
                "title": "Christmas party",
                "datetime": "2022-12-25T16:00:00Z",
                "duration": 7200,
                "description": "Seeing Santa and all that",
                "ownerId": "78e624aa-4cfa-449b-9ceb-e2b12cb8d48a"
            }
        ]
        """

    Scenario: get list of events for the week
        When I send a "GET" request to "http://calendar:8080/events?datetime=2022-12-19T00:00:00Z&period=week"
        Then the response code is 200
        And the response contains list of events
        And the events have IDs
        And events fields are as following:
        """
        [
            {
                "title": "Christmas party",
                "datetime": "2022-12-25T16:00:00Z",
                "duration": 7200,
                "description": "Seeing Santa and all that",
                "ownerId": "12e354aa-4cfa-449b-9ceb-e5b34cb6d84a"
            },
            {
                "title": "Christmas party",
                "datetime": "2022-12-25T16:00:00Z",
                "duration": 7200,
                "description": "Seeing Santa and all that",
                "ownerId": "78e624aa-4cfa-449b-9ceb-e2b12cb8d48a"
            },
            {
                "title": "Last Tuesday before vacation",
                "datetime": "2022-12-20T16:00:00Z",
                "duration": 7200,
                "description": "Enjoying the week before my vacation",
                "ownerId": "12e354aa-4cfa-449b-9ceb-e5b34cb6d84a"
            }
        ]
        """


    Scenario: get list of events for the month
        When I send a "GET" request to "http://calendar:8080/events?datetime=2022-12-01T00:00:00Z&period=month"
        Then the response code is 200
        And the response contains list of events
        And the events have IDs
        And events fields are as following:
        """
        [
            {
                "title": "Christmas party",
                "datetime": "2022-12-25T16:00:00Z",
                "duration": 7200,
                "description": "Seeing Santa and all that",
                "ownerId": "12e354aa-4cfa-449b-9ceb-e5b34cb6d84a"
            },
            {
                "title": "Christmas party",
                "datetime": "2022-12-25T16:00:00Z",
                "duration": 7200,
                "description": "Seeing Santa and all that",
                "ownerId": "78e624aa-4cfa-449b-9ceb-e2b12cb8d48a"
            },
            {
                "title": "Last Tuesday before vacation",
                "datetime": "2022-12-20T16:00:00Z",
                "duration": 7200,
                "description": "Enjoying the week before my vacation",
                "ownerId": "12e354aa-4cfa-449b-9ceb-e5b34cb6d84a"
            },
            {
                "title": "December Monday",
                "datetime": "2022-12-05T16:00:00Z",
                "duration": 7200,
                "description": "Some Monday in December",
                "ownerId": "12e354aa-4cfa-449b-9ceb-e5b34cb6d84a"
            }
        ]
        """


    Scenario: invalid period
        When I send a "GET" request to "http://calendar:8080/events?datetime=2022-12-01T00:00:01Z&period=year"
        Then the response code is 400

    Scenario: invalid route
            When I send a "GET" request to "http://calendar:8080/whassup"
            Then the response code is 404

    Scenario: invalid method
            When I send a "PUT" request to "http://calendar:8080/events"
            Then the response code is 405