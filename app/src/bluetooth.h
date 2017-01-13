#ifndef _BLUETOOTH_H_
#define _BLUETOOTH_H_

#include "ble/Gap.h"

typedef void (*SendBeaconsFn)(char[]);

void ble_start(SendBeaconsFn fn);

#endif // !_BLUETOOTH_H_
