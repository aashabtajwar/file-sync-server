
func main() {
	port := os.Args[1]
	conn, err := net.Dial("tcp", ":3030")
	if err != nil {
		fmt.Println(err)
	}
	go func() {
		buf := new(bytes.Buffer)
		for {
			var size int64
			binary.Read(conn, binary.LittleEndian, &size)
			n, err := io.CopyN(buf, conn, size)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(buf.Bytes())
			fmt.Printf("received %d bytes\n", n)
		}
	}()
	start(port)
}
