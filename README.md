# mad
[![Build Status](https://travis-ci.org/sirkon/mad.svg?branch=master)](https://travis-ci.org/sirkon/mad)

This is a literate configuration format library for Go inspired by Markdown. 


Read about typical usecases [here](USECASES.md). 
In short, ___mad___ library is not yet another configuration format. It is not really suitable for this role. 
Its is main strenght are cases where pieces of structured data or code are needed and there's not much nesting in data.

## Installation

* the prefered way is via the dep
    ```
    dep ensure -add github.com/sirkon/mad    
    ```
* using go get is not recommended, but still should work
    ```
    go get github.com/sirkon/mad
    ```
 

## Simplest usage

Let we have a 
* some temporary table with preaggregated data and several queries running on this table. Let we call
this temporary table creation a *prepare*. 
* quries may run on a temporary table and they also may use other tables
* actually, it is rather uncommon situation when we need a preaggregated table. In the majority of cases we
    will use existing tables. Thus, the *prepare* process is not mandatory.
* prepare query creates table. There should be a method to delete it. 
    
In final, we have
* optional prepare create and delete queries
* regular queries

Let's express this in a Go structure:

```go
import "github.com/sirkon/mad"
…
var job struct {
	Prepare *struct {
		Create mad.Code `mad:"create"`
		Delete mad.Code `mad:"delete"`
	} `mad:"prepare,syntax=sql"` // syntax will stay in the context of nested fields 
	
	Queries []mad.Code `mad:"queries,syntax=sql"`
}
```

It can be decoded from file as

```go
if err := mad.UnmarshalFile("file.md", &job, mad.NewContext()); err != nil {
	panic(err)
}
```

And now how the markdown text would look like for this example:

````markdown
# prepare

## create
```sql
CREATE TABLE tmp AS SELECT * FROM table
```

## delete
```sql
DROP TABLE tmp;
```


# queries

* total amount of events
```sql
SELECT count(1) FROM tmp;
```

* amount of users who created events
```sql
SELECT uniq(user_id) FROM tmp;
```
````

But it is not you only get these SQL pieces. You also get their position in the Markdown file:

```go
fmt.Println(job.Prepare.Create.Start())
fmt.Println(job.Prepare.Create.Finish())
```

will output something like 

```
10 1
14 1
```



* You can also consume integers, floating points, strings and booleans using `raw` code blocks (see [advanced example](HARD.md)):
* you need `mad.Code` type to consume code block
* you need `mad.Comment` type to consume comment block
* comments (everything, that is not header or fenced code block) are normally ignored (this is a bug if they aren't) 
    except the case where you are using `mad.Comment` 


This simple example is not different from what you can used to with stdlib JSON package (except may be only
fields tagged with `mad:".."` can be filled by Decoder). But this is human oriented format and that's why special
interfaces were introduced:
* It is possible to control fragment order in Markdown files (define your `mad.Decodable` type and write down desired decoding order manually)
* It is possible to use header as a source of information, fulfilling fields with data from header text (define your `mad.Manual` type)
* It is possible to reimplement how scalar types are treated (define your own `mad.Unmarshaler` type)
* It is possible to match by regexp, not just by fixed name. Usually `mad.Manual` is used in these cases as a field value type
* It is possible to match against several syntax types, for instance, use `syntax=python perl` to match against perl or python fenced code block

see at the harder example [here](HARD.md)

## Why not to use YAML with text blocks?

Because

* You don't have same "fenced code blocks" in YAML, thus you don't have syntax highlighting, folding, etc for them
* The human readable reporting was in mind, thus you can get positions of code blocks and headers. YAML parsers available
    for go cannot pass this information down to the user. 
