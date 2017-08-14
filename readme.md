# @scribe
A self documenting CSS specification

### What is Scribe?
Scribe is a system for documenting CSS and auto-generating docs for your CSS frameworks.

### How does Scribe work?
Scribe notation is simple. Simply add a comment block above CSS you'd like to document.

```
/*
@scibe Border colors

<template>
  <div class="{{class}}"></div>
</template>
*/

.border-ok {
  border: green;
}

.border-danger {
  border: red;
}
```

Scribe will iterate over every class below a `@scribe` block until it reaches the end of the file or a new comment block.

Given the example above, Scribe will produce the following HTML:

```html
<div class="border-ok">Lorem ipsum dolor ist</div>
<div class="border-danger">Lorem ipsum dolor ist</div>
```

### Running Scribe
`rm -rf test.html && go run *.go --input=YOUR INPUT DIRECTORY`