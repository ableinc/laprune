# laprune

Run delete queries on one or more databases in a given environment. This is useful for delete test account records for users.

## How to run

```bash
go build -o laprune main.go
touch db.json queries.sql
vim db.json # use db.example.json as reference
vim queries.sql # use config/queries.example.sql as reference
./laprune /path/to/db.json /path/to/queries.sql
```

## Setup on server

**Install cron**
```bash
sudo apt-get update
sudo apt-get install cron
```

**Start cron**
```bash
sudo systemctl start cron
```

**Update permissions for executable**
```bash
chmod +x /path/to/laprune
```

**Open crontab file for editing**
```bash
crontab -e
```

**Add cronjob**
```bash
* * * * * /path/to/laprune
```
Example of running executable everyday at 1am
```bash
0 1 * * * /path/to/laprune
```

**Verify cronjob**
```bash
crontab -l
```

## Notes

1. Name must be unique. Using a duplicate name will cause queries to be combined and unexpected behavior will occur.
2. The script requires you pass both db.json and queries.sql file paths. You can provide the directories to each and
the program will append the file names to the end and validate if the files exist.
