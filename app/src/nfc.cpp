#include "nfc.h"

#include "PN532.h"
#include "PN532_I2C.h"

void waitForTarget(PN532 nfc, AuthFn auth) {
    bool authorised = false;

    while (!authorised) {
        bool success;
        uint8_t uid[7];
        uint8_t uidLength;

        success = nfc.readPassiveTargetID(PN532_MIFARE_ISO14443A, &uid[0], &uidLength);

        if (success) {
            authorised = auth(uid);
        }
    }
}

void nfc_start(I2C i2c, AuthFn auth) {
    PN532_I2C pn532 = PN532_I2C(i2c);
    PN532 nfc = PN532(pn532);

    nfc.begin();
    nfc.SAMConfig();
    nfc.setPassiveActivationRetries(0x00);

    waitForTarget(nfc, auth);
}
