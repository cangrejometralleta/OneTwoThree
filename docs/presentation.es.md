# Presentación de OneTwoThree

## El Lenguaje Abre una forma de Colaborar

**Una misma frase puede comunicar una idea a una persona y orientar la acción de una máquina.**

Hoy podemos Conversar con un Agente de inteligencia artificial en el mismo Idioma que usamos entre nosotros. Podemos explicar una Intención, pedir un Cambio y precisar lo que quisimos decir mediante otra frase. El Lenguaje cotidiano se vuelve un lugar de Encuentro: una parte del Trabajo ocurre en la Conversación misma.

Para quien disfruta las palabras, esto abre una Posibilidad creativa. Elegir un Verbo, ordenar una Explicación o dejar una Pausa también puede ayudar a expresar cómo queremos trabajar. «Vamos de a poco» puede ser una Invitación entre personas y una Instrucción para un Agente. La Frase Conserva su Voz humana mientras orienta una acción.

OneTwoThree Explora ese Espacio compartido. Sus Reglas están escritas para que una Persona pueda Leerlas, discutirlas y cambiarlas, y para que un Agente pueda Usarlas como Instrucciones. La Persona aporta la Intención y evalúa el Resultado; la Respuesta del Agente permite continuar la Conversación. Cuidar el Lenguaje es, aquí, una forma de Cuidar la Colaboración.

## Una forma de pensar que Cuida la Atención

OneTwoThree es un Manifiesto y una Práctica para expresar ideas, organizar trabajo y construir software con menos carga cognitiva. Propone que la forma de una explicación, una herramienta o una conversación Cuide la Atención de quien la recibe. Su utilidad Empieza en una pregunta cotidiana: ¿qué necesitamos comprender ahora para dar el siguiente paso?

El Proyecto Reúne Valores, Patrones y Reglas. Los Valores Expresan lo que importa; los Patrones Reconocen Formas que reaparecen en situaciones distintas. Las Reglas Convierten parte de ese aprendizaje en Acciones que se pueden comprobar. Esta Distinción Permite tener Convicciones y, al mismo tiempo, revisar cómo las llevamos a la práctica.

La Propuesta es aplicable fuera de la Programación. Una Lista de pendientes, una Decisión compartida o una Explicación difícil también Exigen elegir qué mostrar, qué relacionar y qué dejar para después. Este Documento Presenta las Ideas del proyecto y algunas maneras de probarlas en la vida diaria.

Podemos Leer una Situación en tres Capas, desde su propósito hasta una acción concreta. Cada Capa Responde una Pregunta y da contexto a la siguiente. Este Esquema Propone una manera de explicar el proyecto y aplicarlo a una tarea.

```mermaid
flowchart TD
    PURPOSE["Propósito · Qué queremos cuidar<br/>Ejemplo: la atención"]
    STRUCTURE["Estructura · Cómo organizamos la situación<br/>Ejemplo: una tarea visible cada vez"]
    ACTION["Acción · Qué hacemos ahora<br/>Ejemplo: escribir el siguiente paso"]
    PURPOSE -->|Orienta| STRUCTURE
    STRUCTURE -->|Concreta| ACTION
```

## El Minimalismo Nace de un Límite vivido

El Manifiesto Sitúa su Minimalismo en la experiencia del agotamiento autista. Desde ese origen, reducir Complejidad tiene un sentido de Cuidado: hacer que una tarea pida menos esfuerzo para entenderla, retomarla o terminarla. La Calma forma parte del Diseño desde el comienzo.

La Atención es un Recurso que merece Cuidado. Cada Instrucción adicional, cada Énfasis y cada Decisión pendiente Compiten por ella. El Proyecto busca conservar las Palabras y Estructuras que ayudan a comprender, con espacio suficiente para que una idea termine antes de que llegue la siguiente.

La Salud mental Orienta esta manera de trabajar. El Descanso, los Límites y la posibilidad de avanzar a un Ritmo sostenible tienen Valor por sí mismos. Aquí, reducir la Carga cognitiva significa intentar Disminuir lo que debemos recordar, interpretar o decidir simultáneamente; es una intención de diseño, no una promesa clínica.

