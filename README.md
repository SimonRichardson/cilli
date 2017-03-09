# cilli

			ELI5: A path compiler

-----

A reasonable implementation for a path DSL

![](http://static.bbc.co.uk/history/img/ic/640/images/resources/histories/code_breaking.jpg)

-----

### Introduction

Cilli is a library for compiling a path query DSL for selecting items found
within a certain corpus. Cilli uses a Pratt parser with a greedy lexer to
generate a simple AST of expressions. The lexer is greedy by generating tokens
of full types instead of token characters (it does it in one loop instead of
possible multiple steps). The parser uses the idea of both prefix and infix
(which can be viewed as postfix) parslets. Parslets are small pieces of code
that do one job and that is to cosume enough information to build the right
expression. The parselets don't check the conformity of the DSL, that is the job
of the interpreter.

### Example

A simple example could be considered:

```
/event.(@Date=="2017-03-10T23:00:00Z")/colour.(@Red==20)
```

The DSL creates an AST and depending on the interpreter will locate the event
you requested.

-----

### Naming

1. [cilli](http://www.bletchleypark.org.uk/resources/filer.rhtm/683077/chapters+11)

	A cribbing method for obtaining Enigma settings derived from the use of slack procedures by German radio operators. “Cillis” included operators using letters easy to remember and derived from such things as own initials, girlfriend’s name or obscene 4-letter words
