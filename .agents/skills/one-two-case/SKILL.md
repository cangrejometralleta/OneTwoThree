---
name: one-two-case
description: Convert code identifiers or prose into a OneTwoCase suggestion — up to two emphasis capitals per sentence beyond its free initial capital, with the main emphasis chosen by voice and an optional second marking a meaningful relationship. For identifiers, apply the language's case convention plus Verb+Noun+context (≤3 words). Use when the user asks to "OneTwoCase" a name, sentence, or identifier, or to rewrite something in the manifesto's convention.
---

# OneTwoCase

A converter, not a checklist. It takes code or words  
and Returns one suggestion, plus one line saying why.

The checklist is [OneTwoRefactor](../one-two-refactor/SKILL.md).  
Use that while writing a whole Unit. Use this to rename one.  
To hear the pattern under an explanation, use the  
[Dove](../../agents/dove.md) agent.

The canon lives in [OneTwoCase](../../../rules/one-two-case.md).

## How to Convert

1. Classify the Input — prose or code?  
   Prose is a sentence or a comment.  
   Code is an identifier, a function, a variable.

2. Split the prose into Sentences.  
   Each sentence Carries its own budget, and a bullet line counts as one.

3. Leave the first Word capitalized, and do not count it.  
   That capital is free, whatever the word is.

4. Choose the main emphasis by Voice.
   active Spends on the Action,  
   passive Spends on the Entity,  
   a copula Spends on the predicate that completes the definition,  
   an imperative Spends on the Entity, because its verb opened for free.  
   Where two entities compete, take the one carrying the Claim.

   Optionally Connect that emphasis to its Object, context or contrast with a second capital.
   Two additional capitals are the Ceiling, never a quota.

5. Lowercase everything Else.  
   Proper names and acronyms Keep their established Spelling outside the budget.
   A list of names is exempt: five Entities in a row stay five entities.

6. For code, apply the language Convention —  
   PascalCase when exported,  
   camelCase when unexported,  
   snake_case for Python.  
   Then fit the name to Verb + Noun + context, three words at most.

7. Return the Suggestion, and one line naming the voice you read  
   and the word or relationship it chose.

## Examples

### Prose — active

Input:  "this function gets the user account data by id"  
Output: "This function Fetches the user Account by id."\
Why:    Active, so *Fetches* carries the action; *Account* names its object.

### Prose — one is Enough

Input:  "the reading breathes"\
Output: "The reading Breathes."\
Why:    The action carries the whole claim; a second emphasis adds nothing.

### Prose — passive

Input:  "all six paths from the checkpoint were committed and pushed"  
Output: "All six Paths from the checkpoint were committed and pushed."  
Why:    Passive, so the Entity takes it; the actor is gone.

### Prose — copular

Input:  "a checkpoint is a snapshot, not a session closing"  
Output: "A checkpoint is a Snapshot, not a session closing."  
Why:    A copula, so the predicate takes it; *Snapshot* completes the definition.

### Prose — imperative

Input:  "read the long line as the opening"  
Output: "Read the long line as the Opening."  
Why:    Imperative, so the verb opened for free and the Entity took the spend.

### Prose — bold

Input:  "spend the capital where you would raise your voice"  
Output: "**Spend the capital where you would raise your voice.**"  
Why:    Bold is the second Tier, one to a section, never on a word a capital already marks.

### Code — Go exported

Input:  GetUserAccountDataById  
Output: FetchAccount  
Why:    Fetch is the Action, account the entity; context and id were local, so three words held.

### Code — Go unexported

Input:  getUserAccountDataById  
Output: fetchAccount  
Why:    camelCase Stays inside; the same three words, lowercased at the front.

### Code — Python

Input:  get_user_account_data_by_id  
Output: fetch_account  
Why:    snake_case for Python; three words again, id and data were local to the scope.

## When not to Convert

- A name that already Reads as three words in its language.
- A variable that lives in three lines — it may keep one Word,  
  because the scope already says the rest.
- An enumeration of names, a table Column or a diagram label.  
  None of them is a Sentence, so none of them spends.
- An identifier Spends nothing, because it is not a sentence.  
  Case there belongs to the Language, and the language decides.
