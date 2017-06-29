package lib

import (
	"flag"
	"fmt"
	"log"
	"strings"
)

// Input groups all input variables
type Input struct {
	name    string
	listen  string
	cluster string
	debug   bool

	leader    string
	followers []string
}

// String stringifies an Input reference
func (input *Input) String() string {
	return fmt.Sprintf(
		"name: %s listen: %s cluster: %s",
		input.name, input.listen, input.cluster)
}

// ParseArgs handles command line argument passing
func ParseArgs() (*Input, error) {
	in := &Input{}

	flag.StringVar(&in.name, "name", "badkv", "node name")
	flag.StringVar(&in.listen, "listen", "localhost:10001", "listening address and port")
	flag.StringVar(&in.cluster, "cluster", "", "tuple of cluster peers; first peer is leader")
	flag.BoolVar(&in.debug, "debug", false, "debug")
	flag.Parse()

	if in.name == "" {
		return nil, fmt.Errorf("Name is empty")
	}
	if in.listen == "" {
		return nil, fmt.Errorf("Listen is empty")
	}
	if strings.HasPrefix(in.listen, "http://") || strings.HasPrefix(in.listen, "https://") {
		return nil, fmt.Errorf("Listen address is a url")
	}
	if len(strings.Split(in.listen, ":")) != 2 {
		return nil, fmt.Errorf("Listen address is invalid: %s", in.listen)
	}

	if in.cluster != "" {
		// cluster defines node tuples separated by comma
		nodes := strings.Split(in.cluster, ",")
		if len(nodes) != 3 {
			log.Fatalf("Invalid 3-node cluster: %s", in.cluster)
		}

		// a node is equal separated name=addr
		leaderNodeParts := strings.Split(nodes[0], "=")
		if len(leaderNodeParts) != 2 {
			log.Fatalf("Invalid leader node: %s", nodes[0])
		}

		if leaderNodeParts[0] == in.name {
			log.Printf("We are the leader")
			in.followers = make([]string, 2)

			for i := 1; i < 3; i++ {
				nodeParts := strings.Split(nodes[i], "=")
				if len(nodeParts) != 2 {
					log.Fatalf("Invalid cluster node: %s", nodes[i])
				}

				log.Printf("Follower %s", nodeParts)
				in.followers[i-1] = nodeParts[1]
			}
		} else {
			log.Printf("We are a follower. Leader is %s", leaderNodeParts[0])
			in.leader = leaderNodeParts[1]
		}
	}

	return in, nil
}
