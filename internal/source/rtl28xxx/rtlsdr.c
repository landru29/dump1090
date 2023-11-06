#include "rtlsdr.h"
#include <math.h>
#include <malloc.h>

#define REMAINING_BUFFER_SIZE 112 * 2

void rtlsdrCallback(unsigned char *buf, uint32_t len, void *ctx) {
    unsigned char message[14];
    int cursor = 0;

    context *currentCtx = (context*)ctx;

    unsigned char* data = (unsigned char*)malloc(len/2 + currentCtx->remainingLength);

    if ((currentCtx->remainingLength>0) && (currentCtx->remainingData != 0)) {
        memcpy(data, currentCtx->remainingData, currentCtx->remainingLength);

        cursor = currentCtx->remainingLength;
    }

    for(int idx = 0; idx<len/2; idx++) {
        int i = buf[idx*2]-127;
        int q = buf[idx*2+1]-127;

        if (i>127) {
            i = i + (i ^ 0xff);
        }

        if (q>127) {
            q = q + (q ^ 0xff);
        }

        data[idx+cursor] = magnitude[i*129+q];
    }

    //       |   |         |   |
    //       |   |         |   |
    //       |   |         |   |
    //       |   |         |   |
    //       | | | | | | | | | | | | | | | |
    //       0 1 2 3 4 5 6 7 8 9 10

    int idx = 0;
    for(int idx = 0; idx<len/2 - REMAINING_BUFFER_SIZE; idx++)  {
        if (data[idx] < data[idx+1]) {
            continue;
        }

        if (data[idx+1] > data[idx+2]) {
            continue;
        }

        if (data[idx+2] < data[idx+3]) {
            continue;
        }
        if (data[idx+2] < data[idx+4]) {
            continue;
        }
        if (data[idx+2] < data[idx+5]) {
            continue;
        }
        if (data[idx+2] < data[idx+6]) {
            continue;
        }

        if (data[idx+6] > data[idx+7]) {
            continue;
        }

        if (data[idx+7] < data[idx+8]) {
            continue;
        }

        if (data[idx+8] > data[idx+9]) {
            continue;
        }

        if (data[idx+9] < data[idx+10]) {
            continue;
        }
        if (data[idx+9] < data[idx+11]) {
            continue;
        }
        if (data[idx+9] < data[idx+12]) {
            continue;
        }
        if (data[idx+9] < data[idx+13]) {
            continue;
        }
        if (data[idx+9] < data[idx+14]) {
            continue;
        }
        if (data[idx+9] < data[idx+15]) {
            continue;
        }

        memset(message, 0, 14);

        int startOfMessage = idx+16;
        int messageLength = 56;

        for(int index=0; index<messageLength; index++) {
            int byteIndex = index / 8;
            int bitIndex = index % 8;

            if ((index==0) && (data[startOfMessage+index*2] > data[startOfMessage+index*2+1])) {
                messageLength == 112;
            }

            unsigned char bit = (data[startOfMessage+index*2] > data[startOfMessage+index*2+1]);

            message[byteIndex] |= bit << bitIndex;
        }

        goRtlsrdData(message, messageLength / 8, ctx);
    }

    if (currentCtx->remainingData != 0) {
        currentCtx->remainingLength = REMAINING_BUFFER_SIZE;
        memcpy(currentCtx->remainingData, &data[len/2 - REMAINING_BUFFER_SIZE], REMAINING_BUFFER_SIZE);
    }
}

context *newContext(void* goContext) {
    context output;
    output.goContext = goContext;
    output.remainingData = (unsigned char*)malloc(REMAINING_BUFFER_SIZE);
    output.remainingLength = 0;

    return &output;
}

int rtlsdrReadAsync(rtlsdr_dev_t *dev, void *ctx, uint32_t buf_num, uint32_t buf_len) {
    return rtlsdr_read_async(dev, rtlsdrCallback, ctx, buf_num, buf_len);
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
