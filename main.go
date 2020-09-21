package main

func main() {
	a := App{}
	a.Initialize(
		"postgres",
		"postgres",
		"products",
	)
	a.Run(":8001")
}
