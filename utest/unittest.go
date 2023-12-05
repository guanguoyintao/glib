package utest

import (
	"context"
	"fmt"
)

func (f TestFunc) Testing(ctx context.Context, name string, wantErr bool) {
	fmt.Println("-----------------------------------")
	fmt.Printf("test %s start \n", name)
	err := f(ctx)
	if (err != nil) != wantErr {
		fmt.Printf("failed, err is %v\n", err)
	} else {
		fmt.Println("success")
	}
	fmt.Println("\n")
}
