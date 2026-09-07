# Naming

- Functions Follow  
  **Verb + Noun + context** rhythm
- A Name longer than Three Words Suggests unclear Responsibility.
- A Variable that Travels Holds three Words,  
  joined by its Language.
- A Variable that Lives in three Lines  
  Holds one Word, because the Scope Says the rest.
- Three Words Fit in Memory and Survive a Rename.
- A Construct Names the Responsibility, never the Vendor.  
  *(Provisional)* `store` Says what it Does;  
  `storegorm` Says who it Called.
- A Filename may Name the Guest.  
  `store_gorm.go` Tells a Reader where the ORM Lives,  
  and no Caller ever Types it.
- A Language with one Type per File Loses that Seam.  
  There the Vendor Lives in the Comment, and nowhere else.

```text
sumItemPrices      JavaScript, Go unexported
SumItemPrices      Go exported
sum_item_prices    Python
Sum Item Prices    Markdown, OneTwoThreeCase
```
