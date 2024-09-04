## Goals
- Allow database agnostic backup & restore using csv/json.
- Work on both local and hosted databases.
- Work on docker and bare metal setups using login configs.
- Have loader/timer function ideally.

## Snap v1 Step by Step:
1. Select database to backup (postgres,mysql,mongodb).
2. Enter database host type (bare metal, container).
3. Enter connection details(host,port,username,password,timeout*)
4. Begin backup progress loader and save .csv/.json/.bson file locally(setup with s3 in future).
5. Enter local/remote database configs of recipient database.
6. Begin restoring database with progress loader.
7. Complete backup and offer the user the file to download locally.

## Docker Test Commands
-> Backup from container to instance directory
* docker exec -it your_postgres_container pg_dump -U your_postgres_user -d your_database > /path/to/your_exported_file.sql 

-> Copy dump file contents to local file from instance
* ssh -i /Users/mac/awesompeproject centos@myteam.host "cat mtbackup.sql" > /Users/mac/Desktop/mtcom.sql

-> Backup postgres container db using backup sql file
* docker exec -i 6612ff0b4505 psql -U awesome_project -d awesome_project < /Users/mac/Desktop/mtcom.sql 

-> Restore postgres database
* docker exec -i awesomeproject-psql-1 psql -U awesome_project -d awesome_project < /Users/mac/Desktop/awesome_project_mtai_backup.sql 

-> Copy file from server to local machine
* scp -i /Users/mac/bestbooks centos@54.145.225.243:/home/centos/bestbooks_23_02_2024.sql /Users/mac/Desktop/
