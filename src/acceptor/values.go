package main

type contextValueKey int

const(
	contextKeyPublisher contextValueKey = iota
	contextKeyCancel
)