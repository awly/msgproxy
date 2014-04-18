// Package msgproxy provides means for two devices to devices to communicate behind NAT
// This is an evolving package that will be changing a lot over time, don't use it!
//
// First goal: provide a TCP bridge for two clients.
// One client connects, sends bridge name and waits for the other side.
// Other client connects, sends bridge name and gets connected to first client.
// After that all communications are piped between them by the server.
// If any side closes connection, the other connection is closed as well.
package msgproxy
