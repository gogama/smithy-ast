# smithy-ast
Native Go implementation of AWS Labs' [Smithy](https://awslabs.github.io/smithy) language's JSON AST

---

TODO:
- Test.
- Document.
- Finish prelude model including traits and metadata [validators]
- Create the parent smithy-go repo.
  - Put the CONTRIBUTORS etc. there,
  - Add copyright notices.


Where my thinking is at now:
  - This repo ~~needs~~ does not need to support extending new traits.
    The non-builtins like AWS traits can be added in another repo/package
    as a post-process step on the model. (Replace InterfaceNode with
    purpose-built trait if desired).
  - A new repo smithy-validate will contain:
    - Selector language, parsing, evaluation
    - Support for metadata-based validators
    - The ability to traverse an AST and apply all validation rules.
