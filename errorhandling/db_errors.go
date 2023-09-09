package errorhandling

import "fmt"

func RequestBodyReadingError(err error) {
	if err != nil {
		fmt.Println("Error Reading Request Body:\n", err)
	}
}

func JsonMarshallingError(err error) {
	if err != nil {
		fmt.Println("Failed Marshalling:\n", err)
	}
}

func DbConnectionError(err error) {
	if err != nil {
		fmt.Println("Error Connecting to DB:\n", err)
	}
}

func DbInsertError(err error) {
	if err != nil {
		fmt.Println("Error Inserting Row in Table:\n", err)
	}
}

func DbQueryError(err error) {
	if err != nil {
		fmt.Println("Error Querying Rows:\n", err)
	}
}
