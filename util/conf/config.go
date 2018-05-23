package conf

import (
	"bufio"
	"io"
	"os"
	"strings"
)

tye Config struct {
	kvs map[string]string
	strcet string
}


func (this *Config) initConfig(path string) {
}

