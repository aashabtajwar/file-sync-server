package tcpserver

import (
	"bytes"
	"database/sql"
	"fmt"
	"net"

	"github.com/aashabtajwar/th-server/errorhandling"
)

func BroadCastToUsers(fileData *bytes.Buffer, conncetedUsers map[string]net.Conn, metadata map[string]string) {
	// query user_ids from db
	// stream file to only those connected users that are in shared workspace

	db, er := sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/filesync")
	defer db.Close()
	if er != nil {
		fmt.Println("Error Connecting to DB: \n", er)
	}

	var validConnections []net.Conn
	query := fmt.Sprintf("SELECT user_id FROM shared_workspace WHERE workspace_id='%s'", metadata["workspace_id"])
	rows, err := db.Query(query)

	errorhandling.DbQueryError(err)

	for rows.Next() {
		var userId string
		if err := rows.Scan(&userId); err != nil {
			fmt.Println("Db Row Scan Error:\n", err)
		}
		connection, exists := conncetedUsers[userId]
		if exists {
			validConnections = append(validConnections, connection)
		}

	}
	// stream file data and metadata to these connections
	for _, con := range validConnections {

	}
}
