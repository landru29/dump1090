#include "rtlsdr.h"
#include <math.h>
#include <malloc.h>
#include <string.h>
#include <stdio.h>

// ADSB message is 112 bits
// Each bit is encoded on 2 magnitudes (uint16)
// Each magnitude is encoded on 2 bytes.

// ADSB message is 448 bytes.

// I: 1 byte, Q: 1 byte
#define IQ_SIZE                         2

// 10 => 1, 01 => 0
#define MAGNITUDE_ENCODED_BIT_SIZE      2

// Number of all magnitudes for a long message
#define MAGNITUDE_LONG_MSG_SIZE         MODES_LONG_MSG_BITS * MAGNITUDE_ENCODED_BIT_SIZE

// Size in memory of all magnitudes for a long message (each magnitude is uint16)
#define MAGNITUDE_LONG_MSG_BYTE_SIZE    MAGNITUDE_LONG_MSG_SIZE * 2

// Number of bytes for a long message
#define IQ_LONG_MSG_SIZE                MAGNITUDE_LONG_MSG_SIZE * IQ_SIZE

int messageLengthBit[25];
uint16_t magnitude[129*129];


void printValue(uint16_t val) {
    printf("[");
    for(int l = 0; l<val/512; l++) {
        printf("=");
    }
    printf(">\n");
}

uint16_t* computeMagnitudes(unsigned char *byteBuffer, uint32_t byteBufferLength, void *ctx, uint32_t *size)  {
    printf("#####################################\n");
    int startIdx = 0;
    context *currentCtx = (context*)ctx;

    int magnitudeBufferLengthByte = (byteBufferLength * sizeof(uint16_t) / IQ_SIZE) + currentCtx->remainingMagnitudeLengthByte;

    uint16_t* magnitudeBuffer = (uint16_t*)malloc(magnitudeBufferLengthByte);

    if ((currentCtx->remainingMagnitudeLengthByte>0) && (currentCtx->remainingMagnitudeData != 0)) {
        printf("copying %d bytes\n",currentCtx->remainingMagnitudeLengthByte);
        memcpy(magnitudeBuffer, currentCtx->remainingMagnitudeData, currentCtx->remainingMagnitudeLengthByte);

        startIdx = currentCtx->remainingMagnitudeLengthByte / sizeof(uint16_t);
    }

    // computes magnitudes
    for(int idx = 0; idx<byteBufferLength/IQ_SIZE; idx++) {
        int i = byteBuffer[idx*IQ_SIZE];
        int q = byteBuffer[idx*IQ_SIZE+1];

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

        magnitudeBuffer[idx+startIdx] = magnitude[i*129+q];
    }

    *size = magnitudeBufferLengthByte/2;

    return magnitudeBuffer;
}

void rtlsdrProcessRaw(unsigned char *byteBuffer, uint32_t byteBufferLength, void *ctx) {
    unsigned char message[14];
    uint32_t magnitudeCount;
    context *currentCtx = (context*)ctx;

    uint16_t *magnitudeBuffer = computeMagnitudes(byteBuffer, byteBufferLength, ctx, &magnitudeCount);

    uint32_t magnitudeBufferLengthByte = magnitudeCount * 2;

    // Signature detection:
    //       |   |         |   |
    //       |   |         |   |
    //       |   |         |   |
    //       |   |         |   |
    //       | | | | | | | | | | | | | | | |
    //       0 1 2 3 4 5 6 7 8 9 10
    for(int idx = 0; idx<magnitudeCount - MAGNITUDE_LONG_MSG_SIZE; idx++)  {
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

        int startOfMessage = idx+16; // skip the preambule of 8µs
        int messageLength = MODES_SHORT_MSG_BITS;

        // +----------+--------------+-----------+
        // |  DF (5)  | (83) or (27) |  PI (24)  |
        // +----------+--------------+-----------+

        for(int index=0; index<messageLength; index++) {
            int byteIndex = index / 8;
            int bitIndex = index % 8;

            unsigned char bit = (magnitudeBuffer[startOfMessage+index*2] > magnitudeBuffer[startOfMessage+index*2+1]);

            if ((index==0) && (bit==1)) {
                messageLength = MODES_LONG_MSG_BITS;
                // If the first bit of DF is 1, this means the message will be long 112 bits
                // otherwise, the message will be short 56 bits
            }

            message[byteIndex] |= bit << bitIndex;
        }

        goRtlsrdData(message, messageLength / 8, ctx);
    }

    

    // Copy remaining data in the context.
    if (currentCtx->remainingMagnitudeData != 0) {
        currentCtx->remainingMagnitudeLengthByte = MAGNITUDE_LONG_MSG_BYTE_SIZE;
        printf("Copying data from %ld, size %d\n", (magnitudeCount - MAGNITUDE_LONG_MSG_SIZE) * sizeof(uint16_t), currentCtx->remainingMagnitudeLengthByte);
        memcpy(currentCtx->remainingMagnitudeData, &magnitudeBuffer[magnitudeCount - MAGNITUDE_LONG_MSG_SIZE], MAGNITUDE_LONG_MSG_BYTE_SIZE);
    }

    free(magnitudeBuffer);
}

context *newContext(void* goContext) {
    context *output = (context*)malloc(sizeof(context));
    output->goContext = goContext;
    output->remainingMagnitudeData = (uint16_t*)malloc(MAGNITUDE_LONG_MSG_SIZE * sizeof(uint16_t));
    output->remainingMagnitudeLengthByte = 0;

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
