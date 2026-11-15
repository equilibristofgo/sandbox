# Ejemplo 01: Concurrencia Pura

**De Cero a Ebiten: IA sin Magia Negra**
*(Capítulo 1 de 13: La Ingesta de Datos)*

---

### Línea Argumental
Antes de hablar de Machine Learning, predicciones o matrices exóticas, hay que hablar de **datos**. Muchos datos.
Go no fue concebido originalmente como un lenguaje para crear redes neuronales, pero brilla con luz propia cuando se trata del pre-procesamiento paralelo gracias a sus **Goroutines** y **WaitGroups**. Si vamos a alimentar al monstruo de la IA de NYU, lo haremos rápido.

### ¿Qué hace este código?
- Lanza 5 procesos "Worker" concurrentes simulando la ingesta o pre-procesamiento de bloques de datos para entrenamiento.
- Usa canales de sincronización básicos (`sync.WaitGroup`) para asegurar que todos los datos se han preparado antes de dar luz verde al inicio del entrenamiento.
- Dejaremos para mas adelante el uso de channels y CSP (y quizas alguna otra cosa que lleva en el tintero mucho tiempo "[single-flight](https://www.bytesizego.com/blog/single-flight)" ... mas info [aqui](https://go.dev/talks/2013/oscon-dl.slide#43) desde la documentacion de la [libreria](https://github.com/golang/groupcache?tab=readme-ov-file)).

### Siguiente Paso ->
[02-ManualTensor](../02) - Donde damos forma a los datos puros creando nuestra primera representación de un Tensor a mano con Slices base.
