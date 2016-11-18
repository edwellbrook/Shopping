#include "mbed.h"
#include "bluetooth.h"
#include "nfc.h"

DigitalOut led(LED1);
Serial     host(USBTX, USBRX);
I2C        i2c(I2C_SDA0, I2C_SCL0);

void tickerCallback() {
    led = !led;
}

bool auth(uint8_t uid[7]) {
    host.printf("AUTH:%s\r\n", uid);
    return host.getc() == 1;
}

int main() {
    Ticker ticker;
    ticker.attach(tickerCallback, 1);

    host.printf("INIT\r\n");

    while (true) {
        host.printf("INFO:Scanning for NFC card\r\n");
        nfc_start(i2c, auth);
        host.printf("INFO:NFC card found and authorised\r\n");

        host.printf("INFO:Scanning for beacons\r\n");
        ble_start();
        host.printf("INFO:Ending beacon scan\r\n");

        host.printf("INFO:Restarting system...\r\n");
    }
}
