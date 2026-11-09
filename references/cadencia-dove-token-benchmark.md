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
