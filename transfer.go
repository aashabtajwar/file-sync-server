
func broadcastToClients(connections []net.Conn, bytes []bytes.Buffer, conn net.Conn) error {
	// gather the connections
	// then broadcast to all of them

	file, err := os.Stat(path)
	if err != nil {
		panic(err)
	}
	fmt.Println("File stats ", file)
	
	// loop over the connections
	for _, connection := range connections {
		binary.Write(connection, binary.LittleEndian, int64(file.Size()))
		writtenBytes, er := io.CopyN(connection, bytes, int64(file.Size()))
		if er != nil {
			log.Fatal(er)
		
		}
		fmt.Println("Written %d bytes over the network")
	}
}


func (s *Server) authorStream(conn net.Conn, connections []net.Conn) error {
	/*
	* conn is the author connection
	* rest of the members from the 'connections' slice represents the 
	* connected clients
	* first create a buffer
	* keep connection open to incoming data continuously
	* read the size first that is in binary
	* then read the incoming data and write it into buffer
	* once that is done, broadcast the data to the connected clients
	*/
	buffer := new(bytes.Buffer)
	
	// keep connection open
	for {
		var size int64 

		binary.Read(conn, binary.LittleEndian, &size)
		bytesA, err := io.CopyN(buffer, conn, size)

		if err != nil {
			panic("Error ", err)
		}
		
		// write file
		err := os.WriteFile("code.py", buffer.Bytes(), 0644)
		if err != nil {
			log.Fatal(err)
		}
	}

	// broadcast data
	transferError := broadcastToClients(connections []net.Conn, bytes []bytes.Buffer, conn net.Conn)
	if transferError != nil {
		return transferError
	}
	return nil

	
}