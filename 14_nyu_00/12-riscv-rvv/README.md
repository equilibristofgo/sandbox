# Ejemplo 12-riscv-rvv: Redes Neuronales sobre Silicio Libre (CanMV-K230 + RVV 1.0)

**Bloque 5: Del Cloud y la GPU al Edge AI en RISC-V con Go**

---

## 1. De AVX (x86_64) a RISC-V Vector Extension (`RVV 1.0`)

En el ejemplo `02/matrix_avx` y en `11-juego-final/nn/avx` usamos instrucciones SIMD de longitud fija (`Y0..Y7` de 256 bits en x86_64). Si cambias a AVX-512 o a otro procesador, tienes que reescribir el kernel.

El SoC **Kendryte K230** (placa **CanMV-K230**) monta dos núcleos RISC-V **T-Head C908**, donde el núcleo de cómputo (`CPU1`, hasta 1.6 GHz) incluye una unidad vectorial **RVV 1.0 de 128 bits** con filosofía **Vector-Length Agnostic (VLA)**:

| Característica | x86_64 AVX2 (`02/matrix_avx`) | RISC-V RVV 1.0 (`12-riscv-rvv`) |
| :--- | :--- | :--- |
| **Longitud de vector** | Fija en el opcode (`256 bits` / 8 floats) | Dinámica en tiempo de ejecución (`VSETVLI`) |
| **Gestión de colas** | Bucle principal + bucle residual escalar | El mismo bucle procesa el remanente automáticamente |
| **Compilación cruzada** | Compleja con CGO / librerías C | `CGO_ENABLED=0 GOOS=linux GOARCH=riscv64` nativo en Go |

---

## 2. Cómo ejecutar y cross-compilar

### Ejecución local en tu PC (`amd64` con fallback Go):
```bash
go run .
```

### Cross-compilación a binario estático para CanMV-K230 (`riscv64` + `RVV 1.0`):
```bash
CGO_ENABLED=0 GOOS=linux GOARCH=riscv64 GORISCV64=rva22u64 go build -ldflags="-s -w" -o neuron_k230 .
file neuron_k230
# neuron_k230: ELF 64-bit LSB executable, UCB RISC-V, RVC, double-float ABI, statically linked
```

---

## 3. Nota sobre Ensamblador Go 1.24 y RVV 1.0

`vector_riscv64.s` implementa `VectorAddRVV` (`VSETVLI`, `VLE32.V`, `VFADD.VV`, `VSE32.V`) y `DotProductRVV` (`VFMACC.VV`, `VFREDUSUM.VS`). Dado que Go 1.24 declara las constantes RVV pero aún no expone el codificador mnemónico en Plan 9 asm, las instrucciones vectoriales se emiten con `WORD $0x...` junto al mnemónico comentado, compilando limpiamente con `go build` estándar sin dependencias externas de C.
