package main

import (
	"fmt"
	"testing"
)

func TestYAML(t *testing.T) {
	fmt.Println(string(mainYAML()))
}

func TestYAMLSpecs(t *testing.T) {
	fmt.Println(string(specsYAML("test-host-elb")))
}


func TestYAMLJava(t *testing.T) {
	fmt.Println(string(javaYAML("test-host-elb")))
}
