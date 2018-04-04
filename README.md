# Colorized and Formatted JSON
This repository provides an implementation of colorizing and formatting JSON using Go. More details are in Assignment2.pdf.

## Introduction
This program includes two parts:
- Colorization rules:
  - Each group of JSON tokens should have a unique color
  
- Formatting rules:
  - { and } tokens always go on their own line
  - There is a space before and after each :
  - Pairs in curly-braces go on their own lines, and if their values happen to also have curly-brace expressions, then they should be indented further in
  
## How to run
```sh
$ go run a2.go input.json
$ go run a2.go input.json > json.html
```
