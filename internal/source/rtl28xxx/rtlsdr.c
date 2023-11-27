#include "rtlsdr.h"
#include <math.h>
#include <malloc.h>
#include <string.h>
#include <stdio.h>

#define MAGNITUDE_ENCODED_BIT_SIZE             2                                                // 2 uint16_t
#define REMAINING_MAGNITUDE_BUFFER_SIZE        MODES_LONG_MSG_BITS * MAGNITUDE_ENCODED_BIT_SIZE

int messageLengthBit[25];
uint16_t magnitude[129*129];


void printValue(uint16_t val) {
    printf("[");
    for(int l = 0; l<val/512; l++) {
        printf("=");
    }
    printf(">\n");
}

void rtlsdrProcessRaw(unsigned char *byteBuffer, uint32_t byteBufferLength, void *ctx) {
    unsigned char message[14];
    int cursor = 0;
    context *currentCtx = (context*)ctx;

    int magnitudeBufferLength = byteBufferLength + currentCtx->remainingMagnitudeLength;

    uint16_t* magnitudeBuffer = (uint16_t*)malloc(magnitudeBufferLength);

    if ((currentCtx->remainingMagnitudeLength>0) && (currentCtx->remainingMagnitudeData != 0)) {
        memcpy(magnitudeBuffer, currentCtx->remainingMagnitudeData, currentCtx->remainingMagnitudeLength);

        cursor = currentCtx->remainingMagnitudeLength;
    }

    // computes magnitudes
    for(int idx = 0; idx<byteBufferLength/2; idx++) {
        int i = byteBuffer[idx*2];
        int q = byteBuffer[idx*2+1];

        // printf("%02x/%02x  |   ", i, q);

        if (i>127) {
            i = i - 127;
        } else {
            i = 127 - i;
        }

        if (q>127) {
            q = q - 127;
        } else {
            q = 127 - q;
        }

        // printf("%02x/%02x => %04x  ", i, q, magnitude[i*129+q]);
        // printValue(magnitude[i*129+q]);

        magnitudeBuffer[idx+cursor] = magnitude[i*129+q];
    }


    // Signature detection:
    //       |   |         |   |
    //       |   |         |   |
    //       |   |         |   |
    //       |   |         |   |
    //       | | | | | | | | | | | | | | | |
    //       0 1 2 3 4 5 6 7 8 9 10
    int idx = 0;
    for(int idx = 0; idx<magnitudeBufferLength - REMAINING_MAGNITUDE_BUFFER_SIZE; idx++)  {
        printf("[foo] magnitude %04x  ", magnitudeBuffer[idx]);
        printValue(magnitudeBuffer[idx]);

        if (magnitudeBuffer[idx] <= magnitudeBuffer[idx+1]) {
            continue;
        }

        if (magnitudeBuffer[idx+1] >= magnitudeBuffer[idx+2]) {
            continue;
        }

        if (magnitudeBuffer[idx+2] <= magnitudeBuffer[idx+3]) {
            continue;
        }
        if (magnitudeBuffer[idx+2] <= magnitudeBuffer[idx+4]) {
            continue;
        }
        if (magnitudeBuffer[idx+2] <= magnitudeBuffer[idx+5]) {
            continue;
        }
        if (magnitudeBuffer[idx+2] <= magnitudeBuffer[idx+6]) {
            continue;
        }

        if (magnitudeBuffer[idx+6] >= magnitudeBuffer[idx+7]) {
            continue;
        }

        if (magnitudeBuffer[idx+7] <= magnitudeBuffer[idx+8]) {
            continue;
        }

        if (magnitudeBuffer[idx+8] >= magnitudeBuffer[idx+9]) {
            continue;
        }

        if (magnitudeBuffer[idx+9] <= magnitudeBuffer[idx+10]) {
            continue;
        }
        if (magnitudeBuffer[idx+9] <= magnitudeBuffer[idx+11]) {
            continue;
        }
        if (magnitudeBuffer[idx+9] <= magnitudeBuffer[idx+12]) {
            continue;
        }
        if (magnitudeBuffer[idx+9] <= magnitudeBuffer[idx+13]) {
            continue;
        }
        if (magnitudeBuffer[idx+9] <= magnitudeBuffer[idx+14]) {
            continue;
        }
        if (magnitudeBuffer[idx+9] <= magnitudeBuffer[idx+15]) {
            continue;
        }

        uint16_t meanHigh = (uint16_t)(
            (
                (uint32_t)(magnitudeBuffer[idx]) + 
                (uint32_t)(magnitudeBuffer[idx + 2]) + 
                (uint32_t)(magnitudeBuffer[idx + 7]) + 
                (uint32_t)(magnitudeBuffer[idx + 9])
            ) / 4
        );

        if (
            (magnitudeBuffer[idx]/meanHigh>2) || 
            (magnitudeBuffer[idx+2]/meanHigh>2) || 
            (magnitudeBuffer[idx+2]/meanHigh>7) || 
            (magnitudeBuffer[idx+9]/meanHigh>2)
            ) {
            continue;
        }

        printf("[foo] _____________________________________________________________________________________________ good preambule ________________________________________________________________________________________\n");

        for(int k=0; k<16; k++) {
            printValue(magnitudeBuffer[idx+k]);
        }
        printf("\n\n");

        memset(message, 0, 14);

        int startOfMessage = idx+16;
        int messageLength = 56;

        for(int index=0; index<messageLength; index++) {
            int byteIndex = index / 8;
            int bitIndex = index % 8;

            if ((index==0) && (magnitudeBuffer[startOfMessage+index*2] > magnitudeBuffer[startOfMessage+index*2+1])) {
                messageLength = 112;
            }

            unsigned char bit = (magnitudeBuffer[startOfMessage+index*2] > magnitudeBuffer[startOfMessage+index*2+1]);

            message[byteIndex] |= bit << bitIndex;
        }

        goRtlsrdData(message, messageLength / 8, ctx);
    }

    printf("Copying data from %d, size %d\n", magnitudeBufferLength - REMAINING_MAGNITUDE_BUFFER_SIZE, REMAINING_MAGNITUDE_BUFFER_SIZE);

    // Copy remaining data in the context.
    if (currentCtx->remainingMagnitudeData != 0) {
        currentCtx->remainingMagnitudeLength = REMAINING_MAGNITUDE_BUFFER_SIZE;
        memcpy(currentCtx->remainingMagnitudeData, &magnitudeBuffer[magnitudeBufferLength - REMAINING_MAGNITUDE_BUFFER_SIZE], REMAINING_MAGNITUDE_BUFFER_SIZE);
    }

    free(magnitudeBuffer);
}

context *newContext(void* goContext) {
    context *output = (context*)malloc(sizeof(context));
    output->goContext = goContext;
    output->remainingMagnitudeData = (uint16_t*)malloc(REMAINING_MAGNITUDE_BUFFER_SIZE);
    output->remainingMagnitudeLength = 0;

    return output;
}


int rtlsdrReadAsync(rtlsdr_dev_t *dev, void *ctx, uint32_t buf_num, uint32_t buf_len) {
    return rtlsdr_read_async(dev, rtlsdrProcessRaw, ctx, buf_num, buf_len);
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