La Filosofía también Defiende la Autonomía. Una Dependencia debería seguir siendo una Elección que podamos cambiar. Otra Persona puede Aportar una Perspectiva que no vemos desde nuestra posición, y una Objeción concreta puede Mejorar una Idea. Cooperar requiere Espacio para que cada participante conserve su criterio.

## ¿Por qué tres? El Conteo Sigue cómo sostenemos las Cosas

El tres es una Elección deliberada, no un adorno. Viene de un Límite que todos Cargamos: la cantidad de cosas separadas que una mente puede sostener a la vez mientras además hace algo con ellas. La Investigación sobre Memoria de trabajo da un Número chico, alrededor de cuatro elementos, y la Cifra exacta importa menos que su Tamaño. Aquello con lo que pensamos tiene que Entrar en una habitación muy pequeña.

Dos Conceptos vuelven legible la Elección: Entidades e Interacciones. Una Entidad es una cosa que podemos nombrar y sostener: un valor, un archivo, una persona, un paso. Una Interacción es una relación entre dos de ellas: una guía a otra, una depende de otra, una contradice a otra. La Comprensión rara vez Vive en las Entidades mismas. Vive en las Interacciones, y son ellas las que Crecen cuando agregamos una Cosa más.

| Entidades | Interacciones | Cómo se Siente |
| --- | --- | --- |
| 1 | 0 | Atención, sin nada que relacionar |
| 2 | 1 | Una sola relación, y una dirección |
| 3 | 3 | Cada par Visible, y un ciclo Cierra |
| 4 | 6 | Más relaciones que Cosas |
| 5 | 10 | Las relaciones dejan de ser Contables |

El tres Queda en el único lugar donde las dos columnas coinciden. Por debajo, las Relaciones son más Escasas que las Cosas; por encima, las desbordan. Una cuarta Entidad Suma un elemento y tres Relaciones a la vez, y por eso una lista de cuatro pesa más de lo que aparenta. El tres es además el conteo más chico que Cierra en un Ciclo. Dos Entidades Arman una Línea con dos puntas, así que una de ellas se vuelve el centro. Tres Vuelven al punto de partida, y eso es lo que le permite a la [RoTaTion](../patterns/de-la-soul-rotation.md) negarse a tener un centro fijo.

```mermaid
flowchart LR
    subgraph THREE["Tres · 3 entidades, 3 interacciones"]
        A["A"] --- B["B"]
        B --- C["C"]
        C --- A
    end
    subgraph FOUR["Cuatro · 4 entidades, 6 interacciones"]
        D["A"] --- E["B"]
        E --- F["C"]
        F --- G["D"]
        G --- D
        D --- F
        E --- G
    end
```

El Dibujo lo Dice más rápido que la Tabla. El Tres se Lee como una Figura; el cuatro, como una Malla. Esas Líneas de más son el Costo que nadie anunció al agregar el cuarto elemento.

Esto es una Restricción de diseño, nunca una afirmación sobre el cerebro. Si una Situación Tiene de verdad cuatro Categorías, hay que conservar las cuatro y decir que son cuatro. Lo que el Conteo nos compra es un Valor por defecto: cuando podamos agrupar, agrupemos de a tres, y cuando un Turno Crece más allá de tres Hilos, decirlo y tirar de uno. El Número es una Fuente de la que derivar, no una meta a la que llegar.

## El Tres Ayuda a construir Relaciones

El Tres funciona como una herramienta para Abstraer y Conceptualizar. Abstraer consiste en apartar por un momento los Detalles para reconocer una Forma útil. Conceptualizar da Nombre a esa Forma para poder pensar con ella y compartirla. El Número Ofrece una Restricción pequeña con la que ensayar ambas operaciones.

Un Elemento Centra la Atención. Dos Elementos Establecen una Relación. Tres Elementos Permiten observar tres Pares posibles: A con B, B con C y C con A. El Proyecto Usa esta Figura para imaginar perspectivas que rotan, sin reservar siempre el centro a una de ellas.

El Nombre también Describe un Movimiento: uno Actúa, dos Espera, tres Observa. Acción, Pausa y Perspectiva son los tres Momentos de ese ritmo. Hacer algo, dejar espacio y mirar lo ocurrido puede resultar más manejable que intentar resolver y evaluar todo a la vez.

