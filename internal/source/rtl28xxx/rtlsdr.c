#include "rtlsdr.h"
#include <math.h>

int rtlsdrReadAsync(rtlsdr_dev_t *dev, void *ctx, uint32_t buf_num, uint32_t buf_len) {
    return rtlsdr_read_async(dev, goRtlsrdData, ctx, buf_num, buf_len);
}

void initTables() {
    for(int idx = 0; idx<25; idx++) {
        messageLengthBit[idx] = idx<=11 ? MODES_SHORT_MSG_BITS : MODES_LONG_MSG_BITS;
    }

    for (int i = 0; i < 129; i++) {
        for (int q = 0; q < 129; q++) {
            magnitude[i*129+q] = round(sqrt(i*i+q*q)*360);
        }
    }
}
