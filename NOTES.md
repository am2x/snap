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