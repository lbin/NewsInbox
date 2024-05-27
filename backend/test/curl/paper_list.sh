curl --location 'http://127.0.0.1:8099/api/paper/list' \
--header 'Content-Type: application/json' \
--data '{
    "start_time": "2023-01-01T00:00:00Z",
    "end_time": "2024-01-01T00:00:00Z",
    "search_content": "AI Pin",
    "start": 0,
    "limit": 100
}'