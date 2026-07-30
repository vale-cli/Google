# Google

> **NOTE**: This project is neither maintained nor endorsed by Google.

This repository contains a [Vale-compatible](https://github.com/errata-ai/vale) implementation of the [*Google Developer Documentation Style Guide*](https://developers.google.com/style/) ([CC BY 4.0](https://creativecommons.org/licenses/by/4.0/)).

## Getting Started

To get started, add the package to your configuration file (as shown below) and then run `vale sync`.

```ini
StylesPath = styles
MinAlertLevel = suggestion

Packages = Google

[*]
BasedOnStyles = Vale, Google
```

See [Packages](https://vale.sh/docs/keys/packages) for more information.

## Repository Structure

<dl>
  <dt><a href="https://github.com/errata-ai/Google/tree/master/Google"><code>/Google</code></a></dt>
  <dd>The <a href="http://yaml.org/">YAML</a>-based rule implementations that make up our style.</dd>

  <dt><a href="https://github.com/errata-ai/Google/tree/master/fixtures"><code>/fixtures</code></a></dt>
  <dd>The individual unit tests. Each directory should be named after a rule found in <code>/Google</code> and include its own <code>.vale.ini</code> file that isolates its target rule.</dd>

  <dt><a href="https://github.com/errata-ai/Google/tree/master/testdata"><code>/testdata</code></a></dt>
  <dd>The expected Vale output for each fixture directory, one <code>&lt;Rule&gt;.ct</code> file per fixture. We use <a href="https://github.com/google/go-cmdtest">go-cmdtest</a> to run Vale against each fixture and compare its output. Run the suite with <code>go test ./...</code>; regenerate the expectations after an intentional change with <code>go test ./... -update</code>.</dd>

  <dt><a href="https://github.com/errata-ai/Google/tree/master/coverage"><code>/coverage</code></a></dt>
  <dd>How much of the style guide we implement, tracked topic by topic. Each key is a subtopic set to <code>true</code> or <code>false</code>, optionally followed by a comment naming the rules that implement it. Run <code>go test -v -run TestCoverage ./...</code> to print the current figures; the same test fails if a named rule no longer exists, so a topic can't silently claim coverage it has lost.</dd>
</dl>
