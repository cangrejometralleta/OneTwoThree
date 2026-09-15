# Naming

- Functions Follow  
  **Verb + Noun + context** rhythm
- A name longer than three words Suggests unclear responsibility.
- A variable that travels Holds three words,  
  joined by its language.
- A variable that lives in three lines  
  Holds one word, because the scope says the rest.
- Three words Fit in memory and survive a rename.
- A construct Names the responsibility, never the vendor.  
  *(Provisional)* `store` says what it Does;  
  `storegorm` says who it Called.
- A filename may Name the guest.  
  `store_gorm.go` Tells a reader where the ORM lives,  
  and no caller ever types it.
- A language with one type per file Loses that seam.  
  There the vendor Lives in the comment, and nowhere else.

```text
sumItemPrices      JavaScript, Go unexported
SumItemPrices      Go exported
sum_item_prices    Python
Sum Item Prices    Markdown, OneTwoThreeCase
```
