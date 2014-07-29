cmds
====

It's a funny replacement for Make maybe, and maybe more.

Installation
------------

	go get github.com/crackcomm/cmds

Usage
-----

cmds is an application which reads cmds.yaml file reading commands from it and executing particular one

example cmds.yaml

```YAML
google.title:
  - http.request:
      url: http://www.google.pl/
  - html.extract:
      selectors:
        title: head > title

clean:
  - cmd.run:
    cmd: rm -rf build

build:
  - clean
  - cmd.run:
    cmd: gulp build
```

```
$ cmds google.title
Running google.title

  title: Google

```
