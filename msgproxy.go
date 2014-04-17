// Package msgproxy provides server and client implementations for simple message passing protocol.
// msgproxy is asymetrical, meaning that interface for sender and receiver are distinct, although device may
// use both interfaces simultaneously.
// In the future, it's possible to expose a more symmetric interface.
//
// Receivers use plain TCP and have to register before they can receive any messages.
// After registration, receiver can either keep an idle TPC connection open waiting for messages, or can periodically check
// for new messages, without having a persistent connection. In the latter case, new messages get buffered. If number of new messages
// exceeds buffer size, old messages get discarded. Buffered messages are stored in main memory, which means that server
// restart drops all pending messages.
//
// Senders are exposed to HTTP interface and can either send messages that expect response or one-way messages.
// One-way messages are not guaranteed to be delivered (if receiver uses polling model)
//
// As of now, there are no plans to restrict/filter messages between different senders/receivers. Anyone can
// subscribe for messages and anyone can send them.
package msgproxy
