#include "bluetooth.h"
#include "ble/BLE.h"
#include <stdio.h>

SendBeaconsFn sendBeacons;

void scan(const Gap::AdvertisementCallbackParams_t *params) {
    if (params->rssi < -40) {
        return;
    }

    char address[Gap::ADDR_LEN * 2];

    for (int i = Gap::ADDR_LEN - 1; i >= 0; i--) {
        sprintf(&address[i * 2], "%2.2x", params->peerAddr[i]);
    }

    sendBeacons(address);
}

void init(BLE::InitializationCompleteCallbackContext *params) {
    BLE &ble = params->ble;

    ble.gap().setScanParams(400, 400);
    ble.gap().startScan(scan);
}

void ble_start(SendBeaconsFn fn) {
    sendBeacons = fn;

    BLE &ble = BLE::Instance();
    ble.init(init);

    while (true) {
        ble.processEvents();
    }
}
