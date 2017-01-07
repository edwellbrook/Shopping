#include <stdio.h>
#include "mbed.h"
#include "bluetooth.h"
#include "nfc.h"
#include "display.h"

I2C i2c(I2C_SDA0, I2C_SCL0);
Serial host(USBTX, USBRX);

volatile int ready = 0;
volatile int authorised = -1;

bool auth(uint8_t uid[7]) {
    host.printf("AUTH:%s\r\n", uid);

    // block until auth completes
    while (authorised == -1) {}

    bool authed = authorised;
    authorised = -1;

    return authed;
}

void sendBeacons(char beaconId[]) {
    host.printf("SCAN:");
    for (int i = 0; i < 12; i++) {
        host.printf("%c", beaconId[i]);
    }
    host.printf("\r\n");
}

void host_writeln(const char *message) {
    host.printf("%s\r\n", message);
    wait_ms(500); // wait for write to complete
}

void serialInterrupt() {
    __disable_irq();

    int idx = 0;
    char str[5] = {0, 0, 0, 0, 0};

    while (idx < 5) {
        str[idx++] = host.getc();
    }

    // reset system
    if (strncmp(str, "RESET", 5) == 0) {
        NVIC_SystemReset();

    // auth response, failed auth
    } else if (strncmp(str, "AUTH0", 5) == 0) {
        authorised = 0;

    // auth response, succeeded auth
    } else if (strncmp(str, "AUTH1", 5) == 0) {
        authorised = 1;

    // host interface accepted handshake and is ready
    } else if (strncmp(str, "READY", 5) == 0) {
        ready = 1;

    // echo the data back
    } else {
        host.printf("ECHO:%s\r\n", str);
    }

    __enable_irq();
}

int main() {
    host.attach(&serialInterrupt);

    display_message("WAITING FOR HANDSHAKE");
    while (!ready) {}

    // host_writeln("INIT");
    host_writeln("INFO:Scanning for NFC card");
    display_message("PLEASE SCAN YOUR CARD");
    nfc_start(i2c, auth);
    host_writeln("INFO:NFC card found and authorised");

    host_writeln("INFO:Scanning for beacons\r\n");
    ble_start(sendBeacons);
    host_writeln("INFO:Ending beacon scan");

    // char items[][FRAME_WIDTH + 1] = {
    //     "CHEESE",
    //     "TUNA",
    //     "BACON",
    //     "SPAGHETTI",
    //     "BUTTER"
    // };
}
