#ifndef _BLUETOOTH_H_
#define _BLUETOOTH_H_

#include "ble/Gap.h"

typedef void (*SendBeaconFn)(const uint8_t*);

void ble_setup(SendBeaconFn fn);
void ble_ping();

#endif // !_BLUETOOTH_H_
