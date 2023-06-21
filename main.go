package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"golang.org/x/sync/errgroup"
)

func main() {
	amount := 100
	if len(os.Args) >= 2 {
		a, err := strconv.Atoi(os.Args[1])
		if err != nil {
			fmt.Printf("Argument is not a number. Usage: %s [AMOUNT]\n", os.Args[0])
			os.Exit(2)
		}
		amount = a
	}

	if err := run(amount); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}

func run(amount int) error {
	eg := new(errgroup.Group)

	for i := 0; i < amount; i++ {
		eg.Go(func() error {
			resp, err := http.Get("https://oskar.staging.openslides.com")
			if err != nil {
				return fmt.Errorf("sending request: %w", err)
			}
			defer resp.Body.Close()

			if _, err := io.ReadAll(resp.Body); err != nil {
				return fmt.Errorf("read body: %w", err)
			}

			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return err
	}

	fmt.Println("All good. No rate limit")
	return nil
}
