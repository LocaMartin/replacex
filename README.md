<div align="right">
  <img src="https://img.shields.io/badge/Go-1.20-00ADD8?style=plastic&logo=go" alt="Go Version">
</div>

<h1 align="center">replacex</h1>

An enhanced version of `qsreplace`. **replacex** accepts URLs via stdin to quickly replace or append query string values. Literal payload mode deduplicates unique parameter combinations per host and path, while payload file mode runs every payload against every input URL. It also features text replacement utilities.

### Install

```
go install github.com/locamartin/replacex@latest
```

Or [download a release](https://github.com/locamartin/replacex/releases) and put it somewhere in your `$PATH`
(e.g. in /usr/local/bin).

<h3 align="center"><u>Usage</u></h3>

### Replace Query String Values

```yml
# Parameter value replaced (String)
-$ cat urls.txt | replacex newval
https://example.com/path?one=newval&two=newval
https://example.com/pathtwo?one=newval&two=newval
https://example.net/a/path?one=newval&two=newval
```

### Replace from wordlist file (auto detect file) 
```yml
# It auto detects an existing file and runs every payload against every input URL (Payload file)
cat file.txt | replacex payload.txt
https://example.com/path?one=1234&two=1234
https://example.com/pathtwo?one=1+1&two=1+1
https://example.net/a/path?one=dsf&two=dsf
```

### Append to Query String Values

```yml
# Parameter value append (String)
-$ cat urls.txt | replacex -a newval
https://example.com/path?one=1newval&two=2newval
https://example.com/pathtwo?one=1newval&two=2newval
https://example.net/a/path?one=1newval&two=2newval
```
```yml
# Parameter value append (Payload file)
-$ cat urls.txt | replacex -a payloads.txt
https://example.com/path?one=1newvafdfgl&two=2newval121
https://example.com/pathtwo?one=1newvasdfsl&two=2nesdfswval
https://example.net/a/path?one=1newvsfsal&two=2newggfval
```

### Remove Duplicate URL and Parameter Combinations
Omit the argument to `-a` (append) to only output each combination of URL and query string parameters once:

```yml
-$ cat urls.txt | replacex -a 
https://example.com/path?one=1&two=2
https://example.com/pathtwo?one=1&two=2
https://example.net/a/path?one=1&two=2
```

### Replace and append with one flag `-ra` (replace and append at the same time)

```yml
# Replace and append every payload from an auto detected payload file
cat file.txt | replacex -ra payload.txt
https://example.com/path?one=1234&two=1234
https://example.com/pathtwo?one=1+1&two=1+1
https://example.net/a/path?one=dsf&two=dsf
```
### Replace specific word without changing the actual content of the file `-rw` (replace word) `-ww` (with word)

```yml
cat file.txt | replacex -rw "google" -ww "facebook" > file2.txt
```

|file.txt|file2.txt|
|-|-|
|google.com|facebook.com|
|google.in|facebook.in|
|adbcgoogle.br|adbcfacebook.br|
|googleabcd.in|facebookabcd.in|
