package main

import (
	"flag"
	"log"
	"net"

	///"strconv"
	"encoding/binary"
	"strconv"

	"./idgenerator"
	"./workerid"
)

var (
	scheme    *string
	host      *string
	port      *uint
	worker_id *uint
	uri       string
)

func parseArgument() {
	scheme = flag.String("scheme", "tcp", "Communication scheme [tcp | http], default is [tcp]")
	host = flag.String("host", "127.0.0.1", "Listening host, default is [127.0.0.1]")
	port = flag.Uint("port", 11337, "Listening port, default 11337")
	worker_id = flag.Uint("workerId", uint(workerid.DetectWorkerId()), "The unique worker id")

	flag.Parse()
}

func main() {

	parseArgument()

	uri = *host + ":" + strconv.FormatUint(uint64(*port), 10)

	// Listen for incoming connections.
	listener, err := net.Listen(*scheme, uri)

	if err != nil {
		log.Fatalf(err.Error())
	}

	// Close the listener when the application closes.
	defer listener.Close()

	log.Printf("Listening: %s://%s\n", *scheme, uri)
	log.Printf("WorkerId: %d\n", *worker_id)

	id_generator, err := idgenerator.NewIdWorker(int32(*worker_id))

	for {

		// Listen for an incoming connection.
		connection, err := listener.Accept()

		if err != nil {
			log.Fatalln(err.Error())
		}

		// Handle connections in a new go-routine.
		go handleRequest(connection, id_generator)
	}
}

// Handles incoming requests.
func handleRequest(connection net.Conn, worker *idgenerator.IdWorker) {

	// number of ids to generate
	var num_ids uint16

	// read the number of requested ids
	err := binary.Read(connection, binary.LittleEndian, &num_ids)

	if err != nil {
		log.Println(err)
	} else {

		// generate ids
		new_ids, err := worker.NextIds(num_ids)

		if err != nil {
			log.Println(err)
		} else {
			// send ids back
			err = binary.Write(connection, binary.LittleEndian, new_ids)

			if err != nil {
				log.Println(err)
			}
		}
	}

	// Close the connection when you're done with it.
	connection.Close()
}