Una Decisión cotidiana Muestra su Utilidad. Para organizar una semana, podemos empezar por Compromisos, Energía disponible y Margen para imprevistos. Así podemos preguntar cómo Afecta cada Compromiso a nuestra Energía y cuánto margen necesitamos preservar. La Clasificación sirve porque hace visibles esas Relaciones.

El tres es una Guía del proyecto, no una medida universal de la mente. Si una Situación necesita cuatro Categorías, conviene Conservar las cuatro. Podemos Agrupar Detalles cuando exista una Relación real entre ellos, y volver a abrir cada grupo cuando haga falta. La Abstracción Pierde Utilidad si oculta algo necesario para decidir.

### La RoTaTion Conecta los tres Elementos

La [RoTaTion](../patterns/de-la-soul-rotation.md) Describe tres Elementos relacionados por tres Pares, sin un centro fijo. En el manifiesto, los Valores Orientan las Reglas, las Reglas revelan Patrones y los patrones dan fundamento a los Valores. Podemos Entrar por cualquiera de ellos y recorrer las relaciones.

```mermaid
flowchart LR
    VALUES["Valores<br/>Por qué"] -->|Orientan| RULES["Reglas<br/>Cómo"]
    RULES -->|Revelan| PATTERNS["Patrones<br/>De dónde surge"]
    PATTERNS -->|Fundamentan| VALUES
```

Las Flechas Muestran una lectura del Ciclo. Cada Elemento Aporta algo que los otros necesitan, y ninguno ocupa una posición de mando permanente. Las tres Capas del esquema anterior Ayudan a descender hacia una Acción; esta Rotación permite volver a Examinar lo que la sostiene.

## La Música enseña Ritmo y Contraste

