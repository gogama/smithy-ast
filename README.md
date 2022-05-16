# smithy-ast
GoLang implementation of AWS Labs' Smithy language's JSON AST

---

TODO:
- Traits probably need to be modelled as Nodes. Decide this and if yes,
  node them.
- Refactor Model interface globals to support installing additional
  trait/node mappings.
- Test and finalize.


Missing Traits:
- @unitShape - Looks like an annotation trait.
- @trait - https://awslabs.github.io/smithy/1.0/spec/core/model.html#traits
- @suppress

Need to figure out:
- Selectors (language)
- Selectors (how to model them as attributes on traits)
- Validators, how to model built-in validators

Where my thinking is at now:
  - This package needs to support extending new traits, and putting selectors
    ON traits.
  - A new package smithy-validate will contain:
    - Selector language, parsing, evaluation
    - Support for metadata-based validators
    - The ability to traverse an AST and apply all validation rules.
