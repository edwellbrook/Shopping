#include "ble/Gap.h"

typedef void (*SendBeaconsFn)(char[]);

void ble_start(SendBeaconsFn fn);
