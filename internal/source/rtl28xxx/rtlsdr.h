#ifndef __RTL_SDR_H__
#define __RTL_SDR_H__

#include <stdlib.h>
#include <unistd.h>
#include "rtl-sdr.h"

#define MODES_LONG_MSG_BITS 112
#define MODES_SHORT_MSG_BITS 56

#define MODES_LONG_MSG_BYTES (MODES_LONG_MSG_BITS/8)
#define MODES_SHORT_MSG_BYTES (MODES_SHORT_MSG_BITS/8)

extern void goRtlsrdData(unsigned char *buf, uint32_t len, void *ctx);
int rtlsdrReadAsync(rtlsdr_dev_t *dev, void *ctx, uint32_t buf_num, uint32_t buf_len);
void initTables();

int messageLengthBit[25];
uint16_t magnitude[129*129];

#endif