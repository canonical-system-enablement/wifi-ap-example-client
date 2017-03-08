/*
 * Copyright (C) 2017 Canonical Ltd
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License version 3 as
 * published by the Free Software Foundation.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	post := flag.String("d", "", "data to POST")
	help := flag.Bool("h", false, "usage help")
	flag.Parse()

	if *help || len(flag.Args()) != 1 {
		fmt.Fprintln(os.Stderr, os.Args[0], "is an example client application for the wifi-ap snap (https://docs.ubuntu.com/core/en/stacks/network/wifi-ap/docs/reference/rest-api).")
		fmt.Fprintln(os.Stderr, "	It communicates with the wifi-ap snap over the provided unix domain socket.")
		fmt.Fprintln(os.Stderr, "	For example you can call its configuration endpoint like:")
		fmt.Fprintln(os.Stderr, os.Args[0], "[-d data] /uri")
		flag.PrintDefaults()
		os.Exit(0)
	}

	fmt.Println("wifi-ap example client")

	httpc := http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
				path := filepath.Join(os.Getenv("SNAP_COMMON"), "/sockets/control")
				return net.Dial("unix", path)
			},
		},
	}

	var response *http.Response
	var err error
	if len(*post) == 0 {
		response, err = httpc.Get("http://unix" + flag.Args()[0])
	} else {
		response, err = httpc.Post("http://unix"+flag.Args()[0], "application/octet-stream", strings.NewReader(*post))
	}

	if err != nil {
		panic(err)
	}
	io.Copy(os.Stdout, response.Body)
}
