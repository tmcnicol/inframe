# Experiment with live updates from Postgres
**Todo:** this is entirely blocking at the moment just to get it running but moving forward this will be in the background.

### Steps to run the script

1. Create a database named `test`

2. Create a users table

   ```bash
   psql -d test -f queries/create_users.sql
   ```

3. Run the script `go run main.go`

   - There are some errors at the moment but we can fix these.

4. In another window insert, update or delete records into the user table, see the notifications printed on the screen.

### Resources
- https://godoc.org/github.com/lib/pq
- http://coussej.github.io/2015/09/15/Listening-to-generic-JSON-notifications-from-PostgreSQL-in-Go/