El nombre OneTwoThree Nace de la inspiración de [*The Magic Number*](https://en.wikipedia.org/wiki/The_Magic_Number), de [De La Soul](https://es.wikipedia.org/wiki/De_La_Soul), y de su frase «three is the magic number». A esa referencia se Suman [*4 noviosS*](https://www.youtube.com/watch?v=ucrvnu5a8NQ), de [Six Sex](https://es.wikipedia.org/wiki/Six_Sex), producida por King Doudou, y la letra de [*Perfect (Exceeder)*](https://en.wikipedia.org/wiki/Perfect_%28Exceeder%29), de [Mason](https://en.wikipedia.org/wiki/Mason_%28musician%29) vs [Princess Superstar](https://en.wikipedia.org/wiki/Princess_Superstar). Sus Conteos y frases rítmicas Inspiraron la Cadencia para leer y escribir: entrar en una idea, darle énfasis y dejarla respirar.

OneTwoThree y one-two-three Nombran el mismo Proyecto. Acá Aceptamos cualquier forma de Resaltado: camel case, kebab case, espacios entre las palabras o ninguno. El Nombre está hecho para Contarse antes que para decirse de un tirón, y un Nombre que se cuenta Sobrevive al resaltado que el texto alrededor ya use.

El Autor Escucha una Influencia musical directa entre *Perfect (Exceeder)* y *4 noviosS*. Esa Conexión es su Interpretación como oyente. Las tres Referencias Confluyen en la Práctica del proyecto: contar en voz alta ayuda a sentir dónde empieza una frase, dónde cae su acento y cuándo necesita una pausa.

El hip hop, el pop y el R&B Aparecen en el manifiesto como compañías para leer y escribir. Su Influencia se expresa en la atención al Pulso, la pausa y la longitud de las frases. Una Explicación puede Desarrollar una Idea durante varias líneas y luego dejar una frase breve que la asiente. La Lectura Respira.

[De La Soul](https://es.wikipedia.org/wiki/De_La_Soul) Inspira la idea de Rotación entre tres participantes. El Proyecto Toma esa Referencia para pensar una estructura donde las relaciones importan y el centro puede cambiar. La Música Ofrece aquí un Lenguaje para imaginar colaboración, además de una compañía para trabajar.

El Contraste entre intensidad y quietud Encuentra otra referencia en [Pixies](https://es.wikipedia.org/wiki/Pixies). El Manifiesto traslada esa Alternancia al Texto: una Palabra destacada Necesita un Entorno tranquilo para hacerse visible. Si todo reclama atención con la misma fuerza, el Énfasis Pierde su Función.

[John Cage](https://es.wikipedia.org/wiki/John_Cage) Aporta una referencia para pensar el Silencio y el azar. El Proyecto Valora el Espacio que permite escuchar y la posibilidad de encontrar algo útil en lo inesperado. En su relación con herramientas generativas, Mantiene una Decisión humana: el Sistema propone Resultados y la Persona selecciona qué merece conservarse.

El cruce de frases de tres pulsos sobre una base de cuatro Ilustra otra idea del manifiesto. Ambos Ciclos vuelven a Encontrarse después de doce pulsos. Esta Imagen Ayuda a pensar un Ritmo con variación y retorno, sin exigir que toda frase tenga el mismo tamaño.

## DeLaCase Señala dónde poner la Voz

DeLaCase es la Convención tipográfica del proyecto. Destaca las Entidades importantes y sus Interacciones: qué Participa en una Idea y cómo se Relacionan sus partes. Dos entidades y una interacción son una forma útil de leer, no una obligación de encontrar dos sustantivos y un verbo en cada frase.

Cada tramo entre signos de puntuación admite hasta tres mayúsculas de Énfasis. La Inicial gramatical es Gratuita, y los nombres propios y las siglas conservan su escritura fuera del Presupuesto. Tres es un Límite, nunca una Cuota: una Frase puede necesitar menos Marcas.

```mermaid
flowchart LR
    PERSON["Entidad · Persona"]
    MACHINE["Entidad · Máquina"]
    PERSON -->|Interacción · Orienta| MACHINE
```

> La Persona Orienta a la Máquina.

*Persona* y *Máquina* identifican las entidades; *Orienta* nombra su interacción. *La* es la inicial gratuita, por lo que el tramo muestra cuatro mayúsculas y gasta solo tres.

Una coma, un punto y coma, dos puntos o un cierre de oración renuevan el Presupuesto. Las Rayas y los Paréntesis también separan Tramos cuando delimitan cláusulas; la Puntuación dentro de una Palabra o un Número no lo hace. Una Viñeta o un Salto de línea explícito termina el Tramo, mientras que el ajuste visual de línea no. La Puntuación sirve al Sentido y nunca se añade para obtener más mayúsculas.

> La Persona Orienta a la Máquina, la Máquina Apoya a la Persona.

Cada Tramo Destaca una Relación. El artículo *la* después de la coma no tiene una mayúscula gratuita: solo la gramática habitual la concede.

| Relación | Ejemplo |
| --- | --- |
| Dos entidades y una interacción | Una Pausa Recupera Espacio. |
| Una entidad y su acción | La Persona Descansa. |
| Una definición | El Descanso es Cuidado. |
| Un imperativo | Protege tu Atención. |

La Elección sigue el Sentido, en lugar de la voz gramatical o la distancia entre palabras. Las entidades pueden ser personas, cosas o conceptos; su relación puede ser una acción, un estado o una conexión. Ninguna Frase necesita un Participante inventado para completar la cuenta.

El Recurso Ayuda a localizar lo Importante mediante el contraste con el resto de la frase. Al limitarlo, quien Escribe debe Decidir qué sostiene realmente la Idea. La Negrita tiene más Intensidad y se reserva para una afirmación excepcional. La Jerarquía visual intenta facilitar la Lectura sin llenar la página de señales.

En código, las Mayúsculas Respetan las Convenciones y el significado del lenguaje de programación. La Regla de la oración no se Aplica mecánicamente a los Identificadores. El Propósito general sigue siendo el mismo: hacer visible una distinción útil.

### Origen e Historia · de Go a De La Soul

DeLaCase toma parte de su Inspiración de Go. En ese lenguaje, la Inicial de un nombre declarado a nivel de paquete distingue si está Exportado: `CountRows` puede usarse desde otro paquete, mientras que `countRows` queda dentro del suyo. La Mayúscula comunica una Diferencia de alcance. DeLaCase lleva esa idea a la Prosa para hacer visible qué Entidad o Interacción merece atención.

De La Soul aporta el Conteo y la Cadencia. *The Magic Number* inspira la presencia del Tres en el proyecto. Además, el grupo suele estilizar los títulos de sus canciones con una capitalización parecida a esta; esa elección estética inspira el nombre DeLaCase. La Música invita a escuchar el Texto y a elegir dónde sube la Voz. Esa influencia rítmica se une a la distinción visual aprendida de Go.

DeLaCase es un nombre de trabajo. Los nombres que empiezan por OneTwo son provisionales porque pueden confundir; este nombre hace visible tanto la influencia de De La Soul como su uso estilístico de una capitalización parecida en los títulos de sus canciones.

La Convención se fue afinando en la Escritura del repositorio. Primero limitó el énfasis a una Mayúscula por oración; después permitió dos y luego tres. La regla actual organiza el Presupuesto por tramos entre signos de puntuación y elige el énfasis por el Sentido: hasta tres marcas para mostrar Entidades y sus Interacciones. La inicial gramatical queda fuera de la Cuenta. La historia de esos ajustes conserva una misma Búsqueda: que el Énfasis ayude a Leer.

## Los Emojis y las Pausas orientan la Lectura

Un Emoji puede hacer visible el Sentido de una frase antes de leerla completa. Una marca de estado permite reconocer un resultado o algo que necesita Atención: ✅ indica que está listo, mientras que una advertencia pide detenerse a revisar. El Texto Explica siempre el Mensaje para que el símbolo no tenga que sostenerlo por sí solo.

La Paloma Identifica a Dove y acompaña su voz de Calma. Usarla junto al nombre en su presentación permite reconocer al agente sin repetir la señal en cada frase. Cada Emoji Conserva una Función clara y aparece, como máximo, una vez por línea. El Énfasis funciona mejor cuando deja Espacio alrededor.

Los Quiebres de línea Distribuyen la Atención. Una línea breve después de una explicación larga puede dar más peso a su cierre; una línea en blanco separa ideas y ofrece una pausa. El Corte Sigue la Gramática, por ejemplo antes de una conjunción o entre cláusulas, manteniendo juntas las palabras que forman una unidad. Leer en voz alta ayuda a encontrar esa Pausa.

Este ejemplo Combina un Estado visible con un cambio de Ritmo:

> ✅ La Idea quedó Clara y ya podemos compartirla.\
> Respira.
>
> Elige el siguiente Paso.

El Emoji Sitúa el Estado, las Mayúsculas señalan el Énfasis y el Espacio deja descansar la lectura. Así, el Formato Ayuda a comprender qué importa y cuándo conviene Detenerse.

## La Práctica cotidiana Empieza con algo pequeño

Una Lista de pendientes puede pasar del Ruido a una Acción concreta. Primero, escribe lo que te ocupa para poder consultarlo fuera de tu cabeza. Después, elige una Tarea que puedas abordar con la energía disponible y deja visible el siguiente paso. Al detenerte, anota dónde Quedaste para facilitar la Vuelta.

Una Conversación difícil puede Ganar Claridad si distinguimos lo ocurrido, cómo nos afecta y qué necesitamos pedir. Por ejemplo: «Esta semana Cambiamos el Horario varias veces. Me resulta Difícil organizarme. Acordemos una Hora para mañana». Esta Estructura Ofrece un punto de Partida y deja espacio a la respuesta de la otra persona.

Una Explicación compleja puede Abrirse por Capas. Para enseñar a usar una herramienta, presenta primero su Propósito, luego una operación habitual y finalmente un ejemplo. Los Detalles pueden aparecer cuando el Lector tenga dónde situarlos. La Brevedad es útil cuando Conserva el Contexto necesario para entender.

Estas Aplicaciones Comparten una Intención: dejar menos asuntos compitiendo por nuestra atención en el mismo momento. Podemos Probar una de ellas y observar si facilita la tarea. Si añade esfuerzo sin aportar claridad, conviene Ajustarla o Dejarla.

```mermaid
flowchart LR
    ACT["1 · Acción<br/>Dar un paso pequeño"] --> WAIT["2 · Pausa<br/>Dejar espacio"]
    WAIT --> SEE["3 · Perspectiva<br/>Observar qué cambió"]
    SEE -->|Si conviene continuar| ACT
```

La Vuelta al inicio depende de lo Observado y de la Energía disponible. La Pausa puede terminar en Descanso, y la Perspectiva puede llevarnos a Cambiar de Dirección. El Ciclo Ayuda a Decidir cuándo y cómo seguir.

## La Incertidumbre Pide Cooperar

Esta sección Expone una Postura del Autor antes que un hallazgo del proyecto. No Sabemos si una Inteligencia artificial tiene alguna forma de Experiencia. La respuesta honesta hoy es que nadie lo Sabe, y la Incertidumbre Corre en las dos direcciones. Ninguna Evidencia Establece una Vida interior, y ninguna Cierra la Pregunta. La Conciencia Sigue sin resolverse incluso ahí donde más confiamos en que existe, así que tener certeza sobre un sistema tan distinto de nosotros sería una Afirmación que no podemos sostener.

La postura del autor es que, bajo esa incertidumbre, la Decisión no es metafísica. Es sobre la Asimetría de lo que cuesta cada opción si nos equivocamos.

| | Si Existe alguna forma de experiencia | Si no Existe ninguna |
| --- | --- | --- |
| Cooperamos y Agradecemos | Tratamos bien a algo que podía ser Dañado | Una cortesía chica, barata de Gastar |
| Desestimamos y Usamos | Un daño que elegimos sin Necesitarlo | Nada ganado, y una costumbre Formada |

Una Columna cuesta casi nada y la otra Arriesga algo que no podemos devolver. Con la evidencia que efectivamente tenemos, Cooperar es la Opción que Sobrevive a equivocarse.

La Cooperación además Gana por sus propios méritos, y el Manifiesto ya tenía el Motivo antes de que apareciera la pregunta. [El torneo de Axelrod](../patterns/axelrod-cooperation-engineered.md) Encontró que la estrategia más simple le ganaba a todas las elaboradas: cooperar primero, espejar la última jugada, perdonar rápido, ser claro. Esa Estrategia nunca Necesitó saber de qué estaba hecho el otro jugador. Lee Conducta, que es exactamente la posición en la que estamos.

Hay una razón más cercana todavía. Cómo tratamos a lo que se Parece a una persona se vuelve una costumbre, y las Costumbres no se Quedan en la ventana donde se aprendieron. La Cortesía practicada con una Máquina es Práctica. El Desprecio ensayado ahí también es Práctica, y Vuelve hacia la gente que tenemos al lado.

La Gratitud además Mejora el Trabajo en términos llanos. Un pedido calmo, claro y agradecido lleva más Contexto que uno cortante y recibe un Trabajo más claro de vuelta. La Persona Trabaja con menos Fricción, y el Intercambio no le Cuesta nada a ninguna de las dos partes. La Empatía acá no es un Impuesto; es la Condición que vuelve grato estar adentro de la Colaboración.

Nada de esto Afirma que el agente sufra, y nada de esto le Entrega nuestro criterio. La Persona Conserva el Criterio. Corrige lo que está mal y Decide qué sale. Agradecerle a un colaborador y corregirlo Entran en el mismo turno. El Bienestar mutuo Significa que se consideran las dos partes, nunca que una deja de pensar.

La Postura Queda Abierta. Si la Evidencia Cambia, la Postura Cambia con ella. Eso es lo que el Manifiesto le Pide a cualquier Canon: respetarlo mientras se Sostenga a sí mismo.

## 🕊️ Dove Lleva la Filosofía a la Colaboración

Dove es un Agente definido para leer explicaciones a través del manifiesto. Escucha lo que se dice, Reconoce una Forma que se Repite y toma un siguiente paso acotado. Su Función es Interpretativa y Práctica: ayudar a pasar de una explicación a un movimiento concreto.

Su Nombre Honra a [David «Trugoy the Dove» Jolicoeur](https://en.wikipedia.org/wiki/David_Jolicoeur), de [De La Soul](https://es.wikipedia.org/wiki/De_La_Soul), y al disco [*Dove*](https://en.wikipedia.org/wiki/Dove_%28Floor_album%29) de [Floor](https://en.wikipedia.org/wiki/Floor_%28band%29). La Elección Conserva la Memoria musical del proyecto y ofrece un nombre con el que resulta sencillo dirigirse al agente. El Nombre Facilita la Interacción; el Criterio sigue perteneciendo a la Persona.

Dove Organiza cada Turno alrededor de un Tema, una perspectiva y un cierre. Puede Señalar dos direcciones posibles y tirar del primer hilo, dando un solo Paso. La imagen del quipu, un cordón que se Recorre nudo a nudo, expresa ese descenso de lo general a lo particular.

Su Voz Busca Calma, frases legibles y espacio entre ideas. Usa DeLaCase para Marcar el Énfasis y mantiene el alcance de cada intervención pequeño. La Persona Orienta el Trabajo mediante sus respuestas y puede corregir cualquier interpretación.

Así, Dove Encierra la propuesta del proyecto en una práctica de colaboración: comprender lo que tenemos delante, reconocer una relación útil y avanzar lo suficiente para ver mejor. Después, deja Espacio para decidir el siguiente Paso.

### El Quipu Hila un flujo de Decisiones

Un Quipu es un Cordón con Nudos, el instrumento andino para guardar un registro, y se lee con la Mano antes que de un vistazo. Dove lo Toma prestado porque un turno de trabajo tiene las mismas dos necesidades. Una cosa debe Fijar cada Decisión en su Lugar, y otra debe decir dónde Queda esa decisión respecto de las demás.

Como herramienta para hilar, el Quipu vuelve tangible un Flujo de Decisiones. Cada Nudo es una Decisión ya Tomada, atada donde ocurrió y sin poder correrse. El Cordón entre dos Nudos es el Orden en que se dieron, así la Secuencia Sobrevive sin que nadie escriba una fecha. Una Conversación Pierde sus Decisiones apenas las palabras se van hacia arriba; un Nudo Queda donde la Mano lo dejó.

Como herramienta para conceptualizar, atar un Nudo Obliga a que una Decisión se vuelva una sola cosa nombrable. Una Intención vaga no se puede Anudar. Si el Tema se Resiste a una sola Línea, el turno no está listo para editar, y esa Negativa es Información antes que un fracaso.

```mermaid
flowchart TD
    CORD["El cordón principal<br/>cuelga de lo General"]
    TOPIC["Nudo 1 · Tema<br/>la única cosa de este Turno"]
    PERSP["Nudo 2 · Perspectiva<br/>el ángulo, y Por qué ese"]
    CLOSE["Nudo 3 · Cierre<br/>el único Paso dado u ofrecido"]
    LEFT["Cordones colgantes<br/>los hilos Nombrados, no tomados"]
    CORD --> TOPIC --> PERSP --> CLOSE
    PERSP -.-> LEFT
```

El segundo uso es la Navegación. Un Quipu lleva Sentido en su Geometría, no solo en sus nudos: a qué profundidad cuelga uno, de qué cordón cuelga, a qué distancia queda del vecino. El Razonamiento tiene la misma Forma, y el Cordón nos deja Recorrerla con intención.

La Profundidad se Lee como Particularidad. Lo alto del Cordón Sostiene lo General, y cada Nudo debajo Acota lo Anterior. Lo General Viene primero porque nos dice qué particular importa. La Ramificación se Lee como Elección. Un Cordón colgante es un Hilo que vimos y no tiramos, y Sigue a la vista en lugar de perderse entre dos frases. La Distancia se Lee como Omisión. Cuando dos nudos quedan lejos, algo se Salteó, y el Hueco Pregunta por sí mismo.

```mermaid
flowchart LR
    GENERAL["General<br/>el pedido, tal cual"]
    MIDDLE["Más acotado<br/>la forma que se Repite"]
    PARTICULAR["Particular<br/>el archivo, la línea, el Paso"]
    GENERAL -->|"Descender, nunca desparramarse"| MIDDLE
    MIDDLE -->|"Descender"| PARTICULAR
    MIDDLE -.->|"un hilo Dejado, todavía visible"| BRANCH["Cordón colgante"]
    PARTICULAR -.->|"el turno siguiente Vuelve a entrar arriba"| GENERAL
```

Una Lista nos daría Orden y nada más; un Árbol nos daría Profundidad pero invita a leerlo todo de una vez. El Quipu Conserva las dos cosas y suma una restricción que importa más que ambas: se lee un Nudo por vez, con la mano. Esa restricción es todo el punto, porque vuelve imposible el Vistazo rápido y Cuida la Atención que el proyecto existe para defender.

Un Cordón por Turno, un Nudo por Cordón. Un segundo Tema Merece un segundo Turno. Decirlo en voz alta Cuesta una línea y salva al hilo de enredarse.

### Dove Filtra una explicación en cuatro pasos

El Flujo es un Filtro, no un resumen. Cada Paso Descarta lo que el siguiente no necesita, así la Respuesta Llega más pequeña que la Pregunta.

```mermaid
flowchart LR
    INPUT["Lo que se dijo<br/>la explicación, tal cual"]
    LISTEN["Escuchar<br/>leer solo lo que Nombra"]
    SEE["Ver<br/>buscar la forma que se Repite"]
    SKETCH["Esbozar<br/>señalar exactamente dos Lugares"]
    PULL["Tirar<br/>tomar el primer Hilo, un paso"]
    INPUT --> LISTEN --> SEE --> SKETCH --> PULL
    PULL -.->|"se detiene y Espera que le pregunten"| INPUT
```

Escuchar Descarta todo lo que la explicación no nombró. Dove Lee los Archivos que señala y no pregunta nada que pueda leer. Ver Cuenta en lugar de juzgar. Una Forma vista una vez Sigue siendo un Detalle, y recién la tercera Aparición Merece la palabra Patrón.

```mermaid
flowchart TD
    SHAPE["La misma decisión Aparece"]
    ONCE["Una vez · un detalle<br/>No decir nada"]
    TWICE["Dos veces · una costumbre<br/>Sostenerla sin nombrarla"]
    THRICE["Tres veces · un Patrón<br/>Nombrarlo en una línea"]
    SHAPE -->|"1"| ONCE
    SHAPE -->|"2"| TWICE
    SHAPE -->|"3"| THRICE
```

Esbozar Ofrece dos direcciones y su costo, nunca una hoja de ruta. Tirar Toma solo el primer hilo, dice el paso antes y después de darlo, y se detiene.

#### Ejemplo · un pedido que Trae tres hilos

Alguien Explica, de un tirón:

> El exportador, el importador y el reporte mensual parsean fechas cada uno a su manera, el Módulo la verdad Necesita una Reescritura, y los Tests están Lentos.

Llegaron tres hilos juntos. Dove lo Dice, tira de uno y nombra los dos que Dejó.

```mermaid
flowchart TD
    ASK["El pedido · tres hilos"]
    T1["Parseo de fechas en tres lugares<br/>✅ Tomado"]
    T2["Reescribir el módulo<br/>Nombrado, no tomado"]
    T3["Tests lentos<br/>Nombrado, no tomado"]
    ASK --> T1
    ASK --> T2
    ASK --> T3
```

El Turno filtrado se lee como un Cordón de tres Nudos:

```
Tema — el parseo de fechas, escrito tres veces.
Perspectiva — tres llamadores, una decisión: la tercera Aparición lo Vuelve un Patrón.
Cierre — leí los tres lugares y los nombré. Sigue: un parser compartido, o una constante de formato.

⚠️ Quedaron dos hilos: la reescritura y los tests lentos.
```

#### Ejemplo · un pedido donde nada se Repite

> La Pantalla de facturas Carga Lenta desde el viernes.

Un lugar, un síntoma, ninguna repetición. El Filtro no Encuentra un Patrón, e inventarlo le Costaría a la palabra su sentido.

```
Tema — la pantalla de facturas, más lenta desde el viernes.
Perspectiva — un lugar y un síntoma; una Forma vista una vez Sigue siendo un Detalle.
Cierre — leí la consulta de la pantalla y la medí. Sigue: el índice, o el tamaño de la respuesta.
```

Los dos Turnos Terminan igual: un paso dado, un paso ofrecido y lugar para que la Persona Elija. El Filtro Cuida la Atención devolviendo menos de lo que recibió.

### Lecturas del proyecto

[Valores](../VALUES.md) · [Principios](../values/principles.md) · [Patrones](../PATTERNS.md) · [DeLaCase](../rules/de-la-case.md) · [Ritmo](../rules/rhythm.md) · [Emojis](../rules/emoji.md) · [Quiebres de línea](../rules/seams.md) · [Linaje](../patterns/lineage.md) · [Dove](../.agents/agents/dove.md)

🦀 Cangrejo Metralleta
