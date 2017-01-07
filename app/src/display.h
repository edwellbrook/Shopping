#ifndef DISPLAY_H
#define DISPLAY_H

#include "C12832.h"

static const int LINE_HEIGHT = 11;
static const int FRAME_WIDTH = 22;

void display_message(const char msg[]);
void shopping_list_start(const char items[][FRAME_WIDTH + 1], int count);

#endif
