package tcpserver

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"time"
)

func BroadCastToUsers(fileData *bytes.Buffer, conncetedUsers map[string]net.Conn, metadata map[string]string, thisConn net.Conn, dataString string, connections []net.Conn) {
	// query user_ids from db
	// stream file to only those connected users that are in shared workspace
	// fmt.Println("Now broadcasting data to users\n", string(fileData.String()))

	// fmt.Println("received file data\n", fileData.Bytes())

	// fmt.Println("\n Received file metadata\n", metadata)

	// db, er := sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/filesync")
	// defer db.Close()
	// if er != nil {
	// 	fmt.Println("Error Connecting to DB: \n", er)
	// }

	// var validConnections []net.Conn
	// query := fmt.Sprintf("SELECT user_id FROM shared_workspace WHERE workspace_id='%s'", metadata["workspace_id"])
	// rows, err := db.Query(query)

	// errorhandling.DbQueryError(err)

	// for rows.Next() {
	// 	var userId string
	// 	if err := rows.Scan(&userId); err != nil {
	// 		fmt.Println("Db Row Scan Error:\n", err)
	// 	}
	// 	connection, exists := conncetedUsers[userId]
	// 	if exists {
	// 		validConnections = append(validConnections, connection)
	// 	}

	// }
	// stream file data and metadata to connected connections
	// except for the connection that sent them
	fmt.Println(connections)
	for _, con := range connections {
		if con != thisConn {
			// up to this point
			// file and file size
			// metadata and metadata size
			fmt.Println("COMING HERE")
			binary.Write(con, binary.LittleEndian, int64(len(fileData.Bytes())))
			_, err := io.CopyN(con, bytes.NewReader(fileData.Bytes()), int64(len(fileData.Bytes())))
			if err != nil {
				fmt.Println("Error broadcasting data:\n", err)
			}
			// fmt.Printf("Written %d\n", n)

			time.Sleep(100 * time.Millisecond)

			mData := []byte(dataString)
			binary.Write(con, binary.LittleEndian, int64(len(mData)))
			_, err = io.CopyN(con, bytes.NewReader(mData), int64(len(mData)))
			if err != nil {
				fmt.Println("Error broadcasting data:\n", err)
			}
			// fmt.Printf("Written %d\n", n)

		}
	}
}
