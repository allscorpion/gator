# gator

A tool for subscribing to RSS feeds and printing them out in the stdout. 

## Installation

- Requires Go
- Requries Postgres
- Use go install
- Setup a config file at ~/.gatorconfig.json containing your db_url in the following format
```json
{
  "db_url": "postgres://example"
}
```

## Commands

- go run . reset
    - resets all the data
- go run . register <username>
    - registers a user and logs them in
- go run . login <username>
    - login a different user
- go run . addfeed <name> <url>
    - add an rss feed to be looked up
- go run . agg <time>
    - aggregate all the rss feeds you have added for this user into the database every <time> (1s, 1m)
- go run . browse <limit>
    - check out the posts that have shown up in the rss feeds

## Commands For Generating the tables

- sqlc generate - Generates the go code from the SQL queries
- goose postgres <connection_string> up - run the up migration
- goose postgres <connection_string> down - run the down migration

## RSS feeds

- TechCrunch: https://techcrunch.com/feed/
- Hacker News: https://news.ycombinator.com/rss
- Boot.dev Blog: https://blog.boot.dev/index.xml

