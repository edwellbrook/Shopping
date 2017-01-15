#include "bluetooth.h"
#include "ble/BLE.h"
#include <stdio.h>

SendBeaconFn sendBeacon;
BLE *ble;

void scan(const Gap::AdvertisementCallbackParams_t *params) {
    if (params->rssi < -40) {
        return;
    }

    sendBeacon(params->peerAddr);
}

void init(BLE::InitializationCompleteCallbackContext *params) {
    BLE &ble = params->ble;

    ble.gap().setScanParams(400, 400);
    ble.gap().startScan(scan);
}

void ble_setup(SendBeaconFn fn) {
    sendBeacon = fn;

    ble = &BLE::Instance();
    ble->init(init);
}

void ble_ping() {
    ble->processEvents();
}
