#include "display.h"

C12832 lcd(P0_23, P0_25, P0_24, P0_18, P0_22);

DigitalIn joy_up(P0_28);
DigitalIn joy_down(P0_29);
DigitalIn joy_in(P0_15);

void display_message(const char msg[]) {
    lcd.cls();
    lcd.locate(0, 0);
    lcd.printf(msg);
}

void update_display(int cursor, int offset, const char items[][FRAME_WIDTH + 1], bool checklist[]) {
    lcd.cls();

    for (int i = 0; i < 3; i++) {
        int row = LINE_HEIGHT * i;
        int item_idx = offset + i;

        lcd.locate(0, row);

        if (cursor == i) {
            lcd.printf("> %s %s", checklist[item_idx] ? "[X]" : "[ ]", items[item_idx]);
        } else {
            lcd.printf("  %s %s", checklist[item_idx] ? "[X]" : "[ ]", items[item_idx]);
        }
    }
}

void shopping_list_start(const char items[][FRAME_WIDTH + 1], int count) {
    bool check[count] = {false};

    int cursor = 0;
    int offset = 0;

    update_display(cursor, offset, items, check);

    while (true) {
        if (!!joy_down) {
            if (cursor < 2) {
                update_display(++cursor, offset, items, check);
            } else if (cursor == 2 && offset + 3 < count) {
                update_display(cursor, ++offset, items, check);
            }
        }

        if (!!joy_up) {
            if (cursor > 0) {
                update_display(--cursor, offset, items, check);
            } else if (cursor == 0 && offset - 1 >= 0) {
                update_display(cursor, --offset, items, check);
            }
        }

        if (!!joy_in) {
            // toggle item
            check[cursor + offset] = !check[cursor + offset];

            update_display(cursor, offset, items, check);
        }
    }
}
