package workerid

import (
	"encoding/binary"
	"errors"
	"net"

	"../idgenerator"
)

func DetectWorkerId() uint16 {
	var id uint16
	ip, err := privateIPv4()

	if err != nil {
		id = 1
	} else {
		id = uint16(ip2int(ip) & uint32(idgenerator.Max_worker_id))
	}

	return id
}

func privateIPv4() (net.IP, error) {
	as, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}

	for _, a := range as {
		ip_net, ok := a.(*net.IPNet)

		if !ok || ip_net.IP.IsLoopback() {
			continue
		}

		ip := ip_net.IP.To4()

		if isPrivateIPv4(ip) {
			return ip, nil
		}
	}

	return nil, errors.New("no private ip address")
}

func isPrivateIPv4(ip net.IP) bool {
	return ip != nil &&
		(ip[0] == 10 || ip[0] == 172 && (ip[1] >= 16 && ip[1] < 32) || ip[0] == 192 && ip[1] == 168)
}

func ip2int(ip net.IP) uint32 {
	if len(ip) == 16 {
		return binary.BigEndian.Uint32(ip[12:16])
	}

	return binary.BigEndian.Uint32(ip)
}
