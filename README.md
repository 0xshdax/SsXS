# SsxS
Simple Scanner XSS and SSTI are inspired by [freq](https://github.com/takshal/freq) with added features and optimizations.

## Install

```
$ go build -o sxs main.go
$ sudo mv sxs /usr/local/bin/
```

## Usage

```
$ cat target.txt

[...]
http://testphp.vulnweb.com:80/comment.php?pid=FUZZ
http://testphp.vulnweb.com:80/disclaimer.php=FUZZ
http://testphp.vulnweb.com:80/guestbook.php%3Cscript%3Ealert(%22XSS%20ATTACK%20By%20BINUSHACKER%20TEAM%22)%3C/script%3E%3CMARQUEE%20BGCOLOR=FUZZ
http://testphp.vulnweb.com:80/hpp/index.php?pp=FUZZ
http://testphp.vulnweb.com/hpp/params.php?aaaa%2f=FUZZ&p=FUZZ
http://testphp.vulnweb.com/hpp/params.php?aaaa%2f=FUZZ&p=FUZZ&pp=FUZZ
http://testphp.vulnweb.com/hpp/params.php?p=FUZZ&pp=FUZZ&aaaa%2f=FUZZ
http://testphp.vulnweb.com/hpp/params.php?aaaa%2f=FUZZ&p
http://testphp.vulnweb.com/hpp/params.php?p=FUZZ&lt;script&gt;alert(1)&lt;/script&gt
http://testphp.vulnweb.com/hpp/params.php?p=FUZZ&lt;script&gt;alert(1)&lt;/script&gt;
http://testphp.vulnweb.com/hpp/params.php?p=FUZZ&lt;iframe
http://testphp.vulnweb.com/hpp/params.php?p=FUZZ
http://testphp.vulnweb.com/hpp/params.php?p=FUZZ&pp=FUZZ
http://testphp.vulnweb.com/hpp/?pp=FUZZ
http://testphp.vulnweb.com:80/hpp/?pp=FUZZ
http://testphp.vulnweb.com:80/http:/listproducts.php=FUZZ
http://testphp.vulnweb.com:80/index.php%20cat=FUZZ
[...]

$ cat target.txt | sxs
```
