#ifndef __CONSTANT_H__
#define __CONSTANT_H__

#define MODES_LONG_MSG_BITS 112
#define MODES_SHORT_MSG_BITS 56

#define MODES_LONG_MSG_BYTES (MODES_LONG_MSG_BITS/8)
#define MODES_SHORT_MSG_BYTES (MODES_SHORT_MSG_BITS/8)


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

#endif