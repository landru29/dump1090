#ifndef __CONTEXT_H__
#define __CONTEXT_H__

#include <inttypes.h>

typedef struct {
    void *goContext;
    uint16_t *remainingMagnitudeData;
    uint32_t remainingMagnitudeLengthByte;
} context;

context *newContext(void* goContext);

#endif