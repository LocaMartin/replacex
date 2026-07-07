package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/url"
	"os"
	"sort"
	"strings"
)

func main() {
	var appendMode bool
	flag.BoolVar(&appendMode, "a", false, "Append the value instead of replacing it")

	var replaceAppendMode bool
	flag.BoolVar(&replaceAppendMode, "ra", false, "Output replaced and appended values")

	var ignorePath bool
	flag.BoolVar(&ignorePath, "ignore-path", false, "Ignore the path when considering what constitutes a duplicate")

	var replaceWord string
	flag.StringVar(&replaceWord, "rw", "", "Word to replace in stdin")

	var withWord string
	flag.StringVar(&withWord, "ww", "", "Word to replace -rw with")

	flag.Parse()

	visitedFlags := map[string]bool{}
	flag.Visit(func(f *flag.Flag) {
		visitedFlags[f.Name] = true
	})

	if visitedFlags["rw"] || visitedFlags["ww"] {
		if !visitedFlags["rw"] || !visitedFlags["ww"] {
			fmt.Fprintln(os.Stderr, "-rw and -ww must be used together")
			os.Exit(1)
		}
		replaceWords(replaceWord, withWord)
		return
	}

	values := []string{flag.Arg(0)}
	payloadFromFile := false
	if arg := flag.Arg(0); arg != "" && fileExists(arg) {
		payloads, err := readPayloads(arg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to read payload file %s [%s]\n", arg, err)
			os.Exit(1)
		}
		values = payloads
		payloadFromFile = true
	}

	seen := make(map[string]bool)

	// read URLs on stdin, then replace the values in the query string
	// with some user-provided value
	sc := bufio.NewScanner(os.Stdin)
	sc.Buffer(make([]byte, 1024), 1024*1024*10)
	for sc.Scan() {
		u, err := url.Parse(sc.Text())
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse url %s [%s]\n", sc.Text(), err)
			continue
		}

		// Go's maps aren't ordered, but we want to use all the param names
		// as part of the key to output only unique requests. To do that, put
		// them into a slice and then sort it.
		pp := make([]string, 0)
		for p, _ := range u.Query() {
			pp = append(pp, p)
		}
		sort.Strings(pp)

		key := fmt.Sprintf("%s%s?%s", u.Hostname(), u.EscapedPath(), strings.Join(pp, "&"))
		if ignorePath {
			key = fmt.Sprintf("%s?%s", u.Hostname(), strings.Join(pp, "&"))
		}

		if !payloadFromFile {
			// Only output each host + path + params combination once.
			if _, exists := seen[key]; exists {
				continue
			}
			seen[key] = true
		}

		writeURLVariants(u, values, payloadFromFile, appendMode, replaceAppendMode)

	}
	if err := sc.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to read stdin [%s]\n", err)
		os.Exit(1)
	}

}

func replaceWords(old, new string) {
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		fmt.Println(strings.ReplaceAll(sc.Text(), old, new))
	}
	if err := sc.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to read stdin [%s]\n", err)
		os.Exit(1)
	}
}

func fileExists(name string) bool {
	info, err := os.Stat(name)
	return err == nil && !info.IsDir()
}

func readPayloads(name string) ([]string, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var payloads []string
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		payloads = append(payloads, sc.Text())
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}

	return payloads, nil
}

func writeURLVariants(u *url.URL, values []string, allValues bool, appendMode bool, replaceAppendMode bool) {
	if allValues {
		for _, value := range values {
			writeURLVariant(u, value, appendMode, replaceAppendMode)
		}
		return
	}

	value := ""
	if len(values) > 0 {
		value = values[0]
	}
	writeURLVariant(u, value, appendMode, replaceAppendMode)
}

func writeURLVariant(u *url.URL, value string, appendMode bool, replaceAppendMode bool) {
	if replaceAppendMode {
		replaced := *u
		replaceQueryValues(&replaced, value, false)
		fmt.Printf("%s\n", &replaced)

		appended := *u
		replaceQueryValues(&appended, value, true)
		fmt.Printf("%s\n", &appended)
		return
	}

	updated := *u
	replaceQueryValues(&updated, value, appendMode)
	fmt.Printf("%s\n", &updated)
}

func replaceQueryValues(u *url.URL, value string, appendMode bool) {
	qs := url.Values{}
	for param, vv := range u.Query() {
		if appendMode && len(vv) > 0 {
			qs.Set(param, vv[0]+value)
		} else {
			qs.Set(param, value)
		}
	}

	u.RawQuery = qs.Encode()
}
