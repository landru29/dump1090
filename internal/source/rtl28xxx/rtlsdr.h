#ifndef __RTL_SDR_H__
#define __RTL_SDR_H__

#include <stdlib.h>
#include <unistd.h>
#include "rtl-sdr.h"

#define MODES_LONG_MSG_BITS 112
#define MODES_SHORT_MSG_BITS 56

#define MODES_LONG_MSG_BYTES (MODES_LONG_MSG_BITS/8)
#define MODES_SHORT_MSG_BYTES (MODES_SHORT_MSG_BITS/8)


typedef struct {
    void *goContext;
    uint16_t *remainingMagnitudeData;
    uint32_t remainingMagnitudeLength;
} context;

extern int messageLengthBit[25];
extern uint16_t magnitude[129*129];

extern void goRtlsrdData(unsigned char *buf, uint32_t len, void *ctx);

int rtlsdrReadAsync(rtlsdr_dev_t *dev, void *ctx, uint32_t buf_num, uint32_t buf_len);
void rtlsdrProcessRaw(unsigned char *buf, uint32_t len, void *ctx);
void initTables();

context *newContext(void* goContext);

#endif