package file_management

import (
	"fmt"
	"log"
	"os"
)

func Save_Chunk(data []byte, filename string, chunk_part int, username string) {
	dir := fmt.Sprintf("chunks/%s/%s", username, filename)
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	file_name := fmt.Sprintf("chunk_%d_%s", chunk_part, filename)
	filepath := dir + "/" + file_name

	// Write the bytes to file
	err = os.WriteFile(filepath, data, 0644) // 0644 = read/write for owner, read for others
	if err != nil {
		log.Fatal(err)
		return
	}

	return
}

func Get_Chunk(name string, filename string, chunknum int) []byte {
	filepath := fmt.Sprintf("./chunks/%s/%s/chunk_%d_%s", name, filename, chunknum, filename)
	data, err := os.ReadFile(filepath)
	if err != nil {
		log.Println(err)
		return nil
	}
	return data
}
