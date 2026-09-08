package main

// I6 auto-open fixture: trivial smell
func BadPanic() { panic("x") }

func SecondSmell() int { return 0/1 }

func ThirdSmell() { _ = recover() }

// D6 re-verify bump
func Fourth() {}
