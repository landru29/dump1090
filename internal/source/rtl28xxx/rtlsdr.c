#include "rtlsdr.h"
#include <math.h>
#include <malloc.h>
#include <string.h>
#include <stdio.h>
#include "constant.h"
#include "../context/context.h"

#define RAW 0
#define DEBUG 1


int messageLengthBit[25];
uint16_t magnitude[129*129];
int globalIndex = 0;


void printValue(uint16_t val) {
    if (DEBUG) {
        printf("[");
        for(int l = 0; l<val/512; l++) {
            printf("=");
        }
        printf(">\n");

        fflush(stdout);
    }
}

void printRawValue(uint16_t val) {
    printf("%c%c", (val>>8)&0xff, val&0xff);
}

uint16_t* computeMagnitudes(unsigned char *byteBuffer, uint32_t byteBufferLength, void *ctx, uint32_t *size)  {
    int startIdx = 0;
    context *currentCtx = (context*)ctx;

    int magnitudeBufferLengthByte = (byteBufferLength * sizeof(uint16_t) / IQ_SIZE) + currentCtx->remainingMagnitudeLengthByte;

    uint16_t* magnitudeBuffer = (uint16_t*)malloc(magnitudeBufferLengthByte);

    if ((currentCtx->remainingMagnitudeLengthByte>0) && (currentCtx->remainingMagnitudeData != 0)) {
        if ((!RAW) && (DEBUG)) {
            printf("copying %d bytes\n",currentCtx->remainingMagnitudeLengthByte);
        }
        memcpy(magnitudeBuffer, currentCtx->remainingMagnitudeData, currentCtx->remainingMagnitudeLengthByte);

        startIdx = currentCtx->remainingMagnitudeLengthByte / sizeof(uint16_t);
    }

    // computes magnitudes.
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

    *size = magnitudeBufferLengthByte / sizeof(uint16_t);

    return magnitudeBuffer;
}

void rtlsdrProcessRaw(unsigned char *byteBuffer, uint32_t byteBufferLength, void *ctx) {
    unsigned char message[14];
    uint32_t magnitudeCount;
    context *currentCtx = (context*)ctx;

    uint16_t *magnitudeBuffer = computeMagnitudes(byteBuffer, byteBufferLength, ctx, &magnitudeCount);

    uint32_t magnitudeBufferLengthByte = magnitudeCount * 2;

    int idx;

    int limitProcess = magnitudeCount - MAGNITUDE_LONG_MSG_SIZE - PREAMBULE_BIT_SIZE;

    // Signature detection:
    //       |   |         |   |
    //       |   |         |   |
    //       |   |         |   |
    //       |   |         |   |
    //       | | | | | | | | | | | | | | | |
    //       0 1 2 3 4 5 6 7 8 9 10
    for(idx = 0; idx<limitProcess; idx++)  {
        if (RAW) {
            printRawValue(magnitudeBuffer[idx]);
        } else if (DEBUG) {
            printf("%04d (%04d / %04d) magnitude %04x  ", globalIndex + idx, idx, magnitudeCount, magnitudeBuffer[idx]);
            printValue(magnitudeBuffer[idx]);
        }

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
            (magnitudeBuffer[idx+7]/meanHigh>2) || 
            (magnitudeBuffer[idx+9]/meanHigh>2)
            ) {
            continue;
        }

        if ((!RAW) && (DEBUG)) {
            printf("%04d [foo] _____________________________________________________________________________________________ good preambule ________________________________________________________________________________________\n", globalIndex+idx);
  
            for(int k=0; k<16; k++) {
                printf("    %04d (%04d / %04d) magnitude %04x  ", globalIndex +idx + k, idx + k, magnitudeCount, magnitudeBuffer[idx+k]);
                printValue(magnitudeBuffer[idx+k]);
            }

            printf("\n\n");
            fflush(stdout);
        }

        // The preambule seems to be right, prepare the message variable.
        int messageLengthBit = decodeMessage(&magnitudeBuffer[idx + PREAMBULE_BIT_SIZE], message); 

        if ((RAW) || (DEBUG)) {
            for (int j=0; j<messageLengthBit/8; j++) {
                printf("%02X", message[j]);
            }
            printf("\n");
        }

        // No error ?
        if (goRtlsrdData(message, messageLengthBit / 8, ctx) == 0) {
            // jump over the message.
            idx += PREAMBULE_BIT_SIZE + messageLengthBit * 2;
            if ((RAW) || (DEBUG)) printf("Jumping to %04d (%04d + %04d = %04d)\n", idx, PREAMBULE_BIT_SIZE, messageLengthBit * 2, PREAMBULE_BIT_SIZE + messageLengthBit * 2);
        }

        // No error ?
        // if (goRtlsrdData(message, messageLengthBit / 8, ctx) == 0) {
        //     // jump over the message.
        //     idx += PREAMBULE_BIT_SIZE + messageLengthBit * 2;
        //     if ((RAW) || (DEBUG)) printf("Jumping to %04d\n", idx);
        // } else {
            // idx += PREAMBULE_BIT_SIZE + messageLengthBit * 2;
            // if ((RAW) || (DEBUG)) printf("Jumping to %04d (%04d + %04d = %04d)\n", idx, PREAMBULE_BIT_SIZE, messageLengthBit * 2, PREAMBULE_BIT_SIZE + messageLengthBit * 2);
        // }
    }

    // Copy remaining data in the context.
    if (currentCtx->remainingMagnitudeData != 0) {
        currentCtx->remainingMagnitudeLengthByte = sizeof(uint16_t) * (magnitudeCount - idx);
        printf("remaining %04d\n", currentCtx->remainingMagnitudeLengthByte);
        if ((!RAW) && (DEBUG)) {
            printf("Copying data from %ld, size %d\n", idx * sizeof(uint16_t), currentCtx->remainingMagnitudeLengthByte);
            fflush(stdout);
        }
        memcpy(currentCtx->remainingMagnitudeData, &magnitudeBuffer[idx], currentCtx->remainingMagnitudeLengthByte);

        globalIndex += idx;
    }

    free(magnitudeBuffer);
}

int decodeMessage(uint16_t* magnitudeBuffer, char * message) {
    memset(message, 0, 14); 

    int messageLengthBit = MODES_SHORT_MSG_BITS;

    // +----------+--------------+-----------+
    // |  DF (5)  | (83) or (27) |  PI (24)  |
    // +----------+--------------+-----------+

    for(int index=0; index<messageLengthBit; index++) {
        int byteIndex = index / 8;
        int bitIndex = index % 8;

        unsigned char bit = (magnitudeBuffer[index*2] > magnitudeBuffer[index*2+1]);

        // If the first bit of DF is 1, this means the message will be long 112 bits (extended squitter),
        // otherwise, the message will be short 56 bits (normal squitter).
        if ((index==0) && (bit==1)) {
            messageLengthBit = MODES_LONG_MSG_BITS;
        }

        message[byteIndex] |= bit << (7 - bitIndex);
    }

    return messageLengthBit;
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
