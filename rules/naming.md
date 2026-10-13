# Naming

- Name a function with a verb, its object and needed context; use at most three words.
- Keep sibling function names structurally consistent.
- Use one term for one concept throughout the code.
- Use short names only for locals whose scope makes their meaning clear.
- Name a construct for its responsibility, not its vendor. *(Provisional)*
- A filename may name the vendor when it identifies an adapter, such as `store_gorm.go`.
- If a language requires one type per file, put vendor details in the adapter's comment.

```text
sumItemPrices      JavaScript, Go unexported
SumItemPrices      Go exported
sum_item_prices    Python
Sum Item Prices    Markdown, DeLaCase
```
