# Ejemplo 04: Escalando y Fracasando

**De Cero a Ebiten: IA sin Magia Negra**
*(Capítulo 4 de 13: Capas Lineales a Mano)*

---

### Línea Argumental
Teníamos nuestro Perceptrón (`03`) funcionando limpiamente. ¿Qué pasa si queremos hacer Redes Profundas (Deep Learning)? Apilamos capas de neuronas usando lo que sabemos: structs y arrays planos de Go. Aquí es donde descubrimos por qué no construimos cohetes con destornilladores.

### ¿Qué hace este código?
- Intenta crear una clase `Network` que contiene múltiples `Layers` (Capas).
- Intenta hacer Forward-Pass sumando y multiplicando el estado manualmente a través de funciones Lineales y Rectificadores ReLU (ignorando números negativos).
- Aplana la matriz de pesos a un único hilo unidimensional para que sea "rápido" en Go, forzándonos a usar la matemática manual de mapeo: `weights[i*n_in+j]`.

### El Peligro
Este tipo de código causa que la gestión de tamaños bidimensionales en índices se ensucie de forma terrorífica. El programador se distrae de "Aprender" para centrarse en "No exceder el tamaño de Slice".

### Siguiente Paso ->
[05-IndexError](../05) - Las consecuencias catastroficas que esto provoca, lo que finalmente nos empujará a librerías de verdad.
