package ipset

import "fmt"

type Ipset struct{}

func (i *Ipset) Test() {
	fmt.Println("Ipset.Test() called")
}
