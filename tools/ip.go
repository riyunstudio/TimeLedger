package tools

import "net"

// 檢查是否為合法的 IPv4 地址
func (tl *Tools) IsIPv4(ip string) bool {
	parsedIP := net.ParseIP(ip)
	return parsedIP != nil && parsedIP.To4() != nil
}
