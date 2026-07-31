# STAGING.md — El Paso entre el Caos y el Canon

> Este archivo NO es la fuente de verdad del proyecto.
> La canon vive en VALUES / RULES / PATTERNS, que se explican solos;
> README es la puerta que las indexa.
>
> Es la segunda de tres etapas:
>
>     CHAOS.md    →    STAGING.md    →    VALUES / RULES / PATTERNS
>     privado          público            canon
>     la vivencia      el patrón          la creencia
>
> CHAOS.md vive fuera del repositorio y guarda la vivencia entera, con
> nombres y fechas. Acá llega la misma pieza con la persona sacada: lo
> que queda es el patrón, nunca el episodio. Nada entra a este archivo
> sin pasar por esa poda, y esa es toda la garantía de que un repo
> público no termine contando una vida privada.
>
> De acá cada pieza se destila y se borra: una creencia va a VALUES,
> una regla que un agente pueda ejecutar va a RULES, una raíz va a
> PATTERNS. El proceso completo está en el README, bajo
> "How Context Becomes Canon".
>
> Tope: nueve entradas. Con las nueve llenas hay que destilar antes de
> promover una décima. Un límite genera calidad, también acá.
> Este archivo encoge; el canon crece.
>
> Cada entrada lleva su estado: SIN DESTILAR, o ESPECULATIVO cuando la
> conexión todavía es una corazonada y no una raíz.
>
> Los términos acuñados se dejan en su forma original (inglés);
> el relato va en español.
>
> No se llama CLAUDE.md porque Claude Code carga ese nombre como
> instrucción del proyecto. Este archivo no lo es: CLAUDE.md queda
> libre para el puntero corto al canon.


## Historia del Nombre

Lo único de la sección de origen que no aterrizó en ningún lado.

- El proyecto empezó llamándose **"Three"**. Pasó a **OneTwoThree** para
  capturar el movimiento: no es un número, es "una cuenta, un ritmo,
  una anticipación".
- **OneTwoThreeCase** se llamó antes **TriCase**. El renombre siguió al del
  proyecto — el nombre viejo contaba el tres en vez de recorrerlo.

Ambos son la misma decisión aplicada dos veces: nombrar el movimiento y no
la cantidad. Candidato a raíz, si encuentra con qué emparejarse.


## Contexto sin Destilar

Ocupadas 7 de 9.

- **Trabajo perdido**: dos conversaciones del proyecto editaron el manifiesto
  y su salida nunca llegó a `develop`. Una escribió en RULES las secciones
  `Trust`, `Documents`, `Maintenance` y `Languages`, en VALUES las secciones
  `Breath`, `Rhyme` y `Trust`, y un `examples/` en Go, JS, Rust y Python
  resolviendo *calificar la confianza de un cambio según sus commits*. La otra
  reencuadró `Structure` como lente narrativa y sacó una regla de dependencias.
  Recuperar antes de reescribir esas zonas, o se duplica el trabajo.

- **Linter propio**: las reglas ya son casi todas mecánicas — contar beats,
  contar palabras de un nombre, medir anidamiento. Un binario en Go que lea
  RULES.md como su propia config cerraría el círculo: el markdown ya es un AST
  de tres capas, y el proyecto se validaría con su propio documento.
  Es roadmap, no creencia. No sube al canon hasta que exista.

- **examples/before-after**: los ejemplos actuales muestran el ideal. Falta el
  mismo problema escrito mal —un nombre largo, cuatro responsabilidades,
  anidamiento profundo— y después corregido.
  *Debugging is Rapping: mostrá el verso que no rimaba.*
  Aterriza en `examples/`, no en el canon.

- **Presencial sobre escrito**: para intercambios emocionalmente
  significativos se elige la conversación en persona sobre el mensaje.
  Posible raíz sobre ancho de banda del canal — el texto pierde justo lo que
  esos intercambios necesitan. Falta el puente con código.

- **Magic Commander / singleton** (ESPECULATIVO): el formato permite una sola
  copia de cada carta. El límite elimina el draw confiable y obliga al mazo a
  improvisar. Podría ser raíz de "A Limit Generates Quality", o podría ser un
  hobby sin carga conceptual. No subir sin confirmar.

- **Otro tres que aparece solo** (ESPECULATIVO): spotify-player se configura
  con tres archivos — app, theme y keymap. Anotado por si se repite. La
  muestra es chica y el sesgo de confirmación acá es enorme: contar treses
  que ya estabas buscando no es evidencia.

- **Cluster modular con LiteLLM**: objetivo técnico de largo plazo.
  Continuación natural de `The Guest you can Evict`. Esperar a que exista.


## Mapa de Procedencia

De qué conversación del proyecto salió cada cosa (para trazar el origen):

- **Elegancia en el código y la música** — nacimiento del proyecto, los tres
  archivos, OneTwoThreeCase, Dev La Soul, el renombre a OneTwoThree, De La Soul,
  Gorillaz, Grace Hopper, camelCase, Bug Fables.
- **El número 3 como patrón filosófico** — el inventario del 3 por dominio.
- **Narradores poco confiables de nosotros mismos** — OneTwoTree, significado
  distribuido, git-SHA como puntos de confianza, Naur, Lehman, el bajo.
- **Patrones de burnout autista y minimalismo en OneTwoThree** — el vínculo
  carga cognitiva ↔ regla del tres; minimalismo como instinto, no como estilo.
- **Cómo funciona el burnout autista** — origen vivencial del minimalismo
  (contexto personal; detalles privados deliberadamente omitidos aquí).
- **Commit de cambios en equipo** — proyecto local "Carmen Gloria", flujo git.
- **Reglas y estrategias de resolución** — el criterio Rules/Values, beats
  contra newlines, Go validando OneTwoThreeCase, la idea del linter.
- **Narrative structure in code** — el tres como lente y no como límite,
  Receive → Transform → Return, el código como oración declarativa.
- **Mammouth.ai OpenCode compatibility** — el escritorio en tres columnas,
  y la evidencia vivida del lock-in de proveedor.
