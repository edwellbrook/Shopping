#ifndef _DISPLAY_H
#define _DISPLAY_H

#include "C12832.h"

static const int MAX_ITEMS = 12;
static const int LINE_HEIGHT = 11;
static const int FRAME_WIDTH = 22;

void display_message(const char msg[]);
void shopping_list_start(const char items[MAX_ITEMS][FRAME_WIDTH + 1]);

#endif // !_DISPLAY_H_
