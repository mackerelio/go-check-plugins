package checkdns

import (
	"fmt"
	"net"
	"strconv"
)

func Do() {
	nameserver, err := adapterAddress()
	if err != nil {
		fmt.Println(err)
	}
	nameserver = net.JoinHostPort(nameserver, strconv.Itoa(53))
	fmt.Println(nameserver)
}
