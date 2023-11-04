#ifndef __RTL_SDR_H__
#define __RTL_SDR_H__

#include <stdlib.h>
#include <unistd.h>
#include "rtl-sdr.h"

extern void goRtlsrdData(unsigned char *buf, uint32_t len, void *ctx);
int rtlsdrReadAsync(rtlsdr_dev_t *dev, void *ctx, uint32_t buf_num, uint32_t buf_len);

#endif