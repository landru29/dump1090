#include "context.h"
#include "malloc.h"
#include "../rtl28xxx/constant.h"

context *newContext(void* goContext) {
    context *output = (context*)malloc(sizeof(context));
    output->goContext = goContext;
    output->remainingMagnitudeData = (uint16_t*)malloc(MAGNITUDE_LONG_MSG_SIZE * sizeof(uint16_t) + PREAMBULE_BIT_SIZE);
    output->remainingMagnitudeLengthByte = 0;

    printf("Allocating memory: %d\n", MAGNITUDE_LONG_MSG_SIZE * sizeof(uint16_t) + PREAMBULE_BIT_SIZE);

    return output;
}