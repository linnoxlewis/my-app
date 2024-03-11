package main

import (
	"applicationDesignTest/cmd"
	"context"
)

func main() {
	ctx := context.Background()
	cmd.InitApp(ctx)
}
