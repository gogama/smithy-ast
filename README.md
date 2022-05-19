# smithy-ast
GoLang implementation of AWS Labs' Smithy language's JSON AST

---

TODO:
- ~~Refactor Model interface globals to support installing additional
  trait/node mappings.~~ [Decided against this.]
- Finish implementing trait node.
- Test.
- Document.
- Write prelude model AST JSON and minify GZIP it. Or just structurally
  model it, who knows?


Need to figure out:
- ~~Selectors (language)~~
- ~~Selectors (how to model them as attributes on traits)~~
- ~~Validators, how to model built-in validators~~

Where my thinking is at now:
  - This repo ~~needs~~ does not need to support extending new traits.
    The non-builtins like AWS traits can be added in another repo/package
    as a post-process step on the model. (Replace InterfaceNode with
    purpose-built trait if desired).
  - A new repo smithy-validate will contain:
    - Selector language, parsing, evaluation
    - Support for metadata-based validators
    - The ability to traverse an AST and apply all validation rules.
