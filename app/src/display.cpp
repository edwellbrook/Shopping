#include "display.h"

C12832 lcd(P0_23, P0_25, P0_24, P0_18, P0_22);

// cursor position
volatile int cursor = 0;

// list position
volatile int offset = 0;

// shopping list
char list[MAX_ITEMS][FRAME_WIDTH + 1];
volatile bool checklist[MAX_ITEMS] = {false};


void display_message(const char msg[]) {
    lcd.cls();
    lcd.locate(0, 0);
    lcd.printf(msg);
}

void display_update() {
    lcd.cls();

    for (int i = 0; i < 3; i++) {
        int row = LINE_HEIGHT * i;
        int item_idx = offset + i;

        lcd.locate(0, row);

        if (cursor == i) {
            lcd.printf("> %s %s", checklist[item_idx] ? "[X]" : "[ ]", list[item_idx]);
        } else {
            lcd.printf("  %s %s", checklist[item_idx] ? "[X]" : "[ ]", list[item_idx]);
        }
    }
}

void display_cursor_down() {
    if (cursor < 2) {
        cursor += 1;
    } else if (cursor == 2 && offset + 3 < MAX_ITEMS) {
        offset += 1;
    }
}

void display_cursor_up() {
    if (cursor > 0) {
        cursor -= 1;
    } else if (cursor == 0 && offset - 1 >= 0) {
        offset -= 1;
    }
}

void display_toggle_item() {
    checklist[cursor + offset] = !checklist[cursor + offset];
}

void shopping_list_start(const char items[MAX_ITEMS][FRAME_WIDTH + 1]) {
    int i = MAX_ITEMS;
    while (i--) {
        strncpy(list[i], items[i], FRAME_WIDTH + 1);
    }

    display_update();
}
