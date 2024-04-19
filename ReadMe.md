# File Synchronization System Server
Code for the backend of File Synchronization System. For the client, check [here](https://github.com/aashabtajwar/file-sync-client).


# Installation
1. Make sure to have Go in your system. You can install it from [here](https://go.dev/doc/install).
2. Clone this repository
```
git clone git@github.com:aashabtajwar/file-sync-client.git
```
3. Use `go mod tidy`. This will install necessary packages.
4. The server uses MySQL, so make sure you have it installed. Digital Ocean has a great [blog](https://www.digitalocean.com/community/tutorials/how-to-install-mysql-on-ubuntu-22-04) covering on how to install MySQL on Ubuntu 22.04. Alternatively, you can use [MariaDB](https://www.digitalocean.com/community/tutorials/how-to-install-mariadb-on-ubuntu-22-04)
5. Import SQL file `mysql -u <username> -p file-sync < filesync.sql`.
6. Use `go run main.go` or `go build main.go` with `./main` to start the server.

# Features
- File storage and sharing: Store files and share them with other users. This is done using workspaces.
- Real time file synchronization: Shared users will receive real-time updates whenever a file from source computer is updated. The changes are made right in to the shared users' computers. Therefore, they can view the updated versions without visiting and downloading from the browser.

# Further Notes
## Features to Expect in Future
The backend is incomplete and still under development. More features are underway, these include but are not limited to
- Min.io for file storage
- Giving permissions to shared users
- Conflict resolution
- Compress files before storing


## 