#include "mbed.h"
#include "bluetooth.h"
#include "nfc.h"
#include "display.h"

Serial host(USBTX, USBRX);

DigitalOut led(LED1);
I2C        i2c(I2C_SDA0, I2C_SCL0);

void tickerCallback() {
    led = !led;
}

bool auth(uint8_t uid[7]) {
    host.printf("AUTH:%s\r\n", uid);
    return host.getc() == 1;
}

void sendBeacons(char beaconId[]) {
    host.printf("SCAN:");
    for (int i = 0; i < 12; i++) {
        host.printf("%c", beaconId[i]);
    }
    host.printf("\r\n");
}

int main() {
    Ticker ticker;
    ticker.attach(tickerCallback, 1);

    host.printf("INIT\r\n");

    while (true) {
        host.printf("INFO:Scanning for NFC card\r\n");
        display_message("PLEASE SCAN YOUR CARD");
        nfc_start(i2c, auth);
        host.printf("INFO:NFC card found and authorised\r\n");

        host.printf("INFO:Scanning for beacons\r\n");
        display_message("LOADING SHOPPING LIST");
        ble_start(sendBeacons);
        host.printf("INFO:Ending beacon scan\r\n");

        host.printf("INFO:Restarting system...\r\n");
    }

    // char items[][FRAME_WIDTH + 1] = {
    //     "CHEESE",
    //     "TUNA",
    //     "BACON",
    //     "SPAGHETTI",
    //     "BUTTER"
    // };
}
