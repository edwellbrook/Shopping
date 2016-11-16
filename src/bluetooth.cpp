#include "bluetooth.h"

#include "ble/BLE.h"
#include "ble/Gap.h"

void scan(const Gap::AdvertisementCallbackParams_t *params) {
    if (params->type != GapAdvertisingParams::ADV_NON_CONNECTABLE_UNDIRECTED) {
        return;
    }

    for (int i = Gap::ADDR_LEN - 1; i >= 0; i--) {
        // host.printf ("%2.2x", params->peerAddr[i]);
    }

    // host.printf(" Got advertisement (%i) with rssi: %i\r\n", params->advertisingDataLen, params->rssi);
}

void init(BLE::InitializationCompleteCallbackContext *params) {
    BLE &ble = params->ble;

    ble.gap().setScanParams(400, 400);
    ble.gap().startScan(scan);
}

void ble_start() {
    BLE &ble = BLE::Instance();
    ble.init(init);

    while (true) {
        ble.processEvents();
    }
}
