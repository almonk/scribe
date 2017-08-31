# Scribe – A syntax for documenting CSS

## What is Scribe?

Scribe is a markup syntax for CSS which allows an interpreter to render documentation for a given stylesheet.

## Why Scribe?

Documenting CSS is hard. It normally involves a lot of manual maintenance and editing. Scribe automates documentation to let you focus on making things.

* * *

## Syntax

### INTRODUCTION

The simplest possible Scribe block consists of the `@scribe` keyword and `<template>`:

```css
/*
@scribe Text colors

<template>
    <span class="{{class}}">Hello</span>
</template>
*/

.red {
   text-color: red;
}

.blue {
   text-color: blue;
}
```

Running Scribe on this file will produce the following HTML:

```html
<span class="red">Hello</span>
<span class="blue">Hello</span>
```

You'll notice that `{{class}}` in our template is replaced with the CSS classes below the Scribe comment block.

### Multiple scribe blocks

Multiple Scribe blocks can be inserted into a CSS document. Scribe will document **all** the classes following the block until either;

* Another Scribe block is found
* Any CSS comment is found
* End of file is reached

For example; we want to document our green text color on a black background to make it easier to read.

```css
/*
@scribe Text colors

<template>
    <span class="{{class}}">Hello</span>
</template>
*/

.red {
   text-color: red;
}

.blue {
   text-color: blue;
}

/*
@scribe Text colors on dark

<template>
    <div class="bg-black">
        <span class="{{class}}">Hello</span>
    </div>
</template>
*/

.yellow {
    text-color: yellow;
}
```

The HTML produced now will be;

```html
<span class="red">Hello</span>
<span class="blue">Hello</span>
<div class="bg-black">
    <span class="yellow">Hello</span>
</div>
```

### Markdown in scribe

Sometimes it's nice to document things with a sentence or two about how to use the CSS class in question.

Scribe allows you to write Markdown as part of the Scribe block which will be parsed and displayed next to the documented CSS.

Write your markdown between `<md>` tags like so;

```css
/*
@scribe Text colors

<md>
# This is a heading
This is a paragraph

* This
* Is 
* A
* List
</md>

<template>
    <span class="{{class}}">Hello</span>
</template>
*/
```

Be careful not to indent your markdown as indentation has special meaning to the markdown parser.

### Ignoring css selectors

We've documented three of our text colors but we still have another two in our file that we don't want to document yet. Let's stop Scribe from reading them with a `/* @scribe nodoc */` comment block.

```css
/*
@scribe Text colors

<template>
    <span class="{{class}}">Hello</span>
</template>
*/

.red {
   text-color: red;
}

.blue {
   text-color: blue;
}

/*
@scribe Text colors on dark

<template>
    <div class="bg-black">
        <span class="{{class}}">Hello</span>
    </div>
</template>
*/

.yellow {
    text-color: yellow;
}

/* @scribe nodoc */

.green {
    text-color: green;
}

.orange {
    text-color: orange;
}
```




