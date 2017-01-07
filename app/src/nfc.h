#include "mbed.h"
#include "display.h"

typedef bool (*AuthFn)(uint8_t[7]);

void nfc_start(I2C i2c, AuthFn auth);
