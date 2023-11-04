#include "rtlsdr.h"

int rtlsdrReadAsync(rtlsdr_dev_t *dev, void *ctx, uint32_t buf_num, uint32_t buf_len) {
    return rtlsdr_read_async(dev, goRtlsrdData, ctx, buf_num, buf_len);
}


