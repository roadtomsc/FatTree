package main

import (
	"fmt"
	"os"

	yaml "gopkg.in/yaml.v2"
)

// Node is a physical server in fat tree topology
type Node struct {
	ID              string   `yaml:"id"`
	Cores           int      `yaml:"cores"`
	RAM             int      `yaml:"ram"`
	VNFSupport      bool     `yaml:"vnfSupport"`
	NotManagerNodes []string `yaml:"notManagerNodes"`
	Egress          bool     `yaml:"egress"`
	Ingress         bool     `yaml:"ingress"`
}

// Link is a physical link in fat tree topology
type Link struct {
	Source      string `yaml:"source"`
	Destination string `yaml:"destination"`
	Bandwidth   int    `yaml:"bandwidth"`
}

// Config aggrates links and nodes as a physical topology
type Config struct {
	Nodes []Node
	Links []Link
}

func main() {
	var nodes []Node
	var links []Link

	var k int

	if _, err := fmt.Scanf("%d", &k); err != nil {
		return
	}

	if k%2 == 1 {
		return
	}

	// globals
	pods := k
	cores := k * k / 4

	// per pods
	aggregations := k / 2
	edges := k / 2
	servers := k * k / 4

	fmt.Printf("Pods: %d\n", pods)
	fmt.Printf("Cores: %d\n", cores)
	fmt.Printf("Aggregation: %d\n", aggregations)
	fmt.Printf("Edges: %d\n", edges)
	fmt.Printf("Servers: %d\n", servers)

	fmt.Printf("Nodes: %d\n", cores+pods*(servers+edges+aggregations))

	for i := 0; i < cores; i++ {
		nodes = append(nodes, Node{
			ID:              fmt.Sprintf("core-switch-%d", i),
			Cores:           0,
			RAM:             0,
			VNFSupport:      true,
			NotManagerNodes: []string{},
			Egress:          true,
			Ingress:         true,
		})

		for j := 0; j < pods; j++ {
			links = append(links, Link{
				Source:      fmt.Sprintf("core-switch-%d", i),
				Destination: fmt.Sprintf("aggr-switch-%d-%d", j, i/(pods/2)),
				Bandwidth:   40 * 1000,
			})
		}
	}

	for i := 0; i < pods; i++ {
		for j := 0; j < aggregations; j++ {
			nodes = append(nodes, Node{
				ID:              fmt.Sprintf("aggr-switch-%d-%d", i, j),
				Cores:           0,
				RAM:             0,
				VNFSupport:      false,
				NotManagerNodes: []string{},
				Egress:          false,
				Ingress:         false,
			})
		}
		for j := 0; j < edges; j++ {
			nodes = append(nodes, Node{
				ID:              fmt.Sprintf("edge-switch-%d-%d", i, j),
				Cores:           0,
				RAM:             0,
				VNFSupport:      false,
				NotManagerNodes: []string{},
				Egress:          false,
				Ingress:         false,
			})
		}

		for j := 0; j < aggregations; j++ {
			for k := 0; k < edges; k++ {
				links = append(links, Link{
					Source:      fmt.Sprintf("aggr-switch-%d-%d", i, j),
					Destination: fmt.Sprintf("edge-switch-%d-%d", i, k),
					Bandwidth:   40 * 1000,
				})
			}
		}

		// servers in the other pods cannot manage these server
		// so create a list of them and set as not manager nodes
		var nmn []string
		for k := 0; k < pods; k++ {
			if k != i {
				for j := 0; j < servers; j++ {
					nmn = append(nmn, fmt.Sprintf("server-%d-%d", k, j))
				}
			}
		}

		for j := 0; j < servers; j++ {
			nodes = append(nodes, Node{
				ID:              fmt.Sprintf("server-%d-%d", i, j),
				Cores:           20,
				RAM:             100,
				VNFSupport:      true,
				NotManagerNodes: nmn,
				Egress:          false,
				Ingress:         false,
			})
			links = append(links, Link{
				Source:      fmt.Sprintf("edge-switch-%d-%d", i, j/(pods/2)),
				Destination: fmt.Sprintf("server-%d-%d", i, j),
				Bandwidth:   40 * 1000,
			})
		}
	}

	b, err := yaml.Marshal(Config{
		Nodes: nodes,
		Links: links,
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	f, err := os.Create("config.yaml")
	if err != nil {
		fmt.Println(err)
		return
	}
	if _, err := f.Write(b); err != nil {
		fmt.Println(err)
		return
	}
	if err := f.Close(); err != nil {
		return
	}
	fmt.Println(string(b))
}
