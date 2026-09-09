---
name: one-two-three-case
description: Convert any code identifier or prose into a OneTwoThreeCase suggestion — capitalize Entities, Actions and Statuses, lowercase connectors, first word of each sentence capitalized; for identifiers, apply the language's case convention plus Verb+Noun+context (≤3 words). Use when the user asks to "OneTwoThreeCase" a name, sentence, or identifier, or to rewrite something in the manifesto's convention.
---

# OneTwoThreeCase

A Converter, not a Checklist. It takes Code or Words
and returns one Suggestion, plus one Line saying why.

The Checklist is [one-two-three-refactor](../one-two-three-refactor/SKILL.md).
Use that while Writing a whole Unit. Use this to Rename one.
To Hear the Pattern under an Explanation, Use the
[one-two-three-agent](../../agents/one-two-three-agent.md) Agent.

The canon lives in [OneTwoThreeCase](../../../rules/one-two-three-case.md).

## How to Convert

1. Classify the Input — Prose or Code?
   Prose is a Sentence or a Comment.
   Code is an Identifier, a Function, a Variable.

2. Find the Important Words —
   Entities, Actions and Statuses stay Capitalized.
   Connectors and local Names go lowercase.
   The First Word of a Sentence stays Capitalized,
   even when it is a Connector.

3. For Code, apply the Language Convention —
   PascalCase when Exported,
   camelCase when Unexported,
   snake_case for Python.
   Then fit the Name to Verb + Noun + context,
   three Words at most.

4. Return the Suggestion, and one Line
   naming which Words you treated as
   Entity, Action or Status, and which as Connectors.

## Examples

### Prose

Input:  "this function gets the user account data by id"
Output: "This Function Fetches the User Account by ID"
Why:    Function, Fetches, User, Account and ID are the Action and Entities; the and by are Connectors.

### Code — Go exported

Input:  GetUserAccountDataById
Output: FetchAccount
Why:    Fetch is the Action, Account the Entity; context and id were Local, so they were Dropped to hold three Words.

### Code — Go unexported

Input:  getUserAccountDataById
Output: fetchAccount
Why:    camelCase Stays inside; the same three Words, lowercased at the Front.

### Code — Python

Input:  get_user_account_data_by_id
Output: fetch_account
Why:    snake_case for Python; three Words again, id and data were Local to the Scope.

## When not to Convert

- A Name that already Reads as three Words in its Language.
- A Variable that Lives in three Lines — it may Keep one Word,
  because the Scope already Says the rest.
