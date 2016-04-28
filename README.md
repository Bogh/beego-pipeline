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
      "root": "static/",
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

**Sources** configuration supports wild card asset discovery. e.g: css/*.css

**IMPORTANT**: Watch out so that any source file doesn't have the same path/name
as the output file. The app will complain if you mistakenly override a source file with an
output file.


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

## Usage

**Pipeline** will automatically register to template functions `pipeline_css` and
`pipeline_js`. Using this you can load the groups by name in your templates. E.g.:

```
<!doctype html>
<html>
<head>
  <title>Beego pipeline example</title>
  {{ pipeline_css "app" }}
</head>
<body>
  ...

  {{ pipeline_js "app" }}
</body>
</html>
```

If you run in `dev` mode then the assets won't be compressed and they will be added
idividually. Assets that don't need compilation (e.g.: basic `.css` files) and any file
that has no matching compiler.

```
<link href="/static/css/base.css">
<link href="/static/css/mixins.css">
```

In case of a compiled file. If the extension is `.less` a file with the same name
but correct asset extension will be generated. e.g. `app.less` => `app.css`. In
the same folder.


If you run in any other mode the assets will be compressed into the `output` setting
of the group and the file named will be appended with a hash generated from the
file contents. e.g.: `css/app.css` => `app.9fa8429f0f816065.css`.
The output of the helpers would be:
```
<link href="/static/app.9fa8429f0f816065.css">
```

###### Auto generation

**Pipeline** will automatically watch the source files for changes and regenerate
that group in `dev` mode.

## TODO

* Sass compiler
* TypeScript compiler

