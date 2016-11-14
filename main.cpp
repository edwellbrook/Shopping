#include "mbed.h"

#include "ble/BLE.h"
#include "ble/Gap.h"
#include "ble/GapAdvertisingParams.h"

#include "PN532.h"
#include "PN532_I2C.h"

DigitalOut led(LED1);
Serial     host(USBTX, USBRX);
I2C        i2c(I2C_SDA0, I2C_SCL0);

PN532_I2C pn532(i2c);
PN532     nfc(pn532);

void ticker_callback() {
    led = !led;
}

void ble_scan(const Gap::AdvertisementCallbackParams_t *params) {
    if (params->type != GapAdvertisingParams::ADV_NON_CONNECTABLE_UNDIRECTED) {
        return;
    }
    
    for (int i = Gap::ADDR_LEN - 1; i>=0; i--) {
        host.printf ("%2.2x", params->peerAddr[i]);
    }
    
    host.printf(" Got advertisement (%i) with rssi: %i\r\n", params->advertisingDataLen, params->rssi);
}

void ble_init(BLE::InitializationCompleteCallbackContext *params) {
    BLE &ble = params->ble;
    
    ble.gap().setScanParams(400, 400);
    ble.gap().startScan(ble_scan);
}

int main() {
    Ticker ticker;
    ticker.attach(ticker_callback, 1);

    // set up bluetooth
    BLE &ble = BLE::Instance(BLE::DEFAULT_INSTANCE);
    ble.init(ble_init);
    
    // set up nfc
    nfc.begin();
    nfc.SAMConfig();
    
    host.printf("Application started\r\n");
    
    while(true) {
        ble.processEvents();
        
        bool success;
        uint8_t uid[] = { 0, 0, 0, 0, 0, 0, 0 };
        uint8_t uidLength;
        
        success = nfc.readPassiveTargetID(PN532_MIFARE_ISO14443A, &uid[0], &uidLength);
  
        if (success) {
            host.printf("Found a card!\r\n");
            host.printf("UID Length: %d bytes\r\n", uidLength);
            
            host.printf("UID Value: 0x00");
            for (uint8_t i=0; i < uidLength; i++) {
                host.printf("%02X", uid[i]);
            }
            host.printf("\r\n");
        }
    }
}
