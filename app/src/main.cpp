#include <stdio.h>
#include "mbed.h"
#include "bluetooth.h"
#include "nfc.h"
#include "display.h"

I2C i2c(I2C_SDA0, I2C_SCL0);
Serial host(USBTX, USBRX);

InterruptIn helpButton(BUTTON1);

volatile int ready = 0;
volatile int authorised = -1;
char items[MAX_ITEMS][FRAME_WIDTH + 1];

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
    }
    // auth response, declined
    else if (strncmp(str, "AUTH0", 5) == 0) {
        authorised = 0;
    }
    // auth response, accepted
    else if (strncmp(str, "AUTH1", 5) == 0) {
        authorised = 1;
    }
    // load shopping list
    else if (strncmp(str, "LLOAD", 5) == 0) {
        host_writeln("loading shit:");

        int i = 0;
        // load in each item
        while(i < MAX_ITEMS) {
            int j = 0;
            // load in each character of the item name
            while (j < FRAME_WIDTH + 1) {
                char c = host.getc();
                items[i][j] = c;

                host.printf("%c", c);

                // end list item with null terminator
                if (c == 0) {
                    host.printf("\r\n");
                    break;
                }

                j += 1;
            }

            i += 1;
        }

        // tell program the list is ready
        ready = 1;
    }

    __enable_irq();
}

void requestHelp() {
    host_writeln("HELP:-");
}

int main() {
    // interrupt on serial data
    host.attach(&serialInterrupt);

    // interrupt when help button is pressed
    helpButton.rise(&requestHelp);

    // host_writeln("INFO:Scanning for NFC card");
    // display_message("PLEASE SCAN YOUR CARD");
    // nfc_start(i2c, auth);
    // host_writeln("INFO:NFC card found and authorised");
    //
    // host_writeln("INFO:Scanning for beacons\r\n");
    // ble_start(sendBeacons);
    // host_writeln("INFO:Ending beacon scan");


    // wait until list is ready
    while (!ready) {}
    shopping_list_start(items);
}
