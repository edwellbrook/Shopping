#include "nfc.h"

#include "PN532.h"
#include "PN532_I2C.h"

bool authoriseCard() {
    return false;
}

void waitForTarget(PN532 nfc) {
    bool authorised = false;

    while (!authorised) {
        bool success;
        uint8_t uid[7];
        uint8_t uidLength;

        success = nfc.readPassiveTargetID(PN532_MIFARE_ISO14443A, &uid[0], &uidLength);

        if (success) {
            // check and set authorised value
            authorised = authoriseCard();

            // host.printf("Found a card!\r\n");
            // host.printf("UID Length: %d bytes\r\n", uidLength);
            //
            // host.printf("UID Value: 0x00");
            // for (uint8_t i = 0; i < uidLength; i++) {
            //     host.printf("%02X", uid[i]);
            // }
            // host.printf("\r\n");
        }
    }
}

void nfc_start(I2C i2c) {
    PN532_I2C pn532 = PN532_I2C(i2c);
    PN532 nfc = PN532(pn532);
    nfc.begin();
    nfc.SAMConfig();

    waitForTarget(nfc);
}
