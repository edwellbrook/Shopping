#include <stdio.h>
#include "mbed.h"
#include "bluetooth.h"
#include "nfc.h"
#include "display.h"

I2C i2c(I2C_SDA0, I2C_SCL0);
Serial host(USBTX, USBRX);

// hardware buttons
DigitalIn joy_up(P0_28);
DigitalIn joy_down(P0_29);
DigitalIn joy_in(P0_15);
InterruptIn helpButton(BUTTON1);


volatile int state = 0;
volatile int ready = 0;
volatile int authorised = -1;
volatile bool inHelp = false;
char cardId[7];
char items[MAX_ITEMS][FRAME_WIDTH + 1];

bool auth(uint8_t uid[7]) {
    host.printf("AUTH:%s\r\n", uid);

    // block until auth completes
    while (authorised == -1) {}

    bool authed = authorised;
    authorised = -1;


    if (authed) {
        sprintf(cardId, "%s", uid);
    }

    return authed;
}

void sendBeacons(const uint8_t *beacon) {
    char beaconId[Gap::ADDR_LEN];
    memcpy(beaconId, beacon, Gap::ADDR_LEN);

    host.printf("HELP:%s:%s\r\n", cardId, beaconId);
}

void host_writeln(const char *message) {
    host.printf("%s\r\n", message);
    wait_ms(500); // wait for write to complete
}

void load_list() {
    int i = 0;
    // load in each item
    while(i < MAX_ITEMS) {
        int j = 0;
        // load in each character of the item name
        while (j < FRAME_WIDTH + 1) {
            char c = host.getc();
            items[i][j] = c;

            // end list item with null terminator
            if (c == 0) {
                break;
            }
            j += 1;
        }
        i += 1;
    }
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
        load_list();
        ready = 1;
    }

    __enable_irq();
}

void requestHelp() {
    inHelp = !inHelp;

    state = inHelp ? 3 : 2;

    if (state == 2) {
        display_update();
    }
}

int main() {

    // set up bluetooth scanner
    ble_setup(sendBeacons);
    // interrupt on serial data
    host.attach(&serialInterrupt);
    // interrupt when help button is pressed
    helpButton.rise(&requestHelp);

    while (true) {
        switch (state) {

            case 0: { // scanning for card
                host_writeln("INFO:Scanning for NFC card");
                display_message("PLEASE SCAN YOUR CARD");
                nfc_start(i2c, auth);

                state += 1;
                host_writeln("INFO:NFC card found and authorised");
                break;
            }

            case 1: { // loading shopping list
                display_message("LOADING SHOPPING LIST");
                host.printf("LIST:%s\r\n", cardId);
                while (!ready) {}

                state += 1;
                shopping_list_start(items);
                break;
            }

            case 2: { // shopping
                if (!!joy_up) {
                    display_cursor_up();
                    display_update();
                } else if (!!joy_down) {
                    display_cursor_down();
                    display_update();
                } else if (!!joy_in) {
                    display_toggle_item();
                    display_update();
                }
                break;
            }

            case 3: {
                host_writeln("INFO:Scanning for beacons");
                display_message("REQUESTING HELP");
                ble_ping();
                break;
            }
        }
    }
}
