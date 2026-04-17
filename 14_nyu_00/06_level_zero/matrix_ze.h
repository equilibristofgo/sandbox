#ifndef MATRIX_ZE_H
#define MATRIX_ZE_H

#ifdef __cplusplus
extern "C" {
#endif

void matrix_multiply_ze(const float* A, const float* B, float* C, int N);

#ifdef __cplusplus
}
#endif

#endif // MATRIX_ZE_H
