# beego-pipeline

Asset compilation and compression for golang [beego framework](http://beego.me/)
inspired by [django-pipeline](https://django-pipeline.readthedocs.org/).



### Requirements

* [NodeJS](https://nodejs.org/)
* [Less](http://lesscss.org/) - Globally installed `sudo npm -g install lessc`
* [Yuglify](https://github.com/yui/yuglify) - Globally installed `sudo npm -g install yuglify`


## Configuration

The configuration is stored in a `pipeline.json` file located in the `conf/` folder in your
beego application. Example usage:

```json
{
  "css": {
    "base": {
      "root": "",
      "sources": [
        "css/base.css",
        "css/mixins.css"
      ],
      "output": "dist/app.css"
    },
    "output": {
      "sources": [
        "css/base.css",
        "css/less1.less",
        "css/less2.less"
      ],
      "output": "dist/output.css"
    }
  },
  "js": {}
}

```

We can specify the following keys: **css**, **js**. Each mapping inside represents the
output of compiled/compressed/versioned file.

The **root** key can be omited and will default to "static/". This is the location where
your assets reside if they are different from the default.

#### Configuration in the code

You must include the pipeline and any compilers/compressor in the code. At the moment
there is only `yuglify` support for compressor cause this supports both `css` and `js`.

In `main.go` import pipeline and any compiler/compressor that you require:

```golang

package main

import (
  ...
  "github.com/astaxie/beego"
  _ "github.com/bogh/beego-pipeline"
  _ "github.com/bogh/beego-pipeline/compilers/less"
  _ "github.com/bogh/beego-pipeline/compressors/yuglify"
)

func main() {
  beego.Run()
}

```

#### Override command parameters

WIP


## TODO

* Sass compiler
* TypeScript compiler
* Do not compress if in debug mode

