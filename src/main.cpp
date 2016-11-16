#include "mbed.h"
#include "bluetooth.h"
#include "nfc.h"

DigitalOut led(LED1);
Serial     host(USBTX, USBRX);
I2C        i2c(I2C_SDA0, I2C_SCL0);

void tickerCallback() {
    led = !led;
}

int main() {
    Ticker ticker;
    ticker.attach(tickerCallback, 1);

    host.printf("Application started\r\n");

    nfc_start(i2c);
    ble_start();
}
