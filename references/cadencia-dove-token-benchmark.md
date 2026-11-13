# Cadencia Dove and Caveman Summary Benchmark

## Run record

- **Stories:** six public-domain stories; three Spanish and three English.
- **Task prompt:** “Summarize this complete story in exactly three paragraphs. Preserve the main plot events, causal relations, and ending. Use only facts supported by the story. Reply in the story’s language.”
- **Vanilla:** task prompt only; no added style instruction.
- **Caveman:** task prompt plus public Caveman skill rules: answer first, remove filler, keep meaning and facts, use concise sentences, do not break grammar.
- **Dove:** task prompt plus provisional Cadencia Dove: concise, natural, rhythmic phrasing; remove repetition while preserving meaning; use DeLaCase.
- **Model:** GPT-6-based Codex runtime. Exact serving model ID is not exposed.
- **Effort:** not exposed by the runtime; no per-condition effort override was available. All three conditions were drafted in the same active session.
- **Execution limit:** these are style-conditioned drafts by one model in one session, not independent API completions or blind samples. The vanilla condition means no added style layer; it still inherits the session context.
- **Token counting:** [`tiktoken`](https://github.com/openai/tiktoken) 0.14.0 with `o200k_base`. This is a consistent proxy: the installed tokenizer has no GPT-6 model mapping. Counts cover summary bodies only, excluding source text, instructions, headings, and hidden system context. They are not billed-token telemetry.

## Results

The count is the token count of each three-paragraph summary body, including its paragraph breaks.

| Language | Story | Vanilla | Caveman | Dove |
|---|---|---:|---:|---:|
| Spanish | La tortuga gigante | 132 | 116 | 100 |
| Spanish | La abeja haragana | 136 | 102 | 108 |
| Spanish | El loro pelado | 115 | 96 | 95 |
| **Spanish total** |  | **383** | **314** | **303** |
| English | The Gift of the Magi | 99 | 80 | 88 |
| English | The Last Leaf | 98 | 80 | 76 |
| English | The Ransom of Red Chief | 101 | 81 | 84 |
| **English total** |  | **298** | **241** | **248** |
| **All stories** |  | **681** | **555** | **551** |

## Shared sources

- [La tortuga gigante](https://es.wikisource.org/wiki/La_tortuga_gigante)
- [La abeja haragana](https://es.wikisource.org/wiki/La_abeja_haragana)
- [El loro pelado](https://es.wikisource.org/wiki/El_loro_pelado)
- [The Gift of the Magi](https://www.gutenberg.org/ebooks/7256)
- [The Last Leaf](https://www.gutenberg.org/files/3707/3707-h/3707-h.htm)
- [The Ransom of Red Chief](https://www.gutenberg.org/cache/epub/5403/pg5403.html)

## Summaries

### La tortuga gigante

#### Vanilla

```
Un hombre de Buenos Aires enferma y se muda al monte para recuperarse. Allí mata a un tigre que atacaba a una tortuga y, pese a su hambre, cura las heridas del animal.

Más tarde, el hombre cae enfermo. La tortuga lo alimenta y, al saber que solo los remedios de Buenos Aires pueden salvarlo, lo carga en su caparazón y emprende el viaje.

Agotada, cree que aún está lejos, hasta que un ratón le revela que ya llegó. El director del zoológico cura al hombre y se queda con la tortuga; desde entonces, el cazador la visita.
```

#### Caveman

```
Un hombre enferma en Buenos Aires y se muda al monte. Allí mata al tigre que ataca a una tortuga. Aunque tiene hambre, cura sus heridas.

Después, el hombre enferma. La tortuga lo alimenta y lo lleva en su caparazón hacia Buenos Aires para salvarlo con remedios. Cruza selvas, ríos y pantanos.

Agotada, cree que está lejos. Un ratón le dice que ya llegó. El director del zoológico cura al hombre y cuida a la tortuga. El cazador la visita.
```

#### Dove

```
Un Hombre enfermo se muda de Buenos Aires al monte. Allí mata al Tigre que ataca a una Tortuga, pero cura sus heridas.

Cuando el Hombre cae enfermo, la Tortuga lo alimenta y lo lleva hacia Buenos Aires sobre su caparazón. El largo Viaje casi la vence.

Un Ratón le revela que ya llegó. El director del zoológico cura al Hombre y cuida a la Tortuga. El cazador la visita; su Amistad permanece.
```

### La abeja haragana

#### Vanilla

```
Una abeja evita trabajar mientras sus compañeras producen miel. Tras varias advertencias, las guardianas le impiden entrar en la colmena durante una noche fría y lluviosa.

La abeja cae en una cueva y una culebra amenaza con comérsela. La abeja la desafía a una prueba: la culebra hace bailar un trompo, y ella se esconde dentro de una sensitiva cuyas hojas se cierran al contacto.

Al amanecer vuelve a la colmena y cambia su conducta. Desde entonces trabaja y, antes de morir, enseña a las abejas jóvenes que el esfuerzo de todas beneficia a la comunidad.
```

#### Caveman

```
Una abeja evita trabajar. Tras varias advertencias, las guardianas la expulsan durante una noche fría y lluviosa.

En una cueva, una culebra quiere comérsela. Compiten en una prueba: la culebra hace bailar un trompo; la abeja se esconde en una sensitiva.

Al amanecer, la abeja vuelve a la colmena y trabaja. Antes de morir, enseña a las jóvenes que el esfuerzo de cada una beneficia a todas.
```

#### Dove

```
Una Abeja evita trabajar mientras las demás llenan la colmena. Tras las advertencias, las guardianas la expulsan en una noche fría.

En una cueva, una Culebra amenaza con comerla. Aceptan una prueba de ingenio: la culebra hace bailar un trompo; la abeja se oculta en una sensitiva.

Al amanecer, la Abeja vuelve y trabaja. Antes de morir, deja una Enseñanza: el Trabajo de cada una sostiene la Felicidad de todas.
```

### El loro pelado

#### Vanilla

```
Pedrito, un loro domesticado, vive feliz con una familia hasta que un tigre lo ataca en el monte. Escapa, pero pierde sus plumas y se esconde avergonzado en un árbol.

Cuando le crecen plumas nuevas, Pedrito cuenta lo ocurrido y pide al dueño que lo acompañe a cazar al tigre en Paraguay.

Pedrito distrae al tigre con sus frases mientras el hombre apunta y lo mata. De vuelta en casa, el loro se burla de la piel del tigre cada tarde durante el té.
```

#### Caveman

```
Pedrito, un loro domesticado, vive con una familia. En Paraguay, un tigre lo ataca y le arranca las plumas. Avergonzado, Pedrito se esconde.

Cuando recupera sus plumas, cuenta lo ocurrido y convence al dueño de cazar al tigre.

Pedrito distrae al tigre mientras el hombre apunta y lo mata. En casa, el loro se burla de la piel del tigre durante el té.
```

#### Dove

```
Pedrito, un Loro domesticado, vive feliz con una familia. En Paraguay, un Tigre lo ataca y le arranca las plumas. Avergonzado, se esconde.

Al recuperar sus plumas, Pedrito cuenta lo ocurrido y convence al dueño de cazar al Tigre.

Pedrito distrae al Tigre mientras el hombre apunta y lo mata. En casa, se burla de la piel del Tigre durante el té.
```

### The Gift of the Magi

#### Vanilla

```
Della has saved $1.87 for a Christmas gift for Jim. She sells her long hair and buys a chain for his gold watch.

Jim has sold his watch to buy Della the combs she admired. Each gift is unusable at first: Jim has no watch for the chain, and Della has cut her hair.

They set the gifts aside and share their love. The narrator calls them wise gift-givers, like the Magi who brought gifts to the Christ child.
```

#### Caveman

```
Della has $1.87 for Jim's Christmas gift. She sells her hair and buys a chain for his gold watch.

Jim sells his watch to buy Della the combs she wanted. The gifts cannot be used: he has no watch, and she has cut her hair.

They put the gifts away. Their sacrifice shows their love. The narrator calls them wiser than the Magi.
```

#### Dove

```
Della has $1.87 for Jim’s Christmas Gift. She sells her long Hair to buy a gold Chain for his watch.

Jim sells his watch to buy Della the combs she admires. Their Gifts cannot be used: he sold his watch, and she cut her hair.

They set the Gifts aside. Their Sacrifice shows their Love. The narrator calls them the wisest gift-givers, like the Magi.
```

### The Last Leaf

#### Vanilla

```
Sue and Johnsy are artists who share a studio. When Johnsy develops pneumonia, she believes she will die when the last ivy leaf falls from the wall outside.

Sue asks their neighbor Behrman, an old painter, to pose for her drawing. After a storm, one leaf remains; Johnsy regains hope and begins to recover.

Behrman dies of pneumonia. Sue tells Johnsy that he painted the leaf during the storm, creating the masterpiece he had promised to make.
```

#### Caveman

```
Sue and Johnsy share an artists' studio. Johnsy gets pneumonia and believes she will die when the last ivy leaf falls.

Sue asks neighbor Behrman to pose. After a storm, one leaf remains. Johnsy regains hope and begins to recover.

Behrman dies of pneumonia. Sue reveals he painted the leaf on the wall. It was the masterpiece he had planned for years.
```

#### Dove

```
Sue and Johnsy share an artists’ studio. Johnsy gets Pneumonia and believes she will die when the last Ivy Leaf falls.

Sue asks neighbor Behrman to pose. After a Storm, one leaf remains; Johnsy regains Hope and recovers.

Behrman dies of pneumonia. Sue reveals he painted the leaf. It was his long-awaited Masterpiece.
```

### The Ransom of Red Chief

#### Vanilla

```
Sam and Bill kidnap Johnny Dorset, expecting his wealthy father to pay a $2,000 ransom. Johnny enjoys the adventure and calls himself Red Chief.

Johnny torments the men, keeps them awake, and attacks Bill. Sam and Bill become desperate and send a ransom demand to Johnny's father.

Ebenezer replies that the kidnappers must pay him $250 to take Johnny back. They agree and return the boy, who clings to Bill as his father promises to hold him for ten minutes.
```

#### Caveman

```
Sam and Bill kidnap Johnny Dorset for a $2,000 ransom. Johnny enjoys the cave and plays Red Chief.

The boy torments them, attacks Bill, and keeps them awake. The kidnappers demand ransom from his father.

Dorset offers to take Johnny back for $250. The men pay and return him. Johnny clings to Bill; his father promises to hold him for ten minutes.
```

#### Dove

```
Sam and Bill kidnap Johnny Dorset, expecting $2,000 from his wealthy Father. Johnny enjoys the cave and becomes Red Chief.

Johnny torments them, attacks Bill, and keeps them awake. The desperate men demand Ransom from Ebenezer Dorset.

Dorset offers to take Johnny back for $250. They pay and return him; Johnny clings to Bill, so his Father promises to hold him ten minutes.
```

## Second run: Sonnet at high effort, five conditions, with a quality check

The first run found brevity without a quality check and ran in one Codex session. This run uses independent calls, real token counts from the runtime, and a judge for the summaries. Raw outputs and scores are in [results.json](cadencia-dove-benchmark-sonnet-high/results.json) and [judge.json](cadencia-dove-benchmark-sonnet-high/judge.json).

### Run record

- **Date and runtime:** 2026-11-13; Claude Code 2.1.282 with `--print --tools "" --model sonnet --effort high`. Every call was its own session. The usage report names the serving model `claude-sonnet-5`.
- **Conditions:** vanilla (task only); **Caveman**, the public [`skills/caveman/SKILL.md`](https://github.com/JuliusBrussee/caveman/blob/2178ed6495c3c0ecd7d037e60acf3b1d0c554b97/skills/caveman/SKILL.md) as of commit `2178ed6`, 744 words; **Cadencia full**, `cadencia-dove` at commit `2253fd4`, 938 words; **Mini on**, `cadencia-dove-mini` with DeLaCase on by default, about 320 words; **Mini off**, the same text with "No DeLaCase." (or "Sin DeLaCase.") added to the request. Each style text was passed as an instruction ahead of the task, as in the first run, not installed as a skill.
- **Summaries:** the six stories of the first run, in full from Wikisource, with the same three-paragraph prompt, two samples per story and condition (12 per condition).
- **Explanations:** eight short questions, four English and four Spanish, "answer in two short paragraphs", one sample each.
- **Counting:** `output_tokens` from the runtime usage report. **Thinking** is `thinking_tokens` from the same report, billed as output but not shown. **Visible** is output minus thinking. Words are of the visible text.
- **Judge:** Sonnet at high effort, one summary per call, in shuffled order, shown only the story and the summary. It scored four criteria from 0 to 2: main events, causal relations, ending, and faithfulness (no wrong or invented fact), and listed each error.

### Summaries (12 per condition, tokens are means)

| Condition | Visible | Thinking | Total output | Judge score (of 8) | Faithfulness (of 2) | Perfect 8/8 | Exactly three paragraphs |
|---|---:|---:|---:|---:|---:|---:|---:|
| Vanilla | 764 | 4 | 768 | 7.58 | 1.58 | 7/12 | 11/12 |
| Caveman | 416 | 578 | 994 | 7.42 | 1.42 | 5/12 | 11/12 |
| Cadencia full | 577 | 100 | 678 | 7.67 | 1.67 | 8/12 | 12/12 |
| Mini on | 524 | 1112 | 1636 | 7.42 | 1.42 | 5/12 | 12/12 |
| Mini off | 570 | 74 | 644 | 7.75 | 1.75 | 9/12 | 12/12 |

Median total output was 786 (vanilla), 812 (Caveman), 710 (Cadencia full), 1122 (Mini on) and 678 (Mini off) tokens. Events, causality and ending scored 2.00 in every condition, so the differences above come only from faithfulness.

### Explanations (8 per condition, tokens are means)

| Condition | Visible | Thinking | Total output | Words |
|---|---:|---:|---:|---:|
| Vanilla | 315 | 0 | 315 | 149 |
| Caveman | 156 | 67 | 223 | 63 |
| Cadencia full | 183 | 0 | 183 | 88 |
| Mini on | 208 | 671 | 879 | 91 |
| Mini off | 211 | 20 | 231 | 103 |

A separate DeLaCase check ran the same eight questions, two samples each, with the Mini on and off ([delacase-check.json](cadencia-dove-benchmark-sonnet-high/delacase-check.json)). With DeLaCase off, none of 16 answers carried a capital of emphasis. With it on, 15 of 16 did and, counting a two-word name once, 2 of 16 went over the three-per-passage ceiling. The counter is a heuristic, not a judgment.

### What the run shows

- **Caveman writes the shortest visible text:** 46% fewer tokens than vanilla on summaries and 50% on explanations, against 24% to 42% for Cadencia full and 25% to 34% for the Mini. The first run's near tie did not repeat.
- **Visible length is not the bill:** Caveman and Mini on triggered long hidden reasoning at high effort, so their total output exceeded vanilla on summaries. Cadencia full and Mini off stayed at or below vanilla in total.
- **DeLaCase on by default is the costly part:** with it, the Mini's total output on summaries rose from a median of 678 to 1122 tokens. The effort setting matters here; lower effort may change this.
- **Quality moved little:** every summary kept events, causes and endings. Faithfulness ranged from 1.42 to 1.75 over 12 summaries per condition, and the judge flagged slips such as a chain "months" long where the story says weeks, or a bird killed where it is only wounded. Those flags appear in all five conditions.

### Limits

- Twelve summaries and eight answers per condition, one session, one model and one effort setting. A gap of a few tenths of a point is inside the noise.
- The judge is one model of the same family as the writer, and it can see a condition's style, so it is not blind to it. It also flagged small inferences as errors. A second judge or a human read would test it.
- Hidden reasoning varies a lot between runs, so the totals carry wide spread; medians are given beside the means.
- Caveman was passed as an instruction, not loaded as an installed skill, and Cadencia full was measured after it was made standalone.
- No cost saving is claimed: per-call cost was dominated by the story and the fixed context, and no style was cheaper than vanilla.
